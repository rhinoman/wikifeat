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

package services

import (
	"errors"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/emicklei/go-restful"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/couchdb-go"
	"github.com/rhinoman/wikifeat/common/auth"
	"github.com/rhinoman/wikifeat/common/config"
	. "github.com/rhinoman/wikifeat/common/database"
	. "github.com/rhinoman/wikifeat/common/entities"
	"github.com/rhinoman/wikifeat/common/util"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Controller interface {
	Service() *restful.WebService
}

type HatLinks struct {
	Self   *HatLink `json:"self"`
	Create *HatLink `json:"create,omitempty"`
	Update *HatLink `json:"update,omitempty"`
	Delete *HatLink `json:"delete,omitempty"`
}

type HatLink struct {
	Href       string `json:"href"`
	HrefLang   string `json:"hreflang,omitempty"`
	Title      string `json:"title,omitempty"`
	Type       string `json:"type,omitempty"`
	Deprecated string `json:"deprecation,omitempty"`
	Name       string `json:"name,omitempty"`
	Profile    string `json:"profile,omitempty"`
	Templated  bool   `json:"templated,omitempty"`
	Method     string `json:"method,omitempty"` //HTTP method, not standard HAL
}

type BooleanResponse struct {
	Success bool `json:"success"`
}

func ApiVersion() string {
	return config.Service.ApiVersion
}

func ApiPrefix() string {
	return "/api/" + ApiVersion()
}

func PluginPrefix() string {
	return "/plugin"
}

//Create usual links for an index
func GenIndexLinks(userRoles []string, dbName string, uri string) HatLinks {
	links := HatLinks{}
	admin := util.HasRole(userRoles, AdminRole(dbName))
	write := util.HasRole(userRoles, WriteRole(dbName))
	//Generate the self link
	links.Self = &HatLink{Href: uri, Method: "GET"}
	if admin || write {
		links.Create = &HatLink{Href: uri, Method: "POST"}
	}
	return links
}

//Create the basic CRUD links for a resource record
func GenRecordLinks(userRoles []string, dbName string, uri string) HatLinks {
	links := HatLinks{}
	//Admin can be a resource admin OR a site admin/master
	admin := util.HasRole(userRoles, AdminRole(dbName)) ||
		util.HasRole(userRoles, AdminRole(MainDbName())) ||
		util.HasRole(userRoles, MasterRole())
	write := util.HasRole(userRoles, WriteRole(dbName))
	//Generate the self link
	links.Self = &HatLink{Href: uri, Method: "GET"}
	if admin || write {
		links.Update = &HatLink{Href: uri, Method: "PUT"}
		links.Delete = &HatLink{Href: uri, Method: "DELETE"}
	}
	return links
}

//Filter function.  Logs incoming requests
func LogRequest(request *restful.Request, resp *restful.Response,
	chain *restful.FilterChain) {
	method := request.Request.Method
	url := request.Request.URL.String()
	remoteAddr := request.Request.RemoteAddr
	log.Printf("[API] %v : %v %v", remoteAddr, method, url)
	chain.ProcessFilter(request, resp)
}

//Log an error
func LogError(request *restful.Request, resp *restful.Response, err error) {
	method := request.Request.Method
	url := request.Request.URL.String()
	remoteAddr := request.Request.RemoteAddr
	log.Printf("[ERROR] %v : %v : %v %v", err, remoteAddr, method, url)
}

//Filter function, figures out the current user from the auth header
func AuthUser(request *restful.Request, resp *restful.Response,
	chain *restful.FilterChain) {
	//authenticator := auth.NewAuthenticator(config.Auth.Authenticator)
	//Check for a basic auth header
	cAuth, err := GetAuth(request.Request)
	if err != nil && config.Auth.AllowGuest {
		cAuth = &couchdb.BasicAuth{
			Username: "guest",
			Password: "guest",
		}
	} else if err != nil {
		Unauthenticated(request, resp)
		return
	}
	userInfo, err := GetUserFromAuth(cAuth)
	if err != nil {
		Unauthenticated(request, resp)
		return
	}
	cui := &CurrentUserInfo{
		Auth: cAuth,
		User: userInfo,
	}
	request.SetAttribute("currentUser", cui)
	chain.ProcessFilter(request, resp)
}

func GetAuth(request *http.Request) (couchdb.Auth, error) {
	if bAuth := auth.GetBasicAuth(request); bAuth != nil {
		return bAuth, nil
	}
	if wAuth, err := auth.GetAuth(request); err != nil {
		return nil, err
	} else {
		return wAuth, nil
	}
}

//Set Updated auth cookies
func SetAuth(response *restful.Response, cAuth couchdb.Auth) {
	//authenticator := auth.NewAuthenticator(config.Auth.Authenticator)
	//authenticator.SetAuth(response.ResponseWriter, cAuth)
}

//Gets the current user from the header
func GetCurrentUser(request *restful.Request,
	response *restful.Response) *CurrentUserInfo {
	curUser, ok := request.Attribute("currentUser").(*CurrentUserInfo)
	if ok == false || curUser == nil {
		return nil
	} else {
		return curUser
	}
}

//Writes unauthenticated error to response
func Unauthenticated(request *restful.Request, response *restful.Response) {
	LogError(request, response, errors.New("Unauthenticated"))
	response.AddHeader("Content-Type", "text/plain")
	response.WriteErrorString(401, "Unauthenticated")
}

func WriteIllegalRequestError(response *restful.Response) {
	response.AddHeader("Content-Type", "text/plain")
	response.WriteErrorString(http.StatusBadRequest, "Bad Request")
}

//Writes an internal server error
func WriteServerError(err error, response *restful.Response) {
	log.Printf("%v", err)
	response.AddHeader("Content-Type", "text/plain")
	response.WriteErrorString(http.StatusInternalServerError, err.Error())
}

//Attempts to extract the http status code from an error
func getErrorCode(err error) int {
	splitStr := strings.Split(err.Error(), ":")
	if len(splitStr) > 1 {
		code, conerr := strconv.Atoi(splitStr[1])
		if conerr != nil {
			return 0
		} else {
			return code
		}
	} else {
		return 500
	}
}

//Writes and logs errors from the couchdb driver
func WriteError(err error, response *restful.Response) {
	str := err.Error()
	errStrings := strings.Split(str, ":")
	statusCode := 0
	var cErr error
	if len(errStrings) > 1 {
		statusCode, cErr = strconv.Atoi(errStrings[1])
	}
	if cErr != nil || statusCode == 0 {
		statusCode = 500
	}
	//Write the error to the response
	response.WriteErrorString(statusCode,
		http.StatusText(statusCode))
	//Log the error
	log.Printf("%v", err)
}

func WriteBadRequestError(response *restful.Response) {
	log.Printf("400: Bad Request")
	response.WriteErrorString(http.StatusBadRequest, "Bad Request")
}

//Returns the Admin Credentials as a CurrentUserInfo
func GetAdminUser() *CurrentUserInfo {
	return &CurrentUserInfo{
		Auth:  AdminAuth,
		Roles: []string{"admin"},
		User: &User{
			Roles: []string{"admin"},
		},
	}
}
