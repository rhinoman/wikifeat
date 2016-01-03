package auth

import (
	"fmt"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/couchdb-go"
	"net/http"
	"strings"
)

type AuthType int

const (
	Standard AuthType = iota
)

func StringToAuthType(str string) AuthType {
	switch strings.ToUpper(str) {
	case "STANDARD":
		return Standard
	default:
		return Standard
	}
}

type Authenticator interface {
	// Reads Auth information/headers from a request and returns a CouchDB auth
	GetAuth(*http.Request) (couchdb.Auth, error)
	// Updates authentication information in the http response before sending it down
	// Usually, this is used to update Auth cookies, CSRF tokens, etc.
	SetAuth(http.ResponseWriter, couchdb.Auth)
}

func NewAuthenticator(authType AuthType) Authenticator {
	switch authType {
	case Standard:
		return &StandardAuthenticator{}
	default:
		return &StandardAuthenticator{}
	}
}

type AuthError struct {
	ErrorCode int
	Reason    string
}

func (err *AuthError) Error() string {
	return fmt.Sprintf("[Auth Error]:%v: %v", err.ErrorCode, err.Reason)
}
