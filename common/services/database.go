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

package services

import (
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/couchdb-go"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/twinj/uuid"
	"github.com/rhinoman/wikifeat/common/config"
	"github.com/rhinoman/wikifeat/common/entities"
	"log"
	"strconv"
	"strings"
	"time"
)

//Holds the connection to the database
var Connection *couchdb.Connection

//The couchdb Admin credentials
var AdminAuth couchdb.Auth

//The name of the main database
var MainDb string

//The name of the user avatar database
var AvatarDb string

//The name of the users database
var UserPrefix = "org.couchdb.user:"
var UserDbName = "_users"

//The various manager instance pointers
//var UsrMgr = new(UserManager)
//var WikiMgr = new(WikiManager)
//var PgMgr = new(PageManager)

type DesignView struct {
	Map    string `json:"map"`
	Reduce string `json:"reduce,omitempty"`
}

type DesignDocument struct {
	Language string                `json:"language"`
	Views    map[string]DesignView `json:"views"`
	Lists    map[string]string     `json:"lists"`
	Shows    map[string]string     `json:"shows"`
}

type AuthDesignDocument struct {
	Language          string `json:"language"`
	ValidateDocUpdate string `json:"validate_doc_update"`
}

type ViewResponse struct {
	TotalRows int `json:"total_rows"`
	Offset    int `json:"offset"`
}

//Only one master account, please.
func MasterRole() string {
	return "master"
}

func AdminRole(dbName string) string {
	if dbName == MainDbName() {
		return "admin"
	}
	return dbName + ":admin"
}

func WriteRole(dbName string) string {
	return dbName + ":write"
}

func ReadRole(dbName string) string {
	return dbName + ":read"
}

func MainDbName() string {
	return MainDb
}

func AvatarDbName() string {
	return AvatarDb
}

func AllUsersRole() string {
	return "all_users"
}

func BadRequestError() error {
	return &couchdb.Error{
		StatusCode: 400,
		Reason:     "Bad Request",
	}
}

func NotAdminError() error {
	return &couchdb.Error{
		StatusCode: 403,
		Reason:     "Not an admin",
	}
}

func NotFoundError() error {
	return &couchdb.Error{
		StatusCode: 404,
		Reason:     "Not Found",
	}
}

//Initialize the Database Connection
func InitDb() {
	log.Println("Initializing Database Connection")
	timeoutMs, err := strconv.Atoi(config.Database.DbTimeout)
	if err != nil {
		log.Fatalf("Error! %v", err)
	}
	port, err := strconv.Atoi(config.Database.DbPort)
	if err != nil {
		log.Fatalf("Error! %v", err)
	}
	connfun := connectionFunc(config.Database.UseSSL)
	Connection, err = connfun(config.Database.DbAddr,
		port, time.Duration(timeoutMs)*time.Millisecond)
	if err != nil {
		log.Fatalf("Error! %v", err)
	}
	AdminAuth = &couchdb.BasicAuth{
		Username: config.Database.DbAdminUser,
		Password: config.Database.DbAdminPassword,
	}
	MainDb = config.Database.MainDb
	AvatarDb = config.Users.AvatarDb
	//Set DB Configuration options
	err = Connection.SetConfig("couch_httpd_auth",
		"allow_persistent_cookies",
		strconv.FormatBool(config.Auth.PersistentSessions),
		AdminAuth)
	if err != nil {
		log.Fatalf("Error! %v", err)
	}
	err = Connection.SetConfig("couch_httpd_auth",
		"timeout",
		strconv.Itoa(config.Auth.SessionTimeout),
		AdminAuth)
	if err != nil {
		log.Fatalf("Error! %v", err)
	}
	//Enable or Revoke Guest access as appropriate
	if config.Auth.AllowGuest {
		if err = EnableGuest(); err != nil {
			log.Fatalf("Error! %v", err)
		}
	} else {
		if err = DisableGuest(); err != nil {
			log.Fatalf("Error! %v", err)
		}
	}
}

// Enable Guest access by adding read access on the main database
// to the guest user
func EnableGuest() error {
	_, err := Connection.GrantRole("guest", ReadRole(MainDbName()), AdminAuth)
	if err != nil {
		code := checkErrorCode(err)
		if code == 404 {
			return CreateGuestUser()
		} else {
			return err
		}
	}
	return nil
}

func checkErrorCode(err error) int {
	splitStr := strings.Split(err.Error(), ":")
	if len(splitStr) > 1 {
		code, conerr := strconv.Atoi(splitStr[1])
		if conerr != nil {
			return 0
		} else {
			return code
		}
	} else {
		return 0
	}
}

// Revoke Guest access
func DisableGuest() error {
	_, err := Connection.RevokeRole("guest", ReadRole(MainDbName()), AdminAuth)
	if err != nil {
		code := checkErrorCode(err)
		if code == 404 {
			return nil
		} else {
			return err
		}
	}
	return nil
}

func CreateGuestUser() error {
	_, err := Connection.AddUser("guest", "guest",
		[]string{ReadRole(MainDbName())}, AdminAuth)
	return err
}

//Sets initial db security
func setDbSecurity(dbName string) error {
	//set initial database security
	db := Connection.SelectDB(dbName, AdminAuth)
	sec, err := db.GetSecurity()
	if err != nil {
		return err
	}
	//Create our roles and store them to the db security doc
	sec.Admins.Roles = []string{AdminRole(dbName)}
	sec.Members.Roles = []string{ReadRole(dbName), WriteRole(dbName)}
	log.Println("Setting security doc for: " + dbName)
	err = db.SaveSecurity(*sec)
	if err != nil {
		return err
	}
	//Set the write validation function in couchdb
	validationFunc := "function(newDoc, oldDoc, userCtx){" +
		"if((userCtx.roles.indexOf('" + WriteRole(dbName) + "') === -1) &&" +
		"(userCtx.roles.indexOf('" + AdminRole(dbName) + "') === -1) &&" +
		"(userCtx.roles.indexOf('_admin') === -1)){" +
		"throw({forbidden: \"Not Authorized\"});" +
		"}" +
		"}"

	adoc := AuthDesignDocument{
		Language:          "javascript",
		ValidateDocUpdate: validationFunc,
	}
	log.Println("Saving validation function for: " + dbName)
	_, err = db.SaveDesignDoc("_auth", adoc, "")
	if err != nil {
		return err
	}
	return nil
}

func InitMainDatabase() error {
	mainDb := Connection.SelectDB(MainDbName(), AdminAuth)
	getWikis := `
		function(doc) {
			if(doc.type==="wiki_record"){
				emit(doc.name, doc);
			}
		}
	`
	getWikiBySlug := `
		function(doc) {
			if(doc.type==="wiki_record"){
				emit(doc.slug, {wikiRev: doc._rev, wiki_record: doc});
			}
		}
	`
	userWikiList := `
		function(head, req){
			var row;
			var user=req['userCtx']['name'];
			var response={
				total_rows:0,
				offset:0, rows:[]
			};
			while(row=getRow()){
				if(user in row.value.members){
					response.rows.push(row);
				}
			}
			response.total_rows = response.rows.length;
			send(toJSON(response))
		}
	`
	checkUniqueSlug := `
		function(doc){
			if(doc.type==="wiki_record"){
				emit(doc.slug, 1);
			}
		}
	`
	gw := DesignView{Map: getWikis}
	gwbs := DesignView{Map: getWikiBySlug}
	cus := DesignView{Map: checkUniqueSlug, Reduce: "_count"}
	ddoc := DesignDocument{
		Language: "javascript",
		Views: map[string]DesignView{"getWikis": gw, "getWikiBySlug": gwbs,
			"checkUniqueSlug": cus},
		Lists: map[string]string{"userWikiList": userWikiList},
	}
	_, err := mainDb.SaveDesignDoc("wiki_query", ddoc, "")
	return err
}

//Select connection function (SSL or not)
func connectionFunc(ssl bool) func(string, int, time.Duration) (*couchdb.Connection, error) {
	if ssl {
		return couchdb.NewSSLConnection
	} else {
		return couchdb.NewConnection
	}
}

func GenUuid() string {
	return uuid.Formatter(uuid.NewV4(), uuid.Clean)
}

//Create a database
func CreateDb(dbName string) error {
	log.Println("Creating Database " + dbName)
	err := Connection.CreateDB(dbName, AdminAuth)
	if err != nil {
		log.Println("ERROR: Couldn't Create Database "+dbName+" -", err)
		return err
	}
	err = setDbSecurity(dbName)
	if err != nil {
		DeleteDb(dbName)
		log.Println("ERROR: Couldn't set security for Db "+dbName+" -", err)
		return err
	}
	return nil
}

//Delete a database
func DeleteDb(dbName string) error {
	//TODO: Cleanup roles associated with this database
	log.Println("Deleting Database " + dbName)
	err := Connection.DeleteDB(dbName, AdminAuth)
	if err != nil {
		log.Println("ERROR: Couldn't Delete Database "+dbName+" -", err)
		return err
	}
	return nil
}

//Get the current user record from auth header
func GetUserFromAuth(auth couchdb.Auth) (*entities.User, error) {
	authInfo, err := Connection.GetAuthInfo(auth)
	if err != nil {
		return nil, err
	}
	userDoc := new(entities.User)
	userDb := Connection.SelectDB(UserDbName, auth)
	_, err = userDb.Read(UserPrefix+authInfo.UserCtx.Name, userDoc, nil)
	if err != nil {
		return nil, err
	}
	return userDoc, nil
}
