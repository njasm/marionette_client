package marionette_client

import (
	"testing"
	"time"
)

const (
	TARGET_URL        = "http://www.abola.pt/"
	ID_SELECTOR       = "clubes-hp"
	CSS_SELECTOR_LI   = "li"
	ID_SELECTOR_INPUT = "topo_txtPesquisa"
	TIMEOUT           = 10000 // milliseconds
)

var client *Client

func init() {
	client = NewClient()
	client.Transport(&MarionetteTransport{})
	RunningInDebugMode = true
}

func TestNewSession(t *testing.T) {
	err := client.Connect("", 0)
	if err != nil {
		t.Error(err)
	}
	t.Log("got here")
	r, err := client.NewSession("", nil)
	if err != nil {
		t.Error(err)
	}

	t.Log(r.Value)
}

// working
func TestGetSessionID(t *testing.T) {
	if client.SessionId != client.SessionID() {
		t.Fatalf("SessionId differ...")
	}

	t.Log("session is : ", client.SessionId)
}

func TestGetPage(t *testing.T) {
	r, err := client.Navigate(TARGET_URL)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func TestCurrentUrl(t *testing.T) {
	url, err := client.CurrentUrl()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	if url != TARGET_URL {
		t.Fatalf("Current Url %v not equal to target url %v", url, TARGET_URL)
	}

}

func TestGetCookies(t *testing.T) {
	r, err := client.Cookies()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func TestGetCookie(t *testing.T) {
	r, err := client.Cookie("abolaCookie")
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
	r, err := client.Capabilities()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.BrowserName)
}

// working
func TestScreenshot(t *testing.T) {
	_, err := client.Screenshot()
	if err != nil {
		t.Fatal(err)
	}

	//this print ise a problem for travis builds, since it can surpass the 4 MB of log size.
	// don't print the base64 encoded image.
	//println(base64encoded)
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
	r, err := client.Logs()
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
	r, err := client.Context()
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
	r, err := client.SetScriptTimeout(TIMEOUT)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func TestSetPageTimout(t *testing.T) {
	r, err := client.SetPageTimeout(TIMEOUT)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func TestSetSearchTimout(t *testing.T) {
	r, err := client.SetSearchTimeout(TIMEOUT)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func TestPageSource(t *testing.T) {
	_, err := client.PageSource()
	if err != nil {
		t.Fatalf("%#v", err)
	}
}

func TestExecuteScriptWithoutFunction(t *testing.T) {
	script := "return (document.readyState == 'complete');"
	args := []interface{}{}
	r, err := client.ExecuteScript(script, args, 1000, false)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func TestExecuteScript(t *testing.T) {
	script := "function testMyGoMarionetteClient() { return 'yes'; } return testMyGoMarionetteClient();"
	args := []interface{}{}
	r, err := client.ExecuteScript(script, args, 1000, false)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func TestExecuteScriptWithArgs(t *testing.T) {
	script := "function testMyGoMarionetteClientArgs(a, b) { return a + b; }; return testMyGoMarionetteClientArgs(arguments[0], arguments[1]);"
	args := []interface{}{1, 3}
	r, err := client.ExecuteScript(script, args, 1000, false)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func TestGetTitle(t *testing.T) {
	title, err := client.Title()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(title)

}
func TestFindElement(t *testing.T) {
	element, err := client.FindElement(By(ID), ID_SELECTOR)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(element.Id())
	t.Log(element.Enabled())
	t.Log(element.Selected())
	t.Log(element.Displayed())
	t.Log(element.TagName())
	t.Log(element.Text())
	t.Log(element.Attribute("id"))
	t.Log(element.CssValue("text-decoration"))
	rect, err := element.Rect()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(rect)

	// size
	w, h, err := element.Location()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log("width: %f, height: %f", w, h)

	//location
	x, y, err := element.Size()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log("x: %f, y: %f", x, y)

	collection, err := element.FindElements(By(CSS_SELECTOR), CSS_SELECTOR_LI)
	if 18 != len(collection) {
		t.FailNow()
	}

	t.Logf("%T %#V", collection, collection)
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

func TestCurrentChromeWindowHandle(t *testing.T) {
	r, err := client.CurrentChromeWindowHandle()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func TestWindowHandles(t *testing.T) {
	w, err := client.CurrentWindowHandle()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(w)

	r, err := client.WindowHandles()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	for _, w := range r {
		err := client.SwitchToWindow(w)
		if err != nil {
			t.Fatalf("%#v", err)
		}

		time.Sleep(time.Duration(time.Second))
	}

	// return to original window.
	client.SwitchToWindow(w)
}

func TestNavigatorMethods(t *testing.T) {
	client.SetContext(Context(CONTENT))
	url1 := "https://www.google.pt/"
	url2 := "https://www.bing.com/"

	client.Navigate(url1)
	sleep := time.Duration(2) * time.Second
	time.Sleep(sleep)

	client.Navigate(url2)
	time.Sleep(sleep)

	client.Back()
	client.Refresh()
	time.Sleep(sleep)

	firstUrl, err := client.CurrentUrl()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	if firstUrl != url1 {
		t.Fatalf("Expected url %v - received url %v", url1, firstUrl[0:len(url1)])
	}

	client.Forward()
	secondUrl, err := client.CurrentUrl()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	if secondUrl != url2 {
		t.Fatalf("Expected url %v - received url %v", url2, secondUrl[:len(url2)])
	}
}


func TestWait(t *testing.T) {
	client.SetContext(Context(CONTENT))
	client.Navigate("http://www.w3schools.com/ajax/tryit.asp?filename=tryajax_get")

	timeout := time.Duration(10) * time.Second
	condition := ElementIsPresent(By(ID), "stackH")
	ok, v, err := Wait(client).For(timeout).Until(condition)

	if err != nil || !ok {
		t.Fatalf("%#v", err)
	}

	v.Click()

	err = client.SwitchToFrame(By(ID), "iframeResult")
	if err != nil {
		t.Fatalf("%#v", err)
	}

	e, err := client.FindElement(By(TAG_NAME), "button")
	if err != nil {
		t.Fatal("%#v", err)
	}

	e.Click()
}

func TestAlert(t *testing.T) {
	client.Get("http://www.abola.pt")
	var text string = "marionette is cool or what?"
	var script string = "alert('" + text + "');"
	args := []interface{}{}
	r, err := client.ExecuteScript(script, args, 1000, false)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	err = client.AcceptDialog()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	r, err = client.ExecuteScript(script, args, 1000, false)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	err = client.DismissDialog()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func TestNotPresent(t *testing.T) {
	client.SwitchToParentFrame()
	client.ActiveFrame()

	timeout := time.Duration(10) * time.Second
	condition := ElementIsNotPresent(By(ID), "non-existing-element")
	ok, _, _ := Wait(client).For(timeout).Until(condition)

	if !ok {
		t.Fatal("Element Was Found in ElementIsNotPresent test.")
	}
}

func TestWindowSize(t *testing.T) {
	w, h, err := client.WindowSize()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Logf("w: %v, h: %v", w, h)

	var newW float32 = w / 2
	var newH float32 = h / 2

	w, h, err = client.SetWindowSize(newW, newH)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Logf("new w: %v, new h: %v", w, h)
}

// working - if called before other tests all hell will break loose
func TestCloseWindow(t *testing.T) {
	r, err := client.CloseWindow()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

// working - if called before other tests all hell will break loose
func TestDeleteSession(t *testing.T) {
	err := client.DeleteSession()
	if err != nil {
		t.Fatalf("%#v", err)
	}
}


// working
//func TestQuitApplication(t *testing.T) {
//	r, err := client.QuitApplication()
//	if err != nil {
//		t.Fatalf("%#v", err)
//	}
//
//	t.Log(r.Value)
//}

