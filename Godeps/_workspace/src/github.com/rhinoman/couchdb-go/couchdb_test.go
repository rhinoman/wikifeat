package couchdb_test

import (
	"bytes"
	//"encoding/json"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/couchdb-go"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/twinj/uuid"
	"io/ioutil"
	"math/rand"
	//"net/http"
	"strconv"
	"testing"
	"time"
)

var timeout = time.Duration(500 * time.Millisecond)
var unittestdb = "unittestdb"
var server = "127.0.0.1"
var numDbs = 1
var adminAuth = &couchdb.BasicAuth{Username: "adminuser", Password: "password"}

type TestDocument struct {
	Title string
	Note  string
}

type ViewResult struct {
	Id  string       `json:"id"`
	Key TestDocument `json:"key"`
}

type ViewResponse struct {
	TotalRows int          `json:"total_rows"`
	Offset    int          `json:"offset"`
	Rows      []ViewResult `json:"rows,omitempty"`
}

type MultiReadResponse struct {
	TotalRows int            `json:"total_rows"`
	Offset    int            `json:"offset"`
	Rows      []MultiReadRow `json:"rows"`
}

type MultiReadRow struct {
	Id  string       `json:"id"`
	Key string       `json:"key"`
	Doc TestDocument `json:"doc"`
}

type ListResult struct {
	Id  string       `json:"id"`
	Key TestDocument `json:"key"`
	//Value string       `json:"value"`
}

type ListResponse struct {
	TotalRows int          `json:"total_rows"`
	Offset    int          `json:"offset"`
	Rows      []ListResult `json:"rows,omitempty"`
}

type View struct {
	Map    string `json:"map"`
	Reduce string `json:"reduce,omitempty"`
}

type DesignDocument struct {
	Language string            `json:"language"`
	Views    map[string]View   `json:"views"`
	Lists    map[string]string `json:"lists"`
}

func getUuid() string {
	theUuid := uuid.NewV4()
	return uuid.Formatter(theUuid, uuid.Clean)
}

func getConnection(t *testing.T) *couchdb.Connection {
	conn, err := couchdb.NewConnection(server, 5984, timeout)
	if err != nil {
		t.Logf("ERROR: %v", err)
		t.Fail()
	}
	return conn
}

/*func getAuthConnection(t *testing.T) *couchdb.Connection {
	auth := couchdb.Auth{Username: "adminuser", Password: "password"}
	conn, err := couchdb.NewConnection(server, 5984, timeout)
	if err != nil {
		t.Logf("ERROR: %v", err)
		t.Fail()
	}
	return conn
}*/

func createTestDb(t *testing.T) string {
	conn := getConnection(t)
	dbName := unittestdb + strconv.Itoa(numDbs)
	err := conn.CreateDB(dbName, adminAuth)
	errorify(t, err)
	numDbs += 1
	return dbName
}

func deleteTestDb(t *testing.T, dbName string) {
	conn := getConnection(t)
	err := conn.DeleteDB(dbName, adminAuth)
	errorify(t, err)
}

func genRandomText(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func createLotsDocs(t *testing.T, db *couchdb.Database) {
	for i := 0; i < 10; i++ {
		id := getUuid()
		note := "purple"
		if i%2 == 0 {
			note = "magenta"
		}
		testDoc := TestDocument{
			Title: "TheDoc -- " + strconv.Itoa(i),
			Note:  note,
		}
		_, err := db.Save(testDoc, id, "")
		errorify(t, err)
	}
}

func errorify(t *testing.T, err error) {
	if err != nil {
		t.Logf("ERROR: %v", err)
		t.Fail()
	}
}

func TestPing(t *testing.T) {
	conn := getConnection(t)
	pingErr := conn.Ping()
	errorify(t, pingErr)
}

func TestBadPing(t *testing.T) {
	conn, err := couchdb.NewConnection("unpingable", 1234, timeout)
	errorify(t, err)
	pingErr := conn.Ping()
	if pingErr == nil {
		t.Fail()
	}
}

func TestGetDBList(t *testing.T) {
	conn := getConnection(t)
	dbList, err := conn.GetDBList()
	errorify(t, err)
	if len(dbList) <= 0 {
		t.Logf("No results!")
		t.Fail()
	} else {
		for i, dbName := range dbList {
			t.Logf("Database %v: %v\n", i, dbName)
		}
	}
}

func TestCreateDB(t *testing.T) {
	conn := getConnection(t)
	err := conn.CreateDB("testcreatedb", adminAuth)
	errorify(t, err)
	//try to create it again --- should fail
	err = conn.CreateDB("testcreatedb", adminAuth)
	if err == nil {
		t.Fail()
	}
	//now delete it
	err = conn.DeleteDB("testcreatedb", adminAuth)
	errorify(t, err)
}

func TestSave(t *testing.T) {
	dbName := createTestDb(t)
	conn := getConnection(t)
	//Create a new document
	theDoc := TestDocument{
		Title: "My Document",
		Note:  "This is my note",
	}
	db := conn.SelectDB(dbName, nil)
	theId := getUuid()
	//Save it
	t.Logf("Saving first\n")
	rev, err := db.Save(theDoc, theId, "")
	errorify(t, err)
	t.Logf("New Document ID: %s\n", theId)
	t.Logf("New Document Rev: %s\n", rev)
	t.Logf("New Document Title: %v\n", theDoc.Title)
	t.Logf("New Document Note: %v\n", theDoc.Note)
	if theDoc.Title != "My Document" ||
		theDoc.Note != "This is my note" || rev == "" {
		t.Fail()
	}
	//Now, let's try updating it
	theDoc.Note = "A new note"
	t.Logf("Saving again\n")
	rev, err = db.Save(theDoc, theId, rev)
	errorify(t, err)
	t.Logf("Updated Document Id: %s\n", theId)
	t.Logf("Updated Document Rev: %s\n", rev)
	t.Logf("Updated Document Title: %v\n", theDoc.Title)
	t.Logf("Updated Document Note: %v\n", theDoc.Note)
	if theDoc.Note != "A new note" {
		t.Fail()
	}
	deleteTestDb(t, dbName)
}

func TestAttachment(t *testing.T) {
	dbName := createTestDb(t)
	conn := getConnection(t)
	//Create a new document
	theDoc := TestDocument{
		Title: "My Document",
		Note:  "This one has attachments",
	}
	db := conn.SelectDB(dbName, nil)
	theId := getUuid()
	//Save it
	t.Logf("Saving document\n")
	rev, err := db.Save(theDoc, theId, "")
	errorify(t, err)
	t.Logf("New Document Id: %s\n", theId)
	t.Logf("New Document Rev: %s\n", rev)
	t.Logf("New Document Title: %v\n", theDoc.Title)
	t.Logf("New Document Note: %v\n", theDoc.Note)
	//Create some content
	content := []byte("THIS IS MY ATTACHMENT")
	contentReader := bytes.NewReader(content)
	//Now Add an attachment
	uRev, err := db.SaveAttachment(theId, rev, "attachment", "text/plain", contentReader)
	errorify(t, err)
	t.Logf("Updated Rev: %s\n", uRev)
	//Now try to read it
	theContent, err := db.GetAttachment(theId, uRev, "text/plain", "attachment")
	errorify(t, err)
	defer theContent.Close()
	theBytes, err := ioutil.ReadAll(theContent)
	errorify(t, err)
	t.Logf("how much data: %v\n", len(theBytes))
	data := string(theBytes[:])
	if data != "THIS IS MY ATTACHMENT" {
		t.Fail()
	}
	t.Logf("The data: %v\n", data)
	//Now delete it
	dRev, err := db.DeleteAttachment(theId, uRev, "attachment")
	errorify(t, err)
	t.Logf("Deleted revision: %v\n", dRev)
	deleteTestDb(t, dbName)
}

func TestRead(t *testing.T) {
	dbName := createTestDb(t)
	conn := getConnection(t)
	db := conn.SelectDB(dbName, nil)
	//Create a test doc
	theDoc := TestDocument{
		Title: "My Document",
		Note:  "Time to read",
	}
	emptyDoc := TestDocument{}
	//Save it
	theId := getUuid()
	_, err := db.Save(theDoc, theId, "")
	errorify(t, err)
	//Now try to read it
	rev, err := db.Read(theId, &emptyDoc, nil)
	errorify(t, err)
	t.Logf("Document Id: %v\n", theId)
	t.Logf("Document Rev: %v\n", rev)
	t.Logf("Document Title: %v\n", emptyDoc.Title)
	t.Logf("Document Note: %v\n", emptyDoc.Note)
	deleteTestDb(t, dbName)
}

func TestMultiRead(t *testing.T) {
	dbName := createTestDb(t)
	conn := getConnection(t)
	db := conn.SelectDB(dbName, nil)
	//Create a test doc
	theDoc := TestDocument{
		Title: "My Document",
		Note:  "Time to read",
	}
	//Save it
	theId := getUuid()
	_, err := db.Save(theDoc, theId, "")
	errorify(t, err)
	//Create another test doc
	theOtherDoc := TestDocument{
		Title: "My Other Document",
		Note:  "TIme to unread",
	}
	//Save it
	otherId := getUuid()
	_, err = db.Save(theOtherDoc, otherId, "")
	errorify(t, err)
	//Now, try to read them
	readDocs := MultiReadResponse{}
	keys := []string{theId, otherId}
	err = db.ReadMultiple(keys, &readDocs)
	errorify(t, err)
	t.Logf("\nThe Docs! %v", readDocs)
	if len(readDocs.Rows) != 2 {
		t.Errorf("Should be 2 results!")
	}
	deleteTestDb(t, dbName)
}

func TestCopy(t *testing.T) {
	dbName := createTestDb(t)
	conn := getConnection(t)
	db := conn.SelectDB(dbName, nil)
	//Create a test doc
	theDoc := TestDocument{
		Title: "My Document",
		Note:  "Time to read",
	}
	emptyDoc := TestDocument{}
	//Save it
	theId := getUuid()
	rev, err := db.Save(theDoc, theId, "")
	errorify(t, err)
	//Now copy it
	copyId := getUuid()
	copyRev, err := db.Copy(theId, "", copyId)
	errorify(t, err)
	t.Logf("Document Id: %v\n", theId)
	t.Logf("Document Rev: %v\n", rev)
	//Now read the copy
	_, err = db.Read(copyId, &emptyDoc, nil)
	errorify(t, err)
	t.Logf("Document Title: %v\n", emptyDoc.Title)
	t.Logf("Document Note: %v\n", emptyDoc.Note)
	t.Logf("Copied Doc Rev: %v\n", copyRev)
	deleteTestDb(t, dbName)
}

func TestDelete(t *testing.T) {
	dbName := createTestDb(t)
	conn := getConnection(t)
	db := conn.SelectDB(dbName, nil)
	//Create a test doc
	theDoc := TestDocument{
		Title: "My Document",
		Note:  "Time to read",
	}
	theId := getUuid()
	rev, err := db.Save(theDoc, theId, "")
	errorify(t, err)
	//Now delete it
	newRev, err := db.Delete(theId, rev)
	errorify(t, err)
	t.Logf("Document Id: %v\n", theId)
	t.Logf("Document Rev: %v\n", rev)
	t.Logf("Deleted Rev: %v\n", newRev)
	if newRev == "" || newRev == rev {
		t.Fail()
	}
	deleteTestDb(t, dbName)
}

func TestUser(t *testing.T) {
	dbName := createTestDb(t)
	conn := getConnection(t)
	//Save a User
	t.Logf("AdminAuth: %v\n", adminAuth)
	rev, err := conn.AddUser("turd.ferguson",
		"password", []string{"loser"}, adminAuth)
	errorify(t, err)
	t.Logf("User Rev: %v\n", rev)
	if rev == "" {
		t.Fail()
	}
	//check user can access db
	db := conn.SelectDB(dbName, &couchdb.BasicAuth{"turd.ferguson", "password"})
	theId := getUuid()
	docRev, err := db.Save(&TestDocument{Title: "My doc"}, theId, "")
	errorify(t, err)
	t.Logf("Granting role to user")
	//check session info
	authInfo, err := conn.GetAuthInfo(&couchdb.BasicAuth{"turd.ferguson", "password"})
	errorify(t, err)
	t.Logf("AuthInfo: %v", authInfo)
	if authInfo.UserCtx.Name != "turd.ferguson" {
		t.Errorf("UserCtx name wrong: %v", authInfo.UserCtx.Name)
	}
	//grant a role
	rev, err = conn.GrantRole("turd.ferguson",
		"fool", adminAuth)
	errorify(t, err)
	t.Logf("Updated Rev: %v\n", rev)
	//read the user
	userData := couchdb.UserRecord{}
	rev, err = conn.GetUser("turd.ferguson", &userData, adminAuth)
	errorify(t, err)
	if len(userData.Roles) != 2 {
		t.Error("Not enough roles")
	}
	t.Logf("Roles: %v", userData.Roles)
	//check user can access db
	docRev, err = db.Save(&TestDocument{Title: "My doc"}, getUuid(), docRev)
	errorify(t, err)

	//revoke a role
	rev, err = conn.RevokeRole("turd.ferguson",
		"loser", adminAuth)
	errorify(t, err)
	t.Logf("Updated Rev: %v\n", rev)
	//read the user
	rev, err = conn.GetUser("turd.ferguson", &userData, adminAuth)
	errorify(t, err)
	if len(userData.Roles) != 1 {
		t.Error("should only be 1 role")
	}
	t.Logf("Roles: %v", userData.Roles)
	dRev, err := conn.DeleteUser("turd.ferguson", rev, adminAuth)
	errorify(t, err)
	t.Logf("Del User Rev: %v\n", dRev)
	if rev == dRev || dRev == "" {
		t.Fail()
	}
	deleteTestDb(t, dbName)
}

func TestSecurity(t *testing.T) {
	conn := getConnection(t)
	dbName := createTestDb(t)
	db := conn.SelectDB(dbName, adminAuth)

	members := couchdb.Members{
		Users: []string{"joe, bill"},
		Roles: []string{"code monkeys"},
	}
	admins := couchdb.Members{
		Users: []string{"bossman"},
		Roles: []string{"boss"},
	}
	security := couchdb.Security{
		Members: members,
		Admins:  admins,
	}
	err := db.SaveSecurity(security)
	errorify(t, err)
	err = db.AddRole("sales", false)
	errorify(t, err)
	err = db.AddRole("uberboss", true)
	errorify(t, err)
	sec, err := db.GetSecurity()
	t.Logf("Security: %v\n", sec)
	if sec.Admins.Users[0] != "bossman" {
		t.Fail()
	}
	if sec.Admins.Roles[0] != "boss" {
		t.Fail()
	}
	if sec.Admins.Roles[1] != "uberboss" {
		t.Errorf("\nAdmin Roles nto right! %v\n", sec.Admins.Roles[1])
	}
	if sec.Members.Roles[1] != "sales" {
		t.Errorf("\nRoles not right! %v\n", sec.Members.Roles[1])
	}
	errorify(t, err)
	err = db.RemoveRole("sales")
	errorify(t, err)
	err = db.RemoveRole("uberboss")
	errorify(t, err)
	//try removing a role that ain't there
	err = db.RemoveRole("WHATROLE")
	errorify(t, err)
	sec, err = db.GetSecurity()
	t.Logf("Secuirty: %v\n", sec)
	if len(sec.Members.Roles) > 1 {
		t.Errorf("\nThe Role was not removed: %v\n", sec.Members.Roles)
	} else if sec.Members.Roles[0] == "sales" {
		t.Errorf("\nThe roles are all messed up: %v\n", sec.Members.Roles)
	}
	if len(sec.Admins.Roles) > 1 {
		t.Errorf("\nThe Admin Role was not removed: %v\n", sec.Admins.Roles)
	}
	deleteTestDb(t, dbName)
}

func TestSessions(t *testing.T) {
	conn := getConnection(t)
	dbName := createTestDb(t)
	defer deleteTestDb(t, dbName)
	//Save a User
	t.Logf("AdminAuth: %v\n", adminAuth)
	rev, err := conn.AddUser("turd.ferguson",
		"password", []string{"loser"}, adminAuth)
	errorify(t, err)
	t.Logf("User Rev: %v\n", rev)
	defer conn.DeleteUser("turd.ferguson", rev, adminAuth)
	if rev == "" {
		t.Fail()
	}
	//Create a session for the user
	cookieAuth, err := conn.CreateSession("turd.ferguson", "password")
	errorify(t, err)
	//sleep
	time.Sleep(time.Duration(2 * time.Second))
	//Create something
	db := conn.SelectDB(dbName, cookieAuth)
	theId := getUuid()
	docRev, err := db.Save(&TestDocument{Title: "The test doc"}, theId, "")
	errorify(t, err)
	t.Logf("Document Rev: %v", docRev)
	t.Logf("Updated Auth: %v", cookieAuth.GetUpdatedAuth()["AuthSession"])
	//Delete the user session
	err = conn.DestroySession(cookieAuth)
	errorify(t, err)
}

func TestSetConfig(t *testing.T) {
	conn := getConnection(t)
	err := conn.SetConfig("couch_httpd_auth", "timeout", "30", adminAuth)
	errorify(t, err)
}

func TestDesignDocs(t *testing.T) {
	conn := getConnection(t)
	dbName := createTestDb(t)
	db := conn.SelectDB(dbName, adminAuth)
	createLotsDocs(t, db)

	view := View{
		Map: "function(doc) {\n  if (doc.Note === \"magenta\"){\n    emit(doc)\n  }\n}",
	}
	views := make(map[string]View)
	views["find_all_magenta"] = view
	lists := make(map[string]string)

	lists["getList"] =
		`function(head, req){
			var row;
			var response={
				total_rows:0,
				offset:0, 
				rows:[]
			};
			while(row=getRow()){
				response.rows.push(row);
			}
			response.total_rows = response.rows.length;
			send(toJSON(response))
		}`

	ddoc := DesignDocument{
		Language: "javascript",
		Views:    views,
		Lists:    lists,
	}
	rev, err := db.SaveDesignDoc("colors", ddoc, "")
	errorify(t, err)
	if rev == "" {
		t.Fail()
	} else {
		t.Logf("Rev of design doc: %v\n", rev)
	}
	result := ViewResponse{}
	//now try to query the view
	err = db.GetView("colors", "find_all_magenta", &result, nil)
	errorify(t, err)
	if len(result.Rows) != 5 {
		t.Logf("docList length: %v\n", len(result.Rows))
		t.Fail()
	} else {
		t.Logf("Results: %v\n", result.Rows)
	}
	listResult := ListResponse{}
	err = db.GetList("colors", "getList", "find_all_magenta", &listResult, nil)
	if err != nil {
		t.Logf("ERROR: %v", err)
	}
	errorify(t, err)
	if len(listResult.Rows) != 5 {
		t.Logf("List Result: %v\n", listResult)
		t.Logf("docList length: %v\n", len(listResult.Rows))
		t.Fail()
	} else {
		t.Logf("List Results: %v\n", listResult)
	}
	deleteTestDb(t, dbName)

}

//Test for a specific situation I've been having trouble with
func TestAngryCouch(t *testing.T) {

	testDoc1 := TestDocument{
		Title: "Test Doc 1",
		Note:  genRandomText(8000),
	}
	testDoc2 := TestDocument{
		Title: "Test Doc 2",
		Note:  genRandomText(1000),
	}

	dbName := createTestDb(t)
	defer deleteTestDb(t, dbName)
	conn := getConnection(t)
	db := conn.SelectDB(dbName, nil)
	//client := &http.Client{}
	id1 := getUuid()
	id2 := getUuid()
	/*req, err := http.NewRequest(
		"PUT",
		"http://localhost:5984/"+dbName+"/"+id1,
		bytes.NewReader(testBody1),
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Expect", "100-continue")*/
	rev, err := db.Save(testDoc1, id1, "")
	//resp, err := client.Do(req)
	errorify(t, err)
	//defer resp.Body.Close()
	t.Logf("Doc 1 Rev: %v\n", rev)
	errorify(t, err)
	rev, err = db.Save(testDoc2, id2, "")
	t.Logf("Doc 2 Rev: %v\n", rev)
	errorify(t, err)
}
