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

package wiki_service_test

import (
	"encoding/json"
	"github.com/rhinoman/couchdb-go"
	. "github.com/rhinoman/wikifeat/common/entities"
	"github.com/rhinoman/wikifeat/users/user_service"
	"github.com/rhinoman/wikifeat/wikis/wiki_service/wikit"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func beforePageTest(t *testing.T) error {
	setup()
	user := User{
		UserName: "John.Smith",
		Password: "password",
	}
	registration := user_service.Registration{
		NewUser: user,
	}
	_, err := um.SetUp(&registration)
	if err != nil {
		return err
	}
	return nil
}

func TestPageCRUD(t *testing.T) {
	err := beforePageTest(t)
	jsAuth := &couchdb.BasicAuth{
		Username: "John.Smith",
		Password: "password",
	}
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	theUser := User{}
	_, err = grabUser("John.Smith", &theUser, jsAuth)
	if err != nil {
		t.Error(err)
	}
	defer afterTest(&theUser)
	//Create a wiki
	curUser := getCurUser(jsAuth)
	wikiId := getUuid()
	pageId := getUuid()
	sPageId := getUuid()
	pageSlug := ""
	wikiRecord := WikiRecord{
		Name:        "Cafe Project",
		Description: "Wiki for the Cafe Project",
	}
	_, err = wm.Create(wikiId, &wikiRecord, curUser)
	if err != nil {
		t.Error(err)
	}
	defer wm.Delete(wikiId, curUser)
	Convey("Given a Page with some basic content", t, func() {
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

		Convey("When the pages are saved", func() {
			rev, err := pm.Save(wikiId, &page, pageId, "", curUser)
			sRev, sErr := pm.Save(wikiId, &sPage, sPageId, "", curUser)
			pageSlug = page.Slug
			Convey("The revision should be set and the error should be nil", func() {
				So(rev, ShouldNotEqual, "")
				So(err, ShouldBeNil)
				So(sRev, ShouldNotEqual, "")
				So(sErr, ShouldBeNil)
			})
		})
		Convey("When the Page is Read", func() {
			rPage := wikit.Page{}
			wikiId, rev, err := pm.ReadBySlug(wikiRecord.Slug, pageSlug, &rPage, curUser)
			Convey("The revision should be set and the error should be nil", func() {
				So(wikiId, ShouldNotBeNil)
				So(rev, ShouldNotEqual, "")
				So(err, ShouldBeNil)
			})
			Convey("The Html content should be correct", func() {
				content = rPage.Content
				So(content.Formatted, ShouldEqual,
					"<h1>About</h1>\n<h2>About the project</h2>\n\n")
			})
			Convey("The LastEditor should be correct", func() {
				So(rPage.LastEditor, ShouldEqual, "John.Smith")
			})
		})
		Convey("When the Page is Updated", func() {
			rPage := wikit.Page{}
			rev, _ := pm.Read(wikiId, pageId, &rPage, curUser)
			content = wikit.PageContent{
				Raw: "About Cafe Project\n=\n",
			}
			rPage.Content = content
			//rPage.Title = "About Cafe"
			rPage = jsonifyPage(rPage)
			rev, err := pm.Save(wikiId, &rPage, pageId, rev, curUser)
			Convey("The revision should be set and the error should be nil", func() {
				So(rev, ShouldNotBeNil)
				So(err, ShouldBeNil)
			})

		})
		Convey("When the Page history is requested", func() {
			hist, err := pm.GetHistory(wikiId, pageId, 0, curUser)
			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("History should be complete", func() {
				So(len(hist.Rows), ShouldEqual, 2)
				for _, hvr := range hist.Rows {
					t.Logf("history item: %v", hvr)
					So(hvr.Value.Editor, ShouldEqual, "John.Smith")
				}
			})
		})
		Convey("When the Page index is requested", func() {
			index, err := pm.Index(wikiId, curUser)
			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Index should contain 2 items", func() {
				So(len(index), ShouldEqual, 2)
			})

		})
		Convey("When the Page child index is requested", func() {
			index, err := pm.ChildIndex(wikiId, pageId, curUser)
			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Index should have one item", func() {
				So(len(index), ShouldEqual, 1)
			})
		})
		Convey("When the Page breadcrumbs are requested for a root-level page", func() {
			crumbs, err := pm.GetBreadcrumbs(wikiId, pageId, curUser)
			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Length should be 1", func() {
				So(len(crumbs), ShouldEqual, 1)
			})
		})
		Convey("When the Page breadcrumbs are requested for a child page", func() {
			crumbs, err := pm.GetBreadcrumbs(wikiId, sPageId, curUser)
			Convey("Error should be nil", func() {
				So(err, ShouldBeNil)
			})
			Convey("Length should be 2", func() {
				So(len(crumbs), ShouldEqual, 2)
			})
		})
		Convey("When the Page is Deleted", func() {
			rPage := wikit.Page{}
			rev, err := pm.Read(wikiId, pageId, &rPage, curUser)
			t.Logf("Page Rev: %v", rev)
			if err != nil {
				t.Error(err)
			}
			dRev, err := pm.Delete(wikiId, pageId, rev, curUser)
			t.Logf("Del Rev: %v", dRev)
			Convey("Error should be nil", func() {
				So(rev, ShouldNotEqual, "")
				So(err, ShouldBeNil)
			})

		})
	})
}

func jsonifyPage(page wikit.Page) wikit.Page {
	resultPage := wikit.Page{}
	ePage, _ := json.Marshal(page)
	json.Unmarshal(ePage, &resultPage)
	return resultPage
}
