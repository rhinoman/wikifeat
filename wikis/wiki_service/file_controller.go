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
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type FileController struct{}

type fileLinks struct {
	HatLinks
	SaveAttachment *HatLink `json:"saveContent,omitempty"`
	GetAttachment  *HatLink `json:"getContent,omitempty"`
}

type FileResponse struct {
	Links fileLinks  `json:"_links"`
	File  wikit.File `json:"file"`
}

type FileIndexResponse struct {
	Links         HatLinks      `json:"_links"`
	TotalRows     int           `json:"totalRows"`
	PageNum       int           `json:"offset"`
	FileIndexList FileIndexList `json:"_embedded"`
}

type FileIndexList struct {
	List []FileIndexItem `json:"ea:file"`
}

type FileIndexItem struct {
	Links fileLinks  `json:"_links"`
	Entry wikit.File `json:"file"`
}

var fileUri = "/{wiki-id}/files"

//Define routes
func (fc FileController) AddRoutes(ws *restful.WebService) {

	ws.Route(ws.GET(fileUri).To(fc.index).
		Doc("Get list of files in this wiki").
		Operation("index").
		Param(ws.PathParameter("wiki-id", "Wiki identifier").DataType("string")).
		Param(ws.QueryParameter("pageNum", "Page Number").DataType("integer")).
		Param(ws.QueryParameter("numPerPage", "Number of records to return").DataType("integer")).
		Writes(FileIndexResponse{}))

	ws.Route(ws.POST(fileUri).To(fc.create).
		Doc("Create a new file record").
		Operation("create").
		Reads(wikit.File{}).
		Param(ws.PathParameter("wiki-id", "Wiki identifier").DataType("string")).
		Writes(FileResponse{}))

	ws.Route(ws.GET(fileUri + "/{file-id}").To(fc.read).
		Doc("Reads a File Record").
		Operation("read").
		Param(ws.PathParameter("wiki-id", "Wiki identifier").DataType("string")).
		Param(ws.PathParameter("file-id", "File identifier").DataType("string")).
		Writes(FileResponse{}))

	ws.Route(ws.PUT(fileUri + "/{file-id}").To(fc.update).Doc("Updates a file").
		Operation("update").
		Reads(wikit.File{}).
		Param(ws.PathParameter("wiki-id", "Wiki identifier").DataType("string")).
		Param(ws.PathParameter("file-id", "File identifier").DataType("string")).
		Param(ws.HeaderParameter("If-Match", "File Revision").DataType("string")).
		Writes(FileResponse{}))

	ws.Route(ws.DELETE(fileUri + "/{file-id}").To(fc.del).
		Doc("Delete a File").
		Operation("del").
		Param(ws.PathParameter("wiki-id", "Wiki identifier").DataType("string")).
		Param(ws.PathParameter("file-id", "File identifier").DataType("string")).
		Writes(BooleanResponse{}))

	ws.Route(ws.POST(fileUri + "/{file-id}/content").
		Consumes("multipart/form-data").To(fc.saveContent).
		Operation("saveContent").
		Param(ws.PathParameter("wiki-id", "Wiki identifier").DataType("string")).
		Param(ws.PathParameter("file-id", "File identifier").DataType("string")).
		Param(ws.HeaderParameter("If-Match", "File Revision").DataType("string")).
		Param(ws.FormParameter("file-data", "The File").DataType("file")).
		Writes(BooleanResponse{}))

	ws.Route(ws.GET(fileUri + "/{file-id}/content").To(fc.getContent).
		Operation("getContent").Produces("application/octet-stream").
		Param(ws.PathParameter("wiki-id", "Wiki identifier").DataType("string")).
		Param(ws.PathParameter("file-id", "File identifier").DataType("string")).
		Param(ws.QueryParameter("attName", "Attachment Name").DataType("string")).
		Param(ws.QueryParameter("download", "Download File").DataType("boolean")))
}

func (fc FileController) genFileUri(wikiId, fileId string) string {
	theUri := ApiPrefix() + "/wikis" + fileUri + "/" + fileId
	return strings.Replace(theUri, "{wiki-id}", wikiId, 1)
}

//Get file index
func (fc FileController) index(request *restful.Request,
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
		if ln, err := strconv.Atoi(limitString); err != nil {
			log.Printf("Error: %v", err)
			WriteIllegalRequestError(response)
			return
		} else {
			limit = ln
		}
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
	wikiId := request.PathParameter("wiki-id")
	fivr, err := new(FileManager).Index(wikiId, pageNum, limit, curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	fir := fc.genFileIndexResponse(curUser, wikiId, fivr)
	SetAuth(response, curUser.Auth)
	response.WriteEntity(fir)
}

//Create a new file record
func (fc FileController) create(request *restful.Request,
	response *restful.Response) {
	curUser := GetCurrentUser(request, response)
	if curUser == nil {
		Unauthenticated(request, response)
		return
	}
	wikiId := request.PathParameter("wiki-id")
	theFile := new(wikit.File)
	err := request.ReadEntity(theFile)
	if err != nil {
		WriteServerError(err, response)
		return
	}
	fileId := GenUuid()
	rev, err := new(FileManager).SaveFileRecord(wikiId, theFile, fileId, "", curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	response.AddHeader("ETag", rev)
	response.WriteHeader(http.StatusCreated)
	fr := fc.genRecordResponse(curUser, wikiId, fileId, theFile)
	SetAuth(response, curUser.Auth)
	response.WriteEntity(fr)
}

//Reads a File Record
func (fc FileController) read(request *restful.Request,
	response *restful.Response) {
	curUser := GetCurrentUser(request, response)
	if curUser == nil {
		Unauthenticated(request, response)
		return
	}
	wikiId, fileId := fc.getPathParameters(request)
	if wikiId == "" || fileId == "" {
		WriteBadRequestError(response)
		return
	}
	file := wikit.File{}
	rev, err := new(FileManager).ReadFileRecord(wikiId, &file, fileId, curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	response.AddHeader("ETag", rev)
	fr := fc.genRecordResponse(curUser, wikiId, fileId, &file)
	SetAuth(response, curUser.Auth)
	response.WriteEntity(fr)
}

//Updates a file record
func (fc FileController) update(request *restful.Request,
	response *restful.Response) {
	curUser := GetCurrentUser(request, response)
	if curUser == nil {
		Unauthenticated(request, response)
		return
	}
	wikiId, fileId := fc.getPathParameters(request)
	rev := request.HeaderParameter("If-Match")
	if wikiId == "" || fileId == "" || rev == "" {
		WriteBadRequestError(response)
		return
	}
	theFile := new(wikit.File)
	err := request.ReadEntity(theFile)
	if err != nil {
		WriteServerError(err, response)
		return
	}
	rev, err = new(FileManager).SaveFileRecord(wikiId, theFile, fileId, rev, curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	response.AddHeader("ETag", rev)
	fr := fc.genRecordResponse(curUser, wikiId, fileId, theFile)
	SetAuth(response, curUser.Auth)
	response.WriteEntity(fr)
}

//Deletes a file
func (fc FileController) del(request *restful.Request,
	response *restful.Response) {
	curUser := GetCurrentUser(request, response)
	if curUser == nil {
		Unauthenticated(request, response)
		return
	}
	wikiId, fileId := fc.getPathParameters(request)
	if wikiId == "" || fileId == "" {
		WriteBadRequestError(response)
		return
	}
	rev, err := new(FileManager).DeleteFile(wikiId, fileId, curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	SetAuth(response, curUser.Auth)
	response.WriteEntity(BooleanResponse{Success: true})
	response.AddHeader("ETag", rev)
}

//Saves a file's attachment content
func (fc FileController) saveContent(request *restful.Request,
	response *restful.Response) {
	curUser := GetCurrentUser(request, response)
	if curUser == nil {
		Unauthenticated(request, response)
		return
	}
	//Get the file data
	theFile, header, err := request.Request.FormFile("file-data")
	if err != nil {
		WriteError(err, response)
	}
	wikiId, fileId := fc.getPathParameters(request)
	rev := request.HeaderParameter("If-Match")
	if wikiId == "" || fileId == "" || rev == "" {
		WriteBadRequestError(response)
		return
	}
	attType := header.Header.Get("Content-Type")
	attName := header.Filename
	rev, err = new(FileManager).SaveFileAttachment(wikiId, fileId, rev,
		attName, attType, theFile, curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	SetAuth(response, curUser.Auth)
	response.WriteEntity(BooleanResponse{Success: true})
	response.AddHeader("ETag", rev)
}

//Fetches a file's attachment content
func (fc FileController) getContent(request *restful.Request,
	response *restful.Response) {
	fm := new(FileManager)
	curUser := GetCurrentUser(request, response)
	if curUser == nil {
		Unauthenticated(request, response)
		return
	}
	wikiId, fileId := fc.getPathParameters(request)
	attName := request.QueryParameter("attName")
	if wikiId == "" || fileId == "" || attName == "" {
		WriteBadRequestError(response)
		return
	}
	//First, read the file record
	fileRecord := wikit.File{}
	rev, err := fm.ReadFileRecord(wikiId, &fileRecord, fileId, curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	att, ok := fileRecord.Attachments[attName]
	if ok == false {
		WriteBadRequestError(response)
		return
	}
	attType := att.MimeType
	attSize := att.Length
	reader, err := fm.GetFileAttachment(wikiId, fileId, rev,
		attType, attName, curUser)
	if err != nil {
		WriteError(err, response)
		return
	}
	defer reader.Close()
	SetAuth(response, curUser.Auth)
	response.AddHeader("Content-Type", attType)
	response.AddHeader("Content-Length", strconv.Itoa(attSize))
	response.AddHeader("ETag", rev)
	if download, err := strconv.ParseBool(request.QueryParameter("download")); err != nil {
	} else if download {
		response.AddHeader("Content-Disposition", "attachment; filename=\""+attName+"\"")
	}
	if bytesWritten, err := io.Copy(response.ResponseWriter, reader); err != nil {
		WriteError(err, response)
	} else {
		log.Printf("Downloaded File: " + attName + ", " +
			strconv.FormatInt(bytesWritten, 10) + " bytes written")
	}
}

func (fc FileController) getPathParameters(request *restful.Request) (string, string) {
	return request.PathParameter("wiki-id"), request.PathParameter("file-id")
}

//Generate the File Response
func (fc FileController) genRecordResponse(curUser *CurrentUserInfo,
	wikiId, fileId string, file *wikit.File) FileResponse {
	file.Id = fileId
	return FileResponse{
		Links: fc.genFileRecordLinks(curUser.User.Roles, "wiki_"+wikiId,
			fc.genFileUri(wikiId, fileId)),
		File: *file,
	}
}

//Generate an individual index record for the response
func (fc FileController) genIndexEntryRecordResponse(curUser *CurrentUserInfo,
	wikiId, fileId string, fileEntry *wikit.File) FileIndexItem {
	return FileIndexItem{
		Links: fc.genFileRecordLinks(curUser.User.Roles, "wiki_"+wikiId,
			fc.genFileUri(wikiId, fileId)),
		Entry: *fileEntry,
	}
}

//Generate the file index response
func (fc FileController) genFileIndexResponse(curUser *CurrentUserInfo, wikiId string,
	fivr *wikit.FileIndexViewResponse) FileIndexResponse {
	var fiis []FileIndexItem
	for _, row := range fivr.Rows {
		row.Value.Id = row.Id
		frr := fc.genIndexEntryRecordResponse(curUser,
			wikiId, row.Id, &row.Value)
		fiis = append(fiis, frr)
	}
	theUri := ApiPrefix() + "/wikis/" + wikiId + "/files"
	return FileIndexResponse{
		Links:     GenIndexLinks(curUser.User.Roles, "wiki_"+wikiId, theUri),
		TotalRows: fivr.TotalRows,
		PageNum:   fivr.Offset,
		FileIndexList: FileIndexList{
			List: fiis,
		},
	}
}

func (fc FileController) genFileRecordLinks(userRoles []string,
	wikiDb string, uri string) fileLinks {
	links := fileLinks{}
	admin := util.HasRole(userRoles, AdminRole(wikiDb))
	write := util.HasRole(userRoles, WriteRole(wikiDb))
	links.Self = &HatLink{Href: uri, Method: "GET"}
	links.GetAttachment = &HatLink{Href: uri + "/content", Method: "GET"}
	if admin || write {
		links.Update = &HatLink{Href: uri, Method: "PUT"}
		links.Delete = &HatLink{Href: uri, Method: "DELETE"}
		links.SaveAttachment = &HatLink{Href: uri + "/content", Method: "PUT"}
	}
	return links
}
