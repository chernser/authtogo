package sessions

import (
	"errors"

	"net/http"
	"time"

	"github.com/valyala/fasthttp"

	"github.com/chernser/authtogo/auth"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/ksuid"
)

const authCookieKey = "auth-session"

// SessionManagerImpl is implementation of session manager able to work with fasthttp context
type SessionManagerImpl struct {
	sessionsStorage auth.Storage
}

// StartSession generates new session id and registers appropriate internal structures
func (manager *SessionManagerImpl) StartSession(context interface{}) error {
	ctx, ok := context.(http.ResponseWriter)
	if !ok {
		log.Error().Msgf("Failed to start session")
		return errors.New("Failed to start session")
	}

	authToken, _ := ksuid.NewRandom()
	authCookie := &http.Cookie{Name: authCookieKey, Expires: time.Now().Add(time.Minute * 2),
		HttpOnly: true, Value: authToken.String(), Path: "/",
	}

	http.SetCookie(ctx, authCookie)
	log.Info().Msgf("Starting session with id %s", authToken.String())
	session := make(map[string]interface{})
	session["created_at"] = 0
	manager.sessionsStorage.Insert(authToken.String(), session)
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
	log.Debug().Msg("Is Authenticated")
	ctx, ok := context.(*http.Request)
	if !ok {
		log.Error().Msgf("Unsupported context")
		return false
	}

	authCookie, err := ctx.Cookie(authCookieKey)
	if err != nil {
		log.Error().Err(err)
		return false
	}

	if authCookie == nil {
		log.Debug().Msg("No authentication cookie")
		return false
	}

	log.Info().Msgf("Security cookie: %s", authCookie.Value)
	_, exists := manager.sessionsStorage.Get(string(authCookie.Value), []string{"created_at"})
	if !exists {
		return false
	}

	return true
}

// NewSessionManager - creates new session manager
func NewSessionManager(sessionStorage auth.Storage) auth.SessionManager {
	manager := &SessionManagerImpl{
		sessionsStorage: sessionStorage,
	}
	return manager
}
