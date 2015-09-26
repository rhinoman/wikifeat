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
	"bytes"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/daaku/go.httpgzip"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/gorilla/mux"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/couchdb-go"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/golang.org/x/net/html"
	"github.com/rhinoman/wikifeat/common/config"
	"github.com/rhinoman/wikifeat/common/entities"
	"github.com/rhinoman/wikifeat/common/services"
	"github.com/rhinoman/wikifeat/frontend/fserv"
	"log"
	"net/http"
	"os"
	"path"
)

var webAppDir string
var indexHtml string

// Start doing the HTTP server/router thing
func Start() {
	webAppDir = config.Frontend.WebAppDir
	finishIndex()
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
	ar := r.PathPrefix("/api/" + config.Service.ApiVersion).Subrouter()
	handleApiRoutes(ar)
	// Handle Plugin requests
	pr := r.PathPrefix("/app/plugin").Subrouter()
	handlePluginRoutes(pr)
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
func finishIndex() {
	log.Print("Finishing index file")
	filename := webAppDir + "/index.html"
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
	auth, err := services.GetAuth(r)
	if err != nil && config.Auth.AllowGuest {
		auth = &couchdb.BasicAuth{
			Username: "guest",
			Password: "guest",
		}
	} else if err != nil {
		return nil, err
	}
	userInfo, err := services.GetUserFromAuth(auth)
	if err != nil {
		return nil, err
	}
	cui := &entities.CurrentUserInfo{
		Auth: auth,
		User: userInfo,
	}
	return cui, nil
}

// Log a request
func LogRequest(r *http.Request) {
	log.Printf("[Wikifeat-Frontend] %s  %s", r.Method, r.URL.String())
}

// Log an error
func LogError(err error) {
	log.Printf("[Wikifeat-Frontend] ERROR: %v", err.Error())
}
