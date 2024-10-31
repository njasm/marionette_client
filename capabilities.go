package marionette_client

type Capabilities struct {
	BrowserName                   string
	BrowserVersion                string
	PlatformName                  string
	PlatformVersion               string
	SpecificationLevel            uint
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

/*
{"capabilities":
  {
	"browserName":"firefox",
	"browserVersion":"53.0.3",
	"platformName":"linux",
	"platformVersion":"4.8.12-040812-generic",
	"pageLoadStrategy":"normal",
	"acceptInsecureCerts":false,
	"timeouts":{
		"implicit":0,
		"pageLoad":300000,
		"script":30000
	},
	"rotatable":false,
	"specificationLevel":0,
	"moz:processID":2004,
	"moz:profile":"/home/travis/.mozilla/firefox/594k4686.default",
	"moz:accessibilityChecks":false
  }
}**/
