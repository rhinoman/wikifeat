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
package wikit

import (
	"errors"
	. "github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/couchdb-go"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/go-slugification"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/twinj/uuid"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Wiki struct {
	db       *Database
	wikiName string
}

func getUuid() string {
	theUuid := uuid.NewV4()
	return uuid.Formatter(theUuid, uuid.Clean)
}

//Returns a Wiki struct for the requested wiki
func SelectWiki(connection *Connection, wikiName string, auth Auth) *Wiki {
	return &Wiki{
		db:       connection.SelectDB(wikiName, auth),
		wikiName: wikiName,
	}
}

//Creates a new Wiki with name wikiName.
//owner specifies a username, the first admin.
//Returns an error, if any
func CreateWiki(conn *Connection,
	adminAuth Auth, owner string, wikiName string) error {

	//if something goes wrong, delete the database
	cleanup := func() {
		conn.DeleteDB(wikiName, adminAuth)
	}
	//Create a new database for the wiki
	err := conn.CreateDB(wikiName, adminAuth)
	if err != nil {
		return err
	}
	//...and initialize it
	db := conn.SelectDB(wikiName, adminAuth)
	err = InitDb(db, wikiName)
	if err != nil {
		//Uncreate the database
		cleanup()
		return err
	}
	//Grant the owner the admin role
	_, err = conn.GrantRole(owner, wikiName+":admin", adminAuth)
	if err != nil {
		cleanup()
		return err
	}
	return nil
}

//Retrieves the latest copy of the Page specified by id
//Returns the page and the revision of the page, or an error
func (wiki *Wiki) ReadPage(id string, page *Page) (string, error) {
	rev, err := wiki.db.Read(id, &page, nil)
	if err != nil {
		return "", err
	} else {
		page.Id = id
		return rev, nil
	}
}

// Retrieves multiple pages in one request, given an array of page ids
func (wiki *Wiki) ReadMultiplePages(ids []string, mpr *MultiPageResponse) error {
	return wiki.db.ReadMultiple(ids, mpr)
}

//Gets page by Slug
func (wiki *Wiki) ReadPageBySlug(slug string, page *Page) (string, error) {
	response := SlugViewResponse{}
	err := wiki.db.GetView("wikit", "getPageBySlug", &response, SetKey(slug))
	if err != nil {
		return "", err
	}
	if len(response.Rows) > 0 {
		*page = response.Rows[0].Value.Page
		page.Id = response.Rows[0].Id
		return response.Rows[0].Value.Rev, nil
	} else {
		return "", errors.New("[Error]:404: Page not found")
	}
}

// Get a page's lineage information.
func (wiki *Wiki) GetLineage(pageId string, page *Page) ([]string, error) {
	if parent := page.Parent; parent != "" {
		//Get the parent page so we can extract its lineage
		parentPage := Page{}
		if _, err := wiki.ReadPage(parent, &parentPage); err == nil {
			parentLineage := parentPage.Lineage
			theLineage := append(parentLineage, pageId)
			return theLineage, nil
		} else {
			return []string{}, err
		}
	} else {
		return []string{pageId}, nil
	}
}

// Checks for duplicate page slugs
// We're trying to enforce a unique constraint on page slugs
func (wiki *Wiki) CheckForDuplicateSlug(slug string) error {
	params := SetKey(slug)
	response := KVResponse{}
	params.Add("group", "true")
	err := wiki.db.GetView("wikit", "checkUniqueSlug", &response, params)
	if err != nil {
		return err
	} else if len(response.Rows) <= 0 {
		return nil
	}
	theRecord := response.Rows[0]
	if theRecord.Value > 1 {
		return &Error{
			StatusCode: 409,
			Reason:     "Duplicate Page slug found",
		}
	} else {
		return nil
	}
}

// Validation of Page
func (page Page) Validate() error {
	err := &Error{
		StatusCode: 400,
	}
	if page.DocType != "page" {
		err.Reason = "Type must be page"
		return err
	}
	if page.Title == "" || len(page.Title) > 128 {
		err.Reason = "Page title invalid"
		return err
	}
	if page.Owner == "" || page.LastEditor == "" {
		err.Reason = "Page must have an owner and editor"
		return err
	}
	if page.OwningPage == "" {
		err.Reason = "Owning Page not set"
		return err
	}
	return nil
}

// Validation of Comment
func (comment Comment) Validate() error {
	err := &Error{
		StatusCode: 400,
	}
	if comment.Author == "" {
		err.Reason = "Comment has no author!"
		return err
	}
	if comment.OwningPage == "" {
		err.Reason = "Comment must have an owning page!"
		return err
	}
	return nil
}

func (wiki *Wiki) createPage(page *Page, id string, editor string) (string, error) {
	//Must be a new document eh, if not, it will error out on the write
	page.DocType = "page"
	page.Timestamp = time.Now().UTC()
	page.Owner = editor
	page.LastEditor = editor
	page.Slug = slugification.Slugify(page.Title)
	//A new document is its own owner.
	page.OwningPage = id
	//Set the lineage
	if lineage, err := wiki.GetLineage(id, page); err == nil {
		page.Lineage = lineage
	}
	if err := page.Validate(); err != nil {
		return "", err
	}
	sRev, err := wiki.db.Save(page, id, "")
	//Now, check for duplicate slug
	err = wiki.CheckForDuplicateSlug(page.Slug)
	if err != nil {
		//Delete the page we just created
		if _, dErr := wiki.db.Delete(id, sRev); dErr != nil {
			return "", dErr
		} else {
			return "", err
		}
	} else {
		return sRev, nil
	}
}

func (wiki *Wiki) updatePage(page *Page, id string, rev string, editor string) (string, error) {
	//Updating an existing document
	//Read in the original page
	rPage := Page{}
	_, err := wiki.ReadPage(id, &rPage)
	if err != nil {
		return "", err
	}
	//Validate the page
	page.DocType = "page"
	page.LastEditor = editor
	page.Timestamp = time.Now().UTC()
	page.Slug = slugification.Slugify(page.Title)
	page.OwningPage = id
	page.Owner = rPage.Owner
	if err = page.Validate(); err != nil {
		return "", err
	}

	//copy the current document
	copyId := getUuid()
	copyRev, err := wiki.db.Copy(id, rev, copyId)
	if err != nil {
		return "", err
	}
	//Clear the slug field in the copied page
	copyPage := new(Page)
	_, err = wiki.ReadPage(copyId, copyPage)
	if err != nil {
		wiki.DeletePage(copyId, copyRev)
		return "", err
	}
	//Remeber this slug just in case
	prevSlug := copyPage.Slug
	copyPage.Slug = ""
	wiki.db.Save(copyPage, copyId, copyRev)
	//now save
	rev, err = wiki.db.Save(page, id, rev)
	dsErr := wiki.CheckForDuplicateSlug(page.Slug)
	if dsErr != nil {
		//delete the file copy
		wiki.db.Delete(copyId, copyRev)
		//restore the previous page
		copyPage.Slug = prevSlug
		if _, cpErr := wiki.db.Save(copyPage, id, rev); cpErr != nil {
			return "", cpErr
		}
		return "", dsErr
	}
	if err != nil {
		//Need to undo the file copy.
		//wasted space, otherwise
		wiki.db.Delete(copyId, copyRev)
		return "", err
	} else {
		return rev, err
	}
}

//Saves a page (creates or updates), creating history entries as necessary
//Returns the revision of the new page or an error
func (wiki *Wiki) SavePage(page *Page, id string, rev string, editor string) (string, error) {
	if rev == "" {
		return wiki.createPage(page, id, editor)
	} else {
		return wiki.updatePage(page, id, rev, editor)
	}
}

//Delete a Page
func (wiki *Wiki) DeletePage(id string, rev string) (string, error) {
	//Fetch the document's history first
	history, err := wiki.GetHistory(id, 1, 0)
	if err != nil {
		return "", err
	}
	//Now delete the document
	dRev, err := wiki.db.Delete(id, rev)
	if err != nil {
		return "", err
	}
	//Now cleanup old versions
	for _, entry := range history.Rows {
		//These are fire-and-forget calls.
		//If they fail, it's not the end of the world, we just
		//have some wasted space in the database.
		//And they really shouldn't fail unless someone was mucking around
		//with the historical documents.
		go wiki.db.Delete(entry.Value.DocumentId, entry.Value.DocumentRev)
	}
	return dRev, nil
}

//Gets a History List
//Returns the History List or an error
func (wiki *Wiki) GetHistory(documentId string, pageNum int,
	numPerPage int) (*HistoryViewResponse, error) {
	response := HistoryViewResponse{}
	theKeys := SetKeys([]string{documentId, "{}"}, []string{documentId})
	/* This function gets the "count" by calling the reduce function on the getHistory
	 * couch view function
	 */
	countChan := make(chan int)
	//Grab the count concurrently
	go wiki.getCountForView("wikit", "getHistory", documentId, countChan)
	if numPerPage != 0 {
		theKeys.Add("limit", strconv.Itoa(numPerPage))
	}
	skip := numPerPage * (pageNum - 1)
	if skip > 0 {
		theKeys.Add("skip", strconv.Itoa(skip))
	}
	theKeys.Add("reduce", "false")
	err := wiki.db.GetView("wikit", "getHistory", &response, theKeys)
	if err != nil {
		return nil, err
	} else if len(response.Rows) <= 0 {
		return nil, nil
	} else {
		//Set total rows to the count
		response.TotalRows = <-countChan
		return &response, nil
	}
}

//Gets a list of a page's child pages
func (wiki *Wiki) GetChildPageIndex(pageId string) (PageIndex, error) {
	response := PageIndexViewResponse{}
	theKeys := SetKey(pageId)
	theKeys.Add("reduce", "false")
	err := wiki.db.GetView("wikit", "getChildPageIndex", &response, theKeys)
	if err != nil {
		return nil, err
	} else if len(response.Rows) <= 0 {
		return nil, nil
	} else {
		index := response.Rows
		return index, nil
	}
}

//Gets a list of all pages in a wiki.
//Returns the Page List or an error
func (wiki *Wiki) GetPageIndex() (PageIndex, error) {
	response := PageIndexViewResponse{}
	params := url.Values{}
	params.Add("reduce", "false")
	err := wiki.db.GetView("wikit", "getIndex", &response, &params)
	if err != nil {
		return nil, err
	} else {
		index := response.Rows
		return index, nil
	}
}

//Gets a list of all files in a wiki
//Retus the File List or an error
func (wiki *Wiki) GetFileIndex(pageNum int,
	numPerPage int) (*FileIndexViewResponse, error) {
	response := FileIndexViewResponse{}
	params := url.Values{}
	if numPerPage != 0 {
		params.Add("limit", strconv.Itoa(numPerPage))
	}
	skip := numPerPage * (pageNum - 1)
	if skip > 0 {
		params.Add("skip", strconv.Itoa(skip))
	}
	params.Add("reduce", "false")
	err := wiki.db.GetView("wikit", "getFileIndex", &response, &params)
	if err != nil {
		return nil, err
	} else {
		return &response, nil
	}
}

//Saves a File record
func (wiki *Wiki) SaveFileRecord(file *File, fileId string, rev string,
	uploadedBy string) (string, error) {
	file.DocType = "file"
	file.UploadedBy = uploadedBy
	file.Timestamp = time.Now().UTC()
	return wiki.db.Save(file, fileId, rev)
}

//Reads a File record
func (wiki *Wiki) GetFileRecord(fileId string, file *File) (string, error) {
	if rev, err := wiki.db.Read(fileId, file, nil); err != nil {
		return "", err
	} else {
		file.Id = fileId
		return rev, nil
	}
}

//Deletes a file record
func (wiki *Wiki) DeleteFileRecord(fileId string, rev string) (string, error) {
	return wiki.db.Delete(fileId, rev)
}

//Save file attachment
func (wiki *Wiki) SaveFileAttachment(fileId, fileRev, attName, attType string,
	attContent io.Reader) (string, error) {
	return wiki.db.SaveAttachment(fileId, fileRev, attName, attType, attContent)
}

//Get file attachment
func (wiki *Wiki) GetFileAttachment(fileId, fileRev,
	attType string, attName string) (io.ReadCloser, error) {
	return wiki.db.GetAttachment(fileId, fileRev, attType, attName)
}

//Get file attachment by proxy
func (wiki *Wiki) GetFileAttachmentByProxy(fileId, fileRev,
	attType string, attName string, r *http.Request, w http.ResponseWriter) error {
	return wiki.db.GetAttachmentByProxy(fileId, fileRev, attType, attName, r, w)
}

//---Comments Stuff-----//

// Get a comment by Id
func (wiki *Wiki) ReadComment(commentId string, comment *Comment) (string, error) {
	if rev, err := wiki.db.Read(commentId, comment, nil); err != nil {
		return "", err
	} else {
		comment.Id = commentId
		return rev, nil
	}
}

// Save a comment
func (wiki *Wiki) SaveComment(comment *Comment, id string,
	rev string, pageId string, author string) (string, error) {
	if rev == "" {
		return wiki.createComment(comment, id, pageId, author)
	} else {
		return wiki.updateComment(comment, id, rev)
	}
}

// Create a new comment
func (wiki *Wiki) createComment(comment *Comment, id string,
	pageId string, author string) (string, error) {
	comment.DocType = "comment"
	nowTime := time.Now().UTC()
	comment.CreatedTime = nowTime
	comment.ModifiedTime = nowTime
	comment.Author = author
	comment.OwningPage = pageId
	if err := comment.Validate(); err != nil {
		return "", err
	} else {
		return wiki.db.Save(comment, id, "")
	}
}

// Update a comment
func (wiki *Wiki) updateComment(comment *Comment, id string,
	rev string) (string, error) {
	//First, read the comment out of the database
	readComment := Comment{}
	if _, err := wiki.db.Read(id, &readComment, nil); err != nil {
		return "", err
	} else {
		readComment.ModifiedTime = time.Now().UTC()
		readComment.Content = comment.Content
		if err := readComment.Validate(); err != nil {
			return "", err
		} else {
			return wiki.db.Save(&readComment, id, rev)
		}
	}
}

// Delete a comment
func (wiki *Wiki) DeleteComment(id string, rev string) (string, error) {
	return wiki.db.Delete(id, rev)
}

// Get All Comments for a page
func (wiki *Wiki) GetCommentsForPage(pageId string, pageNum int,
	numPerPage int) (*CommentIndexViewResponse, error) {
	response := CommentIndexViewResponse{}
	theKeys := SetKeys([]string{pageId, "{}"}, []string{pageId})
	// This function gets the "count" by calling the reduce function
	countChan := make(chan int)
	//Grab the count concurrently
	go wiki.getCountForView("wikit_comments", "getCommentsForPage", pageId, countChan)
	if numPerPage != 0 {
		theKeys.Add("limit", strconv.Itoa(numPerPage))
	}
	skip := numPerPage * (pageNum - 1)
	if skip > 0 {
		theKeys.Add("skip", strconv.Itoa(skip))
	}
	theKeys.Add("reduce", "false")
	err := wiki.db.GetView("wikit_comments", "getCommentsForPage", &response, theKeys)
	if err != nil {
		return nil, err
	} else {
		//Set total rows to the count
		response.TotalRows = <-countChan
		return &response, nil
	}
}

// Assumes the Reduce function for the view is "_count"
// Result of _count is written to the channel 'c'
func (wiki *Wiki) getCountForView(ddoc string, view string, key string, c chan int) {
	params := SetKeys([]string{key, "{}"}, []string{key})
	params.Add("group", "true")
	params.Add("group_level", "1")
	params.Add("reduce", "true")
	var kvResp struct {
		Rows []struct {
			Keys  []string `json:"key"`
			Value int      `json:"value"`
		} `json:"rows"`
	}
	err := wiki.db.GetView(ddoc, view, &kvResp, params)
	if err != nil {
		log.Printf("Error counting history: %v", err)
		c <- 0
	} else if len(kvResp.Rows) == 0 {
		c <- 0
	} else {
		rec := kvResp.Rows[0]
		c <- rec.Value
	}
}
