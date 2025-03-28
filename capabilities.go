package marionette_client

import (
	"encoding/json"
	"strings"
)

// Capabilities is a struct that represents the capabilities of a browser - Webdriver capabilities.
type Capabilities struct {
	BrowserName               string         `json:"browserName"`
	BrowserVersion            string         `json:"browserVersion"`
	PlatformName              string         `json:"platformName"`
	PlatformVersion           string         `json:"platformVersion"`
	AcceptInsecureCerts       bool           `json:"acceptInsecureCerts"`
	PageLoadStrategy          string         `json:"PageLoadStrategy"`
	Proxy                     *Proxy         `json:"proxy"`
	SetWindowRect             bool           `json:"setWindowRect"`
	Timeouts                  *Timeouts      `json:"timeouts"`
	StrictFileInteractability bool           `json:"strictFileInteractability"`
	UnhandledPromptBehavior   string         `json:"unhandledPromptBehavior"`
	UserAgent                 string         `json:"userAgent"`
	MozExtensions             map[string]any `json:"moz:extension"`
}

func (c *Capabilities) UnmarshalJSON(data []byte) error {
	var d map[string]any
	err := json.Unmarshal(data, &d)
	if err != nil {
		return err
	}

	// loop over all keys of 'd' and set the corresponding field of 'c'
	for k, v := range d {
		// use reflection to find the field in 'c' that matches the key 'k'
		switch k {
		case "browserName":
			c.BrowserName = v.(string)
		case "browserVersion":
			c.BrowserVersion = v.(string)
		case "platformName":
			c.PlatformName = v.(string)
		case "platformVersion":
			c.PlatformVersion = v.(string)
		case "acceptInsecureCerts":
			c.AcceptInsecureCerts = v.(bool)
		case "PageLoadStrategy":
			c.PageLoadStrategy = v.(string)
		case "proxy":
			p, ok := v.(*Proxy)
			if !ok {
				// if unable to cast it to *Proxy, it's because it is an empty json object at the moment
				continue
			}
			c.Proxy = p
		case "setWindowRect":
			c.SetWindowRect = v.(bool)
		case "timeouts":
			tBytes, err := json.Marshal(v)
			if err != nil {
				continue
			}

			t := &Timeouts{}
			err = json.Unmarshal(tBytes, t)
			if err != nil {
				continue
			}

			c.Timeouts = t
		case "strictFileInteractability":
			c.StrictFileInteractability = v.(bool)
		case "unhandledPromptBehavior":
			c.UnhandledPromptBehavior = v.(string)
		case "userAgent":
			c.UserAgent = v.(string)
		default:
			// moz:extension is a special case, it's a map of strings to anything
			if strings.Index(k, "moz:") != 0 {
				continue
			}

			if c.MozExtensions == nil {
				c.MozExtensions = make(map[string]any)
			}

			c.MozExtensions[k] = v
		}
	}

	return nil
}

type Timeouts struct {
	Implicit int `json:"implicit,omitempty"`
	PageLoad int `json:"pageLoad,omitempty"`
	Script   int `json:"script,omitempty"`
}

// ProxyType is an enumeration of the types of proxies available.
type ProxyType string

const (
	// Direct connection - no proxy in use.
	Direct     ProxyType = "direct"
	Manual               = "manual"
	Autodetect           = "autodetect"
	System               = "system"
	PAC                  = "pac"
)

// Proxy defines configuration for proxies in the browser. see: https://www.w3.org/TR/webdriver2/#proxy
type Proxy struct {
	Type          ProxyType `json:"proxyType"`
	AutoconfigURL string    `json:"proxyAutoconfigUrl,omitempty"`

	// The following are used when Type is set to Manual.
	//
	// Note that in Firefox, connections to localhost are not proxied by default,
	// even if a proxy is set. This can be overridden via a preference setting.
	FTP           string   `json:"ftpProxy,omitempty"`
	HTTP          string   `json:"httpProxy,omitempty"`
	SSL           string   `json:"sslProxy,omitempty"`
	SOCKS         string   `json:"socksProxy,omitempty"`
	SOCKSVersion  int      `json:"socksVersion,omitempty"`
	SOCKSUsername string   `json:"socksUsername,omitempty"`
	SOCKSPassword string   `json:"socksPassword,omitempty"`
	NoProxy       []string `json:"noProxy,omitempty"`

	// The W3C draft spec includes port fields as well. According to the
	// specification, ports can also be included in the above addresses. However,
	// in the Geckodriver implementation, the ports must be specified by these
	// additional fields.
	HTTPPort  int `json:"httpProxyPort,omitempty"`
	SSLPort   int `json:"sslProxyPort,omitempty"`
	SocksPort int `json:"socksProxyPort,omitempty"`
}
