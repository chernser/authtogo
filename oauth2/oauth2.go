package oauth2

import (
	"net/http"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"	
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
	"github.com/rs/zerolog/log"
)

type OAuth2Server struct {
	SrvImpl *server.Server
	RequestHandler fasthttp.RequestHandler 
}

func (srv *OAuth2Server) HandleOauth2Authorize(w http.ResponseWriter, r *http.Request) {	
	
	switch string(r.URL.Path) {
	case "/auth/oauth2/authorize":		
		log.Debug().Msg("Handling OAuth2 authorize request")
		err := srv.SrvImpl.HandleAuthorizeRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	case "/auth/oauth2/token":	
		log.Debug().Msg("Handling OAuth2 token request")
		srv.SrvImpl.HandleTokenRequest(w, r)
	}
}

func InitOAuth2Server() (*OAuth2Server) {
	log.Info().Msg("Init OAuth2 Server")
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
	server := OAuth2Server{}
	server.SrvImpl = srvImpl
	server.RequestHandler = fasthttpadaptor.NewFastHTTPHandlerFunc(server.HandleOauth2Authorize)

	return &server
}