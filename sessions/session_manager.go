package sessions

import (
	"errors"

	"net/http"

	"github.com/valyala/fasthttp"

	"../auth"
)

const authCookieKey = "authS"

// SessionManagerImpl is implementation of session manager able to work with fasthttp context
type SessionManagerImpl struct {
	sessionsStorage auth.Storage
}

// StartSession generates new session id and registers appropriate internal structures
func (manager *SessionManagerImpl) StartSession(context interface{}) error {
	ctx, ok := context.(*fasthttp.RequestCtx)
	if !ok {
		return errors.New("Failed to start session")
	}

	authCookie := &fasthttp.Cookie{}
	authCookie.SetSecure(true)
	authCookie.SetHTTPOnly(true)
	authCookie.SetKey(authCookieKey)
	authCookie.SetMaxAge(180000)
	ctx.Response.Header.SetCookie(authCookie)
	ctx.Response.SetStatusCode(http.StatusNoContent)
	return nil
}

// InvalidateSession removes session by id from all auth related stores
func (manager *SessionManagerImpl) InvalidateSession(context interface{}) error {
	ctx, ok := context.(*fasthttp.RequestCtx)
	if !ok {
		return errors.New("Failed to invalidate session")
	}

	ctx.Response.Header.DelClientCookie(authCookieKey)
	ctx.Response.SetStatusCode(http.StatusNoContent)

	return nil
}

// IsAuthenticated returns true if context contains information about authenticated session
func (manager *SessionManagerImpl) IsAuthenticated(context interface{}) bool {
	ctx, ok := context.(*fasthttp.RequestCtx)
	if !ok {
		return false
	}

	authCookie := ctx.Request.Header.Cookie(authCookieKey)
	if authCookie != nil {
		return false
	}

	_, exists := manager.sessionsStorage.Get(string(authCookie), []string{"created_at"})
	if !exists {
		return false
	}

	return false
}

// NewSessionManager - creates new session manager
func NewSessionManager(sessionStorage auth.Storage) auth.SessionManager {
	manager := &SessionManagerImpl{
		sessionsStorage: sessionStorage,
	}
	return manager
}
