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
	"github.com/gorilla/mux"
	"github.com/rhinoman/wikifeat/common/registry"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func handleApiRoutes(ar *mux.Router) {
	ar.PathPrefix("/users").HandlerFunc(userHandler)
	ar.PathPrefix("/wikis").HandlerFunc(wikiHandler)
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	if endpoint, err := registry.GetServiceLocation("users"); err != nil {
		log.Println("No Available User Services!")
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		reverseProxy(endpoint, w, r)
	}
}

func wikiHandler(w http.ResponseWriter, r *http.Request) {
	if endpoint, err := registry.GetServiceLocation("wikis"); err != nil {
		log.Println("No Available Wiki Services!")
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		reverseProxy(endpoint, w, r)
	}
}

func reverseProxy(endpoint string,
	w http.ResponseWriter,
	r *http.Request) {
	target, err := url.Parse(endpoint)
	if err != nil {
		//I have no idea
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	rp := httputil.NewSingleHostReverseProxy(target)
	rp.ServeHTTP(w, r)
}
