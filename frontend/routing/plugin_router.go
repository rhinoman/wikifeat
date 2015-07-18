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
	"github.com/daaku/go.httpgzip"
	"github.com/gorilla/mux"
	"github.com/rhinoman/wikifeat/frontend/fserv"
	"net/http"
)

type PluginListResponse struct {
	EnabledPlugins []fserv.PluginData `json:"enabledPlugins"`
}

// Route stuff
func handlePluginRoutes(pr *mux.Router) {
	pr.StrictSlash(true).HandleFunc("/", getPluginList).Methods("GET")
	pr.PathPrefix("/{plugin-name}/resource/").
		Methods("GET").
		HandlerFunc(servePluginResource)
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

// Serves a file from a plugin's directory
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
