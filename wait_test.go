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
	var errorMsg = "the Error message."
	timeout := time.Duration(5) * time.Second
	condition := func(c Finder) (bool, *WebElement, error) {
		return false, nil, errors.New(errorMsg)
	}
	_, _, err := Wait(client).For(timeout).Until(condition)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != errorMsg {
		t.Fatalf("Expected error msg %v, got %v", errorMsg, err.Error())
	}
}

func UntilConditionNeverOccuredTest(t *testing.T) {
	timeout := 20 * time.Millisecond
	condition := func(c Finder) (bool, *WebElement, error) {
		return false, nil, nil
	}
	_, _, err := Wait(client).For(timeout).Until(condition)

	if err == nil {
		t.Fatal("Element Was Found in ElementIsNotPresent test.")
	}
}

func WaitForUntilIntegrationTest(t *testing.T) {
	_, err := client.SetContext(Content)
	if err != nil {
		t.Fatalf("failed to set context: %#v", err)
	}

	_, err = navigateLocal("ul.html")
	if err != nil {
		t.Fatalf("failed to navigate: %#v", err)
	}

	timeout := time.Duration(10) * time.Second
	condition := ElementIsPresent(Id, "delayed-frame-link")
	ok, v, err := Wait(client).For(timeout).Until(condition)
	if err != nil || !ok {
		t.Fatalf("%#v", err)
	}

	v.Click()

	err = client.SwitchToFrame(Id, "test-frame")
	if err != nil {
		t.Fatalf("%#v", err)
	}

	e, err := client.FindElement(TagName, "button")
	if err != nil {
		t.Fatalf("%#v", err)
	}

	e.Click()
}
