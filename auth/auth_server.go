package auth

import (
	"github.com/valyala/fasthttp"
	//	"github.com/buaazp/fasthttprouter"
)

// AuthServer describes interface to underlying authentication server implementation
type AuthServer interface {
	//  Registers route for method and path to handler
	//  Used by authentication modules to register their routes
	RegisterRoute(Method string, Path string, Handler fasthttp.RequestHandler)
}

// SessionManager - manages sessions. Thats it.
type SessionManager interface {
	// StartSession generates new session id and registers appropriate internal structures
	StartSession(context interface{})

	// InvalidateSession removes session by id from all auth related stores
	InvalidateSession(id string)

	// IsAuthenticated returns true if context contains information about authenticated session
	IsAuthenticated(context interface{}) bool
}
