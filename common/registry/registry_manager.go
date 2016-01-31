/*
 *  Licensed to Wikifeat under one or more contributor license agreements.
 *  See the LICENSE.txt file distributed with this work for additional information
 *  regarding copyright ownership.
 *
 *  Redistribution and use in source and binary forms, with or without
 *  modification, are permitted provided that the following conditions are met:
 *
 *  * Redistributions of source code must retain the above copyright notice,
 *  this list of conditions and the following disclaimer.
 *  * Redistributions in binary form must reproduce the above copyright
 *  notice, this list of conditions and the following disclaimer in the
 *  documentation and/or other materials provided with the distribution.
 *  * Neither the name of Wikifeat nor the names of its contributors may be used
 *  to endorse or promote products derived from this software without
 *  specific prior written permission.
 *
 *  THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 *  AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 *  IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 *  ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
 *  LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 *  CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 *  SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 *  INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 *  CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 *  ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 *  POSSIBILITY OF SUCH DAMAGE.
 */

// Manages registry cache and such
package registry

import (
	"errors"
	"fmt"
	etcd "github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/coreos/etcd/client"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/rhinoman/wikifeat/common/config"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"
)

type nodeCache struct {
	sync.RWMutex
	m map[string][]*etcd.Node
}

var serviceCache = nodeCache{m: make(map[string][]*etcd.Node)}
var pluginCache = nodeCache{m: make(map[string][]*etcd.Node)}

var random = rand.New(rand.NewSource(time.Now().Unix()))

var client *etcd.Client
var kapi etcd.KeysAPI
var EtcdPrefix = "/wikifeat"
var UsersLocation = EtcdPrefix + "/services/users/"
var WikisLocation = EtcdPrefix + "/services/wikis/"
var NotificationsLocation = EtcdPrefix + "/services/notifications/"
var FrontEndLocation = EtcdPrefix + "/services/frontend/"
var AuthLocation = EtcdPrefix + "/services/auth/"
var PluginsLocation = EtcdPrefix + "/plugin/"
var ttl time.Duration

func protocolString() string {
	if config.Service.UseSSL {
		return "https://"
	} else {
		return "http://"
	}
}

func GetEtcdKeyAPI() etcd.KeysAPI {
	return kapi
}

func hostUrl() string {
	return protocolString() + config.Service.DomainName +
		":" + config.Service.Port
}

func Init(serviceName, registryLocation string) error {
	log.Print("Initializing registry connection.")
	nodeId := config.Service.NodeId
	ttl = time.Duration(config.ServiceRegistry.EntryTTL) * time.Second
	cfg := etcd.Config{
		Endpoints: []string{config.Service.RegistryLocation},
		Transport: etcd.DefaultTransport,
	}
	client, err := etcd.New(cfg)
	if err != nil {
		log.Fatal(err)
		return err
	}
	kapi = etcd.NewKeysAPI(client)
	log.Print("Registering " + serviceName + " Service node at " + hostUrl())
	if _, err := kapi.Set(context.Background(), registryLocation+nodeId, hostUrl(),
		&etcd.SetOptions{TTL: ttl}); err != nil {
		fmt.Println(err)
		log.Fatal(err)
		return err
	}
	fetchServiceLists()
	fetchPluginNodes()
	go sendHeartbeat(registryLocation)
	go updateServiceCache()
	return nil
}

func sendHeartbeat(registryLocation string) {
	nodeId := config.Service.NodeId
	for {
		time.Sleep(time.Duration(config.ServiceRegistry.EntryTTL/2) * time.Second)
		if _, err := kapi.Set(context.Background(), registryLocation+nodeId,
			hostUrl(), &etcd.SetOptions{TTL: ttl}); err != nil {
			log.Print("Can't send Heartbeat to registry! - " + err.Error())
		}
	}
}

func updateServiceCache() {
	cri := config.ServiceRegistry.CacheRefreshInterval
	for {
		time.Sleep(time.Duration(cri) * time.Second)
		fetchServiceLists()
		fetchPluginNodes()
	}
}

func getServiceNodes(serviceLocation string) ([]*etcd.Node, error) {
	ctx, _ := context.WithTimeout(context.Background(), 7*time.Second)
	if resp, err := kapi.Get(ctx, serviceLocation,
		&etcd.GetOptions{Recursive: true}); err != nil {
		return nil, err
	} else {
		return processResponse(resp)
	}

}

//Reads 'plugin' services from the service registry
func fetchPluginNodes() {
	ppn, err := getServiceNodes(PluginsLocation)
	if err != nil {
		//Need to avoid spamming the logs if you have no plugins installed
		//log.Println("Error fetching plugin nodes: " + err.Error())
	}
	for _, node := range ppn {
		if node.Dir {
			splitKey := strings.Split(node.Key, "/")
			pluginName := splitKey[len(splitKey)-1]
			pluginCache.m[pluginName] = node.Nodes
		}
	}
}

// Loads the latest services from Etcd
func fetchServiceLists() {
	// First, fetch the core services
	userNodes, err := getServiceNodes(UsersLocation)
	if err != nil {
		log.Println("Error fetching user services: " + err.Error())
	}
	wikiNodes, err := getServiceNodes(WikisLocation)
	if err != nil {
		log.Println("Error fetching wiki services: " + err.Error())
	}
	notificationNodes, err := getServiceNodes(NotificationsLocation)
	if err != nil {
		log.Println("Error fetching notificaiton services: " + err.Error())
	}
	frontendNodes, err := getServiceNodes(FrontEndLocation)
	if err != nil {
		log.Println("Error fetching frontend services: " + err.Error())
	}
	authNodes, err := getServiceNodes(AuthLocation)
	if err != nil {
		log.Println("Error fetching auth services: " + err.Error())
	}
	serviceCache.Lock()
	defer serviceCache.Unlock()
	serviceCache.m["users"] = userNodes
	serviceCache.m["wikis"] = wikiNodes
	serviceCache.m["notifications"] = notificationNodes
	serviceCache.m["frontend"] = frontendNodes
	serviceCache.m["auth"] = authNodes
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
	if max := len(serviceCache.m[serviceName]); max == 0 {
		return "", errors.New("No " + serviceName + " services listed!")
	} else {
		index := 0
		if max > 1 {
			index = random.Intn(max)
		}
		return getEndpointFromNode(serviceCache.m[serviceName][index]), nil
	}
	return "", nil
}

//Get a plugin node for use
func GetPluginLocation(pluginName string) (string, error) {
	pluginCache.RLock()
	defer pluginCache.RUnlock()
	if max := len(pluginCache.m[pluginName]); max == 0 {
		return "", errors.New("No " + pluginName + " plugins listed!")
	} else {
		index := 0
		if max > 1 {
			index = random.Intn(max)
		}
		return getEndpointFromNode(pluginCache.m[pluginName][index]), nil
	}
	return "", nil
}
