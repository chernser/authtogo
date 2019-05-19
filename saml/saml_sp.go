package saml

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"net/url"

	"github.com/crewjam/saml/samlsp"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"

	"../auth"
)

type SamlSPServerConfig struct {
	Key  string
	Cert string
}

type SamlSPServer struct {
	Impl           *samlsp.Middleware
	RequestHandler fasthttp.RequestHandler
}

func InitSamlSPServer(aServer auth.AuthServer, config *SamlSPServerConfig) *SamlSPServer {
	log.Info().Msg("Init SAML SP Server")

	keyPair, err := tls.LoadX509KeyPair(config.Cert, config.Key)
	if err != nil {
		panic(err) // TODO handle error
	}
	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
	if err != nil {
		panic(err) // TODO handle error
	}

	idpMetadataURL, err := url.Parse("https://www.testshib.org/metadata/testshib-providers.xml")
	if err != nil {
		panic(err) // TODO handle error
	}

	rootURL, err := url.Parse("http://localhost:8088")
	if err != nil {
		panic(err) // TODO handle error
	}

	samlSP, _ := samlsp.New(samlsp.Options{
		URL:               *rootURL,
		Key:               keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate:       keyPair.Leaf,
		IDPMetadataURL:    idpMetadataURL,
		AllowIDPInitiated: true,
	})

	server := &SamlSPServer{
		Impl:           samlSP,
		RequestHandler: fasthttpadaptor.NewFastHTTPHandlerFunc(samlSP.ServeHTTP),
	}

	return server
}
