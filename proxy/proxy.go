package proxy

import (
	"encoding/json"
	"time"
)

// Proxy represents a proxy server configuration in Uptime Kuma.
// Proxies can be used by monitors to route requests through HTTP/HTTPS/SOCKS proxy servers.
type Proxy struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"userId"`
	Protocol    string    `json:"protocol"`
	Host        string    `json:"host"`
	Port        int       `json:"port"`
	Auth        bool      `json:"auth"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	Active      bool      `json:"active"`
	Default     bool      `json:"default"`
	CreatedDate time.Time `json:"-"`
}

func (p Proxy) GetID() int64 {
	return p.ID
}

func (p Proxy) String() string {
	return formatProxy(p)
}

func (p *Proxy) UnmarshalJSON(data []byte) error {
	aux := &struct {
		ID          int64  `json:"id"`
		UserID      int64  `json:"userId"`
		Protocol    string `json:"protocol"`
		Host        string `json:"host"`
		Port        int    `json:"port"`
		Auth        int    `json:"auth"`
		Username    string `json:"username"`
		Password    string `json:"password"`
		Active      int    `json:"active"`
		Default     int    `json:"default"`
		CreatedDate string `json:"createdDate"`
	}{}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	p.ID = aux.ID
	p.UserID = aux.UserID
	p.Protocol = aux.Protocol
	p.Host = aux.Host
	p.Port = aux.Port
	p.Auth = aux.Auth != 0
	p.Username = aux.Username
	p.Password = aux.Password
	p.Active = aux.Active != 0
	p.Default = aux.Default != 0

	if aux.CreatedDate != "" {
		// Try multiple time formats
		formats := []string{
			time.RFC3339,
			"2006-01-02 15:04:05",
			"2006-01-02T15:04:05",
		}

		var parseErr error
		for _, format := range formats {
			t, err := time.Parse(format, aux.CreatedDate)
			if err == nil {
				p.CreatedDate = t
				return nil
			}

			parseErr = err
		}

		return parseErr
	}

	return nil
}

func (p Proxy) MarshalJSON() ([]byte, error) {
	auth := 0
	if p.Auth {
		auth = 1
	}

	active := 0
	if p.Active {
		active = 1
	}

	defaultVal := 0
	if p.Default {
		defaultVal = 1
	}

	createdDate := ""
	if !p.CreatedDate.IsZero() {
		createdDate = p.CreatedDate.Format(time.RFC3339)
	}

	return json.Marshal(&struct {
		ID          int64  `json:"id"`
		UserID      int64  `json:"userId"`
		Protocol    string `json:"protocol"`
		Host        string `json:"host"`
		Port        int    `json:"port"`
		Auth        int    `json:"auth"`
		Username    string `json:"username"`
		Password    string `json:"password"`
		Active      int    `json:"active"`
		Default     int    `json:"default"`
		CreatedDate string `json:"createdDate"`
	}{
		ID:          p.ID,
		UserID:      p.UserID,
		Protocol:    p.Protocol,
		Host:        p.Host,
		Port:        p.Port,
		Auth:        auth,
		Username:    p.Username,
		Password:    p.Password,
		Active:      active,
		Default:     defaultVal,
		CreatedDate: createdDate,
	})
}

// Config represents the configuration for creating or updating a proxy.
// It includes an optional ID field for updates and an ApplyExisting flag
// to apply the proxy to all existing monitors upon creation.
type Config struct {
	ID            int64  `json:"id,omitempty"`
	Protocol      string `json:"protocol"`
	Host          string `json:"host"`
	Port          int    `json:"port"`
	Auth          bool   `json:"auth"`
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	Active        bool   `json:"active"`
	Default       bool   `json:"default"`
	ApplyExisting bool   `json:"applyExisting"`
}
