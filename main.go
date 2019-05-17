package main

import (	
	// "net/http"

	"github.com/valyala/fasthttp"
	"github.com/buaazp/fasthttprouter"
	// "github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	
	"./oauth2"
	"./saml"
)

type AuthServer struct {
	OAuth2Handler *oauth2.OAuth2Server
	OAuth2AuthorizePath string
	OAuth2TokenPath string
	SamlSPServer *saml.SamlSPServer
}

func (srv *AuthServer) handleOAuthAuthorize(ctx *fasthttp.RequestCtx) {
	log.Debug().Msgf("Handling request: %s", ctx.Path())
	srv.OAuth2Handler.RequestHandler(ctx)
}

func (srv *AuthServer) handleOAuthTokenRequest(ctx *fasthttp.RequestCtx) {
	log.Debug().Msgf("Handling request: %s", ctx.Path())
	srv.OAuth2Handler.RequestHandler(ctx)
}

func main() {
	log.Info().Msg("Starting Auth Service")
	
	samlSpConfig := &saml.SamlSPServerConfig{
		Cert: "./myservice.cert",
		Key: "./myservice.key",
	}

	OAuth2AuthorizePath := "/auth/oauth2/authorize"
	OAuth2TokenPath := "/auth/oauth2/token"
	
	srvCtx := &AuthServer{		
		OAuth2Handler: oauth2.InitOAuth2Server(),
		SamlSPServer: saml.InitSamlSPServer(samlSpConfig),
	}

	router := fasthttprouter.New()
	router.POST(OAuth2AuthorizePath, srvCtx.handleOAuthAuthorize)
	router.POST(OAuth2TokenPath, srvCtx.handleOAuthTokenRequest)

	fasthttp.ListenAndServe(":8088", router.Handler)
}