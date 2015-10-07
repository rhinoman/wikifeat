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

// Main Entry point
package main

import (
	"flag"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/wikifeat-common/config"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/wikifeat-common/registry"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/wikifeat-common/services"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/gopkg.in/natefinch/lumberjack.v2"
	"github.com/rhinoman/wikifeat/frontend/fserv"
	"github.com/rhinoman/wikifeat/frontend/routing"
	"log"
)

func main() {
	// Get Command line arguments
	configFile := flag.String("config", "config.ini", "config file to load")
	flag.Parse()
	// Load configuration
	config.LoadConfig(*configFile)
	fserv.LoadPluginData(config.Frontend.PluginDir + "/plugins.ini")
	services.InitDb()
	// Set up the core logger
	log.SetOutput(&lumberjack.Logger{
		Filename:   config.Logger.LogFile,
		MaxSize:    config.Logger.MaxSize,
		MaxBackups: config.Logger.MaxBackups,
		MaxAge:     config.Logger.MaxAge,
	})
	registry.Init("Frontend", registry.FrontEndLocation)
	routing.Start()
}
