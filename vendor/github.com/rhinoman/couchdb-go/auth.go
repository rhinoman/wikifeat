package couchdb

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

//Basic interface for Auth
type Auth interface {
	//Adds authentication headers to a request
	AddAuthHeaders(*http.Request)
	//Extracts Updated auth info from Couch Response
	UpdateAuth(*http.Response)
	//Sets updated auth (headers, cookies, etc.) in an http response
	//For the update function, the map keys are cookie and/or header names
	GetUpdatedAuth() map[string]string
	//Purely for debug purposes.  Do not call, ever.
	DebugString() string
}

//HTTP Basic Authentication support
type BasicAuth struct {
	Username string
	Password string
}

//Pass-through Auth header
type PassThroughAuth struct {
	AuthHeader string
}

//Cookie-based auth (for sessions)
type CookieAuth struct {
	AuthToken        string
	UpdatedAuthToken string
}

//Proxy authentication
type ProxyAuth struct {
	Username  string
	Roles     []string
	AuthToken string
}

//Adds Basic Authentication headers to an http request
func (ba *BasicAuth) AddAuthHeaders(req *http.Request) {
	authString := []byte(ba.Username + ":" + ba.Password)
	header := "Basic " + base64.StdEncoding.EncodeToString(authString)
	req.Header.Set("Authorization", string(header))
}

//Use if you already have an Authentication header you want to pass through to couchdb
func (pta *PassThroughAuth) AddAuthHeaders(req *http.Request) {
	req.Header.Set("Authorization", pta.AuthHeader)
}

//Adds session token to request
func (ca *CookieAuth) AddAuthHeaders(req *http.Request) {
	authString := "AuthSession=" + ca.AuthToken
	req.Header.Set("Cookie", authString)
	req.Header.Set("X-CouchDB-WWW-Authenticate", "Cookie")
}

func (pa *ProxyAuth) AddAuthHeaders(req *http.Request) {
	req.Header.Set("X-Auth-CouchDB-Username", pa.Username)
	rolesString := strings.Join(pa.Roles, ",")
	req.Header.Set("X-Auth-CouchDB-Roles", rolesString)
	if pa.AuthToken != "" {
		req.Header.Set("X-Auth-CouchDB-Token", pa.AuthToken)
	}
}

//Update Auth Data
//If couchdb generates a new token, place it in a separate field so that
//it is available to an application

//do nothing for basic auth
func (ba *BasicAuth) UpdateAuth(resp *http.Response) {}

//Couchdb returns updated AuthSession tokens
func (ca *CookieAuth) UpdateAuth(resp *http.Response) {
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "AuthSession" {
			ca.UpdatedAuthToken = cookie.Value
		}
	}
}

//do nothing for pass through
func (pta *PassThroughAuth) UpdateAuth(resp *http.Response) {}

//do nothing for proxy auth
func (pa *ProxyAuth) UpdateAuth(resp *http.Response) {}

//Get Updated Auth
//Does nothing for BasicAuth
func (ba *BasicAuth) GetUpdatedAuth() map[string]string {
	return nil
}

//Does nothing for PassThroughAuth
func (pta *PassThroughAuth) GetUpdatedAuth() map[string]string {
	return nil
}

//Set AuthSession Cookie
func (ca *CookieAuth) GetUpdatedAuth() map[string]string {
	am := make(map[string]string)
	if ca.UpdatedAuthToken != "" {
		am["AuthSession"] = ca.UpdatedAuthToken
	}
	return am
}

//do nothing for Proxy Auth
func (pa *ProxyAuth) GetUpdatedAuth() map[string]string {
	return nil
}

//Return a Debug string

func (ba *BasicAuth) DebugString() string {
	return fmt.Sprintf("Username: %v, Password: %v", ba.Username, ba.Password)
}

func (pta *PassThroughAuth) DebugString() string {
	return fmt.Sprintf("Authorization Header: %v", pta.AuthHeader)
}

func (ca *CookieAuth) DebugString() string {
	return fmt.Sprintf("AuthToken: %v, Updated AuthToken: %v",
		ca.AuthToken, ca.UpdatedAuthToken)
}

func (pa *ProxyAuth) DebugString() string {
	return fmt.Sprintf("Username: %v, Roles: %v, AuthToken: %v",
		pa.Username, pa.Roles, pa.AuthToken)
}

//TODO: Add support for other Authentication methods supported by Couch:
//OAuth, etc.
