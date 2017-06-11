package marionette_client

type Capabilities struct {
	BrowserName                   string
	BrowserVersion                string
	PlatformName                  string
	PlatformVersion               string
	SpecificationLevel            string
	RaisesAccessibilityExceptions bool
	Rotatable                     bool
	AcceptSslCerts                bool
	TakesElementScreenshot        bool
	TakesScreenshot               bool
	Proxy                         interface{}
	Platform                      string
	XULappId                      string
	AppBuildId                    string
	Device                        string
	Version                       string
	Command_id                    uint32
}
