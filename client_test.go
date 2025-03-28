package marionette_client

import (
	"fmt"
	"os"
	"testing"
	"time"
)

const (
	TestdataFolder   = "testdata"
	WwwFolder        = "html"
	TargetUrl        = "https://www.abola.pt/"
	CssSelectorTagTd = "td"
	Timeout          = 5000 // milliseconds
)

var client *Client

func navigateLocal(page string) (*Response, error) {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(pwd)

	var schema = "file://" + pwd + "/" + TestdataFolder + "/" + WwwFolder + "/"
	return client.Navigate(schema + page)
}

func init() {
	client = NewClient()
	client.Transport(&MarionetteTransport{})
	RunningInDebugMode = true
}

// we don't want parallel execution we need sequence.
func TestInit(t *testing.T) {
	t.Run("sequence", func(t *testing.T) {
		t.Run("NewSessionTest", NewSessionTest)
		t.Run("GetSessionIDTest", GetSessionIDTest)
		t.Run("GetPageTest", GetPageTest)
		t.Run("UrlTest", UrlTest)

		t.Run("AddCookiesTest", AddCookieTest)
		t.Run("GetCookiesTest", GetCookiesTest)
		t.Run("DeleteCookieTest", DeleteCookieTest)
		t.Run("DeleteAllCookiesTest", DeleteAllCookiesTest)

		t.Run("GetSessionCapabilitiesTest", GetSessionCapabilitiesTest)
		t.Run("ScreenshotTest", ScreenshotTest)

		t.Run("SetContextTest", SetContextTest)
		t.Run("GetContextTest", GetContextTest)

		t.Run("GetActiveElementTest", GetActiveElementTest)

		t.Run("GetPageSourceTest", GetPageSourceTest)

		t.Run("SetScriptTimoutTest", SetScriptTimoutTest)
		t.Run("SetPageTimoutTest", SetPageTimoutTest)
		t.Run("SetSearchTimoutTest", SetSearchTimoutTest)
		t.Run("GetTimeoutsTest", GetTimeoutsTest)

		t.Run("PageSourceTest", PageSourceTest)

		t.Run("ExecuteScriptWithoutFunctionTest", ExecuteScriptWithoutFunctionTest)
		t.Run("ExecuteScriptTest", ExecuteScriptTest)
		t.Run("ExecuteScriptWithArgsTest", ExecuteScriptWithArgsTest)

		t.Run("ExecuteAsyncScriptWithArgsTest", ExecuteAsyncScriptWithArgsTest)

		t.Run("GetTitleTest", GetTitleTest)

		t.Run("FindElementTest", FindElementTest)

		t.Run("SendKeysTest", SendKeysTest)
		t.Run("FindElementsTest", FindElementsTest)

		t.Run("NewWindowTest", NewWindowTest)
		t.Run("WindowHandlesTest", WindowHandlesTest)
		t.Run("CloseWindowTest", CloseWindowTest)
		t.Run("SwitchToParentFrameTest", SwitchToParentFrameTest)

		t.Run("NavigatorMethodsTest", NavigatorMethodsTest)

		t.Run("WaitForUntilIntegrationTest", WaitForUntilIntegrationTest)

		t.Run("WindowRectTest", WindowRectTest)

		t.Run("PromptTest", PromptTest)
		t.Run("AlertTest", AlertTest)

		// test expected.go
		t.Run("NotPresentTest", NotPresentTest)

		t.Run("DeleteSessionTest", DeleteSessionTest)

		// test QuitApplication
		t.Run("NewSessionTest", NewSessionTest)
		t.Run("QuitTest", QuitTest)
	})
}

/*********/
/* tests */
/*********/

func NewSessionTest(t *testing.T) {
	err := client.Connect("", 0)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	r, err := client.NewSession("", nil)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func GetSessionIDTest(t *testing.T) {
	if client.SessionId != client.SessionID() {
		t.Fatalf("SessionId differ...")
	}

	t.Log("session is : ", client.SessionId)
}

func GetPageTest(t *testing.T) {
	r, err := client.Navigate(TargetUrl)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func UrlTest(t *testing.T) {
	url, err := client.Url()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	if url != TargetUrl {
		t.Fatalf("Current Url %v not equal to target url %v", url, TargetUrl)
	}

}
func AddCookieTest(t *testing.T) {
	c := Cookie{
		Name:  "test-cookie",
		Value: "test-value",
	}

	r, err := client.AddCookie(c)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func GetCookiesTest(t *testing.T) {
	r, err := client.GetCookies()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r)
}

func DeleteCookieTest(t *testing.T) {
	//confirm cookie: test-cookie exists
	r, err := client.GetCookies()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	found := false
	for _, cookie := range r {
		if cookie.Name == "test-cookie" {
			found = true
			break
		}
	}

	// if not found error
	if !found {
		t.Fatal("Cookie 'test-cookie' not found in browser")
	}

	// delete it
	_, err = client.DeleteCookie("test-cookie")
	if err != nil {
		t.Fatalf("%#v", err)
	}

	// assert
	r, err = client.GetCookies()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	for _, cookie := range r {
		if cookie.Name == "test-cookie" {
			t.Fatal("Cookie 'test-cookie' was not deleted in browser")
		}
	}
}

func DeleteAllCookiesTest(t *testing.T) {
	// set browser in a controlled webpage
	_, _ = client.Navigate("http://example.com")

	// clear all visible cookies now
	cookies, err := client.GetCookies()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	for _, c := range cookies {
		_, _ = client.DeleteCookie(c.Name)
	}

	// add some dummy cookies
	cookies = append([]Cookie{},
		Cookie{Name: "test-cookie1", Value: "test-value1"},
		Cookie{Name: "test-cookie2", Value: "test-value2"},
	)
	for _, c := range cookies {
		_, _ = client.AddCookie(c)
	}

	// test those
	time.Sleep(time.Second)
	newCookies, err := client.GetCookies()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	if len(newCookies) != 2 {
		t.Fatalf("total number of cookies still dont match: %#v", newCookies)
	}

	err = client.DeleteAllCookies()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	time.Sleep(time.Second)
	cookies, err = client.GetCookies()
	t.Logf("Current cookies: %#v", cookies)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	if len(cookies) != 0 {
		t.Logf("%#v", cookies)
		t.Fatalf("Cookie jar should be empty but it's not")
	}

	// reset url for next tests
	_, _ = client.Navigate(TargetUrl)
}

//func TestConnectWithActiveConnection(t *testing.T) {
//	err := client.Connect("", 0)
//	if err == nil {
//		t.Fatalf("%#v", err)
//	}
//
//	t.Log("No Error..")
//}

func GetSessionCapabilitiesTest(t *testing.T) {
	r, err := client.GetCapabilities()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	if r.BrowserName != "firefox" {
		t.Fatal("Capabilities: Browser Name doesn't have the expected 'firefox' name")
	}

	t.Log(r)
}

func ScreenshotTest(t *testing.T) {
	_, err := client.Screenshot()
	if err != nil {
		t.Fatal(err)
	}

	//this print ise a problem for travis builds, since it can surpass the 4 MB of log size.
	// don't print the base64 encoded image.
	//println(base64encoded)
}

func SetContextTest(t *testing.T) {
	r, err := client.SetContext(Context(Chrome))
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)

	r, err = client.SetContext(Context(Content))
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func GetContextTest(t *testing.T) {
	r, err := client.Context()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func GetActiveElementTest(t *testing.T) {
	_, _ = navigateLocal("form.html")
	r, err := client.GetActiveElement()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	// there's always an active element?
	t.Logf("%#v", r)

	form, err := client.FindElement(Name, "optional")
	if err != nil {
		t.Fatalf("%#v", err)
	}

	// click on a *other* form element to activate
	e, _ := client.FindElement(Id, "email")
	e.Click()

	// assert now
	r, err = form.GetActiveElement()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	if r.Attribute("id") != "email" || r.TagName() != "input" {
		t.Fatalf("%#v", err)
	}

	_, _ = client.Navigate(TargetUrl)
}

func GetPageSourceTest(t *testing.T) {
	r, err := client.SetContext(Context(Chrome))
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)

	r, err = client.SetContext(Context(Content))
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func SetScriptTimoutTest(t *testing.T) {
	r, err := client.SetScriptTimeout(Timeout)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func SetPageTimoutTest(t *testing.T) {
	r, err := client.SetPageLoadTimeout(Timeout)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func SetSearchTimoutTest(t *testing.T) {
	r, err := client.SetImplicitTimout(Timeout)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func GetTimeoutsTest(t *testing.T) {
	r, err := client.GetTimeouts()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	if r["pageLoad"] != Timeout {
		t.Fatalf("pageLoad TIMEOUT value, was expected to be: %#v", Timeout)
	}

	t.Log(r)
}

func PageSourceTest(t *testing.T) {
	_, err := client.PageSource()
	if err != nil {
		t.Fatalf("%#v", err)
	}
}

func ExecuteScriptWithoutFunctionTest(t *testing.T) {
	script := "return (document.readyState == 'complete');"
	args := []any{}
	r, err := client.ExecuteScript(script, args, Timeout, false)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func ExecuteScriptTest(t *testing.T) {
	script := "function testMyGoMarionetteClient() { return 'yes'; } return testMyGoMarionetteClient();"
	args := []any{}
	r, err := client.ExecuteScript(script, args, Timeout, false)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func ExecuteScriptWithArgsTest(t *testing.T) {
	script := "function testMyGoMarionetteClientArgs(a, b) { return a + b; }; return testMyGoMarionetteClientArgs(arguments[0], arguments[1]);"
	args := []any{1, 3}
	r, err := client.ExecuteScript(script, args, Timeout, false)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func ExecuteAsyncScriptWithArgsTest(t *testing.T) {
	script := "function testMyGoMarionetteClientArgs(a, b) { return a + b; }; " +
		"let resolve = arguments[arguments.length - 1]; " +
		"let result = testMyGoMarionetteClientArgs(arguments[0], arguments[1]);" +
		"resolve(result);"

	args := []any{3, 3}
	r, err := client.ExecuteAsyncScript(script, args, false)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func GetTitleTest(t *testing.T) {
	title, err := client.Title()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(title)

}

func FindElementTest(t *testing.T) {
	_, err := navigateLocal("table.html")
	if err != nil {
		t.Fatalf("%#v", err)
	}

	element, err := client.FindElement(Id, "the-table")
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
	t.Log(element.Property("id"))
	t.Log(element.CssValue("text-decoration"))
	rect, err := element.Rect()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	// id attribute and property must be equal
	attId := element.Attribute("id")
	propId := element.Property("id")
	if attId != propId {
		t.Fatalf("Missmatch values from Attribute and Property 'id': first is %#v, former is %#v", attId, propId)
	}

	t.Log(rect)

	// location
	point, err := element.Location()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Logf("x: %f, y: %f", point.X, point.Y)

	//size
	size, err := element.Size()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Logf("w: %f, h: %f", size.Width, size.Height)

	// screenshot of node element
	_, err = element.Screenshot()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	collection, err := element.FindElements(CssSelector, CssSelectorTagTd)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	if 3 > len(collection) {
		t.FailNow()
	}

	t.Logf("%T %#v", collection, collection)

	el, err := element.FindElement(CssSelector, CssSelectorTagTd)
	if el == nil || err != nil {
		t.FailNow()
	}

}

func SendKeysTest(t *testing.T) {
	_, err := navigateLocal("form.html")
	if err != nil {
		t.Fatalf("%#v", err)
	}

	e, err := client.FindElement(Id, "email")
	if err != nil {
		t.Fatalf("%#v", err)
	}

	var test = "teste@example.com"
	err = e.SendKeys(test)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	/* FIXME: Text is not yet set. investigate.
	time.Sleep(time.Second * 5)
	if e.Text() != test {
		t.Fatalf("Elements text is not: %#v, it's: %#v", test, e.Text())
	}
	time.Sleep(time.Second * 10)
	*/

	e.Clear()
	if e.Text() != "" {
		t.Fatalf("Elements text should be empty. found: %#v", e.Text())
	}
}

func FindElementsTest(t *testing.T) {
	_, err := navigateLocal("ul.html")
	if err != nil {
		t.Fatalf("failed to navigate local: %#v", err)
	}

	elements, err := client.FindElements(CssSelector, CssSelectorTagTd)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(len(elements))
}

func NewWindowTest(t *testing.T) {
	r, err := client.GetWindowHandles()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	var count = len(r)
	var expectedCount = count + 1
	_, err = client.NewWindow(true, "tab", false)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	r, err = client.GetWindowHandles()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	var current = len(r)
	if len(r) != expectedCount {
		t.Fatalf("Number of Tab windows does not match, expected: %#v, got: %#v", expectedCount, current)
	}
}

func CloseWindowTest(t *testing.T) {
	r, err := client.GetWindowHandles()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	var current = len(r)
	if current <= 1 {
		t.Fatalf("Expected more then one window availiable, got: %#v", current)
	}

	_, err = client.CloseWindow()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	// return to browsing context
	r, err = client.GetWindowHandles()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	if len(r) >= 1 {
		err = client.SwitchToWindow(r[0])
		if err != nil {
			t.Fatalf("%#v", err)
		}
	}
}

func WindowHandlesTest(t *testing.T) {
	w, err := client.GetWindowHandle()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(w)

	r, err := client.GetWindowHandles()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	for _, w := range r {
		err := client.SwitchToWindow(w)
		if err != nil {
			t.Fatalf("%#v", err)
		}

		time.Sleep(time.Second)
	}

	// return to original window.
	err = client.SwitchToWindow(w)
	if err != nil {
		t.Fatalf("%#v", err)
	}
}

func SwitchToParentFrameTest(t *testing.T) {
	err := client.SwitchToParentFrame()
	if err != nil {
		t.Fatalf("%#v", err)
	}
}

func NavigatorMethodsTest(t *testing.T) {
	_, err := client.SetContext(Context(Content))
	if err != nil {
		t.Fatalf("failed to set context: %#v", err)
	}

	url1 := "https://www.google.pt/"
	url2 := "https://www.mercedes-benz.com/en/"

	_, err = client.Navigate(url1)
	if err != nil {
		t.Fatalf("failed to navigate url1: %#v", err)
	}

	sleep := time.Duration(2) * time.Second
	time.Sleep(sleep)

	_, err = client.Navigate(url2)
	if err != nil {
		t.Fatalf("failed to navigate url2: %#v", err)
	}

	time.Sleep(sleep)

	err = client.Back()
	if err != nil {
		t.Fatalf("failed to navigate back btn: %#v", err)
	}

	err = client.Refresh()
	if err != nil {
		t.Fatalf("failed to refresh: %#v", err)
	}
	time.Sleep(sleep)

	firstUrl, err := client.Url()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	if firstUrl != url1 {
		t.Fatalf("Expected url %v - received url %v", url1, firstUrl)
	}

	err = client.Forward()
	if err != nil {
		t.Fatalf("failed to navigate forward btn: %#v", err)
	}

	secondUrl, err := client.Url()
	if err != nil {
		t.Fatalf("failed to get current url: %#v", err)
	}

	if secondUrl != url2 {
		t.Fatalf("Expected url %v - received url %v", url2, secondUrl)
	}
}

func PromptTest(t *testing.T) {
	_, err := navigateLocal("ul.html")
	if err != nil {
		t.Fatalf("failed to navigate local: %#v", err)
	}

	var text = "marionette is cool or what - prompt?"
	var script = "prompt('" + text + "');"
	args := []any{}

	r, err := client.ExecuteScript(script, args, Timeout, false)
	if err != nil {
		t.Fatalf("failed to execute script: %#v", err)
	}

	err = client.SendAlertText("yeah!")
	if err != nil {
		t.Fatalf("failed to send alert text: %#v", err)
	}

	time.Sleep(time.Duration(5) * time.Second)

	err = client.AcceptAlert()
	if err != nil {
		t.Fatalf("failed to accept alert: %#v", err)
	}

	t.Log(r.Value)
}

func AlertTest(t *testing.T) {
	_, err := navigateLocal("table.html")
	if err != nil {
		t.Fatalf("failed to navigate local: %#v", err)
	}

	var text = "marionette is cool or what?"
	var script = "alert('" + text + "');"
	args := []any{}
	r, err := client.ExecuteScript(script, args, Timeout, false)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	textFromdialog, err := client.TextFromAlert()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	if textFromdialog != text {
		t.Fatalf("Text in dialog differ. expected: %v, textfromdialog: %v", text, textFromdialog)
	}

	time.Sleep(time.Duration(5) * time.Second)

	err = client.DismissAlert()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func WindowRectTest(t *testing.T) {
	expectedRect := WindowRect{X: 0, Y: 0, Width: 600, Height: 800}
	err := client.SetWindowRect(expectedRect)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	actualRect, _ := client.GetWindowRect()

	t.Logf("w: %v, h: %v", actualRect.Width, actualRect.Height)

	if expectedRect.Width != actualRect.Width || expectedRect.Height != actualRect.Height {
		t.Fatalf("Size differs. expected: %v, actual: %v", expectedRect, *actualRect)
	}

	_, err = client.MinimizeWindow()
	if err != nil {
		t.Fatal("Unable to Minimize window")
	}

	// FIXME: On CI MinimizeWindow command return the last WindowRect struct
	// this might be cos the lack of a window manager? - Need further investigation

	// sizes are expected to differ
	//if mr.Width == actualRect.Width || mr.Height == actualRect.Height {
	//	t.Fatalf("Size DOES NOT differs. actual: %v, mr: %v", actualRect, mr)
	//}

	_, err = client.MaximizeWindow()
	if err != nil {
		t.Fatal("Unable to Maximize window")
	}

	// FIXME: On CI MaximizeWindow command return the last WindowRect struct
	// this might be cos the lack of a window manager? - Need further investigation

	// sizes are expected to differ once again
	//if wr.Width == mr.Width || wr.Height == mr.Height {
	//	t.Fatalf("Size DOES NOT differs. wr: %v, mr: %v", wr, mr)
	//}

	_, err = client.FullscreenWindow()
	if err != nil {
		t.Fatal("Unable to Fullscreen window")
	}
}

// working - if called before other tests all hell will break loose
func DeleteSessionTest(t *testing.T) {
	err := client.DeleteSession()
	if err != nil {
		t.Fatalf("%#v", err)
	}
}

func QuitTest(t *testing.T) {
	r, err := client.Quit()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}
