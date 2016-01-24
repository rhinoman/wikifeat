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

package wikit

import (
	"fmt"
	. "github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/couchdb-go"
	"net/url"
	"reflect"
	"strings"
)

type DesignDocument struct {
	Language string          `json:"language"`
	Views    map[string]View `json:"views"`
}

type AuthDesignDocument struct {
	Language          string `json:"language"`
	ValidateDocUpdate string `json:"validate_doc_update"`
}

type View struct {
	Map    string `json:"map"`
	Reduce string `json:"reduce,omitempty"`
}

var wikiViews = map[string]View{
	"getHistory": {
		Map: `
		function(doc) {
			if(doc.type==="page"){
				var owningPage = doc.owningPage || doc.owning_page;
				emit([owningPage, doc.timestamp],
					{documentId: doc._id,
					 documentRev: doc._rev,
					 editor: doc.editor,
					 contentSize: doc.content.raw.length}
				);
			}
		}`,
		Reduce: "_count",
	},
	"getIndex": {
		Map: `
			function(doc){
				if(doc.type==="page"){
					var owningPage = doc.owningPage || doc.owning_page;
					if(doc._id === owningPage){
						emit(doc.title, {
							id: doc._id,
							slug: doc.slug,
							title: doc.title,
							owner: doc.owner,
							editor: doc.editor,
							timestamp: doc.timestamp
						});
					}
				}
			}`,
		Reduce: "_count",
	},
	"getFileIndex": {
		Map: `
			function(doc){
				if(doc.type==="file"){
					emit(doc.name, doc);
				}
			}`,
		Reduce: "_count",
	},
	"getImageFileIndex": {
		Map: `
			function(doc){
				if(doc.type==="file"){
					const att=doc._attachments;
					const contentType=att[Object.keys(att)[0]].content_type;
					if(contentType.substring(0,6)==="image/"){
						emit(doc.name,doc);
					}
				}
			}`,
		Reduce: "_count",
	},
	"getPageBySlug": {
		Map: `
			function(doc){
				if(doc.type==="page"){
					emit(doc.slug, {pageRev: doc._rev, page: doc});
				}
			}`,
	},
	"getChildPageIndex": {
		Map: `
			function(doc) {
				if(doc.type==="page"){
					var owningPage = doc.owningPage || doc.owning_page;
				 	if(doc._id === doc.owningPage){
						emit(doc.parent, {
							id: doc._id,
							slug: doc.slug,
							title: doc.title,
							owner: doc.owner,
							editor: doc.editor,
							timestamp: doc.timestamp
						});
					}
				}
			}`,
		Reduce: "_count",
	},
	"getDescendants": {
		Map: `
			function(doc) {
				for (var i in doc.lineage) {
					emit([doc.lineage[i], doc.lineage], {
						id: doc._id,
						slug: doc.slug,
						title: doc.title,
						owner: doc.owner,
						editor: doc.editor,
						timestamp: doc.timestamp
					});
				}
			}
		`,
	},
	"checkUniqueSlug": {
		Map: `
			function(doc){
				if(doc.type==="page"){
					emit(doc.slug, 1);
				}
			}
		`,
		Reduce: "_count",
	},
}

var commentViews = map[string]View{
	"getCommentsForPage": {
		Map: `
			function(doc){
				if(doc.type==="comment"){
					var owningPage = doc.owningPage || doc.owning_page;
					emit([owningPage, doc.createdTime], doc);
				}
			}`,
		Reduce: "_count",
	},
}

//Populate a database with views, etc.
func InitDb(db *Database, wikiName string) error {
	ddoc := DesignDocument{
		Language: "javascript",
		Views:    wikiViews,
	}
	desRev, err := db.SaveDesignDoc("wikit", ddoc, "")
	if err != nil {
		return err
	}
	comment_ddoc := DesignDocument{
		Language: "javascript",
		Views:    commentViews,
	}
	_, err = db.SaveDesignDoc("wikit_comments", comment_ddoc, "")
	if err != nil {
		db.Delete("_design/wikit", desRev)
		return err
	}
	//setup up roles
	sec, err := db.GetSecurity()
	if err != nil {
		return err
	}
	adminRole := wikiName + ":admin"
	writeRole := wikiName + ":write"
	readRole := wikiName + ":read"
	//Wiki Admin and also 'site admin' and 'master' accounts shall have
	//admin privileges
	sec.Admins.Roles = []string{adminRole, "admin", "master"}
	sec.Members.Roles = []string{readRole, writeRole}
	err = db.SaveSecurity(*sec)
	if err != nil {
		return err
	}
	//save the validation document
	validator := createValidator(wikiName, writeRole, adminRole)
	adoc := AuthDesignDocument{
		Language:          "javascript",
		ValidateDocUpdate: validator,
	}
	_, err = db.SaveDesignDoc("_auth", adoc, "")
	if err != nil {
		return err
	}
	return nil
}

func createValidator(wikiName string, writeRole string,
	adminRole string) string {

	validationFunc := "function(newDoc, oldDoc, userCtx){" +
		"if((userCtx.roles.indexOf('" + writeRole + "') == -1) &&" +
		"(userCtx.roles.indexOf('" + adminRole + "') == -1) &&" +
		"(userCtx.roles.indexOf('admin') == -1) &&" +
		"(userCtx.roles.indexOf('master') == -1) &&" +
		"(userCtx.roles.indexOf('_admin') == -1)){" +
		"throw({forbidden: \"Not Authorized\"});" +
		"}" +
		"}"

	return validationFunc
}

func quotifyString(str string) string {
	return "\"" + str + "\""
}

func SetKey(keyval string) *url.Values {
	return &url.Values{"key": []string{"\"" + keyval + "\""}}
}

func SetKeys(startKeys []string, endKeys []string) *url.Values {
	applyQuotes := func(strList []string) {
		for i, v := range strList {
			if v != "{}" {
				strList[i] = quotifyString(v)
			}
		}
	}
	applyQuotes(startKeys)
	applyQuotes(endKeys)
	return &url.Values{"descending": []string{"true"},
		"startkey": []string{
			"[" + strings.Join(startKeys, ",") + "]"},
		"endkey": []string{
			"[" + strings.Join(endKeys, ",") + "]"},
	}
}

func StructToMap(data interface{}) (map[string]interface{}, error) {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	if dataVal.Kind() != reflect.Struct {
		return nil, fmt.Errorf("Not a struct!")
	}
	typ := dataVal.Type()
	result := make(map[string]interface{})
	for i := 0; i < dataVal.NumField(); i++ {
		field := dataVal.Field(i)
		fieldName := typ.Field(i).Name
		fmt.Printf("\nfieldName : %v, field: %v", fieldName, field)
		kind := field.Kind()
		switch {
		case kind >= reflect.Int && kind <= reflect.Int64:
			result[fieldName] = field.Int()
		case kind >= reflect.Uint && kind <= reflect.Uint64:
			result[fieldName] = field.Uint()
		case kind >= reflect.Float32 && kind <= reflect.Float64:
			result[fieldName] = field.Float()
		case kind == reflect.String:
			result[fieldName] = field.String()
		case kind >= reflect.Complex64 && kind <= reflect.Complex128:
			result[fieldName] = field.Complex()
		default:
			if field.CanInterface() {
				result[fieldName] = field.Interface()
			}
		}
	}
	return result, nil
}
