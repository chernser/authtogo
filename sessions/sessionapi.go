package sessions

import (
	"encoding/json"

	"mime"
	"net/http"

	"../auth"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

type SessionAPI struct {
	sessionManager  auth.SessionManager
	secretStorage   auth.Storage
	secretValidator auth.SecretsValidator
}

func (api *SessionAPI) handleLoginRequest(ctx *fasthttp.RequestCtx) {
	contentType := string(ctx.Request.Header.ContentType())
	log.Debug().Msgf("This is content type: %s", contentType)

	mediaType, _, _ := mime.ParseMediaType(contentType)
	log.Debug().Msgf("This is parsed content type: %s", mediaType)

	var login string
	var secret []byte
	switch mediaType {
	case "application/json":
		var loginRequest map[string]string
		json.Unmarshal(ctx.PostBody(), &loginRequest)
		log.Debug().Msgf("This is login request: %s", loginRequest)
		login = loginRequest["login"]
		secret = []byte(loginRequest["secret"])
	case "application/x-www-form-urlencoded":
		fallthrough
	case "multipart/form-data":
		login = string(ctx.FormValue("login"))
		secret = ctx.FormValue("secret")
		log.Debug().Msgf("This is urlencoded request: %s:%s", login, secret)

	default:
		log.Warn().Msgf("Request with unsupported content type %s", contentType)
		errInvalidRequest(ctx)
		return
	}

	isValid, err := api.secretValidator.IsValidSecret("password", login, secret)

	if err != nil {
		log.Error().Msgf("Error while validating secret: %s", err)
		errInvalidRequest(ctx)
		return
	}

	if !isValid {
		errInvalidSecret(ctx)
		return
	}

	err = api.sessionManager.StartSession(ctx)
	if err != nil {
		ctx.SetContentType("text/plain")
		ctx.SetBody([]byte("Server Error"))
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		return
	}
}

func errInvalidRequest(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/plain")
	ctx.SetBody([]byte("Invalid Request"))
	ctx.Response.SetStatusCode(http.StatusBadRequest)
}

func errInvalidSecret(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/plain")
	ctx.SetBody([]byte("Authentication Failed"))
	ctx.Response.SetStatusCode(http.StatusUnauthorized)
}

func (api *SessionAPI) handleLogoutRequest(ctx *fasthttp.RequestCtx) {
	api.sessionManager.InvalidateSession(ctx)
}

// SetupSessionAPI - setup session API handlers
func SetupSessionAPI(aServer auth.AuthServer, sessionManager auth.SessionManager, secretValidator auth.SecretsValidator) *SessionAPI {

	api := &SessionAPI{sessionManager: sessionManager,
		secretStorage:   aServer.GetSecretsStorage(),
		secretValidator: secretValidator}

	aServer.RegisterRoute("POST", "/auth/sessions", api.handleLoginRequest)
	aServer.RegisterRoute("DELETE", "/auth/sessions/current", api.handleLogoutRequest)
	aServer.RegisterRoute("POST", "/auth/sessions/logout", api.handleLogoutRequest)
	aServer.RegisterRoute("GET", "/auth/sessions/logout", api.handleLogoutRequest)

	return api
}
