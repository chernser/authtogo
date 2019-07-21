package login

import (
	"net/http"

	"github.com/valyala/fasthttp"
	//	"github.com/valyala/fasthttp/fasthttpadaptor"

	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp/fasthttpadaptor"

	"github.com/chernser/authtogo/auth"
)

type LoginPage struct {
	aServer        auth.AuthServer
	sessionManager auth.SessionManager
}

func (loginPage *LoginPage) handleWebFormLogin(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("Handling WebFormLogin")

	var password = r.FormValue("password")
	var userId = r.FormValue("userId")

	if password == "" || userId == "" {
		http.Redirect(w, r, "/auth/login/form.html?error=400;msg=InvalidRequest", 302)
		return
	}

}

func InitLogin(aServer auth.AuthServer, sessionManager auth.SessionManager) error {
	log.Info().Msg("InitLogin")
	var loginPage = &LoginPage{aServer, sessionManager}
	aServer.RegisterRoute("POST", "/auth/login/form",
		fasthttpadaptor.NewFastHTTPHandlerFunc(loginPage.handleWebFormLogin))

	loginPageHTMLHandler := &fasthttp.FS{
		Root: "./static_assets/",
		PathRewrite: func(ctx *fasthttp.RequestCtx) []byte {
			return []byte("/login_form.html")
		},
	}

	aServer.RegisterRoute("GET", "/auth/login/form.html", loginPageHTMLHandler.NewRequestHandler())
	return nil
}
