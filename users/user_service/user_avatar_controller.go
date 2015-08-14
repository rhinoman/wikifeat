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
	. "github.com/rhinoman/wikifeat/common/entities"
	. "github.com/rhinoman/wikifeat/common/services"
	"github.com/rhinoman/wikifeat/common/util"
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
		Operation("create").
		Reads(UserAvatar{}).
		Param(ws.PathParameter("user-id", "User id").DataType("string")).
		Writes(AvatarResponse{}))
}

func (ac AvatarController) genAvatarUri(userId string) string {
	theUri := ApiPrefix() + "/users/" + avatarUri + "/avatar"
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
	links.GetLargeAvatar = &HatLink{Href: uri + "/large", Method: "GET"}
	links.GetThumbnailAvatar = &HatLink{Href: uri + "/thumbnail", Method: "GET"}
	if admin || write {
		links.Update = &HatLink{Href: uri, Method: "PUT"}
		links.Delete = &HatLink{Href: uri, Method: "DELETE"}
		links.SaveImage = &HatLink{Href: uri + "/content", Method: "PUT"}
	}
	return links
}
