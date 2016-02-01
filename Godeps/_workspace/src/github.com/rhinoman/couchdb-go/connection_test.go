package couchdb

import (
	"net/http"
	"net/url"
	"testing"
)

type couchWelcome struct {
	Couchdb string      `json:"couchdb"`
	Uuid    string      `json:"uuid"`
	Version string      `json:"version"`
	Vendor  interface{} `json:"vendor"`
}

var serverUrl = "http://127.0.0.1:5984"
var couchReply couchWelcome

func TestUrlBuilding(t *testing.T) {
	params := url.Values{}
	params.Add("Hello", "42")
	params.Add("crazy", "\"me&bo\"")
	stringified, err := buildParamUrl(params, "theDb", "funny?chars")
	if err != nil {
		t.Fail()
	}
	t.Logf("The URL: %s\n", stringified)
	//make sure everything is escaped
	if stringified != "/theDb/funny%3Fchars?Hello=42&crazy=%22me%26bo%22" {
		t.Fail()
	}
}

func TestConnection(t *testing.T) {
	client := &http.Client{}
	c := connection{
		url:    serverUrl,
		client: client,
	}
	resp, err := c.request("GET", "/", nil, nil, nil)
	if err != nil {
		t.Logf("Error: %v\n", err)
		t.Fail()
	} else if resp == nil {
		t.Fail()
	} else {
		jsonError := parseBody(resp, &couchReply)
		if jsonError != nil {
			t.Fail()
		} else {
			if resp.StatusCode != 200 ||
				couchReply.Couchdb != "Welcome" {
				t.Fail()
			}
			t.Logf("STATUS: %v\n", resp.StatusCode)
			t.Logf("couchdb: %v", couchReply.Couchdb)
			t.Logf("uuid: %v", couchReply.Uuid)
			t.Logf("version: %v", couchReply.Version)
			t.Logf("vendor: %v", couchReply.Vendor)
		}
	}
}

func TestBasicAuth(t *testing.T) {
	client := &http.Client{}
	auth := BasicAuth{Username: "adminuser", Password: "password"}
	c := connection{
		url:    serverUrl,
		client: client,
	}
	resp, err := c.request("GET", "/", nil, nil, &auth)
	if err != nil {
		t.Logf("Error: %v", err)
		t.Fail()
	} else if resp == nil {
		t.Logf("Response was nil")
		t.Fail()
	}
}

func TestProxyAuth(t *testing.T) {
	client := &http.Client{}
	pAuth := ProxyAuth{
		Username: "adminuser",
		Roles:    []string{"admin", "master", "_admin"},
	}
	c := connection{
		url:    serverUrl,
		client: client,
	}
	resp, err := c.request("GET", "/", nil, nil, &pAuth)
	if err != nil {
		t.Logf("Error: %v", err)
		t.Fail()
	} else if resp == nil {
		t.Logf("Response was nil")
		t.Fail()
	}
}

func TestBadAuth(t *testing.T) {
	client := &http.Client{}
	auth := BasicAuth{Username: "notauser", Password: "what?"}
	c := connection{
		url:    serverUrl,
		client: client,
	}
	resp, err := c.request("GET", "/", nil, nil, &auth)
	if err == nil {
		t.Fail()
	} else if resp.StatusCode != 401 {
		t.Logf("Wrong Status: %v", resp.StatusCode)
		t.Fail()
	}
}
