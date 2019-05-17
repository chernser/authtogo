package auth

import (
	"github.com/valyala/fasthttp"
	"github.com/buaazp/fasthttprouter"
)

type AuthServer struct {
	router *fasthttprouter.Router
}

func (srv *AuthServer) RegisterRoute(method string, path string, handler fasthttp.RequestHandler) {

	switch method {
	case "POST":
		srv.router.POST(path, handler)
	case "GET":
		srv.router.GET(path, handler)
	}	
}

func New() (*AuthServer) {

	server := &AuthServer{
		router: fasthttprouter.New(),
	}

	return server
}