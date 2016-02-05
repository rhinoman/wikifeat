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

package auth_service

import "github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/emicklei/go-restful"
import . "github.com/rhinoman/wikifeat/common/services"
import (
	"github.com/rhinoman/wikifeat/common/auth"
	"net/http"
)

type AuthController struct{}

var authWebService *restful.WebService

func (ac AuthController) authUri() string {
	return ApiPrefix() + "/auth"
}

func (ac AuthController) sessionUri(sessionId string) string {
	return ac.authUri() + "/sessions/" + sessionId
}

func (ac AuthController) Service() *restful.WebService {
	return authWebService
}

//Define routes
func (ac AuthController) Register(container *restful.Container) {
	authWebService = new(restful.WebService)
	authWebService.Filter(LogRequest)
	authWebService.
		Path(ac.authUri()).
		Doc("Manage Authentication and Sessions").
		ApiVersion(ApiVersion()).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	authWebService.Route(authWebService.GET("").To(ac.getAuth).
		Doc("Get authentication").
		Operation("getAuth").
		Writes(auth.WikifeatAuth{}))

	authWebService.Route(authWebService.POST("/session").To(ac.create).
		Doc("Create a Session").
		Operation("create").
		Reads(auth.UserLoginCredentials{}).
		Writes(BooleanResponse{}))

	authWebService.Route(authWebService.DELETE("/session").To(ac.del).
		Doc("Destroy a Session").
		Operation("del").
		Writes(BooleanResponse{}))

	container.Add(authWebService)

}

// Creates a session and sets a cookie in the http response
func (ac AuthController) create(request *restful.Request,
	response *restful.Response) {
	credentials := auth.UserLoginCredentials{}
	err := request.ReadEntity(&credentials)
	if err != nil {
		Unauthenticated(request, response)
		return
	}
	am := new(AuthManager)
	sess, err := am.Create(credentials.Username,
		credentials.Password, credentials.AuthType)
	if err != nil {
		Unauthenticated(request, response)
		return
	}
	authCookie := ac.genAuthCookie(sess)
	response.Header().Add("Set-Cookie", authCookie.String())
	response.WriteEntity(BooleanResponse{Success: true})
}

// Destroys a session
func (ac AuthController) del(request *restful.Request,
	response *restful.Response) {
	am := new(AuthManager)
	//Get the session id
	sessId, err := am.GetSessionId(request.Request)
	if err != nil {
		Unauthenticated(request, response)
		return
	}
	//Read the session
	sess, err := am.ReadSession(sessId)
	if err != nil {
		Unauthenticated(request, response)
		return
	}
	if err = am.Destroy(sess); err != nil {
		WriteError(err, response)
	} else {
		//Clear the session cookie
		theCookie := http.Cookie{
			Name:     "AuthSession",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
		}
		response.AddHeader("Set-Cookie", theCookie.String())
		response.WriteEntity(BooleanResponse{Success: true})
	}
}

// Takes a session Id and returns an Auth object
func (ac AuthController) getAuth(request *restful.Request,
	response *restful.Response) {
	req := request.Request
	am := new(AuthManager)
	wfAuth, err := am.GetAuth(req)
	if err != nil {
		WriteServerError(err, response)
		return
	}
	response.WriteEntity(wfAuth)
}

func (ac AuthController) genAuthCookie(session *auth.Session) http.Cookie {
	return http.Cookie{
		Name:     "AuthSession",
		Value:    session.Id,
		Path:     "/",
		HttpOnly: true,
	}

}
