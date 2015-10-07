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
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/microcosm-cc/bluemonday"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/go-commonmark"
	. "github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/wikifeat-common/entities"
	. "github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/wikifeat-common/services"
	"github.com/rhinoman/wikifeat/wikis/wiki_service/wikit"
	"log"
	"regexp"
)

type PageManager struct{}

type Breadcrumb struct {
	Name   string `json:"name"`
	PageId string `json:"pageId"`
	WikiId string `json:"wikiId"`
	Parent string `json:"parent"`
}

func wikiDbString(wikiId string) string {
	return "wiki_" + wikiId
}

//Gets a list of pages for a given wiki
func (pm *PageManager) Index(wiki string,
	curUser *CurrentUserInfo) (wikit.PageIndex, error) {
	auth := curUser.Auth
	theWiki := wikit.SelectWiki(Connection, wikiDbString(wiki), auth)
	return theWiki.GetPageIndex()
}

//Gets a list of child pages for a given document
func (pm *PageManager) ChildIndex(wiki string, pageId string,
	curUser *CurrentUserInfo) (wikit.PageIndex, error) {
	auth := curUser.Auth
	theWiki := wikit.SelectWiki(Connection, wikiDbString(wiki), auth)
	return theWiki.GetChildPageIndex(pageId)
}

//Gets a list of breadcrumbs for the current page
func (pm *PageManager) GetBreadcrumbs(wiki string, pageId string,
	curUser *CurrentUserInfo) ([]Breadcrumb, error) {
	thePage := wikit.Page{}
	if _, err := pm.Read(wiki, pageId, &thePage, curUser); err == nil {
		crumbs := []Breadcrumb{}
		response := wikit.MultiPageResponse{}
		theWiki := wikit.SelectWiki(Connection, wikiDbString(wiki), curUser.Auth)
		if szLineage := len(thePage.Lineage); szLineage > 1 {
			lineage := thePage.Lineage[0 : szLineage-1]
			if err = theWiki.ReadMultiplePages(lineage, &response); err != nil {
				return nil, err
			}
		}
		//Add the current page to the end of the list
		currentPageRow := wikit.MultiPageRow{
			Id:  pageId,
			Doc: thePage,
		}
		rows := append(response.Rows, currentPageRow)
		for _, row := range rows {
			theDoc := row.Doc
			parent := ""
			if len(theDoc.Lineage) >= 2 {
				parent = theDoc.Lineage[len(theDoc.Lineage)-2]
			}
			crumb := Breadcrumb{
				Name:   theDoc.Title,
				PageId: row.Id,
				WikiId: wiki,
				Parent: parent,
			}
			crumbs = append(crumbs, crumb)
		}
		return crumbs, nil
	} else {
		return nil, err
	}

}

//Creates or Updates a page
//Returns the revision number, if successful
func (pm *PageManager) Save(wiki string, page *wikit.Page,
	pageId string, pageRev string, curUser *CurrentUserInfo) (string, error) {
	auth := curUser.Auth
	theUser := curUser.User
	//Read the content from the page
	//parse the markdown to Html
	out := make(chan string)
	//Convert (Sanitized) Markdown to HTML
	go processMarkdown(page.Content.Raw, out)
	page.Content.Formatted = <-out
	//Store the thing, if you have the auth
	theWiki := wikit.SelectWiki(Connection, wikiDbString(wiki), auth)
	return theWiki.SavePage(page, pageId, pageRev, theUser.UserName)
}

//Read a page
//Pass an empty page to hold the data. returns the revision
func (pm *PageManager) Read(wiki string, pageId string,
	page *wikit.Page, curUser *CurrentUserInfo) (string, error) {
	auth := curUser.Auth
	theWiki := wikit.SelectWiki(Connection, wikiDbString(wiki), auth)
	return theWiki.ReadPage(pageId, page)
}

// Read a page by its slug.
// Assume the wiki Id passed in is a slug also
// Returns the WikiId, the Page Rev, and an error
func (pm *PageManager) ReadBySlug(wikiSlug string, pageSlug string,
	page *wikit.Page, curUser *CurrentUserInfo) (string, string, error) {
	// Need to get the true wiki Id from the slug
	auth := curUser.Auth
	mainDbName := MainDbName()
	mainDb := Connection.SelectDB(mainDbName, auth)
	response := WikiSlugViewResponse{}
	err := mainDb.GetView("wiki_query",
		"getWikiBySlug",
		&response,
		wikit.SetKey(wikiSlug))
	if err != nil {
		return "", "", err
	}
	if len(response.Rows) > 0 {
		wikiId := response.Rows[0].Id
		theWiki := wikit.SelectWiki(Connection, wikiDbString(wikiId), auth)
		pageRev, err := theWiki.ReadPageBySlug(pageSlug, page)
		return wikiId, pageRev, err
	} else {
		return "", "", NotFoundError()
	}
}

//Delete a page.  Returns the revision, if successful
func (pm *PageManager) Delete(wiki string, pageId string,
	pageRev string, curUser *CurrentUserInfo) (string, error) {
	auth := curUser.Auth
	theWiki := wikit.SelectWiki(Connection, wikiDbString(wiki), auth)
	//Load the page
	thePage := wikit.Page{}
	if _, err := theWiki.ReadPage(pageId, &thePage); err != nil {
		return "", err
	} else if thePage.OwningPage != pageId {
		//Thou shalt not delete historical revisions
		return "", BadRequestError()
	}
	//check if this is a 'home page'
	wm := WikiManager{}
	wr := WikiRecord{}
	if wRev, err := wm.Read(wiki, &wr, curUser); err != nil {
		return "", err
	} else if wr.HomePageId == pageId {
		//This is a home page, so clear the Wiki Record's home page Id
		wr.HomePageId = ""
		_, err = wm.Update(wiki, wRev, &wr, curUser)
		if err != nil {
			return "", err
		}
	}
	return theWiki.DeletePage(pageId, pageRev)
}

//Gets the history for this page
func (pm *PageManager) GetHistory(wiki string, pageId string,
	limit int, curUser *CurrentUserInfo) (*wikit.HistoryViewResponse, error) {
	auth := curUser.Auth
	theWiki := wikit.SelectWiki(Connection, wikiDbString(wiki), auth)
	return theWiki.GetHistory(pageId, limit)
}

//Converts markdown text to html
func processMarkdown(mdText string, out chan string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Parsing Markdown failed: ", err)
			out <- ""
		}
	}()
	//Remove harmful HTML from Raw Markdown Text
	p := getSanitizerPolicy()
	document := commonmark.ParseDocument(mdText, 0)
	htmlString := document.RenderHtml(commonmark.CMARK_OPT_DEFAULT)
	document.Free()
	out <- p.Sanitize(htmlString)
}

func getSanitizerPolicy() *bluemonday.Policy {
	p := bluemonday.UGCPolicy()
	p.AllowAttrs("data-plugin").
		Matching(regexp.MustCompile(`[\p{L}\p{N}\s\-_',:\[\]!\./\\\(\)&]*`)).Globally()
	p.AllowAttrs("data-id").
		Matching(regexp.MustCompile(`[\p{L}\p{N}\s\-_',:\[\]!\./\\\(\)&]*`)).Globally()
	return p
}
