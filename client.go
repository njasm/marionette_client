package marionette_client

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Context string

const (
	CONTEXT_CHROME         Context = "chrome"
	CONTEXT_CONTENT        Context = "content"
	MARIONETTE_PROTOCOL_V2         = 2
	MARIONETTE_PROTOCOL_V3         = 3
)

type session struct {
	SessionId string
}

type Client struct {
	session
	transport Transporter
}

func NewClient() *Client {
	return &Client{
		session{},
		&MarionetteTransport{},
	}
}

func (c *Client) Transport(t Transporter) {
	c.transport = t
}

func (c *Client) GetSessionID() string {
	return c.SessionId
}

func (c *Client) Connect(host string, port int) error {
	return c.transport.Connect(host, port)
}

// Protocol commands
// NOT A COMMAND
//func (c *Client) GetMarionetteID() (*response, error) {
//	response, err := c.transport.Send("getMarionetteID", nil)
//	if err != nil {
//		return nil, err
//	}
//
//	return response, nil
//}

//func (c *Client) SayHello() (*response, error) {
//	response, err := c.transport.Send("sayHello", nil)
//	if err != nil {
//		return nil, err
//	}
//
//	return response, nil
//}

func (c *Client) NewSession(sessionId string, cap *Capabilities) (*response, error) {
	data := map[string]interface{}{
		"sessionId":    sessionId,
		"capabilities": cap,
	}

	response, err := c.transport.Send("newSession", data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(response.Value), &c)
	if err != nil {
		return nil, err
	}

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
	buf, err := c.transport.Send("getSessionCapabilities", nil)
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
	response, err := c.transport.Send("log", map[string]string{"value": message, "level": level})
	if err != nil {
		return nil, err
	}

	return response, nil
}

//  Return all logged messages.
func (c *Client) GetLogs() (*response, error) {
	response, err := c.transport.Send("getLogs", nil)
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
func (c *Client) SetContext(value Context) (*response, error) {
	response, err := c.transport.Send("setContext", map[string]string{"value": string(value)})
	if err != nil {
		return nil, err
	}

	return response, nil
}

//  Gets the context of the server, either "chrome" or "content".
func (c *Client) GetContext() (*response, error) {
	response, err := c.transport.Send("getContext", nil)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) ExecuteScript(script string, args []interface{}, timeout uint, newSandbox bool) (*response, error) {
	parameters := map[string]interface{}{}
	parameters["scriptTimeout"] = timeout
	parameters["script"] = script
	parameters["args"] = args

	parameters["newSandbox"] = newSandbox

	response, err := c.transport.Send("executeScript", parameters)
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
	return timeouts(&c.transport, "script", milliseconds)
}

// Set timeout for searching for elements.
//
// param number ms
//     Search timeout in milliseconds.
func (c *Client) SetSearchTimeout(milliseconds int) (*response, error) {
	return timeouts(&c.transport, "implicit", milliseconds)
}

// Set timeout for page loading.
//
// param number ms
//     Search timeout in milliseconds.
func (c *Client) SetPageTimeout(milliseconds int) (*response, error) {
	return timeouts(&c.transport, "", milliseconds)
}

// Set timeout for page loading, searching, and scripts.
//
// param string type
//     Type of timeout.
// param number ms
//     Timeout in milliseconds.
func timeouts(transport *Transporter, typ string, milliseconds int) (*response, error) {
	if typ != "implicit" && typ != "script" {
		typ = ""
	}

	response, err := (*transport).Send("timeouts", map[string]interface{}{"type": typ, "ms": milliseconds})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) WindowHandles() ([]string, error) {
	r, err := c.transport.Send("getWindowHandles", nil)
	if err != nil {
		return nil, errors.New(r.ResponseError.Message)
	}

	var d []string
	err = json.Unmarshal([]byte(r.Value), &d)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (c *Client) SwitchToWindow(name string) error {
	r, err := c.transport.Send("switchToWindow", map[string]interface{}{"name": name})
	if err != nil {
		return errors.New(r.ResponseError.Message)
	}

	return nil
}

func (c *Client) CloseWindow() (*response, error) {
	r, err := c.transport.Send("close", nil)
	if err != nil {
		return nil, errors.New(r.ResponseError.Message)
	}

	return r, nil
}

////////////////
// NAVIGATION //
////////////////

func (c *Client) Get(url string) (*response, error) {
	r, err := c.transport.Send("get", map[string]string{"url": url})
	if err != nil {
		return nil, err
	}

	return r, nil
}

// get title
func (c *Client) GetTitle() (string, error) {
	r, err := c.transport.Send("getTitle", map[string]string{})
	if err != nil {
		return "", err
	}

	var d = map[string]string{}
	err = json.Unmarshal([]byte(r.Value), &d)
	if err != nil {
		return "", err
	}

	return d["value"], nil
}

// get current url
func (c *Client) CurrentUrl() error {
	r, err := c.transport.Send("goBack", nil)
	if err != nil {
		return errors.New(r.ResponseError.Message)
	}

	return nil
}

// refresh
func (c *Client) Refresh() error {
	r, err := c.transport.Send("refresh", nil)
	if err != nil {
		return errors.New(r.ResponseError.Message)
	}

	return nil
}

// back
func (c *Client) Back() error {
	r, err := c.transport.Send("goBack", nil)
	if err != nil {
		return errors.New(r.ResponseError.Message)
	}

	return nil
}

// forward
func (c *Client) Forward() error {
	r, err := c.transport.Send("goForward", nil)
	if err != nil {
		return errors.New(r.ResponseError.Message)
	}

	return nil
}

//////////////////
// WEB ELEMENTS //
//////////////////

func isElementEnabled(c *Client, id string) bool {
	r, err := c.transport.Send("isElementEnabled", map[string]interface{}{"id": id})
	if err != nil {
		return false
	}

	return strings.Contains(r.Value, "\"value\":true")
}

func isElementSelected(c *Client, id string) bool {
	r, err := c.transport.Send("isElementSelected", map[string]interface{}{"id": id})
	if err != nil {
		return false
	}

	return strings.Contains(r.Value, "\"value\":true")
}

func isElementDisplayed(c *Client, id string) bool {
	r, err := c.transport.Send("isElementDisplayed", map[string]interface{}{"id": id})
	if err != nil {
		return false
	}

	return strings.Contains(r.Value, "\"value\":true")
}

func getElementTagName(c *Client, id string) string {
	r, err := c.transport.Send("getElementTagName", map[string]interface{}{"id": id})
	if err != nil {
		return ""
	}

	var d = map[string]string{}
	json.Unmarshal([]byte(r.Value), &d)

	return d["value"]
}

func getElementText(c *Client, id string) string {
	r, err := c.transport.Send("getElementText", map[string]interface{}{"id": id})
	if err != nil {
		return ""
	}

	var d = map[string]string{}
	json.Unmarshal([]byte(r.Value), &d)

	return d["value"]
}

func getElementAttribute(c *Client, id string, name string) string {
	r, err := c.transport.Send("getElementAttribute", map[string]interface{}{"id": id, "name": name})
	if err != nil {
		return ""
	}

	var d = map[string]string{}
	json.Unmarshal([]byte(r.Value), &d)

	return d["value"]
}

func getElementCssPropertyValue(c *Client, id string, property string) string {
	r, err := c.transport.Send("getElementValueOfCssProperty", map[string]interface{}{"id": id, "propertyName": property})
	if err != nil {
		return ""
	}

	var d = map[string]string{}
	json.Unmarshal([]byte(r.Value), &d)

	return d["value"]
}

func getElementRect(c *Client, id string) map[string]interface{} {
	r, err := c.transport.Send("getElementRect", map[string]interface{}{"id": id})
	if err != nil {
		return nil
	}

	var d = map[string]interface{}{}
	json.Unmarshal([]byte(r.Value), &d)

	return d
}

func clickElement(c *Client, id string) {
	r, err := c.transport.Send("clickElement", map[string]interface{}{"id": id})
	if err != nil {
		return
	}

	var d = map[string]interface{}{}
	json.Unmarshal([]byte(r.Value), &d)

	//return d
}

func sendKeysToElement(c *Client, id string, keys string) {
	slice := make([]string, 0)
	for _, v := range keys {
		slice = append(slice, fmt.Sprintf("%c", v))
	}

	r, err := c.transport.Send("sendKeysToElement", map[string]interface{}{"id": id, "value": slice})
	if err != nil {
		return
	}

	var d = map[string]interface{}{}
	json.Unmarshal([]byte(r.Value), &d)

	//return d
}

func clearElement(c *Client, id string) {
	r, err := c.transport.Send("clearElement", map[string]interface{}{"id": id})
	if err != nil {
		return
	}

	var d = map[string]interface{}{}
	json.Unmarshal([]byte(r.Value), &d)

	//return d
}

// Find elements using the indicated search strategy.
//
// param string using
//     Indicates which search method to use.
// param string value
//     Value the client is looking for.
func (c *Client) FindElements(by string, value string, startNode *string) ([]*webElement, error) {
	return findElements(c, by, value, startNode)
}

func findElements(c *Client, by string, value string, startNode *string) ([]*webElement, error) {
	var params map[string]interface{}
	if startNode == nil || *startNode == "" {
		params = map[string]interface{}{"using": by, "value": value}
	} else {
		params = map[string]interface{}{"using": by, "value": value, "element": *startNode}
	}

	response, err := c.transport.Send("findElements", params)
	if err != nil {
		return nil, err
	}

	var d []map[string]string
	err = json.Unmarshal([]byte(response.Value), &d)
	if err != nil {
		return nil, err
	}

	var e []*webElement
	for _, v := range d {
		e = append(e, &webElement{c: c, id: v["element-6066-11e4-a52e-4f735466cecf"]})
	}

	return e, nil

	//return string(buf), nil
}

// Find an element using the indicated search strategy.
//
// @param {string} using
//     Indicates which search method to use.
// @param {string} value
//     Value the client is looking for.
func (c *Client) FindElement(by string, value string, startNode *string) (*webElement, error) {
	return findElement(c, by, value, startNode)
}

func findElement(c *Client, by string, value string, startNode *string) (*webElement, error) {
	var params map[string]interface{}
	if startNode == nil || *startNode == "" {
		params = map[string]interface{}{"using": by, "value": value}
	} else {
		params = map[string]interface{}{"using": by, "value": value, "element": *startNode}
	}

	response, err := c.transport.Send("findElement", params)
	if err != nil {
		return nil, err
	}

	var e = &webElement{c: c}
	err = json.Unmarshal([]byte(response.Value), &e)
	if err != nil {
		return nil, err
	}

	return e, nil
}

/////////////
// DIALOGS //
/////////////

func (c *Client) DismissDialog() (*responseError, bool, error) {
	ok := false
	r, err := c.transport.Send("dismissDialog", nil)
	if err != nil {
		return nil, ok, err
	}

	if r.ResponseError != nil {
		return r.ResponseError, ok, nil
	}

	return nil, true, nil
}

func (c *Client) AcceptDialog() (*responseError, bool, error) {
	ok := false
	r, err := c.transport.Send("acceptDialog", nil)
	if err != nil {
		return nil, ok, err
	}

	if r.ResponseError != nil {
		return r.ResponseError, ok, nil
	}

	return nil, true, nil
}

///////////////////////
// DISPOSE TEAR DOWN //
///////////////////////

func (c *Client) QuitApplication() (*response, error) {
	r, err := c.transport.Send("quitApplication", map[string]string{"flags": "eForceQuit"})
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (c *Client) Screenshot() (*response, error) {
	r, err := c.transport.Send("takeScreenshot", map[string]string{})
	if err != nil {
		return nil, err
	}

	return r, nil
	/*
		var d = map[string]string{}
		err = json.Unmarshal([]byte(r.Value), &d)
		if err != nil {
			return "", err
		}

		return d["value"], nil
	*/
}
