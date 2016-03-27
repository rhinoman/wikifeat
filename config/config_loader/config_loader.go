package config_loader

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
/**
 * Load the configuration from file into etcd
 */

import (
	etcd "github.com/coreos/etcd/client"
	"github.com/rhinoman/wikifeat/common/config"
	. "github.com/rhinoman/wikifeat/common/database"
	"golang.org/x/net/context"
	"log"
	"reflect"
	"strconv"
)

var kapi etcd.KeysAPI

//Initialize our etcd connection
func InitRegistry() {
	log.Print("Initializing registry connection.")
	cfg := etcd.Config{
		Endpoints: []string{config.Service.RegistryLocation},
		Transport: etcd.DefaultTransport,
	}
	client, err := etcd.New(cfg)
	if err != nil {
		log.Fatal(err)
		return
	}
	kapi = etcd.NewKeysAPI(client)
}

//Perform some initialization on the CouchDB database
func InitDatabase() {
	InitDb()
	SetupDb()
}

//Load the configuration into etcd
func SetConfig() {
	log.Println("Setting service registry config")
	setConfigItems(config.ServiceRegistry, config.RegistryConfigLocation)
	log.Println("Setting database config")
	setConfigItems(config.Database, config.DbConfigLocation)
	log.Println("Setting logger config")
	setConfigItems(config.Logger, config.LogConfigLocation)
	log.Println("Setting auth config")
	setConfigItems(config.Auth, config.AuthConfigLocation)
	log.Println("Setting notifications config")
	setConfigItems(config.Notifications, config.NotificationsConfigLocation)
	log.Println("Setting users config")
	setConfigItems(config.Users, config.UsersConfigLocation)
	log.Println("Setting frontend config")
	setConfigItems(config.Frontend, config.FrontendConfigLocation)
}

//Clear the configuration in etcd
func ClearConfig() {
	kapi.Delete(context.Background(), config.ConfigPrefix,
		&etcd.DeleteOptions{
			Recursive: true,
		})
}

func setConfigItems(configStruct interface{}, configLocation string) {
	cfg := reflect.ValueOf(configStruct)
	for i := 0; i < cfg.NumField(); i++ {
		key := cfg.Type().Field(i).Name
		entry := cfg.Field(i).Interface()
		cfgVal := entryToString(entry)
		log.Printf("Setting Key: %v, Value: %v", key, cfgVal)
		if err := setConfigEntry(configLocation+key, cfgVal); err != nil {
			log.Printf("Error setting config "+key+": %v", err)
		}
	}

}

func setConfigEntry(key string, value string) error {
	_, err := kapi.Set(context.Background(), key, value, nil)
	return err
}

func entryToString(entry interface{}) string {
	field := reflect.ValueOf(entry)
	kind := field.Kind()
	switch {
	case kind >= reflect.Int && kind <= reflect.Int64:
		return strconv.FormatInt(field.Int(), 10)
	case kind >= reflect.Uint && kind <= reflect.Uint64:
		return strconv.FormatUint(field.Uint(), 10)
	case kind == reflect.Bool:
		return strconv.FormatBool(field.Bool())
	case kind >= reflect.Float32 && kind <= reflect.Float64:
		return strconv.FormatFloat(field.Float(), 'E', -1, 64)
	case kind == reflect.String:
		return field.String()
	default:
		return ""

	}
}
