package marionette_client

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
)

const (
	MARIONETTE_PROTOCOL_V2 = 2
	MARIONETTE_PROTOCOL_V3 = 3

	WEBDRIVER_ELEMENT_KEY = "element-6066-11e4-a52e-4f735466cecf"
)

var RunningInDebugMode bool = false

type session struct {
	SessionId    string
	Capabilities Capabilities
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

func (c *Client) SessionID() string {
	return c.SessionId
}

func (c *Client) Connect(host string, port int) error {
	return c.transport.Connect(host, port)
}

// Capabilities Send the current session's capabilities to the client.
// Capabilities informs the client of which WebDriver features are
// supported by Firefox and Marionette.  They are immutable for the
// length of the session.
// The return value is an immutable map of string keys
// ("capabilities") to values, which may be of types boolean,
// numerical or string.
func (c *Client) Capabilities() (*Capabilities, error) {
	if c.session.SessionId != "" {
		return &c.session.Capabilities, nil
	}
	return &Capabilities{}, nil
}

/////////////
// SESSION //
/////////////

// NewSession create new session
func (c *Client) NewSession(sessionId string, cap *Capabilities) (*Response, error) {
	data := map[string]interface{}{
			"sessionId":    sessionId,
			"capabilities": cap,
	}

	var response *Response

	response, err := c.transport.Send("WebDriver:NewSession", data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(response.Value), &c)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// DeleteSession Marionette currently only accepts a session id, so if
// we call delete session can also close the TCP Connection
func (c *Client) DeleteSession() error {
	_, err := c.transport.Send("WebDriver:DeleteSession", nil)
	if err != nil {
		return err
	}

	return c.transport.Close()
}

// SetScriptTimeout Set the timeout for asynchronous script execution.
// Deprecated
func (c *Client) SetScriptTimeout(milliseconds int) (*Response, error) {
	return timeouts(&c.transport, "script", milliseconds)
}

// SetSearchTimeout Set timeout for searching for elements.
// Deprecated
func (c *Client) SetSearchTimeout(milliseconds int) (*Response, error) {
	return timeouts(&c.transport, "implicit", milliseconds)
}

// SetPageTimeout Set timeout for page loading.
// Deprecated
func (c *Client) SetPageTimeout(milliseconds int) (*Response, error) {
	return timeouts(&c.transport, "pageLoad", milliseconds)
}

// Set timeout for page loading, searching, and scripts.
//
// param string type
//     Type of timeout.
// param number ms
//     Timeout in milliseconds.
// Deprecated
func timeouts(transport *Transporter, typ string, milliseconds int) (*Response, error) {
	r, err := (*transport).Send("timeouts", map[string]interface{}{"type": typ, "ms": milliseconds})
	if err != nil {
		return nil, err
	}

	return r, nil
}

////////////////
// NAVIGATION //
////////////////

// Navigate open url
func (c *Client) Navigate(url string) (*Response, error) {
	r, err := c.transport.Send("WebDriver:Navigate", map[string]string{"url": url})
	if err != nil {
		return nil, err
	}

	return r, nil
}

// Title get title
func (c *Client) Title() (string, error) {
	r, err := c.transport.Send("WebDriver:GetTitle", map[string]string{})
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

// Url get current url
func (c *Client) Url() (string, error) {
	r, err := c.transport.Send("WebDriver:GetCurrentURL", nil)
	if err != nil {
		return "", err
	}

	var url map[string]string
	err = json.Unmarshal([]byte(r.Value), &url)
	if err != nil {
		return "", err
	}

	return url["value"], nil
}

// Refresh refresh
func (c *Client) Refresh() error {
	_, err := c.transport.Send("WebDriver:Refresh", nil)
	if err != nil {
		return err
	}

	return nil
}

// Back go back in navigation history
func (c *Client) Back() error {
	_, err := c.transport.Send("WebDriver:Back", nil)
	if err != nil {
		return err
	}

	return nil
}

// Forward go forward in navigation history
func (c *Client) Forward() error {
	_, err := c.transport.Send("WebDriver:Forward", nil)
	if err != nil {
		return err
	}

	return nil
}

// SetContext Sets the context of the subsequent commands to be either "chrome" or "content".
// Must be one of "chrome" or "content" only.
func (c *Client) SetContext(value Context) (*Response, error) {
	response, err := c.transport.Send("Marionette:SetContext", map[string]string{"value": fmt.Sprint(value)})
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Context Gets the context of the server, either "chrome" or "content".
func (c *Client) Context() (*Response, error) {
	response, err := c.transport.Send("Marionette:GetContext", nil)
	if err != nil {
		return nil, err
	}

	return response, nil
}

/////////////////////
// WINDOWS HANDLES //
/////////////////////

// CurrentWindowHandle returns the current window ID
func (c *Client) CurrentWindowHandle() (string, error) {
	r, err := c.transport.Send("WebDriver:GetWindowHandle", nil)
	if err != nil {
		return "", err
	}

	var d map[string]string
	err = json.Unmarshal([]byte(r.Value), &d)
	if err != nil {
		return "", err
	}
	return d["value"], nil
}

// CurrentChromeWindowHandle returns the current chrome window ID
//"getChromeWindowHandle": GeckoDriver.prototype.getChromeWindowHandle,
//"getCurrentChromeWindowHandle": GeckoDriver.prototype.getChromeWindowHandle,
func (c *Client) CurrentChromeWindowHandle() (*Response, error) {
	r, err := c.transport.Send("WebDriver:GetCurrentChromeWindowHandle", nil)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// WindowHandles return array of window ID currently opened
func (c *Client) WindowHandles() ([]string, error) {
	r, err := c.transport.Send("WebDriver:GetWindowHandles", nil)
	if err != nil {
		return nil, err
	}

	var d []string
	err = json.Unmarshal([]byte(r.Value), &d)
	if err != nil {
		return nil, err
	}

	return d, nil
}

// SwitchToWindow switch to specific window.
func (c *Client) SwitchToWindow(name string) error {
	_, err := c.transport.Send("WebDriver:SwitchToWindow", map[string]interface{}{"focus": true, "handle": name})
	if err != nil {
		return err
	}

	return nil
}

// GetWindowRect gets window position and size
func (c *Client) GetWindowRect() (rect *WindowRect, err error) {
	r, err := c.transport.Send("WebDriver:GetWindowRect", nil)
	if err != nil {
		return nil, err
	}

	rect = new(WindowRect)
	err = json.Unmarshal([]byte(r.Value), &rect)
	if err != nil {
		return nil, err
	}

	return
}

// SetWindowRect sets window position and size
func (c *Client) SetWindowRect(rect WindowRect) error {
	_, err := c.transport.Send("WebDriver:SetWindowRect", map[string]interface{}{"x": rect.X, "y": rect.Y, "width": math.Floor(rect.Width), "height": math.Floor(rect.Height)})
	if err != nil {
		return err
	}
	return nil
}

// MaximizeWindow maximizes window.
func (c *Client) MaximizeWindow() error {
	_, err := c.transport.Send("WebDriver:MaximizeWindow", nil)
	if err != nil {
		return err
	}

	return nil
}

// CloseWindow closes current window.
func (c *Client) CloseWindow() (*Response, error) {
	r, err := c.transport.Send("WebDriver:CloseWindow", nil)
	if err != nil {
		return nil, err
	}

	return r, nil
}

////////////
// FRAMES //
////////////

// SwitchToFrame switch to frame - strategies: By(ID), By(NAME) or name only.
func (c *Client) SwitchToFrame(by By, value string) error {

	//with current marionette implementation we have to find the element first and send the switchToFrame
	//command with the UUID, else it wont work.
	//https://bugzilla.mozilla.org/show_bug.cgi?id=1143908
	frame, err := c.FindElement(by, value)
	if err != nil {
		return err
	}

	_, err = c.transport.Send("WebDriver:SwitchToFrame", map[string]interface{}{"element": frame.Id(), "focus": true})
	if err != nil {
		return err
	}

	return nil
}

// SwitchToParentFrame switch to parent frame
func (c *Client) SwitchToParentFrame() error {
	_, err := c.transport.Send("WebDriver:SwitchToParentFrame", nil)
	if err != nil {
		return err
	}

	return nil
}

/////////////
// COOKIES //
/////////////

// Cookies Get all cookies
func (c *Client) Cookies() (*Response, error) {
	r, err := c.transport.Send("WebDriver:GetCookies", nil)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// Cookie Get cookie by name
func (c *Client) Cookie(name string) (*Response, error) {
	r, err := c.transport.Send("WebDriver:GetCookies", map[string]interface{}{"name": name})
	if err != nil {
		return nil, err
	}

	return r, nil
}

//////////////////
// WEB ELEMENTS //
//////////////////

func isElementEnabled(c *Client, id string) bool {
	r, err := c.transport.Send("WebDriver:IsElementEnabled", map[string]interface{}{"id": id})
	if err != nil {
		return false
	}

	return strings.Contains(r.Value, "\"value\":true")
}

func isElementSelected(c *Client, id string) bool {
	r, err := c.transport.Send("WebDriver:IsElementSelected", map[string]interface{}{"id": id})
	if err != nil {
		return false
	}

	return strings.Contains(r.Value, "\"value\":true")
}

func isElementDisplayed(c *Client, id string) bool {
	r, err := c.transport.Send("WebDriver:IsElementDisplayed", map[string]interface{}{"id": id})
	if err != nil {
		return false
	}

	return strings.Contains(r.Value, "\"value\":true")
}

func getElementTagName(c *Client, id string) string {
	r, err := c.transport.Send("WebDriver:GetElementTagName", map[string]interface{}{"id": id})
	if err != nil {
		return ""
	}

	var d = map[string]string{}
	json.Unmarshal([]byte(r.Value), &d)

	return d["value"]
}

func getElementText(c *Client, id string) string {
	r, err := c.transport.Send("WebDriver:GetElementText", map[string]interface{}{"id": id})
	if err != nil {
		return ""
	}

	var d = map[string]string{}
	json.Unmarshal([]byte(r.Value), &d)

	return d["value"]
}

func getElementAttribute(c *Client, id string, name string) string {
	r, err := c.transport.Send("WebDriver:GetElementAttribute", map[string]interface{}{"id": id, "name": name})
	if err != nil {
		return ""
	}

	var d = map[string]string{}
	json.Unmarshal([]byte(r.Value), &d)

	return d["value"]
}

func getElementCssPropertyValue(c *Client, id string, property string) string {
	r, err := c.transport.Send("WebDriver:GetElementCSSValue", map[string]interface{}{"id": id, "propertyName": property})
	if err != nil {
		return ""
	}

	var d = map[string]string{}
	json.Unmarshal([]byte(r.Value), &d)

	return d["value"]
}

func getElementRect(c *Client, id string) (*ElementRect, error) {
	r, err := c.transport.Send("WebDriver:GetElementRect", map[string]interface{}{"id": id})
	if err != nil {
		return nil, err
	}

	var d = &ElementRect{}
	err = json.Unmarshal([]byte(r.Value), &d)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func clickElement(c *Client, id string) {
	r, err := c.transport.Send("WebDriver:ElementClick", map[string]interface{}{"id": id})
	if err != nil {
		return
	}

	var d = map[string]interface{}{}
	json.Unmarshal([]byte(r.Value), &d)

	//return d
}

func sendKeysToElement(c *Client, id string, keys string) error {
	//slice := make([]string, 0)
	//for _, v := range keys {
	//	slice = append(slice, fmt.Sprintf("%c", v))
	//}
	//
	//r, err := c.transport.Send("sendKeysToElement", map[string]interface{}{"id": id, "value": slice})
	r, err := c.transport.Send("WebDriver:ElementSendKeys", map[string]interface{}{"id": id, "text": keys})
	if err != nil {
		return err
	}

	var d = map[string]interface{}{}
	json.Unmarshal([]byte(r.Value), &d)

	return nil
}

func clearElement(c *Client, id string) {
	r, err := c.transport.Send("WebDriver:ElementClear", map[string]interface{}{"id": id})
	if err != nil {
		return
	}

	var d = map[string]interface{}{}
	json.Unmarshal([]byte(r.Value), &d)

	//return d
}

// FindElements Find elements using the indicated search strategy.
func (c *Client) FindElements(by By, value string) ([]*WebElement, error) {
	return findElements(c, by, value, nil)
}

func findElements(c *Client, by By, value string, startNode *string) ([]*WebElement, error) {
	var params map[string]interface{}
	if startNode == nil || *startNode == "" {
		params = map[string]interface{}{"using": fmt.Sprint(by), "value": value}
	} else {
		params = map[string]interface{}{"using": fmt.Sprint(by), "value": value, "element": *startNode}
	}

	response, err := c.transport.Send("WebDriver:FindElements", params)
	if err != nil {
		return nil, err
	}

	var d []map[string]string
	err = json.Unmarshal([]byte(response.Value), &d)
	if err != nil {
		return nil, err
	}

	var e []*WebElement
	for _, v := range d {
		e = append(e, &WebElement{c: c, id: v[WEBDRIVER_ELEMENT_KEY]})
	}

	return e, nil

	//return string(buf), nil
}

// FindElement Find an element using the indicated search strategy.
func (c *Client) FindElement(by By, value string) (*WebElement, error) {
	return findElement(c, by, value, nil)
}

func findElement(c *Client, by By, value string, startNode *string) (*WebElement, error) {
	var params map[string]string
	if startNode == nil || *startNode == "" {
		params = map[string]string{"using": fmt.Sprint(by), "value": value}
	} else {
		params = map[string]string{"using": fmt.Sprint(by), "value": value, "element": *startNode}
	}

	response, err := c.transport.Send("WebDriver:FindElement", params)
	if err != nil {
		return nil, err
	}

	var e = &WebElement{c: c}
	err = json.Unmarshal([]byte(response.Value), &e)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func takeScreenshot(c *Client, startNode *string) (string, error) {
	var params map[string]string
	if startNode == nil || *startNode == "" {
		params = map[string]string{}
	} else {
		params = map[string]string{"id": *startNode}
	}

	r, err := c.transport.Send("WebDriver:TakeScreenshot", params)
	if err != nil {
		return "", err
	}

	return r.Value, nil
}

///////////////////////
// DOCUMENT HANDLING //
///////////////////////

// PageSource get page source
func (c *Client) PageSource() (*Response, error) {
	response, err := c.transport.Send("WebDriver:GetPageSource", nil)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// ExecuteScript Execute JS Script
func (c *Client) ExecuteScript(script string, args []interface{}, timeout uint, newSandbox bool) (*Response, error) {
	parameters := map[string]interface{}{}
	parameters["scriptTimeout"] = timeout
	parameters["script"] = script
	parameters["args"] = args

	parameters["newSandbox"] = newSandbox

	response, err := c.transport.Send("WebDriver:ExecuteScript", parameters)
	if err != nil {
		return nil, err
	}

	return response, nil
}

/////////////
// DIALOGS //
/////////////

// DismissAlert dismisses the dialog - like clicking No/Cancel
func (c *Client) DismissAlert() error {
	_, err := c.transport.Send("WebDriver:DismissAlert", nil)
	if err != nil {
		return err
	}

	return nil
}

// AcceptAlert accepts the dialog - like clicking Ok/Yes
func (c *Client) AcceptAlert() error {
	_, err := c.transport.Send("WebDriver:AcceptAlert", nil)
	if err != nil {
		return err
	}

	return nil
}

// TextFromAlert gets text from the dialog
func (c *Client) TextFromAlert() (string, error) {
	r, err := c.transport.Send("WebDriver:GetAlertText", map[string]interface{}{"key": "value"})
	if err != nil {
		return "", err
	}

	var d = map[string]string{}
	json.Unmarshal([]byte(r.Value), &d)

	return d["value"], nil
}

// SendAlertText sends text to a dialog
func (c *Client) SendAlertText(keys string) error {
	_, err := c.transport.Send("WebDriver:SendAlertText", map[string]interface{}{"text": keys})
	if err != nil {
		return err
	}

	return nil
}

///////////////////////
// DISPOSE TEAR DOWN //
///////////////////////

// Quit quits the session and request browser process to terminate.
func (c *Client) Quit() (*Response, error) {
	r, err := c.transport.Send("Marionette:Quit", map[string][]string{"flags": {"eForceQuit"}})
	if err != nil {
		return nil, err
	}

	return r, nil
}

// Screenshot takes a screenshot of the page.
func (c *Client) Screenshot() (string, error) {
	return takeScreenshot(c, nil)
}

func (c *Client) browserVersion() string {
	r, e := c.Capabilities()
	if e != nil {
		return ""
	}

	return r.BrowserVersion
}
