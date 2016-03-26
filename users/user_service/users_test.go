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
package user_service_test

import (
	"fmt"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/couchdb-go"
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/twinj/uuid"
	"github.com/rhinoman/wikifeat/common/config"
	"github.com/rhinoman/wikifeat/common/database"
	"github.com/rhinoman/wikifeat/common/entities"
	"github.com/rhinoman/wikifeat/common/util"
	"github.com/rhinoman/wikifeat/users/user_service"
	"testing"
	"time"
)

var timeout = time.Duration(500 * time.Millisecond)
var server = "127.0.0.1"
var adminAuth = &couchdb.BasicAuth{Username: "adminuser", Password: "password"}

var um = new(user_service.UserManager)

//var pm = new(managers.PageManager)

func setup() {
	config.LoadDefaults()
	database.InitDb()
	database.SetupDb()
}

func getUuid() string {
	theUuid := uuid.NewV4()
	return uuid.Formatter(theUuid, uuid.Clean)
}

func grabUser(id string, user *entities.User, auth couchdb.Auth) (string, error) {
	curUser := getCurUser(auth)
	return um.Read(id, user, curUser)
}

func getCurUser(auth couchdb.Auth) *entities.CurrentUserInfo {
	userDoc, err := database.GetUserFromAuth(auth)
	if err != nil {
		fmt.Printf("\nERROR: %v\n", err)
	}
	return &entities.CurrentUserInfo{
		Auth: auth,
		User: userDoc,
	}

}

func TestUsers(t *testing.T) {
	setup()
	defer database.DeleteDb(database.MainDb)
	smithUser := func() *entities.CurrentUserInfo {
		smithAuth := &couchdb.BasicAuth{Username: "Steven.Smith", Password: "jabberwocky"}
		smith, err := database.GetUserFromAuth(smithAuth)
		if err != nil {
			t.Error(err)
			return &entities.CurrentUserInfo{}
		}
		return &entities.CurrentUserInfo{
			Auth: smithAuth,
			User: smith,
		}
	}
	//after test cleanup
	defer func() {
		getSmithUser := func() *entities.CurrentUserInfo {
			return smithUser()
		}
		um.Delete("Steven.Smith", getSmithUser())
	}()
	//Create a new master user
	user := entities.User{
		UserName: "Steven.Smith",
		Password: "jabberwocky",
		Public: entities.UserPublic{
			LastName:  "Smith",
			FirstName: "Steven",
		},
	}
	registration := user_service.Registration{
		NewUser: user,
	}
	rev, err := um.SetUp(&registration)
	if err != nil {
		t.Error(err)
	}
	if rev == "" {
		t.Error("rev is empty!")
	}
	t.Logf("New user revision: %v", rev)
	//Create a new subuser
	subUser := entities.User{
		UserName: "Sally.Smith",
		Password: "123456",
		Public: entities.UserPublic{
			LastName:  "Smith",
			FirstName: "Sally",
		},
	}
	rev, err = um.Create(&subUser, smithUser())
	if err != nil {
		t.Error(err)
	}
	if rev == "" {
		t.Error("rev is empty!")
	}
	if !util.HasRole(subUser.Roles, "all_users") {
		t.Error("user doesn't have all_users role!")
	}
	//Update user
	auth := couchdb.BasicAuth{Username: "Sally.Smith", Password: "123456"}
	updateUser := entities.User{}
	rev, err = um.Read("Sally.Smith", &updateUser, smithUser())
	if err != nil {
		t.Error(err)
	}
	curUser := entities.CurrentUserInfo{
		User: &updateUser,
		Auth: &auth,
	}

	updateUser.Public.MiddleName = "Marie"
	rev, err = um.Update("Sally.Smith", rev, &updateUser, &curUser)
	if err != nil {
		t.Error(err)
	}
	if rev == "" {
		t.Error("rev is empty!")
	}
	rev, err = um.Read("Sally.Smith", &updateUser, smithUser())
	if err != nil {
		t.Error(err)
	}
	if rev == "" {
		t.Error("rev is empty!")
	}
	//Change password
	auth = couchdb.BasicAuth{Username: "Sally.Smith", Password: "123456"}
	updateUser = entities.User{}
	rev, err = um.Read("Sally.Smith", &updateUser, smithUser())
	if err != nil {
		t.Error(err)
	}
	curUser = entities.CurrentUserInfo{
		User: &updateUser,
		Auth: &auth,
	}
	newPassword := "234567"
	cpr := user_service.ChangePasswordRequest{
		NewPassword: newPassword,
		OldPassword: "123456",
	}
	rev, err = um.ChangePassword("Sally.Smith", rev, &cpr, &curUser)
	if err != nil {
		t.Error(err)
	}
	if rev == "" {
		t.Error("rev is empty!")
	}
	rev, err = um.Read("Sally.Smith", &updateUser, &curUser)
	if err == nil {
		t.Error("Old password should not have worked!")
	}
	curUser.Auth = &couchdb.BasicAuth{Username: "Sally.Smith", Password: "234567"}
	rev, err = um.Read("Sally.Smith", &updateUser, &curUser)
	if err != nil {
		t.Error(err)
	}
	if rev == "" {
		t.Error("rev is empty!")
	}
	//Grant role
	curUser = *smithUser()
	roleRequest := user_service.RoleRequest{
		ResourceType: "main",
		ResourceId:   "",
		AccessType:   "write",
	}
	rev, err = um.GrantRole("Sally.Smith", &roleRequest, &curUser)
	if err != nil {
		t.Error(err)
	}
	if rev == "" {
		t.Error("rev is empty!")
	}
	readUser := new(entities.User)
	_, err = um.Read("Sally.Smith", readUser, smithUser())
	if err != nil {
		t.Error(err)
	}
	t.Logf("Sally's Record: %v", readUser)
	searchRole := "main_" + ":write"
	searchRoleFound := false
	if readUser.Roles[2] == searchRole {
		searchRoleFound = true
	}
	if !searchRoleFound {
		t.Error("Role not found!")
	}
	//Revoke role
	curUser = *smithUser()
	roleRequest = user_service.RoleRequest{
		ResourceType: "main",
		ResourceId:   "",
		AccessType:   "write",
	}
	rev, err = um.RevokeRole("Sally.Smith", &roleRequest, &curUser)
	if err != nil {
		t.Error(err)
	}
	if rev == "" {
		t.Error("rev is empty!")
	}
	readUser = new(entities.User)
	_, err = um.Read("Sally.Smith", readUser, smithUser())
	if err != nil {
		t.Error(err)
	}
	searchRole = "main" + ":write"
	searchRoleFound = false
	if len(readUser.Roles) > 1 && readUser.Roles[1] == searchRole {
		searchRoleFound = true
	}
	if searchRoleFound {
		t.Error("Role should not have been found!")
	}
	//User List
	userList := user_service.UserListQueryResponse{}
	err = um.GetUserList(1, 5, &userList, smithUser())
	if err != nil {
		t.Error(err)
	}
	if !(len(userList.Rows) >= 2) {
		t.Error("User list should be greater than or equal to 2!")
	}
	t.Logf("UserListResponse: %v", userList)
	//User search
	userList = user_service.UserListQueryResponse{}
	err = um.SearchForUsersByName(1, 5, "Smith", &userList, smithUser())
	if err != nil {
		t.Error(err)
	}
	if len(userList.Rows) != 2 {
		t.Errorf("Userlist should be length 2 was %v", len(userList.Rows))
	}
	t.Logf("UserListResponse: %v", userList)
	//Read user
	readUser = &entities.User{}
	rev, err = um.Read(user.UserName, readUser, smithUser())
	if err != nil {
		t.Error(err)
	}
	if rev == "" {
		t.Error("Rev is empty!")
	}
	if readUser.UserName != "Steven.Smith" {
		t.Errorf("username should be Seteven.Smith, was %v", readUser.UserName)
	}
	//User by Roles list
	userList = user_service.UserListQueryResponse{}
	err = um.GetUserListForRole(1, 5, []string{"all_users"},
		&userList, smithUser())
	if err != nil {
		t.Error(err)
	}
	t.Logf("Response: %v", userList)
	//Password reset
	err = um.RequestPasswordReset("Steven.Smith")
	if err.Error() != "No notifications services listed!" {
		t.Error("Error is wrong")
	}
	//Delete user
	rev, err = um.Delete("Sally.Smith", smithUser())
	if err != nil {
		t.Error(err)
	}
	if rev == "" {
		t.Error("Rev is empty!")
	}
	t.Logf("Deleted User rev: %v", rev)
}
