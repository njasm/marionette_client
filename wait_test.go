package marionette_client

import (
	"errors"
	"testing"
	"time"
)

func TestWait(t *testing.T) {
	t.Run("UntilConditionNeverOccuredTest", UntilConditionNeverOccuredTest)
	// FIXME: t.Run("UntilErrorTest", UntilErrorTest)
}

func UntilErrorTest(t *testing.T) {
	var errorMsg string = "the Error message."
	timeout := time.Duration(5) * time.Second
	condition := func(c Finder) (bool, *WebElement, error) {
		return false, nil, errors.New(errorMsg)
	}
	_, _, err := Wait(client).For(timeout).Until(condition)

	if err.Error() != errorMsg {
		t.Fatalf("Expected error msg %v, got %v", errorMsg, err.Error())
	}
}

func UntilConditionNeverOccuredTest(t *testing.T) {
	timeout := time.Duration(11) * time.Minute
	condition := func(c Finder) (bool, *WebElement, error) {
		return false, nil, nil
	}
	_, _, err := Wait(client).For(timeout).Until(condition)

	if err == nil {
		t.Fatal("Element Was Found in ElementIsNotPresent test.")
	}
}

func WaitForUntilIntegrationTest(t *testing.T) {
	client.SetContext(Context(CONTENT))
	client.Navigate("http://www.w3schools.com/xml/tryit.asp?filename=tryajax_get")

	timeout := time.Duration(10) * time.Second
	condition := ElementIsPresent(By(CssSelector), "a.w3-button.w3-bar-item.topnav-icons.fa.fa-rotate")
	ok, v, err := Wait(client).For(timeout).Until(condition)

	if err != nil || !ok {
		t.Fatalf("%#v", err)
	}

	v.Click()

	err = client.SwitchToFrame(By(Id), "iframeResult")
	if err != nil {
		t.Fatalf("%#v", err)
	}

	e, err := client.FindElement(By(TagName), "button")
	if err != nil {
		t.Fatalf("%#v", err)
	}

	e.Click()
}
