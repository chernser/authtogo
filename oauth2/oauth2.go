package oauth2

import (
	"net/http"

	//	"github.com/valyala/fasthttp"
	//	"github.com/valyala/fasthttp/fasthttpadaptor"

	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"

	"../auth"
)

// OAuth2Server describes authentication server for oauth2
type OAuth2Server struct {
	impl           *server.Server
	sessionManager auth.SessionManager
	tokenStore     oauth2.TokenStore
	clientStore    oauth2.ClientStore
}

func (srv *OAuth2Server) handleOauth2Authorize(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("Handling OAuth2 authorize request")

	if srv.sessionManager.IsAuthenticated(r) {
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

// InitOAuth2Server - Initializes OAuth2 module to authorize
// requires auth.AuthServer to register routes and auth.SessionManager to register sessions
func InitOAuth2Server(aServer auth.AuthServer, sessionManager auth.SessionManager) *OAuth2Server {
	log.Info().Msg("Init OAuth2 Server")

	oauthServer := &OAuth2Server{}
	oauthServer.setupTokenStorage(aServer.GetVolatileStorage())
	oauthServer.setupSecretsStorage(aServer.GetSecretsStorage())
	oauthServer.setupImpl()

	aServer.RegisterRoute("POST", "/auth/oauth2/authorize", fasthttpadaptor.NewFastHTTPHandlerFunc(oauthServer.handleOauth2Authorize))
	aServer.RegisterRoute("POST", "/auth/oauth2/token", fasthttpadaptor.NewFastHTTPHandlerFunc(oauthServer.handleOauth2Token))

	return oauthServer
}

func (srv *OAuth2Server) setupImpl() {
	// General setup
	manager := manage.NewDefaultManager()
	// token memory store
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// client memory store
	clientStore := store.NewClientStore()
	clientStore.Set("000000", &models.Client{
		ID:     "000000",
		Secret: "999999",
		Domain: "http://localhost",
	})
	manager.MapClientStorage(clientStore)

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
	// srv.tokenStore = &tokenStoreImpl{
	// 	Storage: volatileStorage,
	// }
}

func (srv *OAuth2Server) setupSecretsStorage(storage auth.Storage) {
	srv.clientStore = &clientStoreImpl{
		Storage: storage,
	}
}

type tokenStoreImpl struct {
	Storage auth.Storage
}

// clientStoreImpl is adapter between oauth2 store and auth store
type clientStoreImpl struct {
	Storage auth.Storage
}

func (store *clientStoreImpl) GetByID(id string) (oauth2.ClientInfo, error) {

	// info := &oauth2.ClientInfo{}

	return nil, nil
}
