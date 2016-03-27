//Package couchdb provides a simple REST client for CouchDB
package couchdb

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Connection struct{ *connection }

type Database struct {
	dbName     string
	connection *Connection
	auth       Auth
}

//Creates a regular http connection.
//Timeout sets the timeout for the http Client
func NewConnection(address string, port int,
	timeout time.Duration) (*Connection, error) {

	url := "http://" + address + ":" + strconv.Itoa(port)
	return createConnection(url, timeout)
}

//Creates an https connection.
//Timeout sets the timeout for the http Client
func NewSSLConnection(address string, port int,
	timeout time.Duration) (*Connection, error) {

	url := "https://" + address + ":" + strconv.Itoa(port)
	return createConnection(url, timeout)
}

func createConnection(rawUrl string, timeout time.Duration) (*Connection, error) {
	//check that the url is valid
	theUrl, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}
	return &Connection{
		&connection{
			url:    theUrl.String(),
			client: &http.Client{Timeout: timeout},
		},
	}, nil

}

//Use to check if database server is alive.
func (conn *Connection) Ping() error {
	resp, err := conn.request("HEAD", "/", nil, nil, nil)
	if err == nil {
		resp.Body.Close()
	}
	return err
}

//DATABASES.
//Return a list of all databases on the server
func (conn *Connection) GetDBList() (dbList []string, err error) {
	resp, err := conn.request("GET", "/_all_dbs", nil, nil, nil)
	if err != nil {
		return dbList, err
	}
	err = parseBody(resp, &dbList)
	return dbList, err
}

//Create a new Database.
func (conn *Connection) CreateDB(name string, auth Auth) error {
	url, err := buildUrl(name)
	if err != nil {
		return err
	}
	resp, err := conn.request("PUT", url, nil, nil, auth)
	if err == nil {
		resp.Body.Close()
	}
	return err
}

//Delete a Database.
func (conn *Connection) DeleteDB(name string, auth Auth) error {
	url, err := buildUrl(name)
	if err != nil {
		return err
	}
	resp, err := conn.request("DELETE", url, nil, nil, auth)
	if err == nil {
		resp.Body.Close()
	}
	return err
}

//Set a CouchDB configuration option
func (conn *Connection) SetConfig(section string,
	option string, value string, auth Auth) error {
	url, err := buildUrl("_config", section, option)
	if err != nil {
		return err
	}
	body := strings.NewReader("\"" + value + "\"")
	resp, err := conn.request("PUT", url, body, nil, auth)
	if err == nil {
		resp.Body.Close()
	}
	return err
}

//Gets a CouchDB configuration option
func (conn *Connection) GetConfigOption(section string,
	option string, auth Auth) (string, error) {
	url, err := buildUrl("_config", section, option)
	if err != nil {
		return "", err
	}
	resp, err := conn.request("GET", url, nil, nil, auth)
	var val interface{}
	parseBody(resp, &val)
	if num, ok := val.(int); ok == true {
		return strconv.Itoa(num), nil
	}
	if str, ok := val.(string); ok == true {
		return str, nil
	}
	return "", nil
}

type UserRecord struct {
	Name     string   `json:"name"`
	Password string   `json:"password,omitempty"`
	Roles    []string `json:"roles"`
	TheType  string   `json:"type"` //apparently type is a keyword in Go :)

}

//Add a User.
//This is a convenience method for adding a simple user to CouchDB.
//If you need a User with custom fields, etc., you'll just have to use the
//ordinary document methods on the "_users" database.
func (conn *Connection) AddUser(username string, password string,
	roles []string, auth Auth) (string, error) {

	userData := UserRecord{
		Name:     username,
		Password: password,
		Roles:    roles,
		TheType:  "user"}
	userDb := conn.SelectDB("_users", auth)
	namestring := "org.couchdb.user:" + userData.Name
	return userDb.Save(&userData, namestring, "")

}

//Grants a role to a user
func (conn *Connection) GrantRole(username string, role string,
	auth Auth) (string, error) {
	userDb := conn.SelectDB("_users", auth)
	namestring := "org.couchdb.user:" + username
	var userData interface{}

	rev, err := userDb.Read(namestring, &userData, nil)
	if err != nil {
		return "", err
	}
	if reflect.ValueOf(userData).Kind() != reflect.Map {
		return "", errors.New("Type Error")
	}
	userMap := userData.(map[string]interface{})
	if reflect.ValueOf(userMap["roles"]).Kind() != reflect.Slice {
		return "", errors.New("Type Error")
	}
	userRoles := userMap["roles"].([]interface{})
	//Check if our role is already in the array, so we don't add it twice
	for _, r := range userRoles {
		if r == role {
			return rev, nil
		}
	}
	userMap["roles"] = append(userRoles, role)
	return userDb.Save(&userMap, namestring, rev)
}

//Revoke a user role
func (conn *Connection) RevokeRole(username string, role string,
	auth Auth) (string, error) {

	userDb := conn.SelectDB("_users", auth)
	namestring := "org.couchdb.user:" + username
	var userData interface{}
	rev, err := userDb.Read(namestring, &userData, nil)
	if err != nil {
		return "", err
	}
	if reflect.ValueOf(userData).Kind() != reflect.Map {
		return "", errors.New("Type Error")
	}
	userMap := userData.(map[string]interface{})
	if reflect.ValueOf(userMap["roles"]).Kind() != reflect.Slice {
		return "", errors.New("Type Error")
	}
	userRoles := userMap["roles"].([]interface{})
	found := false
	for i, r := range userRoles {
		if r == role {
			userRoles = append(userRoles[:i], userRoles[i+1:]...)
			found = true
			break
		}
	}
	userMap["roles"] = userRoles
	if found == false {
		return "", nil
	} else {
		return userDb.Save(&userMap, namestring, rev)
	}
}

type UserContext struct {
	Name  string   `json:"name"`
	Roles []string `json:"roles"`
}

type AuthInfo struct {
	Authenticated          string   `json:"authenticated"`
	AuthenticationDb       string   `json:"authentication_db"`
	AuthenticationHandlers []string `json:"authentication_handlers"`
}

type AuthInfoResponse struct {
	Info    AuthInfo    `json:"info"`
	Ok      bool        `json:"ok"`
	UserCtx UserContext `json:"userCtx"`
}

//Creates a session using the Couchdb session api.  Returns auth token on success
func (conn *Connection) CreateSession(username string,
	password string) (*CookieAuth, error) {
	sessUrl, err := buildUrl("_session")
	if err != nil {
		return &CookieAuth{}, err
	}
	var headers = make(map[string]string)
	body := "name=" + username + "&password=" + password
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	resp, err := conn.request("POST", sessUrl,
		strings.NewReader(body), headers, nil)
	if err != nil {
		return &CookieAuth{}, err
	}
	defer resp.Body.Close()
	authToken := func() string {
		for _, cookie := range resp.Cookies() {
			if cookie.Name == "AuthSession" {
				return cookie.Value
			}
		}
		return ""
	}()
	return &CookieAuth{AuthToken: authToken}, nil
}

//Destroys a session (user log out, etc.)
func (conn *Connection) DestroySession(auth *CookieAuth) error {
	sessUrl, err := buildUrl("_session")
	if err != nil {
		return err
	}
	var headers = make(map[string]string)
	headers["Accept"] = "application/json"
	resp, err := conn.request("DELETE", sessUrl, nil, headers, auth)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

//Returns auth information for a user
func (conn *Connection) GetAuthInfo(auth Auth) (*AuthInfoResponse, error) {
	authInfo := AuthInfoResponse{}
	sessUrl, err := buildUrl("_session")
	if err != nil {
		return nil, err
	}
	var headers = make(map[string]string)
	headers["Accept"] = "application/json"
	resp, err := conn.request("GET", sessUrl, nil, headers, auth)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = parseBody(resp, &authInfo)
	if err != nil {
		return nil, err
	}
	return &authInfo, nil
}

//Fetch a user record
func (conn *Connection) GetUser(username string, userData interface{},
	auth Auth) (string, error) {
	userDb := conn.SelectDB("_users", auth)
	namestring := "org.couchdb.user:" + username
	return userDb.Read(namestring, &userData, nil)
}

//Delete a user.
func (conn *Connection) DeleteUser(username string, rev string, auth Auth) (string, error) {
	userDb := conn.SelectDB("_users", auth)
	namestring := "org.couchdb.user:" + username
	return userDb.Delete(namestring, rev)
}

//Select a Database.
//TODO: Perhaps verify dbName exists in couchdb?
//Or just do the fast thing here and let subsequent queries fail if the user supplies an incorrect dbname.
func (conn *Connection) SelectDB(dbName string, auth Auth) *Database {
	return &Database{
		dbName:     dbName,
		connection: conn,
		auth:       auth,
	}
}

//Save a document to the database.
//If you're creating a new document, pass an empty string for rev.
//If updating, you must specify the current rev.
//Returns the revision number assigned to the doc by CouchDB.
func (db *Database) Save(doc interface{}, id string, rev string) (string, error) {
	url, err := buildUrl(db.dbName, id)
	if err != nil {
		return "", err
	}
	var headers = make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Accept"] = "application/json"
	if id == "" {
		return "", fmt.Errorf("No ID specified")
	}
	if rev != "" {
		headers["If-Match"] = rev
	}
	data, numBytes, err := encodeData(doc)
	if err != nil {
		return "", err
	}
	headers["Content-Length"] = strconv.Itoa(numBytes)
	//Yes, this needs to be here.
	//Yes, I know the Golang http.Client doesn't support expect/continue
	//This is here to work around a bug in CouchDB.  It shouldn't work, and yet it does.
	//See: http://stackoverflow.com/questions/30541591/large-put-requests-from-go-to-couchdb
	//Also, I filed a bug report: https://issues.apache.org/jira/browse/COUCHDB-2704
	//Go net/http needs to support the HTTP/1.1 spec, or CouchDB needs to get fixed.
	//If either of those happens in the future, I can revisit this.
	//Unless I forget, which I'm sure I will.
	if numBytes > 4000 {
		headers["Expect"] = "100-continue"
	}
	resp, err := db.connection.request("PUT", url, data, headers, db.auth)
	if err != nil {
		return "", err
	}
	resp.Body.Close()
	return getRevInfo(resp)
}

//Copies a document into a new... document.
//Returns the revision of the newly created document
func (db *Database) Copy(fromId string, fromRev string, toId string) (string, error) {
	url, err := buildUrl(db.dbName, fromId)
	if err != nil {
		return "", err
	}
	var headers = make(map[string]string)
	headers["Accept"] = "application/json"
	if fromId == "" || toId == "" {
		return "", fmt.Errorf("Invalid request.  Ids must be specified")
	}
	if fromRev != "" {
		headers["If-Match"] = fromRev
	}
	headers["Destination"] = toId
	resp, err := db.connection.request("COPY", url, nil, headers, db.auth)
	if err != nil {
		return "", err
	}
	resp.Body.Close()
	return getRevInfo(resp)
}

//Fetches a document from the database.
//Pass it a &struct to hold the contents of the fetched document (doc).
//Returns the current revision and/or error
func (db *Database) Read(id string, doc interface{}, params *url.Values) (string, error) {
	var headers = make(map[string]string)
	headers["Accept"] = "application/json"
	var url string
	var err error
	if params == nil {
		url, err = buildUrl(db.dbName, id)
	} else {
		url, err = buildParamUrl(*params, db.dbName, id)
	}
	if err != nil {
		return "", err
	}
	resp, err := db.connection.request("GET", url, nil, headers, db.auth)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if err = parseBody(resp, &doc); err != nil {
		return "", err
	}
	return getRevInfo(resp)
}

//Fetches multiple documents in a single request given a set of arbitrary _ids
func (db *Database) ReadMultiple(ids []string, results interface{}) error {
	type RequestBody struct {
		Keys []string `json:"keys"`
	}
	parameters := url.Values{}
	parameters.Set("include_docs", "true")
	url, err := buildParamUrl(parameters, db.dbName, "_all_docs")
	if err != nil {
		return err
	}
	var headers = make(map[string]string)
	reqBody := RequestBody{Keys: ids}
	requestBody, numBytes, err := encodeData(reqBody)
	if err != nil {
		return err
	}
	headers["Content-Type"] = "application/json"
	headers["Content-Length"] = strconv.Itoa(numBytes)
	if numBytes > 4000 {
		headers["Expect"] = "100-continue"
	}
	headers["Accept"] = "application/json"
	if resp, err :=
		db.connection.request("POST", url, requestBody,
			headers, db.auth); err == nil {
		defer resp.Body.Close()
		return parseBody(resp, &results)
	} else {
		return err
	}
}

//Deletes a document.
//Or rather, tells CouchDB to mark the document as deleted.
//Yes, CouchDB will return a new revision, so this function returns it.
func (db *Database) Delete(id string, rev string) (string, error) {
	url, err := buildUrl(db.dbName, id)
	if err != nil {
		return "", err
	}
	var headers = make(map[string]string)
	headers["Accept"] = "application/json"
	headers["If-Match"] = rev
	resp, err := db.connection.request("DELETE", url, nil, headers, db.auth)
	if err != nil {
		return "", err
	}
	resp.Body.Close()
	return getRevInfo(resp)
}

//Saves an attachment.
//docId and docRev refer to the parent document.
//attType is the MIME type of the attachment (ex: image/jpeg) or some such.
//attContent is a byte array containing the actual content.
func (db *Database) SaveAttachment(docId string,
	docRev string, attName string,
	attType string, attContent io.Reader) (string, error) {
	url, err := buildUrl(db.dbName, docId, attName)
	if err != nil {
		return "", err
	}
	var headers = make(map[string]string)
	headers["Accept"] = "application/json"
	headers["Content-Type"] = attType
	headers["If-Match"] = docRev
	headers["Expect"] = "100-continue"

	resp, err := db.connection.request("PUT", url, attContent, headers, db.auth)
	if err != nil {
		return "", err
	}
	resp.Body.Close()
	return getRevInfo(resp)
}

//Gets an attachment.
//Returns an io.Reader -- the onus is on the caller to close it.
//Please close it.
func (db *Database) GetAttachment(docId string, docRev string,
	attType string, attName string) (io.ReadCloser, error) {
	url, err := buildUrl(db.dbName, docId, attName)
	if err != nil {
		return nil, err
	}
	var headers = make(map[string]string)
	headers["Accept"] = attType
	if docRev != "" {
		headers["If-Match"] = docRev
	}
	resp, err := db.connection.request("GET", url, nil, headers, db.auth)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

//Fetches an attachment and proxies the result
func (db *Database) GetAttachmentByProxy(docId string, docRev string,
	attType string, attName string, r *http.Request, w http.ResponseWriter) error {
	path, err := buildUrl(db.dbName, docId, attName)
	if err != nil {
		return err
	}
	var headers = make(map[string]string)
	headers["Accept"] = attType
	if docRev != "" {
		headers["If-Match"] = docRev
	}
	for k, v := range headers {
		r.Header.Set(k, v)
	}
	return db.connection.reverseProxyRequest(w, r, path, db.auth)
}

//Deletes an attachment
func (db *Database) DeleteAttachment(docId string, docRev string,
	attName string) (string, error) {
	url, err := buildUrl(db.dbName, docId, attName)
	if err != nil {
		return "", err
	}
	var headers = make(map[string]string)
	headers["Accept"] = "application/json"
	headers["If-Match"] = docRev
	resp, err := db.connection.request("DELETE", url, nil, headers, db.auth)
	if err != nil {
		return "", err
	}
	resp.Body.Close()
	return getRevInfo(resp)
}

type Members struct {
	Users []string `json:"names,omitempty"`
	Roles []string `json:"roles,omitempty"`
}

type Security struct {
	Members Members `json:"members"`
	Admins  Members `json:"admins"`
}

//Returns the Security document from the database.
func (db *Database) GetSecurity() (*Security, error) {
	url, err := buildUrl(db.dbName, "_security")
	if err != nil {
		return nil, err
	}
	var headers = make(map[string]string)
	sec := Security{}
	headers["Accept"] = "application/json"
	resp, err := db.connection.request("GET", url, nil, headers, db.auth)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = parseBody(resp, &sec)
	if err != nil {
		return nil, err
	}
	return &sec, err
}

//Save a security document to the database.
func (db *Database) SaveSecurity(sec Security) error {
	url, err := buildUrl(db.dbName, "_security")
	if err != nil {
		return err
	}
	var headers = make(map[string]string)
	headers["Accept"] = "application/json"
	data, numBytes, err := encodeData(sec)
	if err != nil {
		return err
	}
	headers["Content-Length"] = strconv.Itoa(numBytes)
	if numBytes > 4000 {
		headers["Expect"] = "100-continue"
	}
	resp, err := db.connection.request("PUT", url, data, headers, db.auth)
	if err == nil {
		resp.Body.Close()
	}
	return err
}

// Security helper function.
// Adds a role to a database security doc.
func (db *Database) AddRole(role string, isAdmin bool) error {
	sec, err := db.GetSecurity()
	if err != nil {
		return err
	}
	roles := func() *[]string {
		if isAdmin {
			return &sec.Admins.Roles
		} else {
			return &sec.Members.Roles
		}
	}
	//Make sure the role isn't already there (couchdb will let you add it twice :/ )
	for _, r := range *roles() {
		if r == role {
			//already there, just return
			return nil
		}
	}
	rolesarr := roles()
	*rolesarr = append(*rolesarr, role)
	return db.SaveSecurity(*sec)
}

// Security helper function.
// Removes a role from a database security doc.
func (db *Database) RemoveRole(role string) error {
	sec, err := db.GetSecurity()
	if err != nil {
		return err
	}
	remove := func(isAdmin bool) bool {
		var rolesPtr *[]string
		if isAdmin {
			rolesPtr = &sec.Admins.Roles
		} else {
			rolesPtr = &sec.Members.Roles
		}
		roles := *rolesPtr
		for i, r := range roles {
			if r == role {
				*rolesPtr = append(roles[:i], roles[i+1:]...)
				return true
			}
		}
		return false
	}
	var removed bool = false
	if removed = remove(false); !removed {
		removed = remove(true)
	}
	if removed {
		return db.SaveSecurity(*sec)
	}
	return nil
}

//Get the results of a view.
func (db *Database) GetView(designDoc string, view string,
	results interface{}, params *url.Values) error {
	var err error
	var url string
	if params == nil {
		url, err = buildUrl(db.dbName, "_design", designDoc, "_view", view)
	} else {
		url, err = buildParamUrl(*params, db.dbName, "_design",
			designDoc, "_view", view)
	}
	if err != nil {
		return err
	}
	var headers = make(map[string]string)
	headers["Accept"] = "application/json"
	resp, err := db.connection.request("GET", url, nil, headers, db.auth)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = parseBody(resp, &results)
	if err != nil {
		return err
	}
	return nil
}

//Get the result of a list operation
//This assumes your list function in couchdb returns JSON
func (db *Database) GetList(designDoc string, list string,
	view string, results interface{}, params *url.Values) error {
	var err error
	var url string
	if params == nil {
		url, err = buildUrl(db.dbName, "_design", designDoc, "_list",
			list, view)
	} else {
		url, err = buildParamUrl(*params, db.dbName, "_design", designDoc,
			"_list", list, view)
	}
	if err != nil {
		return err
	}
	var headers = make(map[string]string)
	headers["Accept"] = "application/json"
	resp, err := db.connection.request("GET", url, nil, nil, db.auth)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = parseBody(resp, &results)
	if err != nil {
		return err
	}
	return nil
}

//Save a design document.
//If creating a new design doc, set rev to "".
func (db *Database) SaveDesignDoc(name string,
	designDoc interface{}, rev string) (string, error) {
	path := "_design/" + name
	newRev, err := db.Save(designDoc, path, rev)
	if err != nil {
		return "", err
	} else if newRev == "" {
		return "", fmt.Errorf("CouchDB returned an empty revision string.")
	}
	return newRev, nil

}
