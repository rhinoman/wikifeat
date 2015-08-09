/**   Copyright (c) 2014-present James Adam.  All rights reserved.
*
* This file is part of WikiFeat.
*
*     WikiFeat is free software: you can redistribute it and/or modify
* it under the terms of the GNU General Public License as published by
* the Free Software Foundation, either version 2 of the License, or
* (at your option) any later version.
*
*     WikiFeat is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU General Public License for more details.
*
*     You should have received a copy of the GNU General Public License
* along with WikiFeat.  If not, see <http://www.gnu.org/licenses/>.
 */

package user_service

// Manager for User Records

import (
	"errors"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/couchdb-go"
	"github.com/rhinoman/wikifeat/common/config"
	. "github.com/rhinoman/wikifeat/common/entities"
	. "github.com/rhinoman/wikifeat/common/services"
	"github.com/rhinoman/wikifeat/common/util"
	"log"
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

//Set up -- Create the first (master) user
//This performs several actions required to set up a new installation
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
	//Who am I?
	theUser := curUser.User
	//Make sure user is an admin of the resource being requested
	resourceDbName := resourceDbName(grantRequest)
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
	if grantRequest.validate() == false {
		return "", BadRequestError()
	}
	//All is well, grant the role
	newRole := grantRequest.roleString()
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
	//Turns out non-admins need to view a list of users, too.
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
