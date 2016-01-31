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
import "github.com/rhinoman/wikifeat/common/auth"

type AuthController struct{}

type SessionResponse struct {
	Links   HatLinks     `json:"_links"`
	Session auth.Session `json:"session"`
}

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

	authWebService.Route(authWebService.POST("").To(ac.create).
		Doc("Create a Session").
		Operation("create").
		Reads(auth.UserLoginCredentials{}).
		Writes(SessionResponse{}))

}

func (ac AuthController) create(request *restful.Request,
	response *restful.Response) {
	credentials := auth.UserLoginCredentials{}
	err := request.ReadEntity(credentials)
	if err != nil {
		WriteServerError(err, response)
	}
	am := new(AuthManager)
	sess, err := am.Create(credentials.Username,
		credentials.Password, credentials.AuthType)
	sr := ac.genSessionResponse(sess)
	response.WriteEntity(sr)
}

func (ac AuthController) genSessionResponse(session *auth.Session) SessionResponse {
	links := HatLinks{}
	links.Self = &HatLink{Href: ac.sessionUri(session.Id)}
	return SessionResponse{
		Links:   links,
		Session: *session,
	}
}
