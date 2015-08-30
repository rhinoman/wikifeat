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
**/
package util

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"strings"
	"time"
)

//Just a package of utility functions that don't really belong
//anywhere else

func ContainsString(findStr string, strArr []string) bool {
	for _, str := range strArr {
		if str == findStr {
			return true
		}
	}
	return false
}

func HasRole(roles []string, role string) bool {
	return ContainsString(role, roles)
}

// Check if this user is an admin of anything.
// Either a site admin, or a wiki admin returns true
func IsAnyAdmin(roles []string) bool {
	for _, str := range roles {
		if str == "admin" || str == "master" ||
			strings.Contains(str, ":admin") {
			return true
		}
	}
	return false
}

func Retry(maxTries int, f func() error) error {
	numTries := 0
	return doRetry(numTries, maxTries, f)
}

func doRetry(numTries int, maxTries int, f func() error) error {
	err := f()
	if err != nil && numTries < maxTries {
		numTries += 1
		return doRetry(numTries, maxTries, f)
	}
	return err
}

//Adds quotes to a string
func QuotifyString(str string) string {
	return "\"" + str + "\""
}

// Adds quotes to strings in an array,
// used for couchdb url parameters
func ApplyQuotes(strList []string) {
	for i, v := range strList {
		if v != "{}" {
			strList[i] = QuotifyString(v)
		}
	}
}

/* Parses an array in a query string.
/*
/* Returns an array
*/
func ParseArrayParam(str string) ([]string, error) {
	//The first and last characters should be '[' and ']'
	if str[0] != '[' || str[len(str)-1] != ']' {
		return []string{}, errors.New("Malformed array string.")
	}
	strArray := strings.Split(str[1:len(str)-1], ",")
	return strArray, nil
}

func GenHashString(data string) string {
	now := time.Now().UTC()
	mac := hmac.New(sha1.New, []byte(now.String()))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}
