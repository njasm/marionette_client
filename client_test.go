package marionette_client

import (
	"testing"
)

const (
	TARGET_URL = "http://www.abola.pt"
	ID_SELECTOR = "clubes-hp"
	CSS_SELECTOR_LI = "li"
	ID_SELECTOR_INPUT = "topo_txtPesquisa"
)
var client *Client

func init() {
	client = NewClient()
	client.Transport(&MarionetteTransport{})
}

func TestNewSession(t *testing.T) {
	err := client.Connect("", 0)
	if err != nil {
		t.Error(err)
	}

	r, err := client.NewSession("", nil)
	if err != nil {
		t.Error(err)
	}

	t.Log(r.Value)
}

// working
func TestGetSessionID(t *testing.T) {
	if client.SessionId != client.GetSessionID() {
		t.Fatalf("SessionId differ...")
	}

	t.Log("session is : ", client.SessionId)
}

func TestGetPage(t *testing.T) {
	r, err := client.Get(TARGET_URL)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func TestGetCookies(t *testing.T) {
	r, err := client.GetCookies()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func TestGetCookie(t *testing.T) {
	r, err := client.GetCookie("abolaCookie")
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

//func TestConnectWithActiveConnection(t *testing.T) {
//	err := client.Connect("", 0)
//	if err == nil {
//		t.Fatalf("%#v", err)
//	}
//
//	t.Log("No Error..")
//}

// working
func TestGetSessionCapabilities(t *testing.T) {
	r, err := client.GetSessionCapabilities()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.BrowserName)
}

// working
func TestScreenshot(t *testing.T) {
	r, err := client.Screenshot()
	if err != nil {
		t.Fatal(err)
	}

	println(r.MessageID)
	println(r.DriverError)
	println(r.Size)
	//this print ise a problem for travis builds, since it can surpass the 4 MB of log size.
	// don't print the base64 encoded image.
	//println(r.Value)
}

// working
func TestLog(t *testing.T) {
	r, err := client.Log("message testing", "warning")
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

// working
func TestGetLogs(t *testing.T) {
	r, err := client.GetLogs()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func TestSetContext(t *testing.T) {
	r, err := client.SetContext(Context(CHROME))
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)

	r, err = client.SetContext(Context(CONTENT))
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func TestGetContext(t *testing.T) {
	r, err := client.GetContext()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func TestGetPageSource(t *testing.T) {
	r, err := client.SetContext(Context(CHROME))
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)

	r, err = client.SetContext(Context(CONTENT))
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)

}

func TestSetScriptTimout(t *testing.T) {
	r, err := client.SetScriptTimeout(1000)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func TestSetPageTimout(t *testing.T) {
	r, err := client.SetPageTimeout(1000)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func TestSetSearchTimout(t *testing.T) {
	r, err := client.SetSearchTimeout(1000)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func TestExecuteScript(t *testing.T) {
	script := "function testMyGoMarionetteClient() { return 'yes'; } testMyGoMarionetteClient();"
	args := []interface{}{}
	r, err := client.ExecuteScript(script, args, 1000, false)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func TestExecuteScriptWithArgs(t *testing.T) {
	script := "function testMyGoMarionetteClientArgs(a, b) { return a + b; }; testMyGoMarionetteClientArgs(arguments[0], arguments[1]);"
	args := []interface{}{1, 3}
	r, err := client.ExecuteScript(script, args, 1000, false)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func TestGetTitle(t *testing.T) {
	title, err := client.GetTitle()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(title)

}
func TestFindElement(t *testing.T) {
	element, err := client.FindElement(By(ID), ID_SELECTOR)
	if err != nil {
		t.Fatalf("%#v", err)
		t.Log(element)
		t.FailNow()
	}

	t.Log(element.Id())
	t.Log(element.IsEnabled())
	t.Log(element.IsSelected())
	t.Log(element.IsDisplayed())
	t.Log(element.TagName())
	t.Log(element.Text())
	t.Log(element.Attribute("id"))
	t.Log(element.CssValue("text-decoration"))
	t.Log(element.Rect())

	collection, err := element.FindElements(By(CSS_SELECTOR), CSS_SELECTOR_LI)
	if 18 != len(collection) {
		t.FailNow()
	}
}

func TestSendKeys(t *testing.T) {
	e, err := client.FindElement(By(ID), ID_SELECTOR_INPUT)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	e.SendKeys("teste")
}

func TestFindElements(t *testing.T) {
	elements, err := client.FindElements(By(CSS_SELECTOR), CSS_SELECTOR_LI)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(len(elements))
}

func TestCurrentWindowHandle(t *testing.T) {
	r, err := client.GetCurrentWindowHandle()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func TestCurrentChromeWindowHandle(t *testing.T) {
	r, err := client.GetCurrentChromeWindowHandle()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func TestWindowHandles(t *testing.T) {
	r, err := client.WindowHandles()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r)
}

// working
//func TestQuitApplication(t *testing.T) {
//	r, err := client.QuitApplication()
//	if err != nil {
//		t.Fatalf("%#v", err)
//		t.FailNow()
//	}
//
//	t.Log(r.ResponseError)
//}
