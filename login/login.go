package login

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"hash"
	"net/http"
	"strings"

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

	if loginPage.sessionManager.IsAuthenticated(r) {
		http.Redirect(w, r, "/auth/", http.StatusTemporaryRedirect)
		return
	}

	var password = r.FormValue("password")
	var userID = r.FormValue("userId")

	if password == "" || userID == "" {
		http.Redirect(w, r, "/auth/login/form.html?error=400;msg=InvalidRequest", http.StatusTemporaryRedirect)
		return
	}

	// map[string]interface{}, bool
	userInfo, exists := loginPage.aServer.GetSecretsStorage().Get(userID, []string{"password_hash", "hash_func", "salt"})
	if exists {
		log.Info().Msgf("User with id '%s' exists", userID)
		var salt = ""
		if userInfo["salt"] != nil {
			salt = userInfo["salt"].(string)
		}
		match, err := doPasswordMatch(password, userInfo["password_hash"].(string), userInfo["hash_func"].(string), salt)
		if match && err == nil {
			loginPage.sessionManager.StartSession(w)
			log.Info().Msgf("Password match. Authenticated user %s", userID)
			http.Redirect(w, r, "/auth/", http.StatusTemporaryRedirect)
			return
		} else if err != nil {
			log.Error().Err(err).Stack().Msg("Failed to check hashes")
			http.Redirect(w, r, "/auth/login/form.html?error=401;msg=InvalidAuthentication", http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/auth/login/form.html?error=401;msg=InvalidAuthentication", http.StatusTemporaryRedirect)
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

func doPasswordMatch(input string, password string, algorithm string, salt string) (bool, error) {

	var hashFunc hash.Hash
	switch algorithm {
	case "SHA384":
		hashFunc = sha512.New384()
	case "SHA256":
		hashFunc = sha256.New()
	default:
		return false, errors.New("Unknown hash algorithm")
	}

	hashFunc.Write([]byte(input))
	hexStr := hex.EncodeToString(hashFunc.Sum(nil))
	log.Info().Msgf("Hashed input(%s) %s : %s ", input, hexStr, password)

	return hexStr == strings.TrimSpace(password), nil
}
