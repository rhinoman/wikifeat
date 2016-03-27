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

package notification_service

import (
	"github.com/emicklei/go-restful"
	. "github.com/rhinoman/wikifeat/common/entities"
	. "github.com/rhinoman/wikifeat/common/services"
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
