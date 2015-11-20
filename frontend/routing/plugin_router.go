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

package routing

import (
	"encoding/json"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/daaku/go.httpgzip"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/gorilla/mux"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/wikifeat-common/registry"
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
