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
	. "github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/smartystreets/goconvey/convey"
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
	Convey("Given a new account registration", t, func() {
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
		Convey("When the new user is registered", func() {
			rev, err := um.SetUp(&registration)
			Convey("The revision should be set and the error should be nil", func() {
				So(err, ShouldBeNil)
				So(rev, ShouldNotEqual, "")
				t.Logf("New user revision: %v", rev)
			})
		})
		Convey("When a new user for the account is created", func() {
			subUser := entities.User{
				UserName: "Sally.Smith",
				Password: "123456",
				Public: entities.UserPublic{
					LastName:  "Smith",
					FirstName: "Sally",
				},
			}
			rev, err := um.Create(&subUser, smithUser())
			Convey("The revision should be set and the error should be nil", func() {
				So(err, ShouldBeNil)
				So(rev, ShouldNotEqual, "")
			})
			Convey("The user should have the appropriate roles", func() {
				So(util.HasRole(subUser.Roles, "all_users"), ShouldEqual, true)
			})
		})
		Convey("When the user is updated", func() {
			auth := couchdb.BasicAuth{Username: "Sally.Smith", Password: "123456"}
			updateUser := entities.User{}
			rev, err := um.Read("Sally.Smith", &updateUser, smithUser())
			So(err, ShouldBeNil)
			curUser := entities.CurrentUserInfo{
				User: &updateUser,
				Auth: &auth,
			}

			updateUser.Public.MiddleName = "Marie"
			rev, err = um.Update("Sally.Smith", rev, &updateUser, &curUser)
			Convey("Error should be nil and the revision should be set", func() {
				So(err, ShouldBeNil)
				So(rev, ShouldNotEqual, "")
			})
			rev, err = um.Read("Sally.Smith", &updateUser, smithUser())
			Convey("The user's password should still be valid", func() {
				So(err, ShouldBeNil)
				So(rev, ShouldNotEqual, "")
			})

		})
		Convey("When the user's password is changed", func() {
			auth := couchdb.BasicAuth{Username: "Sally.Smith", Password: "123456"}
			updateUser := entities.User{}
			rev, err := um.Read("Sally.Smith", &updateUser, smithUser())
			So(err, ShouldBeNil)
			curUser := entities.CurrentUserInfo{
				User: &updateUser,
				Auth: &auth,
			}
			newPassword := "234567"
			cpr := user_service.ChangePasswordRequest{
				NewPassword: newPassword,
				OldPassword: "123456",
			}
			rev, err = um.ChangePassword("Sally.Smith", rev, &cpr, &curUser)
			Convey("Error should be nil and the revision should be set", func() {
				So(err, ShouldBeNil)
				So(rev, ShouldNotEqual, "")
			})
			rev, err = um.Read("Sally.Smith", &updateUser, &curUser)
			Convey("Old password should NOT work", func() {
				So(err, ShouldNotBeNil)
			})
			curUser.Auth = &couchdb.BasicAuth{Username: "Sally.Smith", Password: "234567"}
			rev, err = um.Read("Sally.Smith", &updateUser, &curUser)
			Convey("New password should work", func() {
				So(err, ShouldBeNil)
				So(rev, ShouldNotEqual, "")
			})

		})
		Convey("When a user is granted a role", func() {
			curUser := smithUser()
			roleRequest := user_service.RoleRequest{
				ResourceType: "main",
				ResourceId:   "",
				AccessType:   "write",
			}
			rev, err := um.GrantRole("Sally.Smith", &roleRequest, curUser)
			Convey("The revision should be set and the error should be nil", func() {
				So(err, ShouldBeNil)
				So(rev, ShouldNotEqual, "")
			})
			Convey("User should have her new role", func() {
				readUser := new(entities.User)
				_, err := um.Read("Sally.Smith", readUser, smithUser())
				So(err, ShouldBeNil)
				t.Logf("Sally's Record: %v", readUser)
				searchRole := "main_" + ":write"
				searchRoleFound := false
				if readUser.Roles[2] == searchRole {
					searchRoleFound = true
				}
				So(searchRoleFound, ShouldEqual, true)
			})

		})
		Convey("When a user has a role revoked", func() {
			curUser := smithUser()
			roleRequest := user_service.RoleRequest{
				ResourceType: "main",
				ResourceId:   "",
				AccessType:   "write",
			}
			rev, err := um.RevokeRole("Sally.Smith", &roleRequest, curUser)
			Convey("The revision should be set and the error should be nil", func() {
				So(err, ShouldBeNil)
				So(rev, ShouldNotEqual, "")
			})
			Convey("User should no longer have her role", func() {
				readUser := new(entities.User)
				_, err := um.Read("Sally.Smith", readUser, smithUser())
				So(err, ShouldBeNil)
				searchRole := "main" + ":write"
				searchRoleFound := false
				if len(readUser.Roles) > 1 && readUser.Roles[1] == searchRole {
					searchRoleFound = true
				}
				So(searchRoleFound, ShouldEqual, false)

			})
		})
		Convey("When the user list is requested", func() {
			userList := user_service.UserListQueryResponse{}
			err := um.GetUserList(1, 5, &userList, smithUser())
			Convey("Error should be nil and we should have some results", func() {
				So(err, ShouldBeNil)
				So(len(userList.Rows) >= 2, ShouldBeTrue)
				t.Logf("UserListResponse: %v", userList)
			})
		})
		Convey("When a user search is requested", func() {
			userList := user_service.UserListQueryResponse{}
			err := um.SearchForUsersByName(1, 5, "Smith", &userList, smithUser())
			Convey("Error should be nil and we should have some results", func() {
				So(err, ShouldBeNil)
				So(len(userList.Rows), ShouldEqual, 2)
				t.Logf("UserListResponse: %v", userList)
			})
		})
		Convey("When the user is read", func() {
			readUser := entities.User{}
			rev, err := um.Read(user.UserName, &readUser, smithUser())
			Convey("The revision should be set and the error should be nil", func() {
				So(err, ShouldBeNil)
				So(rev, ShouldNotEqual, "")
			})
			Convey("The user data should be set", func() {
				So(readUser.CreatedAt, ShouldNotBeNil)
				So(readUser.UserName, ShouldEqual, "Steven.Smith")
			})
		})
		Convey("When the user by roles list is requested", func() {
			userList := user_service.UserListQueryResponse{}
			err := um.GetUserListForRole(1, 5, []string{"all_users"},
				&userList, smithUser())
			Convey("The error should be nil and we should have some results", func() {
				So(err, ShouldBeNil)
				t.Logf("Response: %v", userList)
			})
		})
		Convey("When a user password reset is requested", func() {
			err := um.RequestPasswordReset("Steven.Smith")
			Convey("The error should indicate no notifcation services are present", func() {
				So(err.Error(), ShouldEqual, "No notifications services listed!")
			})
		})
		Convey("When the user is deleted", func() {
			rev, err := um.Delete("Sally.Smith", smithUser())
			Convey("The revision should be set and the error should be nil", func() {
				So(err, ShouldBeNil)
				So(rev, ShouldNotEqual, "")
				t.Logf("Deleted User rev: %v", rev)
			})
		})
	})
}
