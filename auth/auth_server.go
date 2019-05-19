package auth

import (
	"github.com/valyala/fasthttp"
	//	"github.com/buaazp/fasthttprouter"
)

type AuthServer interface {
	RegisterRoute(string, string, fasthttp.RequestHandler)
}
