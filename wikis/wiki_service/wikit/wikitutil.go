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
		Map: "function(doc) {" +
			"if(doc.type===\"page\"){" +
			"emit([doc.owning_page, doc.timestamp], " +
			"{documentId: doc._id," +
			"documentRev: doc._rev," +
			"editor: doc.editor, contentSize: doc.content.raw.length}" +
			");" +
			"} }",
		Reduce: "_count",
	},
	"getIndex": {
		Map: `
			function(doc){
				if(doc.type==="page" && doc._id === doc.owning_page){
					emit(doc.title, {
						id: doc._id,
						slug: doc.slug,
						title: doc.title, 
						owner: doc.owner, 
						editor: doc.editor, 
						timestamp: doc.timestamp 
					}); 
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
	"getPageBySlug": {
		Map: "function(doc){" +
			"if(doc.type===\"page\"){" +
			"emit(doc.slug, {pageRev: doc._rev, page: doc});" +
			"}" +
			"}",
	},
	"getChildPageIndex": {
		Map: `
			function(doc) {
				if(doc.type==="page" && doc._id === doc.owning_page){
					emit(doc.parent, {
						id: doc._id,
						slug: doc.slug,
						title: doc.title, 
						owner: doc.owner, 
						editor: doc.editor, 
						timestamp: doc.timestamp
					});
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
					emit([doc.owning_page, doc.created_time], doc);
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
