package oauth2

import (
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"

	"github.com/chernser/authtogo/auth"
)

// OAuth2ServerConfiguration - defines oauth2 module configuration
type OAuth2ServerConfiguration struct {
	TokenStore  auth.Storage
	ClientStore auth.Storage
}

// OAuth2Server describes authentication server for oauth2
type OAuth2Server struct {
	impl           *server.Server
	sessionManager auth.SessionManager
	tokenStore     oauth2.TokenStore
	clientStore    oauth2.ClientStore
}

// handleOAuth2Authorize - handles authentication.
func (srv *OAuth2Server) handleOauth2Authorize(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("Handling OAuth2 authorize request")

	if !srv.sessionManager.IsAuthenticated(r) {
		log.Info().Msg("Request from unauthenticated session")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	} else {
		err := srv.impl.HandleAuthorizeRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
}

func (srv *OAuth2Server) handleOauth2Token(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("Handling OAuth2 token request")
	srv.impl.HandleTokenRequest(w, r)
}

func (srv *OAuth2Server) handleOauth2RefreshToken(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("Handle refresh token request")
}

// InitOAuth2Server - Initializes OAuth2 module to authorize
// requires auth.AuthServer to register routes and auth.SessionManager to register sessions
func InitOAuth2Server(aServer auth.AuthServer, sessionManager auth.SessionManager, config *OAuth2ServerConfiguration) *OAuth2Server {
	log.Info().Msg("Init OAuth2 Server")

	oauthServer := &OAuth2Server{}
	oauthServer.sessionManager = sessionManager
	oauthServer.setupTokenStorage(config.TokenStore)
	oauthServer.setupSecretsStorage(config.ClientStore)
	oauthServer.setupImpl()

	aServer.RegisterRoute("POST", "/auth/oauth2/authorize", fasthttpadaptor.NewFastHTTPHandlerFunc(oauthServer.handleOauth2Authorize))
	aServer.RegisterRoute("POST", "/auth/oauth2/access_token", fasthttpadaptor.NewFastHTTPHandlerFunc(oauthServer.handleOauth2Token))

	return oauthServer
}

func (srv *OAuth2Server) setupImpl() {
	// General setup
	manager := manage.NewDefaultManager()
	// token memory store
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// client memory store
	// clientStore := store.NewClientStore()
	// clientStore.Set("000000", &models.Client{
	// 	ID:     "000000",
	// 	Secret: "999999",
	// 	Domain: "http://localhost",
	// })

	manager.MapClientStorage(srv.clientStore)

	srvImpl := server.NewDefaultServer(manager)
	srvImpl.SetAllowGetAccessRequest(true)
	srvImpl.SetClientInfoHandler(server.ClientFormHandler)
	srvImpl.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Error().Msgf("Internal Error: %s", err.Error())
		return
	})

	srvImpl.SetResponseErrorHandler(func(re *errors.Response) {
		log.Error().Msgf("Response Error: %s", re.Error.Error())
	})

	srv.impl = srvImpl
}

func (srv *OAuth2Server) setupTokenStorage(volatileStorage auth.Storage) {
	srv.tokenStore = &tokenStoreImpl{
		Storage: volatileStorage,
	}
}

func (srv *OAuth2Server) setupSecretsStorage(storage auth.Storage) {
	srv.clientStore = &clientStoreImpl{
		Storage: storage,
	}
}

// tokenStoreImpl is adapter for oauth2 token store
type tokenStoreImpl struct {
	Storage auth.Storage
}

func (store *tokenStoreImpl) Create(info oauth2.TokenInfo) error {

}

// delete the authorization code
func (store *tokenStoreImpl) RemoveByCode(code string) error {
	return nil
}

// use the access token to delete the token information
func (store *tokenStoreImpl) RemoveByAccess(access string) error {
	return nil
}

// use the refresh token to delete the token information
func (store *tokenStoreImpl) RemoveByRefresh(refresh string) error {
	return nil
}

// use the authorization code for token information data
func (store *tokenStoreImpl) GetByCode(code string) (oauth2.TokenInfo, error) {
	return nil, nil
}

// use the access token for token information data
func (store *tokenStoreImpl) GetByAccess(access string) (oauth2.TokenInfo, error) {
	return nil, nil
}

// use the refresh token for token information data
func (store *tokenStoreImpl) GetByRefresh(refresh string) (oauth2.TokenInfo, error) {
	return nil, nil
}

// clientStoreImpl is adapter between oauth2 store and auth store
type clientStoreImpl struct {
	Storage auth.Storage
}

func (store *clientStoreImpl) GetByID(id string) (oauth2.ClientInfo, error) {
	row, exists := store.Storage.Get(id, []string{})
	if !exists {
		log.Warn().Msgf("Failed to fetch client info for %s", id)
		return nil, errors.ErrInvalidClient
	}

	return &models.Client{
			ID:     id,
			Secret: row["secret"].(string),
			Domain: row["domain"].(string),
			UserID: row["userId"].(string),
		},
		nil
}
