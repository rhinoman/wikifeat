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

package wiki_service

import (
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/emicklei/go-restful"
	. "github.com/rhinoman/wikifeat/common/database"
	. "github.com/rhinoman/wikifeat/common/entities"
	. "github.com/rhinoman/wikifeat/common/services"
	"github.com/rhinoman/wikifeat/common/util"
	"log"
	"net/http"
	"strconv"
)

type WikisController struct{}

type wikiLinks struct {
	HatLinks
	PageIndex  *HatLink `json:"index,omitempty"`
	Search     *HatLink `json:"search,omitempty"`
	CreatePage *HatLink `json:"create_page,omitempty"`
}

type WikiRecordResponse struct {
	Links      wikiLinks  `json:"_links"`
	WikiRecord WikiRecord `json:"wiki_record"`
}

type WikiIndexResponse struct {
	Links         HatLinks      `json:"_links"`
	TotalRows     int           `json:"totalRows"`
	PageNum       int           `json:"offset"`
	WikiIndexList WikiIndexList `json:"_embedded"`
}

type WikiIndexList struct {
	List []WikiRecordResponse `json:"ea:wiki"`
}

func (wc WikisController) wikiUri() string {
	return ApiPrefix() + "/wikis"
}

var wikisWebService *restful.WebService

func (wc WikisController) Service() *restful.WebService {
	return wikisWebService
}

//Define routes
func (wc WikisController) Register(container *restful.Container) {
	//pages is a subcontroller
	pc := PagesController{}
	//Files is a subcontroller
	fc := FileController{}
	wikisWebService = new(restful.WebService)
	wikisWebService.Filter(LogRequest).
		Filter(AuthUser).
		ApiVersion(ApiVersion()).
		Path(wc.wikiUri()).
		Doc("Manage Wikis").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	wikisWebService.Route(wikisWebService.GET("").To(wc.index).
		Doc("Get a list of wikis").
		Operation("index").
		Param(wikisWebService.QueryParameter("pageNum", "Page Number").DataType("integer")).
		Param(wikisWebService.QueryParameter("numPerPage", "Number of records to return").DataType("integer")).
		Param(wikisWebService.QueryParameter("memberOnly", "Only show wikis user belongs to").DataType("boolean")).
		Writes(WikiIndexResponse{}))

	wikisWebService.Route(wikisWebService.POST("").To(wc.create).
		Doc("Create a new wiki").
		Operation("create").
		Reads(WikiRecord{}).
		Writes(WikiRecordResponse{}))

	wikisWebService.Route(wikisWebService.GET("/{wiki-id}").To(wc.read).
		Doc("Fetch a Wiki Record").
		Operation("read").
		Param(wikisWebService.PathParameter("wiki-id", "Wiki identifier").DataType("string")).
		Writes(WikiRecordResponse{}))

	wikisWebService.Route(wikisWebService.GET("/slug/{wiki-slug}").To(wc.readBySlug).
		Doc("Fetch a Wiki Record by its slug").
		Operation("readBySlug").
		Param(wikisWebService.PathParameter("wiki-slug", "Wiki Slug").DataType("string")).
		Writes(WikiRecordResponse{}))

	wikisWebService.Route(wikisWebService.PUT("/{wiki-id}").To(wc.update).
		Doc("Update a Wiki Record").
		Operation("update").
		Param(wikisWebService.PathParameter("wiki-id", "Wiki identifier").DataType("string")).
		Param(wikisWebService.HeaderParameter("If-Match", "Revision").DataType("string")).
		Reads(WikiRecord{}).
		Writes(WikiRecordResponse{}))

	wikisWebService.Route(wikisWebService.DELETE("/{wiki-id}").To(wc.del).
		Doc("Delete a Wiki").
		Operation("del").
		Param(wikisWebService.PathParameter("wiki-id", "Wiki identifier").DataType("string")).
		Writes(BooleanResponse{}))

	//Add routes from pages to the wiki controller
	pc.AddRoutes(wikisWebService)
	//Add routes from files to the wiki controller
	fc.AddRoutes(wikisWebService)
	//Add the wiki controller to the container
	container.Add(wikisWebService)
}

func (wc WikisController) genWikiUri(wikiId string) string {
	return wc.wikiUri() + "/" + wikiId
}

//Get a list of wikis
func (wc WikisController) index(request *restful.Request,
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
	var memberOnly bool = false
	memberOnlyString := request.QueryParameter("memberOnly")
	if memberOnlyString != "" {
		var err error
		memberOnly, err = strconv.ParseBool(memberOnlyString)
		if err != nil {
			log.Printf("Error: %v", err)
			WriteIllegalRequestError(response)
			return
		}

	}
	wlr := WikiListResponse{}
	err := new(WikiManager).GetWikiList(pageNum, limit, memberOnly,
		&wlr, curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	wir := wc.genWikiIndexResponse(curUser.User, &wlr)
	SetAuth(response, curUser.Auth)
	response.WriteEntity(wir)
}

//Create a new wiki
func (wc WikisController) create(request *restful.Request,
	response *restful.Response) {
	curUser := GetCurrentUser(request, response)
	if curUser == nil {
		Unauthenticated(request, response)
		return
	}
	theWiki := new(WikiRecord)
	err := request.ReadEntity(theWiki)
	if err != nil {
		WriteBadRequestError(response)
		return
	}
	wikiId := GenUuid()
	rev, err := new(WikiManager).Create(wikiId, theWiki, curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	//Permissions would have changed for new wiki, re-read user record
	theUser, err := GetUserFromAuth(curUser.Auth)
	if err != nil {
		WriteError(err, response)
		return
	}
	response.AddHeader("ETag", rev)
	response.WriteHeader(http.StatusCreated)
	wr := wc.genRecordResponse(theUser, wikiId, theWiki)
	SetAuth(response, curUser.Auth)
	response.WriteEntity(wr)
}

//Read a Wiki Record
func (wc WikisController) read(request *restful.Request,
	response *restful.Response) {
	curUser := GetCurrentUser(request, response)
	if curUser == nil {
		Unauthenticated(request, response)
		return
	}
	wikiId := request.PathParameter("wiki-id")
	theWiki := new(WikiRecord)
	rev, err := new(WikiManager).Read(wikiId, theWiki, curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	response.AddHeader("ETag", rev)
	wr := wc.genRecordResponse(curUser.User, wikiId, theWiki)
	SetAuth(response, curUser.Auth)
	response.WriteEntity(wr)
}

//Fetch a Wiki Record by its slug
func (wc WikisController) readBySlug(request *restful.Request,
	response *restful.Response) {
	curUser := GetCurrentUser(request, response)
	if curUser == nil {
		Unauthenticated(request, response)
		return
	}
	wikiSlug := request.PathParameter("wiki-slug")
	theWiki := new(WikiRecord)
	rev, err := new(WikiManager).ReadBySlug(wikiSlug, theWiki, curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	response.AddHeader("ETag", rev)
	wr := wc.genRecordResponse(curUser.User, theWiki.Id, theWiki)
	SetAuth(response, curUser.Auth)
	response.WriteEntity(wr)

}

//Update a Wiki Record
func (wc WikisController) update(request *restful.Request,
	response *restful.Response) {
	curUser := GetCurrentUser(request, response)
	if curUser == nil {
		Unauthenticated(request, response)
		return
	}
	wikiId := request.PathParameter("wiki-id")
	rev := request.HeaderParameter("If-Match")
	theWR := new(WikiRecord)
	err := request.ReadEntity(theWR)
	if err != nil {
		WriteServerError(err, response)
		return
	}
	uRev, err := new(WikiManager).Update(wikiId, rev, theWR, curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	response.AddHeader("ETag", uRev)
	wr := wc.genRecordResponse(curUser.User, wikiId, theWR)
	SetAuth(response, curUser.Auth)
	response.WriteEntity(wr)
}

//Delete a Wiki Record
func (wc WikisController) del(request *restful.Request,
	response *restful.Response) {
	curUser := GetCurrentUser(request, response)
	if curUser == nil {
		Unauthenticated(request, response)
		return
	}
	wikiId := request.PathParameter("wiki-id")
	err := new(WikiManager).Delete(wikiId, curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	SetAuth(response, curUser.Auth)
	response.WriteEntity(BooleanResponse{Success: true})

}

//Generate a record response
func (wc WikisController) genRecordResponse(curUser *User,
	wikiId string, wikiRecord *WikiRecord) WikiRecordResponse {
	wikiRecord.Id = wikiId
	wrr := WikiRecordResponse{
		Links: wc.genWikiLinks(curUser.Roles,
			wikiId,
			MainDbName(),
			wc.genWikiUri(wikiId)),
		WikiRecord: *wikiRecord,
	}
	return wrr
}

func (wc WikisController) genWikiIndexResponse(curUser *User,
	wlr *WikiListResponse) WikiIndexResponse {
	wir := WikiIndexResponse{}
	wir.TotalRows = wlr.TotalRows
	wir.PageNum = wlr.Offset
	for _, row := range wlr.Rows {
		wrr := wc.genRecordResponse(curUser,
			row.Id, &row.Value)
		wir.WikiIndexList.List = append(wir.WikiIndexList.List, wrr)
	}
	wir.Links = GenIndexLinks(curUser.Roles, MainDbName(),
		wc.wikiUri())
	return wir
}

func (wc WikisController) genWikiLinks(userRoles []string,
	wikiId string, dbName string, uri string) wikiLinks {
	//First, add links for wikiRecord in main db
	links := wikiLinks{}
	//Now check admin rights for wiki db and add links
	wikiDb := "wiki_" + wikiId
	admin := util.HasRole(userRoles, AdminRole(wikiDb)) ||
		util.HasRole(userRoles, AdminRole(MainDbName())) ||
		util.HasRole(userRoles, MasterRole())
	read := util.HasRole(userRoles, ReadRole(wikiDb))
	write := util.HasRole(userRoles, WriteRole(wikiDb))
	pageUri := uri + "/pages"
	links.Self = &HatLink{Href: uri, Method: "GET"}
	if admin || read || write {
		links.PageIndex = &HatLink{Href: pageUri, Method: "GET"}
	}
	if admin || write {
		links.Update = &HatLink{Href: uri, Method: "PUT"}
		links.CreatePage = &HatLink{Href: pageUri, Method: "POST"}
	}
	if admin {
		links.Delete = &HatLink{Href: uri, Method: "DELETE"}
	}
	return links
}
