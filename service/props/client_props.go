package props

import "time"

type ClientProps struct {
	KeepAlive           time.Duration `yaml:"keep-alive" envconfig:"KEEP_ALIVE"`
	IdleConnTimeout     time.Duration `yaml:"idle-timeout" envconfig:"IDLE_TIMEOUT"`
	TLSHandshakeTimeout time.Duration `yaml:"tls-handshake-timeout" envconfig:"TLS_HANDSHAKE_TIMEOUT"`
	Timeout             time.Duration `yaml:"timeout" envconfig:"TIMEOUT"`
	MaxIdleConnections  int           `yaml:"max-idle-connection" envconfig:"MAX_IDLE_CONNECTION"`
	ForceAttemptHTTP2   bool          `yaml:"attempt-http2" envconfig:"ATTEMPT_HTTP2"`
	OAuth               OAuth         `yaml:"oauth" envconfig:"OAUTH"`
	CertAuth            CertAuth      `yaml:"cert-auth"`
	BasicAuth           BasicAuth     `yaml:"basic-auth"`
}

// OAuth contains OAuth credentials details
type OAuth struct {
	Enabled     bool            `yaml:"enabled" envconfig:"ENABLED"`
	SslInsecure bool            `yaml:"ssl-insecure" envconfig:"SSL_INSECURE"`
	AuthCert    string          `yaml:"auth-cert" envconfig:"SSL_CERT"`
	Token       TokenProperties `yaml:"token" envconfig:"TOKEN"`
}

// CertAuth will provide client certificates for the authentication
type CertAuth struct {
	Enabled    bool   `yaml:"enabled" envconfig:"CERT_AUTH_ENABLED"`
	ClientCert string `yaml:"client-cert" envconfig:"CLIENT_CERT"`
	ClientKey  string `yaml:"client-key" envconfig:"CLIENT_KEY"`
}

// BasicAuth contains credentials for performing basic auth
type BasicAuth struct {
	ClientId     string `yaml:"clientId" envconfig:"CLIENT_ID"`
	ClientSecret string `yaml:"clientSecret" envconfig:"CLIENT_SECRET"`
}

// TokenProperties contains client application information and the server's endpoint URLs
type TokenProperties struct {
	// ClientID is the application's ID
	ClientID string `yaml:"clientID" envconfig:"CLIENT_ID"`

	// ClientSecret is the application's secret
	ClientSecret string `yaml:"clientSecret" envconfig:"CLIENT_SECRET"`

	// TokenURL is the resource server's token endpoint
	TokenURL string `yaml:"tokenURL" envconfig:"TOKEN_URL"`

	// Scopes specifies optional requested permissions
	Scopes []string `yaml:"scopes" envconfig:"SCOPES"`
}
