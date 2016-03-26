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

package config_service

import (
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/emicklei/go-restful"
	. "github.com/rhinoman/wikifeat/common/services"
)

/**
  We need to be able to query a few select configuration parameters at runtime.
  Might eventually expand this to allow on-the-fly configuration changes.
*/

type ConfigController struct{}

type ConfigResponse struct {
	Links      HatLinks `json:"_links"`
	ParamName  string   `json:"paramName"`
	ParamValue string   `json:"paramValue"`
}

var configWebService *restful.WebService

func (cc ConfigController) configUri() string {
	return ApiPrefix() + "/config"
}

func (cc ConfigController) Service() *restful.WebService {
	return configWebService
}

//Define routes
func (cc ConfigController) Register(container *restful.Container) {
	configWebService = new(restful.WebService)
	configWebService.Filter(LogRequest)
	configWebService.
		Path(cc.configUri()).
		Doc("Query system configuration").
		ApiVersion(ApiVersion()).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	configWebService.Route(configWebService.GET("/{section}/{parameter}").To(cc.getConfigParam).
		Doc("Get a single configuration parameter value").
		Operation("getConfigParam").
		Param(configWebService.PathParameter("section", "Section").DataType("string")).
		Param(configWebService.PathParameter("parameter", "Parameter").DataType("string")).
		Writes(ConfigResponse{}))

	container.Add(configWebService)
}

// Get an individual config parameter
func (cc ConfigController) getConfigParam(request *restful.Request,
	response *restful.Response) {
	section := request.PathParameter("section")
	param := request.PathParameter("parameter")
	value, err := new(ConfigManager).getConfigParam(section, param)
	if err != nil {
		WriteIllegalRequestError(response)
		return
	}
	configResponse := cc.genConfigResponse(param, value)
	response.WriteEntity(configResponse)
}

func (cc ConfigController) genConfigResponse(paramName, paramValue string) ConfigResponse {
	links := HatLinks{}
	uri := cc.configUri() + "/" + paramName
	links.Self = &HatLink{Href: uri, Method: "GET"}
	return ConfigResponse{
		Links:      links,
		ParamName:  paramName,
		ParamValue: paramValue,
	}
}
