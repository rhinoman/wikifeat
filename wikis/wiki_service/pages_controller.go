/**
* Copyright (c) 2014-present James Adam.  All rights reserved.
*
* This file is part of WikiFeat
*
*     WikiFeat is free software: you can redistribute it and/or modify
* it under the terms of the GNU General Public License as published by
* the Free Software Foundation, either version 2 of the License, or
* (at your option) any later version.
*
*     WikiFeat is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU General Public License for more details.
*
*     You should have received a copy of the GNU General Public License
* along with WikiFeat.  If not, see <http://www.gnu.org/licenses/>.
 */

package wiki_service

import (
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/emicklei/go-restful"
	. "github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/wikifeat-common/entities"
	. "github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/wikifeat-common/services"
	"github.com/rhinoman/wikifeat/wikis/wiki_service/wikit"
	"net/http"
	"strconv"
	"strings"
)

type PagesController struct{}
type PageResponse struct {
	Links HatLinks   `json:"_links"`
	Page  wikit.Page `json:"page"`
}

type PageIndexItem struct {
	Links HatLinks             `json:"_links"`
	Entry wikit.PageIndexEntry `json:"page"`
}

type PageIndexResponse struct {
	Links         HatLinks      `json:"_links"`
	PageIndexList PageIndexList `json:"_embedded"`
}

type PageIndexList struct {
	List []PageIndexItem `json:"ea:page"`
}

type BreadcrumbsResponse struct {
	Crumbs []Breadcrumb `json:"crumbs"`
}

type HistoryResponse struct {
	Links     HatLinks       `json:"_links"`
	TotalRows int            `json:"total_rows"`
	Offset    int            `json:"offset"`
	Entries   HistoryEntries `json:"_embedded"`
}
type HistoryEntries struct {
	EntryList []HistoryEntryResponse `json:"ea:history_entry"`
}

type HistoryEntryResponse struct {
	Links        HatLinks     `json:"_links"`
	HistoryEntry HistoryEntry `json:"history_entry"`
}

type HistoryEntry struct {
	Timestamp   string `json:"timestamp"`
	Editor      string `json:"editor"`
	ContentSize int    `json:"contentSize"`
	DocumentId  string `json:"documentId"`
	DocumentRev string `json:"documentRev"`
}

var pageUri = "/{wiki-id}/pages"

//Define routes
func (pc PagesController) AddRoutes(ws *restful.WebService) {

	ws.Route(ws.GET(pageUri).To(pc.index).
		Doc("Get list of pages in this wiki").
		Operation("index").
		Param(ws.PathParameter("wiki-id", "Wiki identifier").DataType("string")).
		Writes(PageIndexResponse{}))

	ws.Route(ws.GET(pageUri + "/{page-id}/children").To(pc.childIndex).
		Doc("Get list of children for this page").
		Operation("childIndex").
		Param(ws.PathParameter("wiki-id", "Wiki identifier").DataType("string")).
		Param(ws.PathParameter("page-id", "Page identifier").DataType("string")).
		Writes(PageIndexResponse{}))

	ws.Route(ws.GET(pageUri + "/{page-id}/breadcrumbs").To(pc.breadcrumbs).
		Doc("Get a list of this page's ancestry.").
		Operation("breadcrumbs").
		Param(ws.PathParameter("wiki-id", "Wiki identifier").DataType("string")).
		Param(ws.PathParameter("page-id", "Page identifier").DataType("string")).
		Writes(BreadcrumbsResponse{}))

	ws.Route(ws.POST(pageUri).To(pc.create).
		Doc("Create a new Page").
		Operation("create").
		Reads(wikit.Page{}).
		Param(ws.PathParameter("wiki-id", "Wiki identifier").DataType("string")).
		Writes(PageResponse{}))

	ws.Route(ws.GET(pageUri + "/{page-id}").To(pc.read).
		Doc("Reads a Page").
		Operation("read").
		Param(ws.PathParameter("wiki-id", "Wiki identifier").DataType("string")).
		Param(ws.PathParameter("page-id", "Page identifier").DataType("string")).
		Writes(PageResponse{}))

	ws.Route(ws.GET(pageUri + "/{page-id}/history").To(pc.history).
		Doc("Gets a Page's history").
		Operation("history").
		Param(ws.PathParameter("wiki-id", "Wiki identifier").DataType("string")).
		Param(ws.PathParameter("page-id", "Page identifier").DataType("string")).
		Param(ws.QueryParameter("limit", "Number of records to return").DataType("integer")).
		Writes(HistoryResponse{}))

	ws.Route(ws.GET("/slug/{wiki-slug}/pages/{page-slug}").To(pc.readBySlug).
		Doc("Reads a Page by its slug").
		Operation("readBySlug").
		Param(ws.PathParameter("wiki-slug", "Wiki Slug").DataType("string")).
		Param(ws.PathParameter("page-slug", "Page Slug").DataType("string")).
		Writes(PageResponse{}))

	ws.Route(ws.PUT(pageUri + "/{page-id}").To(pc.update).
		Doc("Updates a Page").
		Operation("update").
		Param(ws.PathParameter("wiki-id", "Wiki identifier").DataType("string")).
		Param(ws.PathParameter("page-id", "Page identifier").DataType("string")).
		Param(ws.HeaderParameter("If-Match", "Page revision").DataType("string")).
		Reads(wikit.Page{}).
		Writes(PageResponse{}))

	ws.Route(ws.DELETE(pageUri + "/{page-id}").To(pc.del).
		Doc("Deletes a Page").
		Operation("del").
		Param(ws.PathParameter("wiki-id", "Wiki identifier").DataType("string")).
		Param(ws.PathParameter("page-id", "Page identifier").DataType("string")).
		Param(ws.HeaderParameter("If-Match", "Page revision").DataType("string")).
		Writes(BooleanResponse{}))
}

func (pc PagesController) genPageUri(wikiId string, pageId string) string {
	theUri := ApiPrefix() + "/wikis" + pageUri + "/" + pageId
	return strings.Replace(theUri, "{wiki-id}", wikiId, 1)
}

func (pc PagesController) getIndexResponse(wikiId string,
	curUser *CurrentUserInfo,
	pIndex wikit.PageIndex) PageIndexResponse {

	var indexList []PageIndexItem
	indexResponse := PageIndexResponse{}
	for _, pr := range pIndex {
		pii := PageIndexItem{}
		pid := pr.Id
		pie := pr.Value
		pii.Entry = pie
		pii.Links = GenRecordLinks(curUser.User.Roles,
			"wiki_"+wikiId, pc.genPageUri(wikiId, pid))
		indexList = append(indexList, pii)
	}
	indexResponse.PageIndexList = PageIndexList{List: indexList}
	return indexResponse
}

//Get Page index
func (pc PagesController) index(request *restful.Request,
	response *restful.Response) {
	curUser := GetCurrentUser(request, response)
	if curUser == nil {
		Unauthenticated(request, response)
		return
	}
	wikiId := request.PathParameter("wiki-id")
	pIndex, err := new(PageManager).Index(wikiId, curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	indexResponse := pc.getIndexResponse(wikiId, curUser, pIndex)
	wikiUri := ApiPrefix() + "/wikis/" + wikiId
	indexResponse.Links = GenIndexLinks(curUser.User.Roles,
		"wiki_"+wikiId, wikiUri)
	SetAuth(response, curUser.Auth)
	response.WriteEntity(indexResponse)
}

//Get Child Page index
func (pc PagesController) childIndex(request *restful.Request,
	response *restful.Response) {
	curUser := GetCurrentUser(request, response)
	if curUser == nil {
		Unauthenticated(request, response)
		return
	}
	wikiId := request.PathParameter("wiki-id")
	pageId := request.PathParameter("page-id")
	pIndex, err := new(PageManager).ChildIndex(wikiId, pageId, curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	indexResponse := pc.getIndexResponse(wikiId, curUser, pIndex)
	indexResponse.Links = GenRecordLinks(curUser.User.Roles,
		"wiki_"+wikiId, pc.genPageUri(wikiId, pageId))
	SetAuth(response, curUser.Auth)
	response.WriteEntity(indexResponse)
}

//Get breadcrumbs
func (pc PagesController) breadcrumbs(request *restful.Request,
	response *restful.Response) {
	curUser := GetCurrentUser(request, response)
	if curUser == nil {
		Unauthenticated(request, response)
	}
	wikiId := request.PathParameter("wiki-id")
	pageId := request.PathParameter("page-id")
	breadcrumbs, err := new(PageManager).GetBreadcrumbs(wikiId, pageId, curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	SetAuth(response, curUser.Auth)
	response.WriteEntity(BreadcrumbsResponse{Crumbs: breadcrumbs})
}

//Create a Page
func (pc PagesController) create(request *restful.Request,
	response *restful.Response) {
	curUser := GetCurrentUser(request, response)
	if curUser == nil {
		Unauthenticated(request, response)
		return
	}
	wikiId := request.PathParameter("wiki-id")
	thePage := new(wikit.Page)
	err := request.ReadEntity(thePage)
	if err != nil {
		WriteServerError(err, response)
		return
	}
	pageId := GenUuid()
	rev, err := new(PageManager).Save(wikiId, thePage, pageId, "", curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	response.AddHeader("ETag", rev)
	response.WriteHeader(http.StatusCreated)
	pr := pc.genRecordResponse(curUser, wikiId, pageId, thePage)
	SetAuth(response, curUser.Auth)
	response.WriteEntity(pr)
}

//Read a Page
func (pc PagesController) read(request *restful.Request,
	response *restful.Response) {

	curUser := GetCurrentUser(request, response)
	if curUser == nil {
		Unauthenticated(request, response)
		return
	}
	wikiId := request.PathParameter("wiki-id")
	pageId := request.PathParameter("page-id")
	if wikiId == "" || pageId == "" {
		WriteBadRequestError(response)
		return
	}
	page := wikit.Page{}
	rev, err := new(PageManager).Read(wikiId, pageId, &page, curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	response.AddHeader("ETag", rev)
	pr := pc.genRecordResponse(curUser, wikiId, pageId, &page)
	SetAuth(response, curUser.Auth)
	response.WriteEntity(pr)
}

//Get a page's history
func (pc PagesController) history(request *restful.Request,
	response *restful.Response) {

	curUser := GetCurrentUser(request, response)
	if curUser == nil {
		Unauthenticated(request, response)
		return
	}
	wikiId := request.PathParameter("wiki-id")
	pageId := request.PathParameter("page-id")
	limit, err := strconv.Atoi(request.QueryParameter("limit"))
	if err != nil {
		limit = 50
	}
	if wikiId == "" || pageId == "" {
		WriteBadRequestError(response)
		return
	}
	history, err := new(PageManager).GetHistory(wikiId, pageId, limit, curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	SetAuth(response, curUser.Auth)
	hr := pc.genHistoryResponse(curUser, wikiId, pageId, history)
	response.WriteEntity(hr)
}

//Read a Page by its slug
func (pc PagesController) readBySlug(request *restful.Request,
	response *restful.Response) {

	curUser := GetCurrentUser(request, response)
	if curUser == nil {
		Unauthenticated(request, response)
		return
	}
	wikiSlug := request.PathParameter("wiki-slug")
	pageSlug := request.PathParameter("page-slug")
	if wikiSlug == "" || pageSlug == "" {
		WriteBadRequestError(response)
		return
	}
	page := wikit.Page{}
	wikiId, rev, err := new(PageManager).ReadBySlug(wikiSlug, pageSlug, &page, curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	response.AddHeader("ETag", rev)
	pr := pc.genRecordResponse(curUser, wikiId, page.Id, &page)
	SetAuth(response, curUser.Auth)
	response.WriteEntity(pr)
}

//Update a Page
func (pc PagesController) update(request *restful.Request,
	response *restful.Response) {
	curUser := GetCurrentUser(request, response)
	if curUser == nil {
		Unauthenticated(request, response)
		return
	}
	wikiId := request.PathParameter("wiki-id")
	pageId := request.PathParameter("page-id")
	rev := request.HeaderParameter("If-Match")
	if wikiId == "" || pageId == "" || rev == "" {
		WriteBadRequestError(response)
		return
	}
	thePage := new(wikit.Page)
	err := request.ReadEntity(thePage)
	if err != nil {
		WriteServerError(err, response)
		return
	}
	rev, err = new(PageManager).Save(wikiId, thePage, pageId, rev, curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	response.AddHeader("ETag", rev)
	pr := pc.genRecordResponse(curUser, wikiId, pageId, thePage)
	SetAuth(response, curUser.Auth)
	response.WriteEntity(pr)

}

//Delete a Page
func (pc PagesController) del(request *restful.Request,
	response *restful.Response) {
	curUser := GetCurrentUser(request, response)
	if curUser == nil {
		Unauthenticated(request, response)
		return
	}
	wikiId := request.PathParameter("wiki-id")
	pageId := request.PathParameter("page-id")
	rev := request.HeaderParameter("If-Match")
	if wikiId == "" || pageId == "" || rev == "" {
		WriteBadRequestError(response)
		return
	}
	rev, err := new(PageManager).Delete(wikiId, pageId, rev, curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	SetAuth(response, curUser.Auth)
	response.WriteEntity(BooleanResponse{Success: true})
	response.AddHeader("ETag", rev)
}

func (pc PagesController) genRecordResponse(curUser *CurrentUserInfo,
	wikiId string, pageId string, page *wikit.Page) PageResponse {
	page.Id = pageId
	pr := PageResponse{
		Links: GenRecordLinks(curUser.User.Roles, "wiki_"+wikiId,
			pc.genPageUri(wikiId, pageId)),
		Page: *page,
	}
	return pr
}

//This is gnarly
func (pc PagesController) genHistoryResponse(curUser *CurrentUserInfo,
	wikiId string, pageId string, history *wikit.HistoryViewResponse) HistoryResponse {
	historyUri := pc.genPageUri(wikiId, pageId) + "/history"
	indexLinks := HatLinks{
		Self: &HatLink{Href: historyUri, Method: "GET"},
	}
	var entries []HistoryEntryResponse
	for _, he := range history.Rows {
		values := he.Value
		entries = append(entries, HistoryEntryResponse{
			Links: HatLinks{
				Self: &HatLink{Href: pc.genPageUri(wikiId, he.Id), Method: "GET"},
			},
			HistoryEntry: HistoryEntry{
				Timestamp:   he.Key[1],
				Editor:      values.Editor,
				ContentSize: values.ContentSize,
				DocumentId:  values.DocumentId,
				DocumentRev: values.DocumentRev,
			},
		})
	}
	return HistoryResponse{
		Links:     indexLinks,
		TotalRows: history.TotalRows,
		Offset:    history.Offset,
		Entries: HistoryEntries{
			EntryList: entries,
		},
	}
}
