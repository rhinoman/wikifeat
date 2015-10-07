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

package user_service

import (
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/emicklei/go-restful"
	. "github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/wikifeat-common/entities"
	. "github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/wikifeat-common/services"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/wikifeat-common/util"
	"io"
	"net/http"
	"strings"
)

type AvatarController struct{}

var avatarUri = "/{user-id}/avatar"

type avatarLinks struct {
	HatLinks
	SaveImage          *HatLink `json:"saveImage,omitempty"`
	GetLargeAvatar     *HatLink `json:"getLargeAvatar,omitempty"`
	GetThumbnailAvatar *HatLink `json:"getThumbnailAvatar",omitempty"`
}

type AvatarResponse struct {
	Links        avatarLinks `json:"_links"`
	AvatarRecord UserAvatar  `json:"avatar_record"`
}

//Define routes

func (ac AvatarController) AddRoutes(ws *restful.WebService) {
	ws.Route(ws.POST(avatarUri).To(ac.create).
		Doc("Create a new User Avatar Record").
		Filter(AuthUser).
		Operation("create").
		Reads(UserAvatar{}).
		Param(ws.PathParameter("user-id", "User id").DataType("string")).
		Writes(AvatarResponse{}))

	ws.Route(ws.GET(avatarUri).To(ac.read).
		Doc("Reads a User Avatar Record").
		Filter(AuthUser).
		Operation("read").
		Param(ws.PathParameter("user-id", "User id").DataType("string")).
		Writes(AvatarResponse{}))

	ws.Route(ws.PUT(avatarUri).To(ac.update).
		Doc("Updated a User Avatar Record").
		Filter(AuthUser).
		Operation("update").
		Param(ws.PathParameter("user-id", "User id").DataType("string")).
		Param(ws.HeaderParameter("If-Match", "Avatar Revision").DataType("string")).
		Reads(UserAvatar{}).
		Writes(AvatarResponse{}))

	ws.Route(ws.DELETE(avatarUri).To(ac.del).
		Doc("Delete a User Avatar Record").
		Filter(AuthUser).
		Operation("delete").
		Param(ws.PathParameter("user-id", "User id").DataType("string")).
		Writes(BooleanResponse{}))

	ws.Route(ws.POST(avatarUri + "/image").
		Doc("Saves a User Avatar Image").
		Filter(AuthUser).
		Consumes("multipart/form-data").To(ac.saveImage).
		Operation("saveImage").
		Param(ws.PathParameter("user-id", "User id").DataType("string")).
		Param(ws.HeaderParameter("If-Match", "File Revision").DataType("string")).
		Param(ws.FormParameter("file-data", "The Image File").DataType("string")).
		Writes(BooleanResponse{}))

	ws.Route(ws.GET(avatarUri + "/image").To(ac.getImage).
		Doc("Fetches a User Avatar Image").
		Operation("getAvatar").Produces("image/jpeg").
		Param(ws.PathParameter("user-id", "User id").DataType("string")))

	ws.Route(ws.GET(avatarUri + "/thumbnail").To(ac.getThumb).
		Doc("Fetches a Thumbnail Avatar Image").
		Operation("getThumbnail").Produces("image/jpeg").
		Param(ws.PathParameter("user-id", "User id").DataType("string")))

}

func (ac AvatarController) genAvatarUri(userId string) string {
	theUri := ApiPrefix() + "/users" + avatarUri
	return strings.Replace(theUri, "{user-id}", userId, 1)
}

//Create a new User Avatar Record
func (ac AvatarController) create(request *restful.Request,
	response *restful.Response) {
	curUser := GetCurrentUser(request, response)
	if curUser == nil {
		Unauthenticated(request, response)
		return
	}
	userId := request.PathParameter("user-id")
	ua := new(UserAvatar)
	err := request.ReadEntity(ua)
	if err != nil {
		WriteServerError(err, response)
		return
	}
	rev, err := new(UserAvatarManager).Save(userId, "", ua, curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	response.AddHeader("ETag", rev)
	response.WriteHeader(http.StatusCreated)
	ar := ac.genRecordResponse(curUser, userId, ua)
	SetAuth(response, curUser.Auth)
	response.WriteEntity(ar)
}

//Read a User Avatar Record
func (ac AvatarController) read(request *restful.Request,
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
	ua := new(UserAvatar)
	rev, err := new(UserAvatarManager).Read(userId, ua, curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	response.AddHeader("ETag", rev)
	ar := ac.genRecordResponse(curUser, userId, ua)
	SetAuth(response, curUser.Auth)
	response.WriteEntity(ar)
}

//Updates a User Avatar Record
func (ac AvatarController) update(request *restful.Request,
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
	ua := new(UserAvatar)
	err := request.ReadEntity(ua)
	if err != nil {
		WriteServerError(err, response)
		return
	}
	rev, err = new(UserAvatarManager).Save(userId, rev, ua, curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	response.AddHeader("ETag", rev)
	ar := ac.genRecordResponse(curUser, userId, ua)
	SetAuth(response, curUser.Auth)
	response.WriteEntity(ar)
}

//Deletes a User Avatar Record
func (ac AvatarController) del(request *restful.Request,
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
	rev, err := new(UserAvatarManager).Delete(userId, curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	SetAuth(response, curUser.Auth)
	response.WriteEntity(BooleanResponse{Success: true})
	response.AddHeader("ETag", rev)
}

//Saves an Avatar Image
func (ac AvatarController) saveImage(request *restful.Request,
	response *restful.Response) {
	curUser := GetCurrentUser(request, response)
	if curUser == nil {
		Unauthenticated(request, response)
		return
	}
	//Get the file data
	imageFile, header, err := request.Request.FormFile("file-data")
	if err != nil {
		WriteError(err, response)
	}
	userId := request.PathParameter("user-id")
	rev := request.HeaderParameter("If-Match")
	if userId == "" || rev == "" {
		WriteBadRequestError(response)
		return
	}
	attType := header.Header.Get("Content-Type")
	rev, err = new(UserAvatarManager).SaveImage(userId, rev, attType,
		imageFile, curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	SetAuth(response, curUser.Auth)
	response.WriteEntity(BooleanResponse{Success: true})
	response.AddHeader("ETag", rev)
}

type imageReader func(string) (io.ReadCloser, error)

//Get an Avatar Image
func (ac AvatarController) getAvatar(request *restful.Request,
	response *restful.Response, uam *UserAvatarManager, ir imageReader) {
	userId := request.PathParameter("user-id")
	if userId == "" {
		WriteBadRequestError(response)
		return
	}
	image, err := ir(userId)
	if err != nil {
		WriteError(err, response)
		return
	}
	defer image.Close()
	response.AddHeader("Content-Type", "image/jpeg")
	if _, err := io.Copy(response.ResponseWriter, image); err != nil {
		WriteError(err, response)
	}
}

//Get an Avatar Image (large size)
func (ac AvatarController) getImage(request *restful.Request,
	response *restful.Response) {
	uam := new(UserAvatarManager)
	ac.getAvatar(request, response, uam, uam.GetLargeAvatar)
}

//Get an Avatar Image (thumb size)
func (ac AvatarController) getThumb(request *restful.Request,
	response *restful.Response) {
	uam := new(UserAvatarManager)
	ac.getAvatar(request, response, uam, uam.GetThumbnailAvatar)
}

//Generates a UserAvatar Record Response
func (ac AvatarController) genRecordResponse(curUser *CurrentUserInfo,
	userId string, ua *UserAvatar) AvatarResponse {
	return AvatarResponse{
		Links: ac.genAvatarRecordLinks(curUser.User, userId,
			ac.genAvatarUri(userId)),
		AvatarRecord: *ua,
	}
}

//Generates the Links for a UserAvatar Record Response
func (ac AvatarController) genAvatarRecordLinks(user *User,
	userId string, uri string) avatarLinks {
	links := avatarLinks{}
	userRoles := user.Roles
	admin := util.HasRole(userRoles, AdminRole(MainDbName())) ||
		util.HasRole(userRoles, MasterRole())
	write := user.UserName == userId
	links.Self = &HatLink{Href: uri, Method: "GET"}
	links.GetLargeAvatar = &HatLink{Href: uri + "/image", Method: "GET"}
	links.GetThumbnailAvatar = &HatLink{Href: uri + "/thumbnail", Method: "GET"}
	if admin || write {
		links.Update = &HatLink{Href: uri, Method: "PUT"}
		links.Delete = &HatLink{Href: uri, Method: "DELETE"}
		links.SaveImage = &HatLink{Href: uri + "/image", Method: "PUT"}
	}
	return links
}
