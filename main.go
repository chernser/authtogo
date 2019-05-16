package main

import (
	"./oauth2"
	"github.com/valyala/fasthttp"
	"fmt"
	"net/http"
)

type AuthServer struct {
	OAuth2Handler *oauth2.OAuth2Server
}

func (srv *AuthServer) HandleRequest(ctx *fasthttp.RequestCtx) {
	fmt.Printf("Path is %s\n", ctx.Path())
	switch (string(ctx.Path())) {
	case "/auth/oauth2/authorize":
		fallthrough
	case "/auth/oauth2/token":
		fmt.Printf("Going to handle oauth2\n")
		srv.OAuth2Handler.RequestHandler(ctx)
		break
	default: 
		fmt.Printf("Unknown authentication\n")
		ctx.Error("Unknown authentication mechanis", http.StatusNotFound)
	}
	fmt.Printf("Request handled\n")
}

func main() {
	fmt.Printf("Starting Auth Service\n")
	srvCtx := &AuthServer{
		OAuth2Handler: oauth2.InitOAuth2Server(),

	}
	fasthttp.ListenAndServe(":8088", srvCtx.HandleRequest)
}