package monitor

type HTTP struct {
	Base
	HTTPDetails
}

func (h HTTP) Type() string {
	return h.Base.Type()
}

func (h *HTTP) UnmarshalJSON(data []byte) error {
	panic("// FIXME: not implemented")
}

func (h HTTP) MarshalJSON() ([]byte, error) {
	panic("// FIXME: not implemented")
}

type HTTPDetails struct {
	URL                 string     `json:"url"`
	Timeout             int64      `json:"timeout"`
	ExpiryNotification  bool       `json:"expiryNotification"`
	IgnoreTLS           bool       `json:"ignoreTls"`
	MaxRedirects        int        `json:"maxredirects"`
	AcceptedStatusCodes []string   `json:"accepted_statuscodes"`
	ProxyID             int64      `json:"proxyID"`
	Method              string     `json:"method"`
	HTTPBodyEncoding    string     `json:"httpBodyEncoding"`
	Body                string     `json:"body"`
	Headers             string     `json:"headers"`
	AuthMethod          AuthMethod `json:"authMethod"`
	BasicAuthUser       string     `json:"basic_auth_user"`
	BasicAuthPass       string     `json:"basic_auth_pass"`
	AuthDomain          string     `json:"authDomain"`
	AuthWorkstation     string     `json:"authWorkstation"`
	TLSCert             string     `json:"tlsCert"`
	TLSKey              string     `json:"tlsKey"`
	TLSCa               string     `json:"tlsCa"`
	OAuthAuthMethod     string     `json:"oauth_auth_method"`
	OAuthTokenURL       string     `json:"oauth_token_url"`
	OAuthClientID       string     `json:"oauth_client_id"`
	OAuthClientSecret   string     `json:"oauth_client_secret"`
	OAuthScopes         string     `json:"oauth_scopes"`
}

func (h HTTPDetails) Type() string {
	return "http"
}

// AuthMethod represents the authentication method for monitors
type AuthMethod string

// Auth methods
const (
	AuthMethodNone     AuthMethod = ""
	AuthMethodBasic    AuthMethod = "basic"
	AuthMethodNTLM     AuthMethod = "ntlm"
	AuthMethodMTLS     AuthMethod = "mtls"
	AuthMethodOAtuh2CC AuthMethod = "oauth2-cc"
)
