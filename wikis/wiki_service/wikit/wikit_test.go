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
package wikit_test

import (
	"bytes"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/couchdb-go"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/twinj/uuid"
	. "github.com/rhinoman/wikifeat/wikis/wiki_service/wikit"
	"io/ioutil"
	"strconv"
	"testing"
	"time"
)

const server = "127.0.0.1"
const timeout = time.Duration(500 * time.Millisecond)
const unittestdb = "unittestdb"

var adminAuth = &couchdb.BasicAuth{Username: "adminuser", Password: "password"}
var connection, _ = couchdb.NewConnection(server, 5984, timeout)

type TestContent struct {
	Name     string    `json:"name"`
	Location string    `json:"location"`
	Time     time.Time `json:"time"`
}

type TestPage struct {
	Page
	Name     string
	Location string
	Time     time.Time
}

func getUuid() string {
	return uuid.Formatter(uuid.NewV4(), uuid.Clean)
}

func printError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Error: %v\n", err)
	}
}

func deleteDb(t *testing.T, wikiName string) {
	err := connection.DeleteDB(wikiName, adminAuth)
	printError(t, err)
}

func deleteUser(t *testing.T, username string) {
	userdata := couchdb.UserRecord{}
	rev, _ := connection.GetUser(username, &userdata, adminAuth)
	connection.DeleteUser(username, rev, adminAuth)
}

func createTestWiki(t *testing.T) (string, string) {
	wikiName := "wiki_" + getUuid()
	newSteve := "steve" + getUuid()[1:4]
	connection.AddUser(newSteve, "password", []string{}, adminAuth)
	err := CreateWiki(connection, adminAuth, newSteve, wikiName)
	printError(t, err)
	t.Logf("Created Wiki with id: %v\n", wikiName)
	return wikiName, newSteve
}

func TestCreateWiki(t *testing.T) {
	wikiName, user := createTestWiki(t)
	deleteDb(t, wikiName)
	deleteUser(t, user)
	t.Logf("Wiki id: %v\n", wikiName)
}

func TestPages(t *testing.T) {
	wikiName, user := createTestWiki(t)
	ba := &couchdb.BasicAuth{user, "password"}
	theWiki := SelectWiki(connection, wikiName, ba)
	t.Logf("dbname: %v\n", wikiName)
	t.Logf("user: %v\n", user)
	content := PageContent{
		Raw: "Name: Concert, Location: Town Square",
	}
	page := Page{
		Title:   "The Story",
		Owner:   "Steve",
		Content: content,
	}
	dupPage := Page{
		Title:   "The Story",
		Owner:   "Bill",
		Content: content,
	}
	theId := getUuid()
	rev, err := theWiki.SavePage(&page, theId, "", "steve")
	printError(t, err)
	t.Logf("new doc rev: %v\n", rev)
	//Try to save a dup page
	rev, err = theWiki.SavePage(&dupPage, getUuid(), "", "steve")
	if err == nil {
		t.Errorf("Dup was saved, should not have been.")
	}
	//Save it with a diff title
	dupPage.Title = "Not the Story"
	dupUuid := getUuid()
	rev, err = theWiki.SavePage(&dupPage, dupUuid, "", "steve")
	printError(t, err)
	//Now, create a conflict
	dupPage.Title = "The Story"
	rev, err = theWiki.SavePage(&dupPage, dupUuid, rev, "steve")
	if err == nil {
		t.Errorf("Dup was saved, should not have been.")
	}
	//Read the dup
	rDup := Page{}
	rev, err = theWiki.ReadPage(dupUuid, &rDup)
	printError(t, err)
	if rDup.Title != "Not The Story" && rDup.Slug != "not-the-story" {
		t.Errorf("Title and Slug ain't right\nTitle: %v\nSlug: %v\n",
			rDup.Title, rDup.Slug)
	}
	//now read it
	rPage := Page{}

	rev, err = theWiki.ReadPage(theId, &rPage)
	//now update it
	page.Content.Raw = "This is the first edit"
	rev, err = theWiki.SavePage(&page, theId, rev, "bill")
	printError(t, err)
	t.Logf("first update rev: %v\n", rev)
	//now update it again
	page.Content.Raw = "This is the second edit"
	rev, err = theWiki.SavePage(&page, theId, rev, "sarah")
	printError(t, err)
	t.Logf("second update rev: %v\n", rev)
	//Create another page
	anotherPage := Page{
		Title:  "The Sequel",
		Owner:  "Joe",
		Parent: theId,
		Content: PageContent{
			Raw: "Howdy!\n===\n",
		},
	}
	sId := getUuid()
	sRev, err := theWiki.SavePage(&anotherPage, sId, "", "joe")
	printError(t, err)
	t.Logf("2nd page: %v\n", anotherPage)
	t.Logf("2nd doc rev: %v\n", sRev)
	//Get both pages
	multiPages := MultiPageResponse{}
	err = theWiki.ReadMultiplePages([]string{theId, sId}, &multiPages)
	printError(t, err)
	if len(multiPages.Rows) != 2 {
		t.Errorf("Not enough documents returned!")
	}
	t.Logf("MultiPage 1: %v\n", multiPages.Rows[0].Doc)
	t.Logf("MultiPage 2: %v\n", multiPages.Rows[1].Doc)
	//Get child page Index
	pIndex, err := theWiki.GetChildPageIndex(theId)
	printError(t, err)
	t.Logf("Child Page Index: %v", pIndex)
	if len(pIndex) < 1 {
		t.Errorf("length of child page index should not be 0")
	}
	//Get page index
	index, err := theWiki.GetPageIndex()
	printError(t, err)
	//Get the history
	history, err := theWiki.GetHistory(theId, 1, 0)
	printError(t, err)
	if history.TotalRows != 3 {
		t.Errorf("Wrong number of Rows reported!")
	}
	if len(index) != 3 {
		t.Errorf("len index should be 3, was %v", len(history.Rows))
	}
	for i, pIdx := range index {
		t.Logf("Index %v, %v: %v\n", i, pIdx.Key, pIdx.Value)
	}
	for i, hist := range history.Rows {
		t.Logf("History %v, %v: %v\n", i, hist.Key[1], hist.Value)
	}
	//now delete it
	rev, err = theWiki.DeletePage(theId, rev)
	printError(t, err)
	t.Logf("Deleted rev: %v\n", rev)
	deleteDb(t, wikiName)
	deleteUser(t, user)
}

func TestReadPage(t *testing.T) {
	wikiName, user := createTestWiki(t)
	ba := &couchdb.BasicAuth{user, "password"}
	theWiki := SelectWiki(connection, wikiName, ba)
	t.Logf("dbname: %v\n", wikiName)
	page := Page{
		Title: "The Story",
		Owner: "Steve",
		Content: PageContent{
			Raw: "This is the original",
		},
	}
	theId := getUuid()
	rev, err := theWiki.SavePage(&page, theId, "", "steve")
	printError(t, err)
	t.Logf("new doc rev: %v\n", rev)
	//now update it
	page.Content.Raw = "This is the first edit"
	rev, err = theWiki.SavePage(&page, theId, rev, "bill")
	printError(t, err)
	t.Logf("first update rev: %v\n", rev)
	//now read it
	readPage := Page{}
	readRev, err := theWiki.ReadPage(theId, &readPage)
	printError(t, err)
	t.Logf("read doc: %v\n", readPage)
	t.Logf("read doc rev: %v\n", readRev)
	if readPage.Content.Raw != "This is the first edit" ||
		readPage.Title != "The Story" ||
		readRev == "" {
		t.Fail()
	}
	//now read it by its slug
	slug := readPage.Slug
	readPage = Page{}
	readRev, err = theWiki.ReadPageBySlug(slug, &readPage)
	printError(t, err)
	t.Logf("read doc by slug: %v\n", readPage)
	t.Logf("read doc rev by slug: %v\n", readRev)
	deleteDb(t, wikiName)
	deleteUser(t, user)

}

func TestFiles(t *testing.T) {
	wikiName, user := createTestWiki(t)
	ba := &couchdb.BasicAuth{user, "password"}
	theWiki := SelectWiki(connection, wikiName, ba)
	t.Logf("dbname: %v\n", wikiName)
	file := File{
		Name:        "tps_report.txt",
		Description: "TPS Report 9999",
	}
	theId := getUuid()
	rev, err := theWiki.SaveFileRecord(&file, theId, "", "Steve")
	printError(t, err)
	t.Logf("File Rev: %v\n", rev)
	content := []byte("TPS REPORT 9999\nWidgets widgetized this week: 42\n")
	contentReader := bytes.NewReader(content)
	//Add the attachment to the file record
	uRev, err := theWiki.SaveFileAttachment(theId, rev, "tps_report.txt", "text/plain",
		contentReader)
	printError(t, err)
	t.Logf("Updated Rev: %v\n", uRev)
	//Now try to read it
	theContent, err := theWiki.GetFileAttachment(theId, uRev, "text/plain", "tps_report.txt")
	t.Logf("read content")
	printError(t, err)
	defer theContent.Close()
	theBytes, err := ioutil.ReadAll(theContent)
	printError(t, err)
	data := string(theBytes[:])
	if data != "TPS REPORT 9999\nWidgets widgetized this week: 42\n" {
		t.Errorf("Content was a lie!")
	}
	t.Logf("CONTENT: %v", data)
	//Get the File index
	fileIndex, err := theWiki.GetFileIndex(0, 0)
	printError(t, err)
	if len(fileIndex.Rows) != 1 {
		t.Errorf("File index is wrong length: " + strconv.Itoa(len(fileIndex.Rows)))
	}
	if fileIndex.Rows[0].Key != "tps_report.txt" {
		t.Errorf("File Record has wrong name! " + fileIndex.Rows[0].Key)
	}
	if fileIndex.Rows[0].Value.Description != "TPS Report 9999" {
		t.Errorf("Description is wrong! " + fileIndex.Rows[0].Value.Description)
	}
	dRev, err := theWiki.DeleteFileRecord(theId, uRev)
	printError(t, err)
	if dRev == "" {
		t.Errorf("Deleted Rev blank.")
	}
	deleteDb(t, wikiName)
	deleteUser(t, user)
}

func TestComments(t *testing.T) {
	wikiName, user := createTestWiki(t)
	defer deleteDb(t, wikiName)
	defer deleteUser(t, user)
	ba := &couchdb.BasicAuth{user, "password"}
	theWiki := SelectWiki(connection, wikiName, ba)
	t.Logf("dbnam: %v\n", wikiName)
	page := Page{
		Title: "Page of Commenting",
		Owner: "Steve",
		Content: PageContent{
			Raw: "This is the original",
		},
	}
	pageId := getUuid()
	rev, err := theWiki.SavePage(&page, pageId, "", "steve")
	printError(t, err)
	t.Logf("new doc rev: %v\n", rev)
	//Now let's create a comment
	theId := getUuid()
	comment := Comment{
		Content: PageContent{
			Raw: "This is a commment",
		},
	}
	rev, err = theWiki.SaveComment(&comment, theId, "", pageId, "steve")
	t.Logf("new doc rev: %v\n", rev)
	//now update it
	comment.Content.Raw = "I done had to edit my comment!"
	uRev, err := theWiki.SaveComment(&comment, theId, rev, "", "")
	printError(t, err)
	t.Logf("first update rev: %v\n", uRev)
	//now read it
	readComment := Comment{}
	readRev, err := theWiki.ReadComment(theId, &readComment)
	printError(t, err)
	t.Logf("read comment: %v\n", readComment)
	t.Logf("read doc rev: %v\n", readRev)
	if readRev != uRev {
		t.Errorf("revisions don't match!")
	}
	if readComment.Content.Raw != "I done had to edit my comment!" {
		t.Errorf("Content ain't right!")
	}
	//Create a second comment
	secondComment := Comment{
		Content: PageContent{
			Raw: "This is another comment",
		},
	}
	secondId := getUuid()
	sRev, err := theWiki.SaveComment(&secondComment, secondId, "", pageId, "steve")
	printError(t, err)
	t.Logf("Second comment rev: %v\n", sRev)
	//Create a comment, replying to the first one
	replyComment := Comment{
		Content: PageContent{
			Raw: "This is a reply",
		},
	}
	replyId := getUuid()
	rRev, err := theWiki.SaveComment(&replyComment, replyId, "", pageId, "steve")
	printError(t, err)
	t.Logf("Update rev: %v\n", rRev)
	//Get a list of all comments
	civr, err := theWiki.GetCommentsForPage(pageId, 1, 0)
	printError(t, err)
	t.Logf("CommentList: %v\n", civr)
	if civr.TotalRows != 3 {
		t.Errorf("Total Rows should be 3, but was: %v", civr.TotalRows)
	}
	//now delete it!
	_, err = theWiki.DeleteComment(theId, uRev)
	printError(t, err)
}
