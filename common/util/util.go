/*
 *  Licensed to Wikifeat under one or more contributor license agreements.
 *  See the LICENSE.txt file distributed with this work for additional information
 *  regarding copyright ownership.
 *
 *  Redistribution and use in source and binary forms, with or without
 *  modification, are permitted provided that the following conditions are met:
 *
 *  * Redistributions of source code must retain the above copyright notice,
 *  this list of conditions and the following disclaimer.
 *  * Redistributions in binary form must reproduce the above copyright
 *  notice, this list of conditions and the following disclaimer in the
 *  documentation and/or other materials provided with the distribution.
 *  * Neither the name of Wikifeat nor the names of its contributors may be used
 *  to endorse or promote products derived from this software without
 *  specific prior written permission.
 *
 *  THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 *  AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 *  IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 *  ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
 *  LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 *  CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 *  SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 *  INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 *  CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 *  ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 *  POSSIBILITY OF SUCH DAMAGE.
 */

package util

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

//Returns the directory the executable is running from
func GetExecDirectory() (string, error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}

//Returns the "default" config file name and location
func DefaultConfigLocation() (string, error) {
	dir, err := GetExecDirectory()
	if err != nil {
		return "", err
	} else {
		return path.Join(dir, "config.ini"), nil
	}
}

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

//Generates a random token
func GenToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

//Encodes a struct into json and returns an io.Reader
func EncodeJsonData(data interface{}) (io.Reader, int, error) {
	if data == nil {
		return nil, 0, nil
	}
	if buf, err := json.Marshal(&data); err != nil {
		return nil, 0, err
	} else {
		return bytes.NewReader(buf), len(buf), nil
	}
}

func EncodeBase64Url(data []byte) string {
	b64 := base64.StdEncoding.EncodeToString(data)
	b64 = strings.Replace(b64, "+", "-", -1)
	b64 = strings.Replace(b64, "/", "_", -1)
	b64 = strings.Replace(b64, "=", "", -1)
	return b64
}

//Returns a minimum of two integers
func MinInt(x, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}

//Returns a maximum of two integers
func MaxInt(x, y int) int {
	if x > y {
		return x
	} else {
		return y
	}
}

//Returns the Absolute value of an int
func AbsInt(x int) int {
	if x < 0 {
		return -x
	} else {
		return x
	}
}
