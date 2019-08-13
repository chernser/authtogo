package oauth2

import (
	"context"
	_ "encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/chernser/authtogo/auth"
	"github.com/chernser/authtogo/datastore"
	"github.com/chernser/authtogo/utils"

	"gotest.tools/assert"

	"github.com/buaazp/fasthttprouter"
	_ "github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttputil"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/rs/zerolog/log"
)

var router *fasthttprouter.Router
var netListener *fasthttputil.InmemoryListener

type mockAuthServer struct {
	VolatileStorage auth.Storage
	SecretsStorage  auth.Storage
}

func (mockServer *mockAuthServer) RegisterRoute(method string, path string, handler fasthttp.RequestHandler) {
	router.Handle(method, path, handler)
}

// GetVolatileStorage returns volatile storage for tokens and runtime information
func (mockServer *mockAuthServer) GetVolatileStorage() auth.Storage {
	return mockServer.VolatileStorage
}

// GetSecretsStorage returns storage of secrets
func (mockServer *mockAuthServer) GetSecretsStorage() auth.Storage {
	return mockServer.SecretsStorage
}

type mockSessionManager struct {
	IsAuthorizedSession bool
}

func (sessions *mockSessionManager) StartSession(context interface{}) error {
	return nil
}

// InvalidateSession removes session by id from all auth related stores
func (sessions *mockSessionManager) InvalidateSession(context interface{}) error {
	return nil
}

// IsAuthenticated returns true if context contains information about authenticated session
func (sessions *mockSessionManager) IsAuthenticated(context interface{}) bool {
	log.Info().Msgf("Handling isAuthenticated request. Will return %v", sessions.IsAuthorizedSession)
	return sessions.IsAuthorizedSession
}

func TestMain(m *testing.M) {
	router = fasthttprouter.New()
	netListener = fasthttputil.NewInmemoryListener()
	defer netListener.Close()

	// err := InitAccountsAPI(viper.GetViper(), router)
	// if err != nil {
	// 	fmt.Printf("Failed to init account API %s\n", err)
	// 	os.Exit(1)
	// }
	// apiClient(netListener, router)
	os.Exit(m.Run())
}

func TestAuthorization(t *testing.T) {
	server := &mockAuthServer{}
	server.VolatileStorage = datastore.CreateMemoryStorage()
	server.SecretsStorage = datastore.CreateMemoryStorage()

	row := map[string]interface{}{
		"userId": "11111",
		"domain": "http://localhost",
		"secret": "999999",
	}
	server.SecretsStorage.Insert("11111", row)

	sessMgr := &mockSessionManager{IsAuthorizedSession: true}

	oauthServer := InitOAuth2Server(server, sessMgr)
	assert.Assert(t, oauthServer != nil)

	client, err := utils.ApiClient(netListener, router)
	assert.NilError(t, err, "Failed to get API client")
	clentCfg := &clientcredentials.Config{
		ClientID:     "11111",
		ClientSecret: "999999",
		TokenURL:     "http://localhost/auth/oauth2/token",
	}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, client)
	fmt.Printf("context: %v\n", ctx)
	token, err := clentCfg.Token(ctx)
	assert.NilError(t, err)
	assert.Assert(t, token.Valid())
	assert.Assert(t, token.AccessToken != "")
}
