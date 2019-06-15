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

	// GetVolatileStorage returns volatile storage for tokens and runtime information
	GetVolatileStorage() Storage

	// GetSecretsStorage returns storage of secrets
	GetSecretsStorage() Storage
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

// Storage defines behaviour of authentication storage used for tokens and secrets
// Server uses one storage for tokens what is very volatile information,
// and for accessing secrets what is more stable information
// In general both storages are similar in way accessing data.
type Storage interface {

	// Return map of fields representing the row. second tuple argument is true
	// if value exists, otherwise false
	GetFieldsOfRow(id string, fields []string) (map[string]string, bool)

	// Sets row information
	SetFieldsOfRow(id string, fields map[string]string)
}
