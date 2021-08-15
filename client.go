package marionette_client

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
)

const (
	MARIONETTE_PROTOCOL_V3 = 3
	WEBDRIVER_ELEMENT_KEY  = "element-6066-11e4-a52e-4f735466cecf"
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

// Capabilities informs the client of which WebDriver features are
// supported by Firefox and Marionette. They are immutable for the
// length of the session.
func (c *Client) GetCapabilities() (*Capabilities, error) {
	r, err := c.transport.Send("WebDriver:GetCapabilities", map[string]string{})
	if err != nil {
		return nil, err
	}

	var d = map[string]*Capabilities{}
	err = json.Unmarshal([]byte(r.Value), &d)
	if err != nil {
		return nil, err
	}

	return d["capabilities"], nil
}

// SetScriptTimeout Set the timeout for asynchronous script execution.
func (c *Client) SetScriptTimeout(milliseconds int) (*Response, error) {
	data := map[string]int{"script": milliseconds}
	return c.SetTimeouts(data)
}

// SetImplicitTimout Set timeout for searching for elements.
func (c *Client) SetImplicitTimout(milliseconds int) (*Response, error) {
	data := map[string]int{"implicit": milliseconds}
	return c.SetTimeouts(data)
}

// SetPageLoadTimeout Set timeout for page loading.
func (c *Client) SetPageLoadTimeout(milliseconds int) (*Response, error) {
	data := map[string]int{"pageLoad": milliseconds}
	return c.SetTimeouts(data)
}

// <h4>Timeouts object</h4>
//
// <dl>
// <dt><code>script</code> (number)
// <dd>Determines when to interrupt a script that is being evaluates.
//
// <dt><code>pageLoad</code> (number)
// <dd>Provides the timeout limit used to interrupt navigation of the
//  browsing context.
//
// <dt><code>implicit</code> (number)
// <dd>Gives the timeout of when to abort when locating an element.
// </dl>
func (c *Client) SetTimeouts(data map[string]int) (*Response, error) {
	r, err := c.transport.Send("WebDriver:SetTimeouts", data)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// Get Timeouts Get current set timeouts
func (c *Client) GetTimeouts() (map[string]uint, error) {
	r, err := c.transport.Send("WebDriver:GetTimeouts", map[string]string{})
	if err != nil {
		return nil, err
	}

	var d = map[string]uint{}
	err = json.Unmarshal([]byte(r.Value), &d)
	if err != nil {
		return nil, err
	}

	return d, nil
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

// GetWindowHandle returns the current window ID
func (c *Client) GetWindowHandle() (string, error) {
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

// GetWindowHandles return array of window ID currently opened
func (c *Client) GetWindowHandles() ([]string, error) {
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

// GetChromeWindowHandle returns the current chrome window ID
func (c *Client) GetChromeWindowHandle() (*string, error) {
	//"getChromeWindowHandle": GeckoDriver.prototype.getChromeWindowHandle,
	//"getCurrentChromeWindowHandle": GeckoDriver.prototype.getChromeWindowHandle,
	r, err := c.transport.Send("WebDriver:GetCurrentChromeWindowHandle", nil)
	if err != nil {
		return nil, err
	}

	var d map[string]string
	err = json.Unmarshal([]byte(r.Value), &d)
	if err != nil {
		return nil, err
	}

	t := d["value"]
	return &t, nil
}

// GetChromeWindowHandles return array of chrome window ID
func (c *Client) GetChromeWindowHandles() ([]string, error) {
	r, err := c.transport.Send("WebDriver:GetChromeWindowHandles", nil)
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
func (c *Client) MaximizeWindow() (rect *WindowRect, err error) {
	r, err := c.transport.Send("WebDriver:MaximizeWindow", nil)
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

// MinimizeWindow Synchronously minimizes the user agent window as if the user pressed
// the minimize button.
func (c *Client) MinimizeWindow() (rect *WindowRect, err error) {
	r, err := c.transport.Send("WebDriver:MinimizeWindow", nil)
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

// FullscreenWindow Synchronously sets the user agent window to full screen as if the user
// had done "View > Enter Full Screen"
func (c *Client) FullscreenWindow() (rect *WindowRect, err error) {
	r, err := c.transport.Send("WebDriver:FullscreenWindow", nil)
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

// NewWindow opens a new top-level browsing context window.
//
// param: type string
// Optional type of the new top-level browsing context. Can be one of
// `tab` or `window`. Defaults to `tab`.
//
// param: focus bool
// Optional flag if the new top-level browsing context should be opened
// in foreground (focused) or background (not focused). Defaults to false.
//
// param: private bool
// Optional flag, which gets only evaluated for type `window`. True if the
// new top-level browsing context should be a private window.
// Defaults to false.
//
// return {"handle": string, "type": string}
// Handle and type of the new browsing context.
func (c *Client) NewWindow(focus bool, typ string, private bool) (*Response, error) {
	r, err := c.transport.Send("WebDriver:NewWindow", map[string]interface{}{"focus": focus, "type": typ, "private": private})
	if err != nil {
		return nil, err
	}

	//TODO: would be nice if we could create a Window struct and return that struct instead of the Response object
	return r, nil
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

// AddCookie Adds a cookie
func (c *Client) AddCookie(cookie Cookie) (*Response, error) {
	r, err := c.transport.Send("WebDriver:AddCookie", map[string]interface{}{"cookie": cookie})
	if err != nil {
		return nil, err
	}

	return r, nil
}

// GetCookies Get all cookies
func (c *Client) GetCookies() ([]Cookie, error) {
	r, err := c.transport.Send("WebDriver:GetCookies", nil)
	if err != nil {
		return nil, err
	}

	var cookies []Cookie
	_ = json.Unmarshal([]byte(r.Value), &cookies)

	return cookies, nil
}

// DeleteCookie Deletes cookie by name
func (c *Client) DeleteCookie(name string) (error, error) {
	_, err := c.transport.Send("WebDriver:DeleteCookie", map[string]interface{}{"name": name})
	return err, nil
}

// DeleteAllCookies Delete all cookies
func (c *Client) DeleteAllCookies() error {
	_, err := c.transport.Send("WebDriver:DeleteAllCookies", nil)
	return err
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

func getElementProperty(c *Client, id string, name string) string {
	r, err := c.transport.Send("WebDriver:GetElementProperty", map[string]interface{}{"id": id, "name": name})
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

// GetActiveElement Returns the page's active element.
func (c *Client) GetActiveElement() (*WebElement, error) {
	return getActiveElement(c)
}

func getActiveElement(c *Client) (*WebElement, error) {
	response, err := c.transport.Send("WebDriver:GetActiveElement", nil)
	if err != nil {
		return nil, err
	}

	var e = &WebElement{c: c}
	err = json.Unmarshal([]byte(response.Value), e)
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

// ExecuteAsyncScript Execute JS Script Async
//TODO: Add missing arguments/options
func (c *Client) ExecuteAsyncScript(script string, args []interface{}, newSandbox bool) (*Response, error) {
	parameters := map[string]interface{}{}
	parameters["script"] = script
	parameters["args"] = args

	parameters["newSandbox"] = newSandbox

	response, err := c.transport.Send("WebDriver:ExecuteAsyncScript", parameters)
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
