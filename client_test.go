package marionette_client

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"
)

const (
	TESTDATA_FOLDER   = "testdata"
	WWW_FOLDER        = "html"
	TARGET_URL        = "https://www.abola.pt/"
	ID_SELECTOR       = "clubes-hp"
	CSS_SELECTOR_LI   = "td"
	ID_SELECTOR_INPUT = "topo_txtPesquisa"
	TIMEOUT           = 10000 // milliseconds
)

var client *Client

func navigateLocal(page string) (*Response, error) {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(pwd)

	var schema = "file://" + pwd + "/" + TESTDATA_FOLDER + "/" + WWW_FOLDER + "/"
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

		t.Run("GetCookiesTest", GetCookiesTest)
		t.Run("GetCookieTest", GetCookieTest)

		t.Run("GetSessionCapabilitiesTest", GetSessionCapabilitiesTest)
		t.Run("ScreenshotTest", ScreenshotTest)

		t.Run("SetContextTest", SetContextTest)
		t.Run("GetContextTest", GetContextTest)

		t.Run("GetPageSourceTest", GetPageSourceTest)

		t.Run("SetScriptTimoutTest", SetScriptTimoutTest)
		t.Run("SetPageTimoutTest", SetPageTimoutTest)
		t.Run("SetSearchTimoutTest", SetSearchTimoutTest)

		t.Run("PageSourceTest", PageSourceTest)

		t.Run("ExecuteScriptWithoutFunctionTest", ExecuteScriptWithoutFunctionTest)
		t.Run("ExecuteScriptTest", ExecuteScriptTest)
		t.Run("ExecuteScriptWithArgsTest", ExecuteScriptWithArgsTest)

		t.Run("GetTitleTest", GetTitleTest)

		t.Run("FindElementTest", FindElementTest)

		t.Run("SendKeysTest", SendKeysTest)
		t.Run("FindElementsTest", FindElementsTest)

		t.Run("CurrentChromeWindowHandleTest", CurrentChromeWindowHandleTest)
		t.Run("WindowHandlesTest", WindowHandlesTest)

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
	r, err := client.Navigate(TARGET_URL)
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

	if url != TARGET_URL {
		t.Fatalf("Current Url %v not equal to target url %v", url, TARGET_URL)
	}

}

func GetCookiesTest(t *testing.T) {
	r, err := client.Cookies()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func GetCookieTest(t *testing.T) {
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

func GetSessionCapabilitiesTest(t *testing.T) {
	r, err := client.Capabilities()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.BrowserName)
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

func GetContextTest(t *testing.T) {
	r, err := client.Context()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func GetPageSourceTest(t *testing.T) {
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

func SetScriptTimoutTest(t *testing.T) {
	r, err := client.SetScriptTimeout(TIMEOUT)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func SetPageTimoutTest(t *testing.T) {
	r, err := client.SetPageLoadTimeout(TIMEOUT)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func SetSearchTimoutTest(t *testing.T) {
	r, err := client.SetImplicitTimout(TIMEOUT)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func PageSourceTest(t *testing.T) {
	_, err := client.PageSource()
	if err != nil {
		t.Fatalf("%#v", err)
	}
}

func ExecuteScriptWithoutFunctionTest(t *testing.T) {
	script := "return (document.readyState == 'complete');"
	args := []interface{}{}
	r, err := client.ExecuteScript(script, args, TIMEOUT, false)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func ExecuteScriptTest(t *testing.T) {
	script := "function testMyGoMarionetteClient() { return 'yes'; } return testMyGoMarionetteClient();"
	args := []interface{}{}
	r, err := client.ExecuteScript(script, args, TIMEOUT, false)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func ExecuteScriptWithArgsTest(t *testing.T) {
	script := "function testMyGoMarionetteClientArgs(a, b) { return a + b; }; return testMyGoMarionetteClientArgs(arguments[0], arguments[1]);"
	args := []interface{}{1, 3}
	r, err := client.ExecuteScript(script, args, TIMEOUT, false)
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
	navigateLocal("table.html")
	element, err := client.FindElement(By(ID), "the-table")
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

	collection, err := element.FindElements(By(CSS_SELECTOR), CSS_SELECTOR_LI)
	if 3 > len(collection) {
		t.FailNow()
	}

	t.Logf("%T %#v", collection, collection)

	el, err := element.FindElement(By(CSS_SELECTOR), CSS_SELECTOR_LI)
	if el == nil || err != nil {
		t.FailNow()
	}

}

func SendKeysTest(t *testing.T) {
	navigateLocal("form.html")
	e, err := client.FindElement(By(ID), "email")
	if err != nil {
		t.Fatalf("%#v", err)
	}

	var test string = "teste@example.com"
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
	navigateLocal("ul.html")
	elements, err := client.FindElements(By(CSS_SELECTOR), CSS_SELECTOR_LI)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(len(elements))
}

func CurrentChromeWindowHandleTest(t *testing.T) {
	r, err := client.CurrentChromeWindowHandle()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func WindowHandlesTest(t *testing.T) {
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

		time.Sleep(time.Second)
	}

	// return to original window.
	err = client.SwitchToWindow(w)
	if err != nil {
		t.Fatalf("%#v", err)
	}
}

func NavigatorMethodsTest(t *testing.T) {
	client.SetContext(Context(CONTENT))
	url1 := "https://www.google.pt/"
	url2 := "https://www.mercedes-benz.com/en/"

	client.Navigate(url1)
	sleep := time.Duration(2) * time.Second
	time.Sleep(sleep)

	client.Navigate(url2)
	time.Sleep(sleep)

	client.Back()
	client.Refresh()
	time.Sleep(sleep)

	firstUrl, err := client.Url()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	if firstUrl != url1 {
		t.Fatalf("Expected url %v - received url %v", url1, firstUrl)
	}

	client.Forward()
	secondUrl, err := client.Url()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	if secondUrl != url2 {
		t.Fatalf("Expected url %v - received url %v", url2, secondUrl)
	}
}

func PromptTest(t *testing.T) {
	navigateLocal("ul.html")
	var text string = "marionette is cool or what - prompt?"
	var script string = "prompt('" + text + "');"
	args := []interface{}{}

	r, err := client.ExecuteScript(script, args, TIMEOUT, false)
	if err != nil {
		t.Fatalf("%#v", err)
	}

	err = client.SendAlertText("yeah!")
	if err != nil {
		t.Fatalf("%#v", err)
	}

	time.Sleep(time.Duration(5) * time.Second)

	err = client.AcceptAlert()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

func AlertTest(t *testing.T) {
	navigateLocal("ul.html")
	var text string = "marionette is cool or what?"
	var script string = "alert('" + text + "');"
	args := []interface{}{}
	r, err := client.ExecuteScript(script, args, TIMEOUT, false)
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

	err = client.AcceptAlert()
	if err != nil {
		var version = client.browserVersion()
		if len(version) > 2 {
			major, err := strconv.ParseInt(version[0:2], 10, 0)
			if err == nil && major <= 83 {
				t.Skip("Skipping AcceptAlert for spurious failures")
				return
			}
		}
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

	// FIXME: Position is not changed via SetWindowRect so cannot assert with the expected value here

	if expectedRect.Width != actualRect.Width || expectedRect.Height != actualRect.Height {
		t.Fatalf("Size differs. expected: %v, actual: %v", expectedRect, *actualRect)
	}

	err = client.MaximizeWindow()
	if err != nil {
		t.Fatal("Unable to Maximize window")
	}
}

// working - if called before other tests all hell will break loose
//func TestCloseWindow(t *testing.T) {
//	r, err := client.CloseWindow()
//	if err != nil {
//		t.Fatalf("%#v", err)
//	}
//
//	t.Log(r.Value)
//}

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
