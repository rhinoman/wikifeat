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
	. "github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/wikifeat-common/entities"
	. "github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/wikifeat-common/services"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/wikifeat-common/util"
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

type CommentResponse struct {
	Links   HatLinks      `json:"_links"`
	Comment wikit.Comment `json:"comment"`
}

type CommentIndexResponse struct {
	Links     HatLinks         `json:"_links"`
	TotalRows int              `json:"total_rows"`
	Offset    int              `json:"offset"`
	Entries   CommentIndexList `json:"_embedded"`
}

type CommentIndexList struct {
	List []CommentResponse `json:"ea:comment"`
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
		Param(ws.QueryParameter("pageNum", "Page number for pagination").DataType("integer")).
		Param(ws.QueryParameter("numPerPage", "Number of records to return").DataType("integer")).
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

	ws.Route(ws.POST(pageUri + "/{page-id}/comments").To(pc.createComment).
		Doc("Creates a Comment for this page").
		Operation("create").
		Param(ws.PathParameter("wiki-id", "Wiki identifier").DataType("string")).
		Param(ws.PathParameter("page-id", "Page identifier").DataType("string")).
		Reads(wikit.Comment{}).
		Writes(CommentResponse{}))

	ws.Route(ws.GET(pageUri + "/{page-id}/comments/{comment-id}").To(pc.readComment).
		Doc("Reads a Comment").
		Operation("read").
		Param(ws.PathParameter("wiki-id", "Wiki identifier").DataType("string")).
		Param(ws.PathParameter("page-id", "Page identifier").DataType("string")).
		Param(ws.PathParameter("comment-id", "Comment identifier").DataType("string")).
		Writes(CommentResponse{}))

	ws.Route(ws.PUT(pageUri + "/{page-id}/comments/{comment-id}").To(pc.updateComment).
		Doc("Updates a Comment").
		Operation("update").
		Param(ws.PathParameter("wiki-id", "Wiki identifier").DataType("string")).
		Param(ws.PathParameter("page-id", "Page identifier").DataType("string")).
		Param(ws.PathParameter("comment-id", "Comment identifier").DataType("string")).
		Param(ws.HeaderParameter("If-Match", "Comment Revision").DataType("string")).
		Reads(wikit.Comment{}).
		Writes(CommentResponse{}))

	ws.Route(ws.DELETE(pageUri + "/{page-id}/comments/{comment-id}").To(pc.delComment).
		Doc("Deletes a Comment").
		Operation("delete").
		Param(ws.PathParameter("wiki-id", "Wiki identifier").DataType("string")).
		Param(ws.PathParameter("page-id", "Page identifier").DataType("string")).
		Param(ws.PathParameter("comment-id", "Comment identifier").DataType("string")).
		Writes(BooleanResponse{}))

	ws.Route(ws.GET(pageUri + "/{page-id}/comments").To(pc.commentIndex).
		Doc("Get all comments for this page").
		Operation("index").
		Param(ws.PathParameter("wiki-id", "Wiki identifier").DataType("string")).
		Param(ws.PathParameter("page-id", "Page identifier").DataType("string")).
		Param(ws.QueryParameter("pageNum", "Page number for pagination").DataType("integer")).
		Param(ws.QueryParameter("numPerPage", "Number of records to return").DataType("integer")).
		Writes(CommentIndexResponse{}))

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
	numPerPage, err := strconv.Atoi(request.QueryParameter("numPerPage"))
	if err != nil {
		numPerPage = 50
	}
	pageNum, err := strconv.Atoi(request.QueryParameter("pageNum"))
	if err != nil {
		pageNum = 1
	}
	if wikiId == "" || pageId == "" {
		WriteBadRequestError(response)
		return
	}
	history, err := new(PageManager).GetHistory(wikiId, pageId, pageNum, numPerPage, curUser)
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

//Create a comment
func (pc PagesController) createComment(request *restful.Request,
	response *restful.Response) {
	curUser := GetCurrentUser(request, response)
	if curUser == nil {
		Unauthenticated(request, response)
		return
	}
	wikiId := request.PathParameter("wiki-id")
	pageId := request.PathParameter("page-id")
	theComment := new(wikit.Comment)
	err := request.ReadEntity(theComment)
	if err != nil {
		WriteServerError(err, response)
		return
	}
	commentId := GenUuid()
	rev, err := new(PageManager).SaveComment(wikiId, pageId, theComment,
		commentId, "", curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	response.AddHeader("ETag", rev)
	response.WriteHeader(http.StatusCreated)
	cr := pc.genCommentRecordResponse(curUser, wikiId, pageId, commentId, theComment)
	SetAuth(response, curUser.Auth)
	response.WriteEntity(cr)
}

//Read a comment
func (pc PagesController) readComment(request *restful.Request,
	response *restful.Response) {
	curUser := GetCurrentUser(request, response)
	if curUser == nil {
		Unauthenticated(request, response)
		return
	}
	wikiId := request.PathParameter("wiki-id")
	pageId := request.PathParameter("page-id")
	commentId := request.PathParameter("comment-id")
	if wikiId == "" || pageId == "" || commentId == "" {
		WriteBadRequestError(response)
		return
	}
	comment := wikit.Comment{}
	rev, err := new(PageManager).ReadComment(wikiId, commentId,
		&comment, curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	response.AddHeader("ETag", rev)
	cr := pc.genCommentRecordResponse(curUser, wikiId, pageId, commentId, &comment)
	SetAuth(response, curUser.Auth)
	response.WriteEntity(cr)
}

//Update a comment
func (pc PagesController) updateComment(request *restful.Request,
	response *restful.Response) {
	curUser := GetCurrentUser(request, response)
	if curUser == nil {
		Unauthenticated(request, response)
		return
	}
	wikiId := request.PathParameter("wiki-id")
	pageId := request.PathParameter("page-id")
	commentId := request.PathParameter("comment-id")
	rev := request.HeaderParameter("If-Match")
	if wikiId == "" || pageId == "" || commentId == "" || rev == "" {
		WriteBadRequestError(response)
		return
	}
	theComment := new(wikit.Comment)
	err := request.ReadEntity(theComment)
	if err != nil {
		WriteServerError(err, response)
		return
	}
	pm := new(PageManager)
	rev, err = pm.SaveComment(wikiId, pageId, theComment,
		commentId, rev, curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	response.AddHeader("ETag", rev)
	readComment := wikit.Comment{}
	cr := CommentResponse{}
	if _, err = pm.ReadComment(wikiId, commentId, &readComment, curUser); err != nil {
		cr = pc.genCommentRecordResponse(curUser, wikiId, pageId, commentId, theComment)
	} else {
		cr = pc.genCommentRecordResponse(curUser, wikiId, pageId, commentId, &readComment)
	}
	SetAuth(response, curUser.Auth)
	response.WriteEntity(cr)
}

//Delete a comment
func (pc PagesController) delComment(request *restful.Request,
	response *restful.Response) {
	curUser := GetCurrentUser(request, response)
	if curUser == nil {
		Unauthenticated(request, response)
		return
	}
	wikiId := request.PathParameter("wiki-id")
	pageId := request.PathParameter("page-id")
	commentId := request.PathParameter("comment-id")
	if wikiId == "" || pageId == "" || commentId == "" {
		WriteBadRequestError(response)
		return
	}
	rev, err := new(PageManager).DeleteComment(wikiId, commentId, curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	SetAuth(response, curUser.Auth)
	response.WriteEntity(BooleanResponse{Success: true})
	response.AddHeader("ETag", rev)

}

//Gets all comments for a page
func (pc PagesController) commentIndex(request *restful.Request,
	response *restful.Response) {
	curUser := GetCurrentUser(request, response)
	if curUser == nil {
		Unauthenticated(request, response)
		return
	}
	wikiId := request.PathParameter("wiki-id")
	pageId := request.PathParameter("page-id")
	numPerPage, err := strconv.Atoi(request.QueryParameter("numPerPage"))
	if err != nil {
		numPerPage = 50
	}
	pageNum, err := strconv.Atoi(request.QueryParameter("pageNum"))
	if err != nil {
		pageNum = 1
	}
	if wikiId == "" || pageId == "" {
		WriteBadRequestError(response)
		return
	}
	cList, err := new(PageManager).GetComments(wikiId, pageId, pageNum, numPerPage, curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	SetAuth(response, curUser.Auth)
	cr := pc.genCommentIndexResponse(curUser, wikiId, pageId, cList)
	response.WriteEntity(cr)
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

func (pc PagesController) genCommentRecordResponse(curUser *CurrentUserInfo,
	wikiId string, pageId string, commentId string, comment *wikit.Comment) CommentResponse {
	comment.Id = commentId
	cr := CommentResponse{
		Links: pc.genCommentRecordLinks(curUser, "wiki_"+wikiId,
			pc.genPageUri(wikiId, pageId)+"/comments/"+commentId,
			comment.Author),
		Comment: *comment,
	}
	return cr
}

// Permissions are a bit different for comments, so we need a custom function
func (pc PagesController) genCommentRecordLinks(curUser *CurrentUserInfo,
	dbName string, commentUri string, commentAuthor string) HatLinks {
	links := HatLinks{}
	userRoles := curUser.User.Roles
	admin := util.HasRole(userRoles, AdminRole(dbName)) ||
		util.HasRole(userRoles, AdminRole(MainDbName())) ||
		util.HasRole(userRoles, MasterRole())
	ownComment := commentAuthor == curUser.User.UserName
	//Generate the self link
	links.Self = &HatLink{Href: commentUri, Method: "GET"}
	//Update links
	if admin || ownComment {
		links.Update = &HatLink{Href: commentUri, Method: "PUT"}
		links.Delete = &HatLink{Href: commentUri, Method: "DELETE"}
	}
	return links
}

func (pc PagesController) genCommentIndexResponse(curUser *CurrentUserInfo,
	wikiId string, pageId string, commentList *wikit.CommentIndexViewResponse) CommentIndexResponse {
	commentsUri := pc.genPageUri(wikiId, pageId) + "/comments"
	userRoles := curUser.User.Roles
	dbName := "wiki_" + wikiId
	admin := util.HasRole(userRoles, AdminRole(dbName)) ||
		util.HasRole(userRoles, AdminRole(MainDbName())) ||
		util.HasRole(userRoles, MasterRole())
	write := util.HasRole(userRoles, WriteRole(dbName))
	indexLinks := HatLinks{}
	indexLinks.Self = &HatLink{Href: commentsUri, Method: "GET"}
	if admin || write {
		indexLinks.Create = &HatLink{Href: commentsUri, Method: "POST"}
	}
	var entries []CommentResponse
	for _, com := range commentList.Rows {
		entries = append(entries,
			pc.genCommentRecordResponse(curUser, wikiId, pageId, com.Id, &com.Value))
	}
	return CommentIndexResponse{
		Links:     indexLinks,
		TotalRows: commentList.TotalRows,
		Offset:    commentList.Offset,
		Entries: CommentIndexList{
			List: entries,
		},
	}
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
