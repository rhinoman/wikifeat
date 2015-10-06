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

package routing

import (
	"encoding/json"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/daaku/go.httpgzip"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/gorilla/mux"
	"github.com/rhinoman/wikifeat/common/registry"
	"github.com/rhinoman/wikifeat/frontend/fserv"
	"log"
	"net/http"
)

type PluginListResponse struct {
	EnabledPlugins []fserv.PluginData `json:"enabledPlugins"`
}

// Plugin "Resources" -- static files and such
func handlePluginRoutes(pr *mux.Router) {
	pr.StrictSlash(true).HandleFunc("/", getPluginList).Methods("GET")
	pr.PathPrefix("/{plugin-name}/resource/").
		Methods("GET").
		HandlerFunc(servePluginResource)
}

// Backend plugin routes
func handlePluginBackendRoutes(pr *mux.Router) {
	pr.PathPrefix("/{plugin-name}").HandlerFunc(pluginHandler)
}

// Forward a request to the appropriate plugin
func pluginHandler(w http.ResponseWriter, r *http.Request) {
	pathVars := mux.Vars(r)
	pluginName := pathVars["plugin-name"]
	if endpoint, err := registry.GetPluginLocation(pluginName); err != nil {
		log.Println("No available Plugin Node for: " + pluginName)
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		reverseProxy(endpoint, w, r)
	}
}

// Returns the list of all enabled plugins
// Non-enabled plugins are not included
func getPluginList(w http.ResponseWriter, r *http.Request) {
	_, err := AuthUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}
	LogRequest(r)
	pluginList := fserv.GetEnabledPlugins()
	jsonList, err := json.Marshal(pluginList)
	if err != nil {
		LogError(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonList)
}

// Serves a file from a plugin's resource directory
func servePluginResource(w http.ResponseWriter, r *http.Request) {
	LogRequest(r)
	pathVars := mux.Vars(r)
	pluginName := pathVars["plugin-name"]
	thePlugin, err := fserv.GetPluginData(pluginName)
	if err != nil {
		LogError(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pluginResourceDir := thePlugin.PluginDir
	theHandler := httpgzip.NewHandler(http.StripPrefix("/app/plugin/"+pluginName+"/resource/",
		http.FileServer(http.Dir(pluginResourceDir))))
	theHandler.ServeHTTP(w, r)
}
