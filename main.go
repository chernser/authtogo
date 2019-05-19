package main

import (
	// "net/http"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"

	// "github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"./auth"
	"./oauth2"
	"./saml"
)

type AuthServerImpl struct {
	auth.AuthServer
	router       *fasthttprouter.Router
	OAuth2Server *oauth2.OAuth2Server
	SamlSPServer *saml.SamlSPServer
}

func (aServer *AuthServerImpl) RegisterRoute(method string, path string, handler fasthttp.RequestHandler) {

	log.Info().Msgf("RegisterRoute %s %s", method, path)
	switch method {
	case "POST":
		aServer.router.POST(path, handler)
	case "GET":
		aServer.router.GET(path, handler)
	default:
		log.Error().Msgf("Failed to register route %s %s ", method, path)
	}
}

func main() {
	log.Info().Msg("Starting Auth Service")

	aServer := &AuthServerImpl{}

	aServer.init()
	aServer.readConfig()

	aServer.setupOAuth2Server()
	// aServer.setupSAMLSPServer()

	aServer.start()
}

func (aServer *AuthServerImpl) init() {
	aServer.router = fasthttprouter.New()
}

func (aServer *AuthServerImpl) readConfig() {

}

func (aServer *AuthServerImpl) setupOAuth2Server() {
	aServer.OAuth2Server = oauth2.InitOAuth2Server(aServer)

}

func (aServer *AuthServerImpl) setupSAMLSPServer() {
	samlSpConfig := &saml.SamlSPServerConfig{
		Cert: "./myservice.cert",
		Key:  "./myservice.key",
	}

	aServer.SamlSPServer = saml.InitSamlSPServer(aServer, samlSpConfig)

}

func (aServer *AuthServerImpl) start() {
	fasthttp.ListenAndServe(":8088", aServer.router.Handler)
}
