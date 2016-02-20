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

package config

import (
	etcd "github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/coreos/etcd/client"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/golang.org/x/net/context"
	"log"
	"reflect"
	"strconv"
)

//Config Locations in Etcd
var ConfigPrefix = "/wikifeat/config/"
var DbConfigLocation = ConfigPrefix + "db/"
var LogConfigLocation = ConfigPrefix + "log/"
var AuthConfigLocation = ConfigPrefix + "auth/"
var NotificationsConfigLocation = ConfigPrefix + "notifications/"
var UsersConfigLocation = ConfigPrefix + "users/"
var FrontendConfigLocation = ConfigPrefix + "frontend/"
var RegistryConfigLocation = ConfigPrefix + "registry/"

// Fetch configuration from etcd
// Because golang has no generics, this makes heavy use of reflection :|
func FetchConfig() {
	log.Printf("\nFetching Configuration from %v\n", Service.RegistryLocation)
	// Get an etcd Client
	etcdCfg := etcd.Config{
		Endpoints: []string{Service.RegistryLocation},
		Transport: etcd.DefaultTransport,
	}
	etcdClient, err := etcd.New(etcdCfg)
	if err != nil {
		log.Fatal(err)
		return
	}
	kapi := etcd.NewKeysAPI(etcdClient)

	// Fetch config sections from etcd
	fetchConfigSection(&Database, DbConfigLocation, kapi)
	fetchConfigSection(&Logger, LogConfigLocation, kapi)
	fetchConfigSection(&ServiceRegistry, RegistryConfigLocation, kapi)
	fetchConfigSection(&Auth, AuthConfigLocation, kapi)
	fetchConfigSection(&Notifications, NotificationsConfigLocation, kapi)
	fetchConfigSection(&Users, UsersConfigLocation, kapi)
	fetchConfigSection(&Frontend, FrontendConfigLocation, kapi)
}

func setConfigVal(str string, field reflect.Value) error {
	t := field.Kind()
	switch {
	case t == reflect.String:
		log.Printf("STR: %v\n", str)
		field.SetString(str)
	case t >= reflect.Int && t <= reflect.Int64:
		if x, err := strconv.ParseInt(str, 10, 64); err != nil {
			return err
		} else {
			field.SetInt(x)
		}
	case t >= reflect.Uint && t <= reflect.Uint64:
		if x, err := strconv.ParseUint(str, 10, 64); err != nil {
			return err
		} else {
			field.SetUint(x)
		}
	case t >= reflect.Float32 && t <= reflect.Float64:
		if x, err := strconv.ParseFloat(str, 64); err != nil {
			return err
		} else {
			field.SetFloat(x)
		}
	case t == reflect.Bool:
		if x, err := strconv.ParseBool(str); err != nil {
			return err
		} else {
			field.SetBool(x)
		}
	default:
		return nil
	}
	return nil
}

//Fetches a single config section
func fetchConfigSection(configStruct interface{}, location string, kapi etcd.KeysAPI) {
	cfg := reflect.ValueOf(configStruct).Elem()
	for i := 0; i < cfg.NumField(); i++ {
		key := cfg.Type().Field(i).Name
		resp, getErr := kapi.Get(context.Background(), location+key, nil)
		if getErr != nil {
			log.Printf("Error getting key %v: %v\n", key, getErr)
			continue
		}
		valErr := setConfigVal(resp.Node.Value, cfg.Field(i))
		if valErr != nil {
			log.Printf("Error setting config field %v: %v\n", key, valErr)
		}
	}
}
