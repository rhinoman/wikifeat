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
	"bytes"
	"github.com/daaku/go.httpgzip"
	"github.com/gorilla/mux"
	"github.com/rhinoman/couchdb-go"
	"github.com/rhinoman/wikifeat/common/config"
	"github.com/rhinoman/wikifeat/common/database"
	"github.com/rhinoman/wikifeat/common/entities"
	"github.com/rhinoman/wikifeat/common/services"
	"github.com/rhinoman/wikifeat/common/util"
	"github.com/rhinoman/wikifeat/frontend/fserv"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
)

var webAppDir string
var indexHtml string

// Start doing the HTTP server/router thing
func Start() {
	curDir, err := util.GetExecDirectory()
	if err != nil {
		log.Fatal(err)
	}
	indexFile := path.Join(curDir, "/index.html")
	webAppDir = config.Frontend.WebAppDir
	finishIndex(indexFile)
	r := mux.NewRouter()
	// Serve the login page
	r.HandleFunc("/login", getLogin).Methods("GET")
	// Serve the forgot password page
	r.HandleFunc("/forgot_password", getForgotPassword).Methods("GET")
	// Serve the rest password page
	r.HandleFunc("/reset_password", getResetPassword).Methods("GET")
	// Serve the Main app
	r.StrictSlash(true).HandleFunc("/app", getAppRoot)
	// Fetch the location of the homepage
	r.HandleFunc("/app/home", getHome).Methods("GET")
	// Serve up static assets
	r.PathPrefix("/app/resource/").
		Methods("GET").
		Handler(httpgzip.NewHandler(http.StripPrefix("/app/resource/",
			http.FileServer(http.Dir(webAppDir)))))
	//Handle other routes
	r.PathPrefix("/app/wikis").HandlerFunc(getAppRoot)
	r.PathPrefix("/app/users").HandlerFunc(getAppRoot)
	// Custom extensions (plugins) should use /app/x routes on the front end
	r.PathPrefix("/app/x").HandlerFunc(getAppRoot)
	// Handle API requests
	ar := r.PathPrefix("/api/" + config.ApiVersion).Subrouter()
	handleApiRoutes(ar)
	// Handle Plugin requests for frontend resources
	pr := r.PathPrefix("/app/plugin").Subrouter()
	handlePluginRoutes(pr)
	// Backend Plugin requests
	bpr := r.PathPrefix("/plugin").Subrouter()
	handlePluginBackendRoutes(bpr)
	log.Print("Starting HTTP router")
	if config.Service.UseSSL {
		certFile := config.Service.SSLCertFile
		keyFile := config.Service.SSLKeyFile
		log.Fatal(http.ListenAndServeTLS(":"+config.Service.Port,
			certFile, keyFile, r))
	} else {
		log.Fatal(http.ListenAndServe(":"+config.Service.Port, r))
	}
}

// Opens up the index.html file and adds additional (mostly plugin) elements
func finishIndex(filename string) {
	log.Print("Finishing index file")

	fileReader, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	//Parse the html document
	doc, err := html.Parse(fileReader)
	if err != nil {
		log.Fatal(err)
	}
	//Find the HEAD
	var processNode func(*html.Node)
	processNode = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "head" {
			enabledPlugins := fserv.GetEnabledPlugins()
			for _, plugin := range enabledPlugins {
				name := plugin.Name
				css := plugin.Stylesheet
				if css != "" {
					path := "/app/plugin/" + name + "/resource/" + css
					//Construct a new link css node and add it to the head
					attributes := []html.Attribute{
						html.Attribute{Key: "rel", Val: "stylesheet"},
						html.Attribute{Key: "type", Val: "text/css"},
						html.Attribute{Key: "href", Val: path},
					}
					linkNode := html.Node{
						Type: html.ElementNode,
						Data: "link",
						Attr: attributes,
					}
					node.AppendChild(&linkNode)
				}
			}
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			processNode(c)
		}
	}
	processNode(doc)
	var b bytes.Buffer
	html.Render(&b, doc)
	indexHtml = b.String()
	fileReader.Close()
}

// Get the main app
func getAppRoot(w http.ResponseWriter, r *http.Request) {
	//location := path.Join(webAppDir, "index.html")
	LogRequest(r)
	log.Print("Serving web root")
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(indexHtml))
	//http.ServeFile(w, r, location)
}

// Get the Home Uri
func getHome(w http.ResponseWriter, r *http.Request) {
	homeUri := config.Frontend.Homepage
	jsonString := "{\"home\": \"" + homeUri + "\"}"
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonString))
}

// Serve up the login page
func getLogin(w http.ResponseWriter, r *http.Request) {
	location := path.Join(webAppDir, "login.html")
	LogRequest(r)
	log.Printf("Serving login page from %s", location)
	http.ServeFile(w, r, location)
}

// Serve up the forgot password page
func getForgotPassword(w http.ResponseWriter, r *http.Request) {
	location := path.Join(webAppDir, "forgot_password.html")
	LogRequest(r)
	log.Printf("Serving Forgot Password page from %s", location)
	http.ServeFile(w, r, location)
}

// Serve up the reset password page
func getResetPassword(w http.ResponseWriter, r *http.Request) {
	location := path.Join(webAppDir, "reset_password.html")
	LogRequest(r)
	log.Printf("Serving Reset Password page from %s", location)
	http.ServeFile(w, r, location)
}

// Authenticate a user
// Returns CurrentUserInfo if authenticated, error if not
func AuthUser(r *http.Request) (*entities.CurrentUserInfo, error) {
	cAuth, err := services.GetAuth(r)
	if err == http.ErrNoCookie && config.Auth.AllowGuest {
		cAuth = &couchdb.BasicAuth{
			Username: "guest",
			Password: "guest",
		}
	} else if err != nil {
		return nil, err
	}
	userInfo, err := database.GetUserFromAuth(cAuth)
	if err != nil {
		return nil, err
	}
	cui := &entities.CurrentUserInfo{
		Auth: cAuth,
		User: userInfo,
	}
	return cui, nil
}

// Proxy request to service node
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

// Log a request
func LogRequest(r *http.Request) {
	log.Printf("[Wikifeat-Frontend] %s  %s", r.Method, r.URL.String())
}

// Log an error
func LogError(err error) {
	log.Printf("[Wikifeat-Frontend] ERROR: %v", err.Error())
}
