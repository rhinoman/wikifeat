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

package wiki_service_test

import (
	"encoding/json"
	. "github.com/rhinoman/wikifeat/common/entities"
	"github.com/rhinoman/wikifeat/wikis/wiki_service/wikit"
	"testing"
)

func cleanup(wikiId string){
	wm.Delete(wikiId, curUser)
}

func doPageTest(t *testing.T) {
	//Create a wiki
	wikiId := getUuid()
	pageId := getUuid()
	sPageId := getUuid()
	commentId := getUuid()
	sCommentId := getUuid()
	rCommentId := getUuid()
	pageSlug := ""
	wikiRecord := WikiRecord{
		Name:        "Cafe Project",
		Description: "Wiki for the Cafe Project",
	}
	_, err := wm.Create(wikiId, &wikiRecord, curUser)
	if err != nil {
		t.Error(err)
	}
	defer cleanup(wikiId)
	//Create a page with some markdown
	content := wikit.PageContent{
		Raw:       "About\n=\nAbout the project\n--\n<script type=\"text/javascript\">alert(\"no!\");</script>",
		Formatted: "",
	}
	page := wikit.Page{
		Content: content,
		Title:   "About",
	}
	//page = jsonifyPage(page)
	//Create another page
	sContent := wikit.PageContent{
		Raw:       "Contact\n=\nContact Us\n--\n",
		Formatted: "",
	}
	sPage := wikit.Page{
		Content: sContent,
		Title:   "Contact Us",
		Parent:  pageId,
	}
	//sPage = jsonifyPage(sPage)

	rev, err := pm.Save(wikiId, &page, pageId, "", curUser)
	sRev, sErr := pm.Save(wikiId, &sPage, sPageId, "", curUser)
	pageSlug = page.Slug
	if rev == ""{
		t.Error("rev is empty")
	}
	if err != nil {
		t.Error(err)
	}
	if sRev == "" {
		t.Error("sRev is empty")
	}
	if sErr != nil {
		t.Error(sErr)
	}
	//Read Page
	rPage := wikit.Page{}
	nWikiId, rev, err := pm.ReadBySlug(wikiRecord.Slug, pageSlug, &rPage, curUser)
	if nWikiId == "" {
		t.Error("wikiId is empty")
	}
	if rev == "" {
		t.Error("rev is empty!")
	}
	if err != nil {
		t.Error(err)
	}

	content = rPage.Content
	if content.Formatted != "<h1>About</h1>\n<h2>About the project</h2>\n\n"{
		t.Error("content.Formatted is wrong!")
	}
	if rPage.LastEditor != "John.Smith"{
		t.Error("rPage is not John.Smith!")
	}
	//Update Page
	rPage = wikit.Page{}
	rev, _ = pm.Read(wikiId, pageId, &rPage, curUser)
	content = wikit.PageContent{
		Raw: "About Cafe Project\n=\n",
	}
	rPage.Content = content
	//rPage.Title = "About Cafe"
	rPage = jsonifyPage(rPage)
	rev, err = pm.Save(wikiId, &rPage, pageId, rev, curUser)
	if rev == "" {
		t.Error("Rev is empty!")
	}
	if err != nil {
		t.Error(err)
	}
	//Page history
	hist, err := pm.GetHistory(wikiId, pageId, 1, 0, curUser)
	if err != nil {
		t.Error(err)
	}
	if len(hist.Rows) != 2{
		t.Errorf("Insufficient history length, should be 2 was %v", len(hist.Rows))
	}
	for _, hvr := range hist.Rows {
		t.Logf("history item: %v", hvr)
		if hvr.Value.Editor != "John.Smith"{
			t.Error("Editor is wrong!")
		}
	}
	//Page index
	index, err := pm.Index(wikiId, curUser)
	if err != nil {
		t.Error(err)
	}
	if len(index) != 2 {
		t.Errorf("Index should be length 2, was: %v", len(index))
	}
	//Child index
	index, err = pm.ChildIndex(wikiId, pageId, curUser)
	if err != nil {
		t.Error(err)
	}
	if len(index) != 1{
		t.Errorf("Index should be 1, was %v", len(index))
	}
	//Breadcrumbs
	crumbs, err := pm.GetBreadcrumbs(wikiId, pageId, curUser)
	if err != nil {
		t.Error(err)
	}
	if len(crumbs) != 1{
		t.Errorf("Should be 1 breadcrumb, was %v", len(crumbs))
	}
	crumbs, err = pm.GetBreadcrumbs(wikiId, sPageId, curUser)
	if err != nil{
		t.Error(err)
	}
	if len(crumbs) != 2{
		t.Errorf("Length should be 2, was %v", len(crumbs))
	}
	//Comments
	firstComment := wikit.Comment{
		Content: wikit.PageContent{
			Raw: "This is a comment",
		},
	}
	secondComment := wikit.Comment{
		Content: wikit.PageContent{
			Raw: "This is another comment",
		},
	}
	replyComment := wikit.Comment{
		Content: wikit.PageContent{
			Raw: "This is a reply",
		},
	}
	_, err1 := pm.SaveComment(wikiId, pageId, &firstComment,
		commentId, "", curUser)
	_, err2 := pm.SaveComment(wikiId, pageId, &secondComment,
		sCommentId, "", curUser)
	_, err3 := pm.SaveComment(wikiId, pageId, &replyComment,
		rCommentId, "", curUser)
	if err1 != nil {
		t.Error(err1)
	}
	if err2 != nil{
		t.Error(err2)
	}
	if err3 != nil{
		t.Error(err3)
	}
	//Comment queries
	comments, err := pm.GetComments(wikiId, pageId, 1, 0, curUser)
	if err != nil {
		t.Error(err)
	}
	numComments := len(comments.Rows)
	if numComments != 3 {
		t.Errorf("Wrong number of comments, should be 3 was %v", numComments)
	}
	//Read comment
	//Read the comment to get the revision
	readComment := wikit.Comment{}
	sCommentRev, err := pm.ReadComment(wikiId, sCommentId,
		&readComment, curUser)
	if err != nil{
		t.Error(err)
	}
	if sCommentRev == ""{
		t.Error("sCommentRev is empty")
	}
	t.Logf("Comment rev: %v\n", sCommentRev)
	//Comment deletion
	readComment = wikit.Comment{}
	sCommentRev, err = pm.ReadComment(wikiId, sCommentId,
		&readComment, curUser)
	t.Logf("Comment rev: %v\n", sCommentRev)
	dRev, err := pm.DeleteComment(wikiId, sCommentId, curUser)
	if err != nil{
		t.Error(err)
	}
	if dRev == ""{
		t.Error("dRev is empty!")
	}
	//Delete Page
	rPage = wikit.Page{}
	rev, err = pm.Read(wikiId, pageId, &rPage, curUser)
	t.Logf("Page Rev: %v", rev)
	if err != nil {
		t.Error(err)
	}
	dRev, err = pm.Delete(wikiId, pageId, rev, curUser)
	t.Logf("Del Rev: %v", dRev)
	if err != nil {
		t.Error(err)
	}
	if rev == ""{
		t.Error("dRev is empty!")
	}

}

func jsonifyPage(page wikit.Page) wikit.Page {
	resultPage := wikit.Page{}
	ePage, _ := json.Marshal(page)
	json.Unmarshal(ePage, &resultPage)
	return resultPage
}
