package main

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
 * Configuration Loader
 * This program loads the wikifeat configuration from a file into etcd
 * It also performs some initialization of supporting services.
 */

import (
	"flag"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/emicklei/go-restful"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/gopkg.in/natefinch/lumberjack.v2"
	"github.com/rhinoman/wikifeat/common/config"
	"github.com/rhinoman/wikifeat/common/registry"
	"github.com/rhinoman/wikifeat/common/util"
	"github.com/rhinoman/wikifeat/config/config_loader"
	"github.com/rhinoman/wikifeat/config/config_service"
	"log"
	"net/http"
)

func parseCmdParams() string {
	defaultConfig, err := util.DefaultConfigLocation()
	if err != nil {
		log.Fatalf("Error setting config file: %v", err)
	}
	hostName := flag.String("hostName", "localhost", "The host name for this instance")
	nodeId := flag.String("nodeId", "cfg1", "The node Id for this instance")
	port := flag.String("port", "4140", "The port number for this instance")
	useSSL := flag.Bool("useSSL", false, "use SSL")
	sslCertFile := flag.String("sslCertFile", "", "The SSL certificate file")
	sslKeyFile := flag.String("sslKeyFile", "", "The SSL key file")
	registryLocation := flag.String("registryLocation", "http://localhost:2379", "URL for etcd")
	configFile := flag.String("config", defaultConfig, "config file to load")
	flag.Parse()
	config.Service.DomainName = *hostName
	config.Service.NodeId = *nodeId
	config.Service.Port = *port
	config.Service.UseSSL = *useSSL
	config.Service.SSLCertFile = *sslCertFile
	config.Service.SSLKeyFile = *sslKeyFile
	config.Service.RegistryLocation = *registryLocation
	return *configFile
}

func main() {
	// Get command line arguments
	configFile := parseCmdParams()
	// Load Configuration
	config.LoadConfig(configFile)
	// Set up Logger
	log.SetOutput(&lumberjack.Logger{
		Filename:   config.Logger.LogFile,
		MaxSize:    config.Logger.MaxSize,
		MaxBackups: config.Logger.MaxBackups,
		MaxAge:     config.Logger.MaxAge,
	})
	// Initialize our etcd and couchdb connections
	config_loader.InitRegistry()
	config_loader.InitDatabase()
	// Clear out any old config that may be hanging around
	config_loader.ClearConfig()
	// Set the configuration keys in etcd
	config_loader.SetConfig()
	//Start up the config service
	wsContainer := restful.NewContainer()
	wsContainer.Router(restful.CurlyRouter{})
	//Enable GZIP support
	wsContainer.EnableContentEncoding(true)
	cc := config_service.ConfigController{}
	cc.Register(wsContainer)
	registry.Init("Config", registry.ConfigServiceLocation)
	httpAddr := ":" + config.Service.Port
	if config.Service.UseSSL == true {
		certFile := config.Service.SSLCertFile
		keyFile := config.Service.SSLKeyFile
		log.Fatal(http.ListenAndServeTLS(httpAddr,
			certFile, keyFile, wsContainer))
	} else {
		log.Fatal(http.ListenAndServe(httpAddr, wsContainer))
	}
}
