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

package services

import (
	"errors"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/emicklei/go-restful"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/couchdb-go"
	"github.com/rhinoman/wikifeat/common/config"
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
	auth, err := GetAuth(request.Request)
	if err != nil && config.Auth.AllowGuest {
		auth = &couchdb.BasicAuth{
			Username: "guest",
			Password: "guest",
		}
	} else if err != nil {
		Unauthenticated(request, resp)
		return
	}
	userInfo, err := GetUserFromAuth(auth)
	if err != nil {
		Unauthenticated(request, resp)
		return
	}
	cui := &CurrentUserInfo{
		Auth: auth,
		User: userInfo,
	}
	request.SetAttribute("currentUser", cui)
	chain.ProcessFilter(request, resp)

}

//Authenticates Request
//Returns Auth header from request
func GetAuth(request *http.Request) (couchdb.Auth, error) {
	//What kind of authentication type do we have?
	if request.Header.Get("Authorization") != "" {
		//We have an authorization header, just use it
		return &couchdb.PassThroughAuth{
			AuthHeader: request.Header.Get("Authorization"),
		}, nil
	}
	//Ok, Check for a session cookie
	//TODO: CHECK FOR EXPIRED COOKIES!!!
	var sessionToken string
	var csrfToken string
	for _, cookie := range request.Cookies() {
		if cookie.Name == "AuthSession" {
			sessionToken = cookie.Value
		}
		if cookie.Name == "CsrfToken" {
			csrfToken = cookie.Value
		}
	}
	if sessionToken == "" || csrfToken == "" {
		return nil, errors.New("Unauthenticated")
	}
	//Set the cookie auth
	cookieAuth := &couchdb.CookieAuth{
		AuthToken: sessionToken,
	}
	//Better check for our Csrf Token
	csrf := request.Header.Get("X-Csrf-Token")
	//Csrf token should match the csrf cookie
	if csrfToken == csrf {
		//Request is good!
		return cookieAuth, nil
	} else {
		//CSRF is bad, Go away!
		return nil, errors.New("Unauthenticated")
	}
}

//Set Updated auth cookies
func SetAuth(response *restful.Response, auth couchdb.Auth) {
	authData := auth.GetUpdatedAuth()
	if authData == nil {
		return
	}
	if val, ok := authData["AuthSession"]; ok {
		authCookie := http.Cookie{
			Name:     "AuthSession",
			Value:    val,
			Path:     "/",
			HttpOnly: true,
		}
		//Create a CSRF cookie
		csrfCookie := http.Cookie{
			Name:     "CsrfToken",
			Value:    util.GenHashString(val),
			Path:     "/",
			HttpOnly: false,
		}
		response.AddHeader("Set-Cookie", authCookie.String())
		response.AddHeader("Set-Cookie", csrfCookie.String())
	}
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
