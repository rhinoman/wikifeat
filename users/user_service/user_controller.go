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

package user_service

import (
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/emicklei/go-restful"
	. "github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/wikifeat-common/entities"
	. "github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/wikifeat-common/services"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/wikifeat-common/util"
	"log"
	"net/http"
	"strconv"
)

type UsersController struct{}

type UserResponse struct {
	Links HatLinks `json:"_links"`
	User  User     `json:"user"`
}

type UserListResponse struct {
	Links     HatLinks `json:"_links"`
	TotalRows int      `json:"totalRows"`
	PageNum   int      `json:"offset"`
	UserList  struct {
		List []UserResponse `json:"ea:user"`
	} `json:"_embedded"`
}

func (uc UsersController) userUri() string {
	return ApiPrefix() + "/users"
}

var usersWebService *restful.WebService

func (uc UsersController) Service() *restful.WebService {
	return usersWebService
}

//Define routes
func (uc UsersController) Register(container *restful.Container) {
	//avatars is a subcontroller
	ac := AvatarController{}
	usersWebService = new(restful.WebService)
	usersWebService.Filter(LogRequest)
	usersWebService.
		Path(uc.userUri()).
		Doc("Manage Users").
		ApiVersion(ApiVersion()).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	usersWebService.Route(usersWebService.POST("").To(uc.create).
		Filter(AuthUser).
		Doc("Create a User").
		Operation("create").
		Reads(User{}).
		Writes(UserResponse{}))

	usersWebService.Route(usersWebService.PUT("/{user-id}").To(uc.update).
		Filter(AuthUser).
		Doc("Updates a User").
		Operation("update").
		Param(usersWebService.PathParameter("user-id", "User Name").DataType("string")).
		Param(usersWebService.HeaderParameter("If-Match", "User revision").DataType("string")).
		Reads(User{}).
		Writes(UserResponse{}))

	usersWebService.Route(usersWebService.PUT("/{user-id}/grant_role").To(uc.grant).
		Filter(AuthUser).
		Doc("Grants access to a resource to a user").
		Operation("grant").
		Param(usersWebService.PathParameter("user-id", "User Name").DataType("string")).
		Reads(RoleRequest{}).
		Writes(BooleanResponse{}))

	usersWebService.Route(usersWebService.PUT("/{user-id}/revoke_role").To(uc.revoke).
		Filter(AuthUser).
		Doc("Revokes access to a resource").
		Operation("revoke").
		Param(usersWebService.PathParameter("user-id", "User Name").DataType("string")).
		Reads(RoleRequest{}).
		Writes(BooleanResponse{}))

	usersWebService.Route(usersWebService.PUT("/{user-id}/change_password").
		To(uc.changePassword).
		Filter(AuthUser).
		Doc("Change a user's password").
		Operation("changePassword").
		Param(usersWebService.PathParameter("user-id", "User Name").DataType("string")).
		Param(usersWebService.HeaderParameter("If-Match", "User revision").DataType("string")).
		Reads(ChangePasswordRequest{}).
		Writes(BooleanResponse{}))

	usersWebService.Route(usersWebService.PUT("/{user-id}/request_reset").
		To(uc.requestPasswordReset).
		Doc("Request a password reset (forgot, etc.)").
		Operation("resetPasswordRequest").
		Param(usersWebService.PathParameter("user-id", "User Name").DataType("string")).
		Writes(BooleanResponse{}))

	usersWebService.Route(usersWebService.PUT("/{user-id}/reset_password").
		To(uc.resetPassword).
		Doc("Resets a users' password").
		Operation("resetPassword").
		Param(usersWebService.PathParameter("user-id", "User Name").DataType("string")).
		Reads(ResetTokenRequest{}).
		Writes(BooleanResponse{}))

	usersWebService.Route(usersWebService.GET("/{user-id}").To(uc.read).
		Filter(AuthUser).
		Doc("Gets a User").
		Operation("read").
		Param(usersWebService.PathParameter("user-id", "User Name").DataType("string")).
		Writes(UserResponse{}))

	usersWebService.Route(usersWebService.GET("/current_user").To(uc.readCurrentUser).
		Filter(AuthUser).
		Doc("Gets the currently authenticated user").
		Operation("readCurrentUser").
		Writes(UserResponse{}))

	usersWebService.Route(usersWebService.DELETE("/{user-id}").To(uc.del).
		Filter(AuthUser).
		Doc("Deletes a User").
		Operation("delete").
		Param(usersWebService.PathParameter("user-id", "User Name").DataType("string")).
		Writes(BooleanResponse{}))

	usersWebService.Route(usersWebService.POST("/login").To(uc.login).
		Doc("Creates a new User Session").
		Operation("login").
		Reads(UserLoginCredentials{}).
		Writes(BooleanResponse{}))

	usersWebService.Route(usersWebService.DELETE("/login").To(uc.logout).
		Doc("Destroys User Session").
		Operation("logout").
		Writes(BooleanResponse{}))

	usersWebService.Route(usersWebService.GET("").To(uc.list).
		Filter(AuthUser).
		Doc("Gets a list of users").
		Operation("list").
		Param(usersWebService.QueryParameter("pageNum", "Page Number").DataType("integer")).
		Param(usersWebService.QueryParameter("numPerPage", "Number of records to return").DataType("integer")).
		Param(usersWebService.QueryParameter("forResource", "Return users that have of roles associated with a resource").DataType("string")).
		Writes(UserListResponse{}))
	//Add routes form avatars to the users controller
	ac.AddRoutes(usersWebService)
	//Add the users controller to the container
	container.Add(usersWebService)
}

func (uc UsersController) genUserUri(userId string) string {
	theUri := uc.userUri() + "/" + userId
	return theUri
}

//Create a User
func (uc UsersController) create(request *restful.Request,
	response *restful.Response) {
	curUser := GetCurrentUser(request, response)
	if curUser == nil {
		Unauthenticated(request, response)
		return
	}
	newUser := new(User)
	err := request.ReadEntity(newUser)
	if err != nil {
		WriteServerError(err, response)
		return
	}
	rev, err := new(UserManager).Create(newUser, curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	response.AddHeader("ETag", rev)
	ur := uc.genRecordResponse(curUser, newUser.UserName, newUser)
	SetAuth(response, curUser.Auth)
	response.WriteHeader(http.StatusCreated)
	response.WriteEntity(ur)
}

//Returns the currently authenticated user
func (uc UsersController) readCurrentUser(request *restful.Request,
	response *restful.Response) {
	curUser := GetCurrentUser(request, response)
	if curUser == nil {
		Unauthenticated(request, response)
		return
	}
	readUser := new(User)
	userId := curUser.User.UserName
	rev, err := new(UserManager).Read(userId, readUser, curUser)
	if err != nil {
		WriteServerError(err, response)
		return
	}
	response.AddHeader("ETag", rev)
	ur := uc.genRecordResponse(curUser, readUser.UserName, readUser)
	SetAuth(response, curUser.Auth)
	response.WriteEntity(ur)
}

//Read a User
func (uc UsersController) read(request *restful.Request,
	response *restful.Response) {
	curUser := GetCurrentUser(request, response)
	if curUser == nil {
		Unauthenticated(request, response)
		return
	}
	userId := request.PathParameter("user-id")
	readUser := new(User)
	rev, err := new(UserManager).Read(userId, readUser, curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	response.AddHeader("ETag", rev)
	ur := uc.genRecordResponse(curUser, readUser.UserName, readUser)
	SetAuth(response, curUser.Auth)
	response.WriteEntity(ur)
}

//Update a User
func (uc UsersController) update(request *restful.Request,
	response *restful.Response) {
	curUser := GetCurrentUser(request, response)
	if curUser == nil {
		Unauthenticated(request, response)
		return
	}
	userId := request.PathParameter("user-id")
	rev := request.HeaderParameter("If-Match")
	if userId == "" || rev == "" {
		WriteBadRequestError(response)
		return
	}
	updateUser := new(User)
	err := request.ReadEntity(updateUser)
	if err != nil {
		WriteServerError(err, response)
		return
	}
	rev, err = new(UserManager).Update(userId, rev, updateUser, curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	response.AddHeader("ETag", rev)
	ur := uc.genRecordResponse(curUser, userId, updateUser)
	SetAuth(response, curUser.Auth)
	response.WriteEntity(ur)
}

//Grant a role to a User
func (uc UsersController) grant(request *restful.Request,
	response *restful.Response) {
	curUser := GetCurrentUser(request, response)
	if curUser == nil {
		Unauthenticated(request, response)
		return
	}
	userId := request.PathParameter("user-id")
	if userId == "" {
		WriteBadRequestError(response)
		return
	}
	grantRequest := new(RoleRequest)
	err := request.ReadEntity(grantRequest)
	if err != nil {
		WriteServerError(err, response)
		return
	}
	rev, err := new(UserManager).GrantRole(userId, grantRequest, curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	response.AddHeader("ETag", rev)
	SetAuth(response, curUser.Auth)
	response.WriteEntity(BooleanResponse{Success: true})
}

//Revoke a role from a User
func (uc UsersController) revoke(request *restful.Request,
	response *restful.Response) {
	curUser := GetCurrentUser(request, response)
	if curUser == nil {
		Unauthenticated(request, response)
		return
	}
	userId := request.PathParameter("user-id")
	if userId == "" {
		WriteBadRequestError(response)
		return
	}
	revokeRequest := new(RoleRequest)
	err := request.ReadEntity(revokeRequest)
	if err != nil {
		WriteServerError(err, response)
		return
	}
	rev, err := new(UserManager).RevokeRole(userId, revokeRequest, curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	response.AddHeader("ETag", rev)
	SetAuth(response, curUser.Auth)
	response.WriteEntity(BooleanResponse{Success: true})
}

//Change user password
func (uc UsersController) changePassword(request *restful.Request,
	response *restful.Response) {
	curUser := GetCurrentUser(request, response)
	if curUser == nil {
		Unauthenticated(request, response)
		return
	}
	userId := request.PathParameter("user-id")
	if userId == "" {
		WriteBadRequestError(response)
		return
	}
	rev := request.HeaderParameter("If-Match")
	if rev == "" {
		WriteBadRequestError(response)
		return
	}
	cpr := new(ChangePasswordRequest)
	err := request.ReadEntity(cpr)
	if err != nil {
		WriteServerError(err, response)
		return
	}
	rev, err = new(UserManager).ChangePassword(userId, rev, cpr, curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	response.AddHeader("ETag", rev)
	SetAuth(response, curUser.Auth)
	response.WriteEntity(BooleanResponse{Success: true})
}

//Request password reset (Forgot, etc.)
func (uc UsersController) requestPasswordReset(request *restful.Request,
	response *restful.Response) {
	userId := request.PathParameter("user-id")
	if userId == "" {
		WriteBadRequestError(response)
		return
	}
	err := new(UserManager).RequestPasswordReset(userId)
	if err != nil {
		WriteError(err, response)
		return
	}
	br := BooleanResponse{Success: true}
	response.WriteEntity(br)
}

//Resets password
func (uc UsersController) resetPassword(request *restful.Request,
	response *restful.Response) {
	userId := request.PathParameter("user-id")
	if userId == "" {
		WriteBadRequestError(response)
		return
	}
	tr := new(ResetTokenRequest)
	err := request.ReadEntity(tr)
	if err != nil {
		WriteServerError(err, response)
		return
	}
	err = new(UserManager).ResetPassword(userId, tr)
	if err != nil {
		WriteError(err, response)
		return
	}
	response.WriteEntity(BooleanResponse{Success: true})
}

//Delete a User
func (uc UsersController) del(request *restful.Request,
	response *restful.Response) {
	curUser := GetCurrentUser(request, response)
	if curUser == nil {
		Unauthenticated(request, response)
		return
	}
	userId := request.PathParameter("user-id")
	if userId == "" {
		WriteBadRequestError(response)
		return
	}
	rev, err := new(UserManager).Delete(userId, curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	response.AddHeader("ETag", rev)
	rr := BooleanResponse{Success: true}
	SetAuth(response, curUser.Auth)
	response.WriteEntity(rr)
}

//User login
func (uc UsersController) login(request *restful.Request,
	response *restful.Response) {
	loginCredentials := new(UserLoginCredentials)
	err := request.ReadEntity(loginCredentials)
	if err != nil {
		LogError(request, response, err)
		WriteIllegalRequestError(response)
		return
	}
	cookieAuth, err := new(UserManager).Login(loginCredentials)
	if err != nil {
		LogError(request, response, err)
		WriteError(err, response)
		return
	}
	//Create an Auth cookie
	authCookie := http.Cookie{
		Name:     "AuthSession",
		Value:    cookieAuth.AuthToken,
		Path:     "/",
		HttpOnly: true,
	}
	//Create a CSRF cookie for this session
	//Subsequent requests must include this in a header field
	//X-Csrf-Token
	csrfCookie := http.Cookie{
		Name:     "CsrfToken",
		Value:    util.GenHashString(cookieAuth.AuthToken),
		Path:     "/",
		HttpOnly: false,
	}
	response.AddHeader("Set-Cookie", authCookie.String())
	response.AddHeader("Set-Cookie", csrfCookie.String())
	response.WriteEntity(BooleanResponse{Success: true})
}

//User logout
func (uc UsersController) logout(request *restful.Request,
	response *restful.Response) {
	token := func() string {
		for _, cookie := range request.Request.Cookies() {
			if cookie.Name == "AuthSession" {
				return cookie.Value
			}
		}
		return ""
	}()
	err := new(UserManager).Logout(token)
	if err != nil {
		LogError(request, response, err)
		WriteError(err, response)
		return
	}
	//Because CouchDB doesn't actually destroy the session,
	//best we can do is clear the cookie in the browser.
	//This is apparently "not a bug" :|
	//http://webmail.dev411.com/t/couchdb/dev/141xwf5vb0/jira-created-couchdb-2042-session-not-cleared-after-delete-session-cookie-auth
	theCookie := http.Cookie{
		Name:     "AuthSession",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
	}
	response.AddHeader("Set-Cookie", theCookie.String())
	response.WriteEntity(BooleanResponse{Success: true})
}

//Get user list
func (uc UsersController) list(request *restful.Request,
	response *restful.Response) {
	curUser := GetCurrentUser(request, response)
	if curUser == nil {
		Unauthenticated(request, response)
		return
	}
	var limit int
	var pageNum int
	limitString := request.QueryParameter("numPerPage")
	if limitString == "" {
		limit = 0
	} else {
		ln, err := strconv.Atoi(limitString)
		if err != nil {
			log.Printf("Error: %v", err)
			WriteIllegalRequestError(response)
			return
		}
		limit = ln
	}
	pageNumString := request.QueryParameter("pageNum")
	if pageNumString == "" {
		pageNum = 1
	} else {
		if ln, err := strconv.Atoi(pageNumString); err != nil {
			log.Printf("Error: %v", err)
			WriteIllegalRequestError(response)
		} else {
			pageNum = ln
		}
	}
	ulr := UserListQueryResponse{}
	forResource := request.QueryParameter("forResource")
	var err error
	if forResource != "" {
		//Make sure the user is an admin for the given resource
		if !util.HasRole(curUser.User.Roles, AdminRole(forResource)) &&
			!util.HasRole(curUser.User.Roles, AdminRole(MainDbName())) {
			WriteError(NotAdminError(), response)
			return
		}
		rolesArray := []string{AdminRole(forResource),
			WriteRole(forResource), ReadRole(forResource)}
		err = new(UserManager).GetUserListForRole(
			pageNum, limit, rolesArray, &ulr, curUser)
	} else {
		err = new(UserManager).GetUserList(pageNum, limit, &ulr, curUser)
	}
	if err != nil {
		WriteError(err, response)
		return
	}
	uir := uc.genUserListResponse(curUser, &ulr)
	SetAuth(response, curUser.Auth)
	response.WriteEntity(uir)
}

func (uc UsersController) genUserListResponse(curUser *CurrentUserInfo,
	ulr *UserListQueryResponse) UserListResponse {
	uir := UserListResponse{}
	uir.TotalRows = ulr.TotalRows
	uir.PageNum = ulr.Offset
	for _, row := range ulr.Rows {
		urr := uc.genRecordResponse(curUser, row.Value.UserName, &row.Value)
		uir.UserList.List = append(uir.UserList.List, urr)
	}
	uir.Links = GenIndexLinks(curUser.User.Roles, MainDbName(),
		uc.userUri())
	return uir
}

func (uc UsersController) genRecordResponse(curUser *CurrentUserInfo,
	userId string, user *User) UserResponse {

	user.Id = user.UserName
	return UserResponse{
		Links: uc.genUserRecordLinks(curUser.User.Roles, userId,
			curUser.User.UserName, uc.genUserUri(userId)),
		User: *user,
	}
}

func (uc UsersController) genUserRecordLinks(userRoles []string,
	userId string, curUserId string, uri string) HatLinks {
	links := HatLinks{}
	dbName := MainDbName()
	admin := util.HasRole(userRoles, AdminRole(dbName))
	//Write := util.HasRole(userRoles, WriteRole(dbName))
	self := func() bool {
		if curUserId == userId {
			return true
		} else {
			return false
		}
	}()
	links.Self = &HatLink{Href: uri, Method: "GET"}
	if admin || self {
		links.Update = &HatLink{Href: uri, Method: "PUT"}
	}
	return links
}
