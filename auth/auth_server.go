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
	StartSession(context interface{}) error

	// InvalidateSession removes session by id from all auth related stores
	InvalidateSession(context interface{}) error

	// IsAuthenticated returns true if context contains information about authenticated session
	IsAuthenticated(context interface{}) bool
}

// Storage defines behaviour of authentication storage used for tokens and secrets
// Server uses one storage for tokens what is very volatile information,
// and for accessing secrets what is more stable information
// In general both storages are similar in way accessing data.
type Storage interface {

	// Get returns mapped values for row with rowId
	Get(rowID string, fields []string) (map[string]string, bool)

	// Insert records row under rowID and with provided values. Returns true if success
	Insert(rowID string, values map[string]string) bool

	// Update records new values in the row.
	Update(rowID string, values map[string]string) bool

	// Delete removes record from store.
	Delete(rowID string) bool
}

// SecretsValidator checks secret for validity
// Commonly it is used to validate passwords and certificates
type SecretsValidator interface {
	// IsValidSecret checks what secret belongs to owner
	IsValidSecret(sType string, owner string, secret []byte) (bool, error)
}
