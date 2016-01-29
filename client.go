package marionette_client

import (
	"encoding/json"
)

type Context string

const (
	CONTEXT_CHROME  Context = "chrome"
	CONTEXT_CONTENT Context = "content"
)

type session struct {
	SessionId string
}

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
}

type Client struct {
	session
	transport
}

func NewClient() *Client {
	return &Client{
		session{},
		transport{},
	}
}

func (c *Client) GetSessionID() string {
	return c.SessionId
}

func (c *Client) Connect(host string, port int) error {
	return c.transport.connect(host, port)
}

// Protocol commands
func (c *Client) GetMarionetteID() (*response, error) {
	response, err := c.transport.send("getMarionetteID", nil)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) SayHello() (*response, error) {
	response, err := c.transport.send("sayHello", nil)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) NewSession(sessionId string, cap *Capabilities) (*response, error) {
	data := map[string]interface{}{
		"sessionId":    sessionId,
		"capabilities": cap,
	}

	response, err := c.transport.send("newSession", data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(response.Value), &c)
	return response, nil
}

// Send the current session's capabilities to the client.
// Capabilities informs the client of which WebDriver features are
// supported by Firefox and Marionette.  They are immutable for the
// length of the session.
// The return value is an immutable map of string keys
// ("capabilities") to values, which may be of types boolean,
// numerical or string.

func (c *Client) GetSessionCapabilities() (*Capabilities, error) {
	buf, err := c.transport.send("getSessionCapabilities", nil)
	if err != nil {
		return nil, err
	}

	response := map[string]*Capabilities{"Capabilities": &Capabilities{}}
	err = json.Unmarshal([]byte(buf.Value), &response)
	if err != nil {
		return nil, err
	}

	cap, _ := response["capabilities"]
	return cap, nil
}

// Log message.  Accepts user defined log-level.
//
// param string value
//     Log message.
// param string level
//     Arbitrary log level.

func (c *Client) Log(message string, level string) (*response, error) {
	response, err := c.transport.send("log", map[string]string{"value": message, "level": level})
	if err != nil {
		return nil, err
	}

	return response, nil
}

//  Return all logged messages.
func (c *Client) GetLogs() (*response, error) {
	response, err := c.transport.send("getLogs", nil)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Sets the context of the subsequent commands to be either "chrome" or
// "content".
//
// param string value
//     Name of the context to be switched to.  Must be one of "chrome" or
//     "content".

func (c *Client) SetContext(value string) (*response, error) {
	response, err := c.transport.send("setContext", map[string]string{"value": value})
	if err != nil {
		return nil, err
	}

	return response, nil
}

//  Gets the context of the server, either "chrome" or "content".
func (c *Client) GetContext() (*response, error) {
	response, err := c.transport.send("getContext", nil)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) ExecuteScript(value string) (*response, error) {
	response, err := c.transport.send("executeScript", map[string]string{"value": value})
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Set the timeout for asynchronous script execution.
//
// param number ms
//     Time in milliseconds.

func (c *Client) SetScriptTimeout(milliseconds int) (*response, error) {
	return c.timeouts("script", milliseconds)
}

// Set timeout for searching for elements.
//
// param number ms
//     Search timeout in milliseconds.

func (c *Client) SetSearchTimeout(milliseconds int) (*response, error) {
	return c.timeouts("implicit", milliseconds)
}

// Set timeout for page loading.
//
// param number ms
//     Search timeout in milliseconds.

func (c *Client) SetPageTimeout(milliseconds int) (*response, error) {
	return c.timeouts("", milliseconds)
}

// Set timeout for page loading, searching, and scripts.
//
// param string type
//     Type of timeout.
// param number ms
//     Timeout in milliseconds.

func (c *Client) timeouts(typ string, milliseconds int) (*response, error) {
	if typ != "implicit" && typ != "script" {
		typ = ""
	}

	response, err := c.transport.send("timeouts", map[string]interface{}{"type": typ, "ms": milliseconds})
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Find elements using the indicated search strategy.
//
// param string using
//     Indicates which search method to use.
// param string value
//     Value the client is looking for.

//func (c *Client) FindElements(by string, value string) {
//	buf, err := c.transport.send("findElements", map[string]interface{}{})
//    if err != nil {
//        return "", err
//    }
//
//    return string(buf), nil
//}

func (c *Client) Close() (*response, error) {
	response, err := c.transport.send("close", nil)
	if err != nil {
		return nil, err
	}

	err = c.transport.close()
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) GetStatus() {

}
