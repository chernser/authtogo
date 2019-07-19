package main

import (
	// "net/http"

	"fmt"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"

	"github.com/rs/zerolog/log"

	"github.com/spf13/viper"

	"github.com/chernser/authtogo/auth"
	"github.com/chernser/authtogo/datastore"
	"github.com/chernser/authtogo/oauth2"
	"github.com/chernser/authtogo/saml"
	"github.com/chernser/authtogo/sessions"
)

// Authentication server joining structure
type AuthServerImpl struct {
	auth.AuthServer
	router          *fasthttprouter.Router
	sessionManager  auth.SessionManager
	sessionAPI      *sessions.SessionAPI
	OAuth2Server    *oauth2.OAuth2Server
	SamlSPServer    *saml.SamlSPServer
	volatileStorage auth.Storage
	secretsStorage  auth.Storage
}

// RegisterRoute adds route and handler to internal router
func (aServer *AuthServerImpl) RegisterRoute(method string, path string, handler fasthttp.RequestHandler) {

	wrappedHandler := func(ctx *fasthttp.RequestCtx) {
		// aServer.sessionManager.StartSession(ctx)
		handler(ctx)
	}

	log.Info().Msgf("RegisterRoute %s %s", method, path)
	switch method {
	case "POST":
		aServer.router.POST(path, wrappedHandler)
	case "GET":
		aServer.router.GET(path, wrappedHandler)
	default:
		log.Error().Msgf("Failed to register route %s %s ", method, path)
	}
}

func main() {
	log.Info().Msg("Starting Auth Service")

	aServer := &AuthServerImpl{}

	aServer.init()
	aServer.readConfig()
	aServer.setupSessionManager()

	aServer.setupOAuth2Server()
	aServer.setupSAMLSPServer()

	aServer.start()
}

func (aServer *AuthServerImpl) init() {
	aServer.router = fasthttprouter.New()
}

func (aServer *AuthServerImpl) readConfig() {
	viper.SetConfigFile("authserver.yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("failed to read configuration: %s", err))
	}

	aServer.volatileStorage = datastore.CreateVolatileDataStore(viper.GetViper())
	aServer.secretsStorage = datastore.CreateSecretsDataStore(viper.GetViper())
}

func (aServer *AuthServerImpl) setupSessionManager() {
	aServer.sessionManager = sessions.NewSessionManager(aServer.GetVolatileStorage())
	aServer.sessionAPI = sessions.SetupSessionAPI(aServer, aServer.sessionManager, nil)
}

func (aServer *AuthServerImpl) setupOAuth2Server() {
	viper.SetDefault("oauth2.enabled", true)
	if viper.GetBool("oauth2.enabled") {
		aServer.OAuth2Server = oauth2.InitOAuth2Server(aServer, aServer.sessionManager)
	} else {
		log.Info().Msg("OAuth2 is disabled")
	}
}

func (aServer *AuthServerImpl) setupSAMLSPServer() {
	viper.SetDefault("saml.enabled", false)
	if viper.GetBool("saml.enabled") {
		viper.SetDefault("saml.cert.path", "./saml_server.cert")
		viper.SetDefault("saml.key.path", "./saml_server.key")
		samlSpConfig := &saml.SamlSPServerConfig{
			Cert: viper.GetString("saml.cert.path"),
			Key:  viper.GetString("saml.key.path"),
		}

		aServer.SamlSPServer = saml.InitSamlSPServer(aServer, samlSpConfig)
	} else {
		log.Info().Msg("SAML is disabled")
	}
}

func (aServer *AuthServerImpl) start() {
	viper.SetDefault("http_server.port", ":8088")
	fasthttp.ListenAndServe(viper.GetString("http_server.port"), aServer.router.Handler)
}

// GetVolatileStorage returns configured volatile storage suitable for storing runtime
// data such as tokens
func (aServer *AuthServerImpl) GetVolatileStorage() auth.Storage {
	return aServer.volatileStorage
}

// GetSecretsStorage returns storage where user secret (like keys, password hashes etc)
// are stored. This storage usually composite one what solves integration problems
func (aServer *AuthServerImpl) GetSecretsStorage() auth.Storage {
	return aServer.secretsStorage
}
