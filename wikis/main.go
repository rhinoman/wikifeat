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

// The Wiki service
package main

import (
	"flag"
	"github.com/emicklei/go-restful"
	"github.com/rhinoman/wikifeat/common/config"
	"github.com/rhinoman/wikifeat/common/registry"
	"github.com/rhinoman/wikifeat/common/services"
	"github.com/rhinoman/wikifeat/wikis/wiki_service"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"net/http"
)

func main() {
	// Get command line arguments
	configFile := flag.String("config", "config.ini", "config file to load")
	flag.Parse()
	// Load Configuration
	config.LoadConfig(*configFile)
	// Set up Logger
	log.SetOutput(&lumberjack.Logger{
		Filename:   config.Logger.LogFile,
		MaxSize:    config.Logger.MaxSize,
		MaxBackups: config.Logger.MaxBackups,
		MaxAge:     config.Logger.MaxAge,
	})
	wsContainer := restful.NewContainer()
	wsContainer.Router(restful.CurlyRouter{})
	//Enable Gzip
	wsContainer.EnableContentEncoding(true)
	wc := wiki_service.WikisController{}
	wc.Register(wsContainer)
	services.InitDb()
	registry.Init("Wikis", registry.WikisLocation)
	httpAddr := ":" + config.Service.Port
	server := &http.Server{Addr: httpAddr, Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
