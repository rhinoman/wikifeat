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

package notification_service

import (
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/emicklei/go-restful"
	. "github.com/rhinoman/wikifeat/common/entities"
	. "github.com/rhinoman/wikifeat/common/services"
	//"log"
)

type NotificationsController struct{}

func (nc NotificationsController) notificationUri() string {
	return ApiPrefix() + "/notifications"
}

var notificationsWebService *restful.WebService

func (nc NotificationsController) Service() *restful.WebService {
	return notificationsWebService
}

//Define routes
func (nc NotificationsController) Register(container *restful.Container) {
	notificationsWebService := new(restful.WebService)
	notificationsWebService.Filter(LogRequest).
		ApiVersion(ApiVersion()).
		Path(nc.notificationUri()).
		Doc("Send Notifications").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	notificationsWebService.Route(notificationsWebService.POST("{notification-id}/send").
		To(nc.send).
		Operation("send").
		Param(notificationsWebService.PathParameter("notification-id",
		"Notification Template").DataType("string")).
		Reads(NotificationRequest{}).
		Writes(BooleanResponse{}))

	//Add the notifications controller to the container
	container.Add(notificationsWebService)
}

func (nc NotificationsController) genNotifUri(notifId string) string {
	return nc.notificationUri() + "/" + notifId
}

func (nc NotificationsController) send(request *restful.Request,
	response *restful.Response) {
	notificationId := request.PathParameter("notification-id")
	nr := new(NotificationRequest)
	err := request.ReadEntity(nr)
	if err != nil || !nc.ValidateNotificationRequest(nr) {
		WriteBadRequestError(response)
		return
	}
	err = new(NotificationManager).Send(notificationId, nr)
	if err != nil {
		WriteError(err, response)
		return
	}
	response.WriteEntity(BooleanResponse{Success: true})
}

func (nc NotificationsController) ValidateNotificationRequest(nr *NotificationRequest) bool {
	if nr.To == "" || nr.Subject == "" || nr.Data == nil {
		return false
	} else {
		return true
	}
}
