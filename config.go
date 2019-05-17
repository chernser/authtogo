package main


type AuthServerConfiguration struct {
	OAuth2Enabled bool
	OAuth2AuthorizePath string
	OAuth2TokenPath string

	SamlSPCertificate string
	SamlSPKey string
}

func initConfiguration() (*AuthServerConfiguration) {
	config := &AuthServerConfiguration{}
	

	return config
}