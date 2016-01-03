package auth

import (
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/couchdb-go"
	"github.com/rhinoman/wikifeat/common/util"
	"net/http"
)

/**
 * With a StandardAuthenticator, we're just going against CouchDB directly
 * for Authentication
 */
type StandardAuthenticator struct{}

func (sta *StandardAuthenticator) unauthError() error {
	return &AuthError{
		ErrorCode: 401,
		Reason:    "Unauthenticated",
	}
}

func (sta *StandardAuthenticator) GetAuth(req *http.Request) (couchdb.Auth, error) {

	//What kind of authentication type do we have?
	if req.Header.Get("Authorization") != "" {
		return &couchdb.PassThroughAuth{
			AuthHeader: req.Header.Get("Authorization"),
		}, nil
	}
	//Ok, check for a session cookie
	sessCookie, err := req.Cookie("AuthSession")
	if err != nil {
		return nil, sta.unauthError()
	}
	//TODO: CHECK FOR EXPIRED COOKIES?
	sessionToken := sessCookie.Value
	csrfErr := sta.checkCsrf(req)
	if sessionToken == "" || csrfErr != nil {
		//Bad user, no cookie
		return nil, sta.unauthError()
	}
	//Return the cookie auth
	return &couchdb.CookieAuth{
		AuthToken: sessionToken,
	}, nil
}

func (sta *StandardAuthenticator) SetAuth(rw http.ResponseWriter, cAuth couchdb.Auth) {
	authData := cAuth.GetUpdatedAuth()
	if authData == nil {
		return
	}
	if val, ok := authData["AuthSession"]; ok {
		authCookie := http.Cookie{
			Name:     "AuthSession",
			Value:    val,
			Path:     "/",
			HttpOnly: true,
		}
		//Create a CSRF cookie
		csrfCookie := http.Cookie{
			Name:     "CsrfToken",
			Value:    util.GenHashString(val),
			Path:     "/",
			HttpOnly: false,
		}
		rw.Header().Add("Set-Cookie", authCookie.String())
		rw.Header().Add("Set-Cookie", csrfCookie.String())
	}
}

func (sta *StandardAuthenticator) checkCsrf(request *http.Request) error {
	//Check CSRF token
	csrfCookie, err := request.Cookie("CsrfToken")
	ourCsrf := request.Header.Get("X-Csrf-Token")
	if err != nil || ourCsrf == "" {
		return sta.unauthError()
	} else if token := csrfCookie.Value; token != ourCsrf {
		return sta.unauthError()
	} else {
		return nil
	}
}
