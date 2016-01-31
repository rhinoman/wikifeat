package auth_service_test

import (
	"github.com/rhinoman/wikifeat/Godeps/_workspace/src/github.com/rhinoman/couchdb-go"
	"github.com/rhinoman/wikifeat/auth/auth_service"
	"github.com/rhinoman/wikifeat/common/config"
	"github.com/rhinoman/wikifeat/common/entities"
	"github.com/rhinoman/wikifeat/common/registry"
	"github.com/rhinoman/wikifeat/common/services"
	"github.com/rhinoman/wikifeat/users/user_service"
	"testing"
	"time"
)

var timeout = time.Duration(500 * time.Millisecond)
var um = new(user_service.UserManager)
var user = entities.User{
	UserName: "John.Smith",
	Password: "password",
}

func setup(t *testing.T) {
	config.LoadDefaults()
	config.ServiceRegistry.CacheRefreshInterval = 1000
	services.InitDb()
	//This will cause the registry manager to complain, but we don't
	//really need the service being registered here.
	registry.Init("TestAuth", "/services/test/auth")
	//We need to create a user in order to have any sessions, so
	registration := user_service.Registration{
		NewUser: user,
	}
	_, err := um.SetUp(&registration)
	if err != nil {
		t.Error(err)
	}
}

func afterTest(t *testing.T) {
	auth := &couchdb.BasicAuth{
		Username: "John.Smith",
		Password: "password",
	}
	userDoc, _ := services.GetUserFromAuth(auth)
	curUser := &entities.CurrentUserInfo{
		Auth: auth,
		User: userDoc,
	}
	um.Delete(user.UserName, curUser)
	services.DeleteDb(services.MainDbName())
}

func TestSessions(t *testing.T) {
	setup(t)
	// Test Standard
	am := auth_service.AuthManager{}
	sess, err := am.Create("John.Smith", "password", "standard")
	if err != nil {
		t.Error(err)
	}
	t.Logf("Session: %v", sess)
	defer afterTest(t)

}
