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

// The User Service
package main

import (
	"flag"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/emicklei/go-restful"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/gopkg.in/natefinch/lumberjack.v2"
	"github.com/rhinoman/wikifeat/common/config"
	"github.com/rhinoman/wikifeat/common/registry"
	"github.com/rhinoman/wikifeat/common/services"
	"github.com/rhinoman/wikifeat/users/user_service"
	"log"
	"net/http"
)

func main() {
	// Get command line arguments
	configFile := flag.String("config", "config.ini", "config file to load")
	flag.Parse()
	// Load configuration
	config.LoadConfig(*configFile)
	// Set up logger
	log.SetOutput(&lumberjack.Logger{
		Filename:   config.Logger.LogFile,
		MaxSize:    config.Logger.MaxSize,
		MaxBackups: config.Logger.MaxBackups,
		MaxAge:     config.Logger.MaxAge,
	})
	wsContainer := restful.NewContainer()
	wsContainer.Router(restful.CurlyRouter{})
	//Enable Gzip support
	wsContainer.EnableContentEncoding(true)
	uc := user_service.UsersController{}
	uc.Register(wsContainer)
	services.InitDb()
	registry.Init("Users", registry.UsersLocation)
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
