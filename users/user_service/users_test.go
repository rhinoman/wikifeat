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

package user_service_test

import (
	"fmt"
	"github.com/rhinoman/couchdb-go"
	"github.com/rhinoman/wikifeat/common/config"
	"github.com/rhinoman/wikifeat/common/entities"
	"github.com/rhinoman/wikifeat/common/services"
	"github.com/rhinoman/wikifeat/common/util"
	"github.com/rhinoman/wikifeat/users/user_service"
	"github.com/rhinoman/wikifeat/vendor/github.com/twinj/uuid"
	. "github.com/smartystreets/goconvey/convey"
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
	services.InitDb()
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
	userDoc, err := services.GetUserFromAuth(auth)
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
	defer services.DeleteDb(services.MainDb)
	smithUser := func() *entities.CurrentUserInfo {
		smithAuth := &couchdb.BasicAuth{Username: "Steven.Smith", Password: "jabberwocky"}
		smith, err := services.GetUserFromAuth(smithAuth)
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
