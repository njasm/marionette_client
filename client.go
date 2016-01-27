package marionette_client
import (
    "encoding/json"
)

type Context string

const (
    CHROME  Context = "chorme"
    CONTENT Context = "content"
)

type session struct {
    sessionId       string
}

type Capabilities struct {
    BrowserName                     string
    BrowserVersion                  string
    PlatformName                    string
    PlatformVersion                 string
    SpecificationLevel              string
    RaisesAccessibilityExceptions   bool
    Rotatable                       bool
    AcceptSslCerts                  bool
    TakesElementScreenshot          bool
    TakesScreenshot                 bool
    Proxy                           interface{}
    Platform                        string
    XULappId                        string
    AppBuildId                      string
    Device                          string
    Version                         string
}

type Client struct {
    session
    transport
}

func NewClient() Client {
    return &Client{
        &session{},
        &transport{},
    }
}

func (c *Client) GetSessionID() string {
    if c.session != nil {
        return c.sessionId
    }

    return ""
}

// Protocol commands
func (c *Client) GetMarionetteID() string {
    return c.transport.send("getMarionetteID", nil)
}

func (c *Client) SayHello() string {
    return c.transport.send("sayHello", nil)
}

func (c *Client) NewSession(sessionId string, cap Capabilities) error {

    data := map[string]interface{} {
        "name": "newSession",
        "parameters": {
            "sessionId": sessionId,
            "capabilities" : cap,
        },
    }

    c.transport.send(data);
    return nil
}

/**
 * Send the current session's capabilities to the client.
 *
 * Capabilities informs the client of which WebDriver features are
 * supported by Firefox and Marionette.  They are immutable for the
 * length of the session.
 *
 * The return value is an immutable map of string keys
 * ("capabilities") to values, which may be of types boolean,
 * numerical or string.
 */
func (c *Client) GetSessionCapabilities() (Capabilities, error) {
    buf, err := c.transport.send("getSessionCapabilities", nil)
    if err != nil {
        return nil, err
    }

    response := &Capabilities{}
    err = json.Unmarshal(buf, response)
    if err != nil {
        return nil, err
    }

    return response, nil
}

/**
 * Log message.  Accepts user defined log-level.
 *
 * @param {string} value
 *     Log message.
 * @param {string} level
 *     Arbitrary log level.
 */
func (c *Client) Log(message string, level string) string {
    return c.transport.send("log", map[string]string{"value": message, "level": level})
}

/** Return all logged messages. */
func (c *Client) GetLogs() string {
    return c.transport.send("getLogs", nil)
}

/**
 * Sets the context of the subsequent commands to be either "chrome" or
 * "content".
 *
 * @param {string} value
 *     Name of the context to be switched to.  Must be one of "chrome" or
 *     "content".
 */
func (c *Client) SetContext(value Context) {
    return c.transport.send("setContext", map[string]string{"value": value})
}

/** Gets the context of the server, either "chrome" or "content". */
func (c *Client) GetContext() {
    return c.transport.send("getContext", nil)
}

func (c *Client) ExecuteScript(value Context) {
    return c.transport.send("executeScript", map[string]string{"value": value})
}

/**
 * Set the timeout for asynchronous script execution.
 *
 * @param {number} ms
 *     Time in milliseconds.
 */
func (c *Client) SetScriptTimeout(milliseconds int) {
    return c.timeouts("script", milliseconds)
}

/**
 * Set timeout for searching for elements.
 *
 * @param {number} ms
 *     Search timeout in milliseconds.
 */
func (c *Client) SetSearchTimeout(milliseconds int) {
    return c.timeouts("implicit", milliseconds)
}

/**
 * Set timeout for page loading.
 *
 * @param {number} ms
 *     Search timeout in milliseconds.
 */
func (c *Client) SetPageTimeout(milliseconds int) {
    return c.timeouts(nil, milliseconds)
}

/** internal
 * Set timeout for page loading, searching, and scripts.
 *
 * @param {string} type
 *     Type of timeout.
 * @param {number} ms
 *     Timeout in milliseconds.
 */
func (c *Client) timeouts(typ string, milliseconds int) {
    if typ != "implicit" && typ != "script" {
        typ = nil
    }

    return c.transport.send("timeouts", map[string]interface{}{"type": typ, "ms": milliseconds})
}

/** Single tap. */
/* TODO: grab a firefoxOS device to test and implement this.
GeckoDriver.prototype.singleTap = function(cmd, resp) {
let {id, x, y} = cmd.parameters;

switch (this.context) {
case Context.CHROME:
throw new WebDriverError("Command 'singleTap' is not available in chrome context");

case Context.CONTENT:
this.addFrameCloseListener("tap");
yield this.listener.singleTap(id, x, y);
break;
}
};
*/

/**
 * Find elements using the indicated search strategy.
 *
 * @param {string} using
 *     Indicates which search method to use.
 * @param {string} value
 *     Value the client is looking for.
 */
func (c *Client) FindElements(by string, value string) {
    return c.transport.send("findElements", map[string]interface{} {})
}

func (c *Client) Close() (response, error){
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

func (c *Client) GetStatus()