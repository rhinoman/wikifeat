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

package user_service

// Manager for User Records

import (
	"errors"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/couchdb-go"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/wikifeat-common/config"
	. "github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/wikifeat-common/entities"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/wikifeat-common/registry"
	. "github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/wikifeat-common/services"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/wikifeat-common/util"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type UserLoginCredentials struct {
	UserName string `json:"name"`
	Password string `json:"password"`
}

type UserListQueryResponse struct {
	ViewResponse
	Rows []UserListItem `json:"rows"`
}

type UserListItem struct {
	Id    string `json:"id"`
	Key   string `json:"key"`
	Value User   `json:"value"`
}

type RoleRequest struct {
	//ResourceType: wiki, main, etc.
	ResourceType string `json:"resourceType"`
	//The Uuid of the resource
	ResourceId string `json:"resourceId"`
	//AccessType is read, write, or admin
	AccessType string `json:"accessType"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type ResetTokenRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

func (rr *RoleRequest) roleString() string {
	if rr.ResourceType == "main" && rr.AccessType == "admin" {
		return "admin"
	} else {
		return rr.ResourceType + "_" +
			rr.ResourceId + ":" +
			rr.AccessType
	}
}

func (rr *RoleRequest) validate() bool {
	if rr.AccessType != "read" &&
		rr.AccessType != "write" &&
		rr.AccessType != "admin" {
		return false
	}
	if rr.ResourceType != "wiki" &&
		rr.ResourceType != "main" {
		return false
	}
	return true
}

type Registration struct {
	NewUser User `json:"user"`
}

type UserManager struct{}

// Create the first (master) user and perform various set up operations.
// This functionality was moved to an external script.
// Do NOT use -- It has been left here to facilitate unit tests only.
func (um *UserManager) SetUp(registration *Registration) (string, error) {
	err := um.validateUser(&registration.NewUser)
	if err != nil {
		return "", err
	}
	user := registration.NewUser
	//Create a new user
	//This will be the master user for this installation
	user.Roles = []string{AdminRole(MainDbName()),
		MasterRole(),
		AllUsersRole()}
	user.Type = "user"
	nowTime := time.Now().UTC()
	user.CreatedAt = nowTime
	user.ModifiedAt = nowTime
	userDb := Connection.SelectDB(UserDbName, AdminAuth)
	log.Printf("Registering new user account for: %v", user.UserName)

	namestring := UserPrefix + user.UserName
	//Try to create the main database
	if err := CreateDb(MainDbName()); err != nil {
		return "", err
	}
	if err := InitMainDatabase(); err != nil {
		DeleteDb(MainDbName())
		return "", err
	}
	if uRev, err := userDb.Save(&user, namestring, ""); err != nil {
		DeleteDb(MainDbName())
		return "", err
	} else {
		return uRev, nil
	}
}

//Create a normal user
func (um *UserManager) Create(newUser *User,
	curUser *CurrentUserInfo) (string, error) {
	//Who am I?
	theUser := curUser.User

	//check for admin
	if !util.HasRole(theUser.Roles, AdminRole(MainDbName())) &&
		!util.HasRole(theUser.Roles, MasterRole()) {
		return "", NotAdminError()
	}
	//ok, now create the user
	err := um.validateUser(newUser)
	if err != nil {
		return "", err
	}
	newUser.Roles = []string{ReadRole(MainDbName()), AllUsersRole()}
	newUser.Type = "user"
	namestring := UserPrefix + newUser.UserName
	nowTime := time.Now().UTC()
	newUser.CreatedAt = nowTime
	newUser.ModifiedAt = nowTime

	userDb := Connection.SelectDB(UserDbName, AdminAuth)
	log.Printf("Creating new user account for: %v", newUser.UserName)
	return userDb.Save(newUser, namestring, "")
}

//Delete a user
func (um *UserManager) Delete(id string,
	curUser *CurrentUserInfo) (string, error) {
	//Who am I?
	theUser := curUser.User
	//check for admin
	if !util.HasRole(theUser.Roles, AdminRole(MainDbName())) {
		return "", NotAdminError()
	}
	//pull the user record
	namestring := UserPrefix + id
	readUser := new(User)
	rev, err := um.Read(id, readUser, curUser)
	if err != nil {
		return "", err
	}
	//Do the delete
	theDb := Connection.SelectDB(UserDbName, AdminAuth)
	return theDb.Delete(namestring, rev)
}

//Update a user
func (um *UserManager) Update(id string, rev string, updatedUser *User,
	curUser *CurrentUserInfo) (string, error) {
	//Who am I?
	theUser := curUser.User
	var auth couchdb.Auth
	//make sure Id matches the user object
	if id != updatedUser.UserName {
		return "", BadRequestError()
	}
	//check for admin
	if util.HasRole(theUser.Roles, AdminRole(MainDbName())) ||
		util.HasRole(theUser.Roles, MasterRole()) {
		auth = AdminAuth
	} else {
		auth = curUser.Auth
	}
	userDb := Connection.SelectDB(UserDbName, auth)
	//pull the user record
	namestring := UserPrefix + id
	//readUser := FullUserRecord{}
	var userData interface{}
	_, err := userDb.Read(namestring, &userData, nil)
	if err != nil {
		return "", err
	}
	readUser := userData.(map[string]interface{})
	//Update the user parameters of the read user
	//We do this instead of just pushing the updated user to CouchDB,
	//because the updated user doesn't contain password information,
	//and we don't want to wipe away the password.
	nowTime := time.Now().UTC()
	readUser["modifiedAt"] = nowTime
	readUser["userPublic"] = updatedUser.Public
	readUser["userName"] = updatedUser.UserName
	readUser["password_reset"] = updatedUser.PassResetToken
	//And save the updated user
	uRev, err := userDb.Save(&readUser, namestring, rev)
	if err != nil {
		return "", err
	}
	//Censor the password
	//updatedUser.Password = ""
	return uRev, nil

}

// Change a user's password
func (um *UserManager) ChangePassword(id string, rev string,
	cpr *ChangePasswordRequest, curUser *CurrentUserInfo) (string, error) {
	theUser := curUser.User
	var auth couchdb.Auth
	//check for admin
	isAdmin := util.HasRole(theUser.Roles, AdminRole(MainDbName())) ||
		util.HasRole(theUser.Roles, MasterRole())

	if isAdmin {
		auth = AdminAuth
	} else {
		auth = curUser.Auth
	}
	//Must we validate the old password?
	if !isAdmin || (isAdmin && theUser.UserName == id) {
		//yes, let's use the old password to do this thing
		auth = &couchdb.BasicAuth{Username: id, Password: cpr.OldPassword}
	}
	//validate password
	err := um.validatePassword(cpr.NewPassword)
	if err != nil {
		return "", err
	}
	userDb := Connection.SelectDB(UserDbName, auth)
	//pull the user record
	namestring := UserPrefix + id
	var userData interface{}
	_, err = userDb.Read(namestring, &userData, nil)
	if err != nil {
		return "", err
	}
	readUser := userData.(map[string]interface{})
	//Update JUST the password (and modified time)
	nowTime := time.Now().UTC()
	readUser["modifiedAt"] = nowTime
	readUser["password"] = cpr.NewPassword
	//Save the updated user
	uRev, err := userDb.Save(&readUser, namestring, rev)
	if err != nil {
		return "", err
	}
	return uRev, nil

}

func (um *UserManager) ResetPassword(id string, tr *ResetTokenRequest) error {
	//Read the user
	user := User{}
	userDb := Connection.SelectDB(UserDbName, AdminAuth)
	rev, err := userDb.Read(UserPrefix+id, &user, nil)
	if err != nil {
		return err
	}
	//Check the token
	tok := user.PassResetToken.Token
	expireTime := user.PassResetToken.Expires
	nowTime := time.Now().UTC()
	if tok != tr.Token || nowTime.After(expireTime) {
		return errors.New("[Error]:400: This token is invalid or expired.")
	}
	cpr := ChangePasswordRequest{NewPassword: tr.NewPassword}
	rev, err = um.ChangePassword(id, rev, &cpr, GetAdminUser())
	if err != nil {
		return err
	}
	//Now, expire the token
	user.PassResetToken.Token = ""
	user.PassResetToken.Expires = nowTime
	_, err = um.Update(id, rev, &user, GetAdminUser())
	return err
}

func resourceDbName(rr *RoleRequest) string {
	if rr.ResourceType == "main" {
		return MainDbName()
	}
	return rr.ResourceType + "_" + rr.ResourceId
}

//Grant a role to a user
func (um *UserManager) GrantRole(id string,
	grantRequest *RoleRequest,
	curUser *CurrentUserInfo) (string, error) {

	newRole := grantRequest.roleString()
	//Who am I?
	theUser := curUser.User
	//Make sure user is an admin of the resource being requested
	resourceDbName := resourceDbName(grantRequest)
	if !util.HasRole(theUser.Roles, AdminRole(MainDbName())) &&
		!util.HasRole(theUser.Roles, MasterRole()) &&
		!util.HasRole(theUser.Roles, AdminRole(resourceDbName)) {
		//Not an admin
		return "", NotAdminError()
	}
	//Are we trying to grant the site admin role?
	//If so, curUser must be a master user
	if newRole == AdminRole(MainDbName()) &&
		!util.HasRole(theUser.Roles, MasterRole()) {
		return "", NotAdminError()

	}
	//Fetch the user
	userDb := Connection.SelectDB(UserDbName, AdminAuth)
	readUser := new(User)
	rev, err := userDb.Read(UserPrefix+id, readUser, nil)
	if err != nil {
		return "", err
	}
	//validate the role request
	if grantRequest.validate() == false {
		return "", BadRequestError()
	}
	//All is well, grant the role
	rev, err = Connection.GrantRole(id, newRole, AdminAuth)
	if err != nil {
		return "", err
	}
	return rev, nil
}

//Revoke user access
func (um *UserManager) RevokeRole(id string,
	revokeRequest *RoleRequest,
	curUser *CurrentUserInfo) (string, error) {

	theUser := curUser.User
	//Make sure user is an admin of the resource being requested
	resourceDbName := resourceDbName(revokeRequest)
	if !util.HasRole(theUser.Roles, AdminRole(MainDbName())) &&
		!util.HasRole(theUser.Roles, AdminRole(resourceDbName)) {
		//Not an admin
		return "", NotAdminError()
	}
	//Fetch the user
	userDb := Connection.SelectDB(UserDbName, AdminAuth)
	readUser := new(User)
	rev, err := userDb.Read(UserPrefix+id, readUser, nil)
	if err != nil {
		return "", err
	}
	//validate the role request
	if revokeRequest.validate() == false {
		return "", BadRequestError()
	}
	//ok, revoke the role
	revokeRole := revokeRequest.roleString()
	rev, err = Connection.RevokeRole(id, revokeRole, AdminAuth)
	if err != nil {
		return "", err
	}
	return rev, nil
}

//Get a user
func (um *UserManager) Read(id string, user *User,
	curUser *CurrentUserInfo) (string, error) {
	//Who am I?
	theUser := curUser.User
	readUser := User{}
	if util.HasRole(theUser.Roles, AdminRole(MainDbName())) {
		//This is an admin
		userDb := Connection.SelectDB(UserDbName, AdminAuth)
		rev, err := userDb.Read(UserPrefix+id, &readUser, nil)
		if err != nil {
			return "", err
		} else {
			*user = readUser
			return rev, nil
		}
	} else {
		userDb := Connection.SelectDB(UserDbName, curUser.Auth)
		if rev, err := userDb.Read(UserPrefix+id, &readUser, nil); err != nil {
			return "", err
		} else {
			*user = readUser
			return rev, nil
		}
	}
}

//Request a password reset (forgot password, etc).
func (um *UserManager) RequestPasswordReset(id string) error {
	//Read the user
	user := User{}
	userDb := Connection.SelectDB(UserDbName, AdminAuth)
	rev, err := userDb.Read(UserPrefix+id, &user, nil)
	if err != nil {
		return err
	}
	//Generate a token
	tok := util.GenToken()
	//Set token expiration time
	nowTime := time.Now().UTC()
	//Tokens are good for 4 hours
	hours := time.Duration(4) * time.Hour
	expireTime := nowTime.Add(hours)
	user.PassResetToken.Token = tok
	user.PassResetToken.Expires = expireTime
	//Now save the user
	log.Println("Saving reset token to user document")
	rev, err = um.Update(id, rev, &user, GetAdminUser())
	if err != nil {
		return err
	}
	//Now we need to send an email to the user containing our token
	notifEndpoint, err := registry.GetServiceLocation("notifications")
	if err != nil {
		return err
	}
	nr := NotificationRequest{
		To:      user.Public.Contact.Email,
		Subject: "Reset Password Request",
		Data: map[string]string{
			"user": user.Public.FirstName,
			"uri": "/reset_password?user=" + id +
				"&token=" + tok,
		},
	}
	nrJson, _, err := util.EncodeJsonData(&nr)
	if err != nil {
		return err
	}
	//Assemble the request
	reqUrl := notifEndpoint + "/api/v1/notifications/reset_password/send"
	client := &http.Client{}
	request, err := http.NewRequest("POST", reqUrl, nrJson)
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

//Login a user and create a new session
func (um *UserManager) Login(credentials *UserLoginCredentials) (*couchdb.CookieAuth, error) {
	return Connection.CreateSession(credentials.UserName, credentials.Password)
}

//Logout
func (um *UserManager) Logout(sessionToken string) error {
	return Connection.DestroySession(&couchdb.CookieAuth{AuthToken: sessionToken})
}

//Get list of users
func (um *UserManager) GetUserList(pageNum int, numPerPage int,
	ulr *UserListQueryResponse, curUser *CurrentUserInfo) error {
	//You have to be some sort of admin to do this
	if !util.IsAnyAdmin(curUser.User.Roles) {
		return NotAdminError()
	}
	params := url.Values{}
	if numPerPage != 0 {
		params.Add("limit", strconv.Itoa(numPerPage))
	}
	skip := numPerPage * (pageNum - 1)
	if skip > 0 {
		params.Add("skip", strconv.Itoa(skip))
	}
	userDb := Connection.SelectDB("_users", AdminAuth)
	err := userDb.GetView("user_queries", "listUsers", &ulr, &params)
	if err != nil {
		return err
	}
	return nil
}

//Get list of users by role(s)
func (um *UserManager) GetUserListForRole(pageNum int, numPerPage int,
	roles []string, ulr *UserListQueryResponse, curUser *CurrentUserInfo) error {
	util.ApplyQuotes(roles)
	keyVals := "[" + strings.Join(roles, ",") + "]"
	userDb := Connection.SelectDB("_users", AdminAuth)
	type IntResult struct {
		value int
		err   error
	}
	type UlResult struct {
		value UserListQueryResponse
		err   error
	}
	//These two CouchDB requests may be run concurrently
	c_count := make(chan IntResult)
	c_results := make(chan UlResult)
	go func() {
		cr := CountResponse{}
		params := url.Values{}
		params.Set("keys", keyVals)
		params.Add("group", "true")
		if err := userDb.GetList("user_queries", "browseUsers",
			"usersByRole", &cr, &params); err == nil {
			if len(cr.Rows) > 0 {
				c_count <- IntResult{value: cr.Rows[0].Value, err: nil}
			} else {
				c_count <- IntResult{value: 0, err: nil}
			}
		} else {
			log.Printf("\nError in User Role Count Query: %v\n", err)
			c_count <- IntResult{value: 0, err: err}
		}
	}()
	go func() {
		params := url.Values{}
		params.Set("keys", keyVals)
		if numPerPage != 0 {
			params.Add("limit", strconv.Itoa(numPerPage))
		}
		skip := numPerPage * (pageNum - 1)
		if skip > 0 {
			params.Add("skip", strconv.Itoa(skip))
		}
		qr := UserListQueryResponse{}
		params.Add("reduce", "false")
		if err := userDb.GetList("user_queries", "browseUsers",
			"usersByRole", &qr, &params); err == nil {
			c_results <- UlResult{value: qr, err: nil}
		} else {
			log.Printf("\nError in User Role Query: %v\n", err)
			c_results <- UlResult{value: UserListQueryResponse{}, err: err}
		}
	}()
	results := <-c_results
	count := <-c_count
	if results.err != nil {
		return results.err
	} else if count.err != nil {
		return count.err
	} else {
		*ulr = results.value
		ulr.TotalRows = count.value
	}
	return nil
}

//Simple validation of user data
func (um *UserManager) validateUser(user *User) error {
	var err error
	if len(user.UserName) < 3 {
		err = errors.New("Username invalid")
	}
	err = um.validatePassword(user.Password)
	return err
}

//Validate Password
func (um *UserManager) validatePassword(password string) error {
	if len(password) < config.Auth.MinPasswordLength {
		return errors.New("Password too short")
	} else {
		return nil
	}
}
