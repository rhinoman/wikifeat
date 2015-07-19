package couchdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

//represents a couchdb 'connection'
type connection struct {
	url    string
	client *http.Client
}

//processes a request
func (conn *connection) request(method, path string,
	body io.Reader, headers map[string]string, auth Auth) (*http.Response, error) {

	req, err := http.NewRequest(method, conn.url+path, body)
	//set headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if err != nil {
		return nil, err
	}
	if auth != nil {
		auth.AddAuthHeaders(req)
	}
	resp, err := conn.processResponse(0, req)
	if err == nil && resp != nil && auth != nil {
		auth.updateAuth(resp)
	}
	return resp, err
}

//Returns a result from couchdb directly to a requesting client
//Useful for downloading large files
func (conn *connection) reverseProxyRequest(w http.ResponseWriter,
	r *http.Request, path string, auth Auth) error {
	target, err := url.Parse(conn.url)
	if err != nil {
		return err
	}
	if auth != nil {
		auth.AddAuthHeaders(r)
	}
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, path)
	}
	rp := &httputil.ReverseProxy{Director: director}
	rp.ServeHTTP(w, r)
	return nil
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func (conn *connection) processResponse(numTries int,
	req *http.Request) (*http.Response, error) {
	resp, err := conn.client.Do(req)
	if err != nil {
		errStr := err.Error()
		// Because sometimes couchdb rudely
		// slams the connection shut and we get a race condition.
		// Of course, Go http presents one of two possibilities
		// for error strings, so we check for both.
		if (strings.Contains(errStr, "EOF") ||
			strings.Contains(errStr, "broken connection")) && numTries < 3 {
			//wait a bit and try again
			fmt.Printf("\nERROR! %v\n", errStr)
			time.Sleep(10 * time.Millisecond)
			numTries += 1
			return conn.processResponse(numTries, req)
		} else {
			return nil, err
		}
	} else if resp.StatusCode >= 400 {
		return resp, parseError(resp)
	} else {
		return resp, nil
	}
}

type Error struct {
	StatusCode int
	URL        string
	Method     string
	ErrorCode  string //empty for HEAD requests
	Reason     string //empty for HEAD requests
}

//stringify the error
func (err *Error) Error() string {
	return fmt.Sprintf("[Error]:%v: %v %v - %v %v",
		err.StatusCode, err.Method, err.URL, err.ErrorCode, err.Reason)
}

//extracts rev code from header
func getRevInfo(resp *http.Response) (string, error) {
	if rev := resp.Header.Get("ETag"); rev == "" {
		return "", fmt.Errorf("CouchDB did not return rev info")
	} else {
		return rev[1 : len(rev)-1], nil
	}
}

//unmarshalls a JSON Response Body
func parseBody(resp *http.Response, o interface{}) error {
	err := json.NewDecoder(resp.Body).Decode(&o)
	if err != nil {
		resp.Body.Close()
		return err
	} else {
		return resp.Body.Close()
	}
}

// encodes a struct to JSON and returns an io.Reader,
// the buffer size, and an error (if any)
func encodeData(o interface{}) (io.Reader, int, error) {
	if o == nil {
		return nil, 0, nil
	}
	buf, err := json.Marshal(&o)
	if err != nil {
		return nil, 0, err
	} else {
		return bytes.NewReader(buf), len(buf), nil
	}
}

//Parse a CouchDB error response
func parseError(resp *http.Response) error {
	var couchReply struct{ Error, Reason string }
	if resp.Request.Method != "HEAD" {
		err := parseBody(resp, couchReply)
		if err != nil {
			return fmt.Errorf("Unknown error accessing CouchDB: %v", err)
		}
	}
	return &Error{
		StatusCode: resp.StatusCode,
		URL:        resp.Request.URL.String(),
		Method:     resp.Request.Method,
		ErrorCode:  couchReply.Error,
		Reason:     couchReply.Reason,
	}
}

//smooshes url segments together
func buildString(pathSegments []string) string {
	urlString := ""
	for _, pathSegment := range pathSegments {
		urlString += "/"
		urlString += url.QueryEscape(pathSegment)
	}
	return urlString
}

//Build Url
func buildUrl(pathSegments ...string) (string, error) {
	var Url *url.URL
	urlString := buildString(pathSegments)
	Url, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}
	return Url.String(), nil
}

//Build Url with query arguments
func buildParamUrl(params url.Values, pathSegments ...string) (string, error) {
	var Url *url.URL
	urlString := buildString(pathSegments)
	Url, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}
	Url.RawQuery = params.Encode()
	return Url.String(), nil
}
