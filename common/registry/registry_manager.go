/**  Copyright (c) 2014-present James Adam.  All rights reserved.
*
*		 This file is part of Wikifeat.
*
*    Wikifeat is free software: you can redistribute it and/or modify
*    it under the terms of the GNU General Public License as published by
*    the Free Software Foundation, either version 2 of the License, or
*    (at your option) any later version.
*
*    This program is distributed in the hope that it will be useful,
*    but WITHOUT ANY WARRANTY; without even the implied warranty of
*    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
*    GNU General Public License for more details.
*
*    You should have received a copy of the GNU General Public License
*    along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

// Manages registry cache and such
package registry

import (
	"errors"
	"fmt"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/coreos/go-etcd/etcd"
	"github.com/rhinoman/wikifeat/common/config"
	"log"
	"math/rand"
	"sync"
	"time"
)

var serviceCache = struct {
	sync.RWMutex
	m map[string][]*etcd.Node
}{m: make(map[string][]*etcd.Node)}

var random = rand.New(rand.NewSource(time.Now().Unix()))

var client *etcd.Client
var etcdPrefix = "/wikifeat"
var UsersLocation = etcdPrefix + "/services/users/"
var WikisLocation = etcdPrefix + "/services/wikis/"
var FrontEndLocation = etcdPrefix + "/services/frontend/"

func protocolString() string {
	if config.Service.UseSSL {
		return "https://"
	} else {
		return "http://"
	}
}

func hostUrl() string {
	return protocolString() + config.Service.DomainName +
		":" + config.Service.Port
}

func Init(serviceName, registryLocation string) error {
	log.Print("Initializing registry connection.")
	machines := []string{config.Service.RegistryLocation}
	nodeId := config.Service.NodeId
	ttl := config.ServiceRegistry.EntryTTL
	client = etcd.NewClient(machines)
	log.Print("Registering " + serviceName + " Service node at " + hostUrl())
	if _, err := client.Set(registryLocation+nodeId, hostUrl(), ttl); err != nil {
		fmt.Println(err)
		log.Fatal(err)
		return err
	}
	fetchServiceLists()
	go sendHeartbeat(registryLocation)
	go updateServiceCache()
	return nil
}

func sendHeartbeat(registryLocation string) {
	ttl := config.ServiceRegistry.EntryTTL
	nodeId := config.Service.NodeId
	for {
		time.Sleep(time.Duration(ttl/2) * time.Second)
		log.Print(nodeId + " Sending heartbeat...")
		if _, err := client.Set(registryLocation+nodeId, hostUrl(), ttl); err != nil {
			log.Print("Can't send Heartbeat to registry! - " + err.Error())
		}
	}
}

func updateServiceCache() {
	cri := config.ServiceRegistry.CacheRefreshInterval
	for {
		time.Sleep(time.Duration(cri) * time.Second)
		fetchServiceLists()
	}
}

// Loads the latest services from Etcd
func fetchServiceLists() {
	// First, fetch the core services
	userNodes := []*etcd.Node{}
	resp, err := client.Get(UsersLocation, false, false)
	if err != nil {
		log.Println("Unable to fetch Users Service list from etcd!")
	} else {
		if nodes, err := processResponse(resp); err != nil {
			log.Println("Error fetching user services: " + err.Error())
		} else {
			userNodes = nodes
		}
	}
	wikiNodes := []*etcd.Node{}
	resp, err = client.Get(WikisLocation, false, false)
	if err != nil {
		log.Println("Unable to fetch Wikis Service list from etcd!")
	} else {
		if nodes, err := processResponse(resp); err != nil {
			log.Println("Error fetching wiki services: " + err.Error())
		} else {
			wikiNodes = nodes
		}
	}
	serviceCache.Lock()
	defer serviceCache.Unlock()
	serviceCache.m["users"] = userNodes
	serviceCache.m["wikis"] = wikiNodes
}

//Read nodes from an etcd response
func processResponse(response *etcd.Response) ([]*etcd.Node, error) {
	rootNode := response.Node
	if !rootNode.Dir {
		return nil, errors.New("Not a directory!")
	}
	if len(rootNode.Nodes) == 0 {
		return nil, errors.New("No listed services!")
	}
	return rootNode.Nodes, nil
}

func getEndpointFromNode(node *etcd.Node) string {
	return node.Value
}

//Get a service node for use
func GetServiceLocation(serviceName string) (string, error) {
	serviceCache.RLock()
	defer serviceCache.RUnlock()
	switch serviceName {
	case "users":
		if max := len(serviceCache.m["users"]); max == 0 {
			return "", errors.New("No User services listed!")
		} else {
			index := 0
			if max > 1 {
				index = random.Intn(max - 1)
			}
			return getEndpointFromNode(serviceCache.m["users"][index]), nil
		}
	case "wikis":
		if max := len(serviceCache.m["wikis"]); max == 0 {
			return "", errors.New("No Wiki services listed!")
		} else {
			index := 0
			if max > 1 {
				index = random.Intn(max - 1)
			}
			return getEndpointFromNode(serviceCache.m["wikis"][index]), nil
		}
	}
	return "", nil
}
