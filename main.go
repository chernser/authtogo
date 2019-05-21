package main

import (
	// "net/http"

	"net/http"

	"github.com/google/uuid"

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
	router         *fasthttprouter.Router
	sessionManager auth.SessionManager
	OAuth2Server   *oauth2.OAuth2Server
	SamlSPServer   *saml.SamlSPServer
}

// RegisterRoute adds route and handler to internal router
func (aServer *AuthServerImpl) RegisterRoute(method string, path string, handler fasthttp.RequestHandler) {

	wrappedHandler := func(ctx *fasthttp.RequestCtx) {
		aServer.sessionManager.StartSession(ctx)
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

// Session represents security session
// If authentication is done for a session - IsAuthenticated is true,
// if session is ananymous - than false.
type Session struct {
	Id              string
	IsAuthenticated bool
}

type SessionManagerImpl struct {
	Sessions map[string]Session
}

func (sessions *SessionManagerImpl) StartSession(ctx interface{}) {
	context, isFastHttpCtx := ctx.(fasthttp.RequestCtx)
	if isFastHttpCtx {
		context.SetUserValue("__sessionId", uuid.New().String())

	} else {
		log.Warn().Msg("unsupported context type")
	}
}

func (sessions *SessionManagerImpl) InvalidateSession(id string) {
	delete(sessions.Sessions, id)
}

func (sessions *SessionManagerImpl) IsAuthenticated(ctx interface{}) bool {
	request, isHttpRequest := ctx.(http.Request)
	if isHttpRequest {
		authentication := request.Header.Get("Authentication")
		_, exists := sessions.Sessions[authentication]
		// TODO: add loading session from database
		return exists
	}
	return false
}

func main() {
	log.Info().Msg("Starting Auth Service")

	aServer := &AuthServerImpl{}

	aServer.init()
	aServer.readConfig()
	aServer.setupSessionManager()

	aServer.setupOAuth2Server()
	// aServer.setupSAMLSPServer()

	aServer.start()
}

func (aServer *AuthServerImpl) init() {
	aServer.router = fasthttprouter.New()
}

func (aServer *AuthServerImpl) readConfig() {

}

func (aServer *AuthServerImpl) setupSessionManager() {
	aServer.sessionManager = &SessionManagerImpl{}

}

func (aServer *AuthServerImpl) setupOAuth2Server() {
	aServer.OAuth2Server = oauth2.InitOAuth2Server(aServer, aServer.sessionManager)

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
