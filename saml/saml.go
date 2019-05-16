package saml

import (
	"github.com/rs/zerolog/log"
)

type SamlSPServerConfig struct {

}

type SamlSPServer struct { 

}


type SamlIDPServerConfig struct {

}

type SamlIDPServer struct {

}



func InitSamlSPServer(config *SamlSPServerConfig) (*SamlSPServer) {
	log.Info().Msg("Init SAML SP Server")
	return nil
}


func InitSamlIDPServer(config *SamlIDPServerConfig) (*SamlIDPServer) {
	log.Info().Msg("Init SAML IDP Server")
	return nil
}