package main

import (	
	"net/http"

	"github.com/valyala/fasthttp"
	_ "github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	
	"./oauth2"
)

type AuthServer struct {
	OAuth2Handler *oauth2.OAuth2Server
}

func (srv *AuthServer) HandleRequest(ctx *fasthttp.RequestCtx) {
	log.Debug().Msgf("Handling request: %s", ctx.Path())
	switch (string(ctx.Path())) {
	case "/auth/oauth2/authorize":
		fallthrough
	case "/auth/oauth2/token":
		srv.OAuth2Handler.RequestHandler(ctx)
		break
	default: 
		log.Error().Msg("Unknown auth request")
		ctx.Error("Unknown authentication mechanis", http.StatusNotFound)
	}
	log.Debug().Msg("Request handled")
}

func main() {
	log.Info().Msg("Starting Auth Service")
	
	initConfiguration()
	srvCtx := &AuthServer{
		OAuth2Handler: oauth2.InitOAuth2Server(),

	}
	fasthttp.ListenAndServe(":8088", srvCtx.HandleRequest)
}