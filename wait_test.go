package marionette_client

import (
	"errors"
	"testing"
	"time"
)

func TestWait(t *testing.T) {
	t.Run("UntilConditionNeverOccuredTest", UntilConditionNeverOccuredTest)
	t.Run("UntilErrorTest", UntilErrorTest)
	t.Run("UntilSuccessTest", UntilSuccessTest)
	t.Run("ForBoundsTest", ForBoundsTest)
}

func UntilErrorTest(t *testing.T) {
	var errorMsg = "the Error message."
	timeout := time.Second
	condition := func(c Finder) (bool, *WebElement, error) {
		return false, nil, errors.New(errorMsg)
	}
	_, _, err := Wait(fakeFinder{}).For(timeout).Until(condition)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != errorMsg {
		t.Fatalf("Expected error msg %v, got %v", errorMsg, err.Error())
	}
}

func UntilSuccessTest(t *testing.T) {
	expected := &WebElement{id: "found"}
	ok, element, err := Wait(fakeFinder{}).For(0).Until(func(Finder) (bool, *WebElement, error) {
		return true, expected, nil
	})
	if err != nil || !ok || element != expected {
		t.Fatalf("unexpected result: ok=%v element=%#v err=%v", ok, element, err)
	}
}

func ForBoundsTest(t *testing.T) {
	waiter := Wait(fakeFinder{})
	if waiter.d != time.Duration(1) {
		t.Fatalf("unexpected default duration %v", waiter.d)
	}
	if waiter.For(25*time.Millisecond).d != 25*time.Millisecond {
		t.Fatalf("valid duration was not retained: %v", waiter.d)
	}
	if waiter.For(-time.Millisecond).d != time.Second {
		t.Fatalf("negative duration did not use fallback: %v", waiter.d)
	}
	if waiter.For(10*time.Minute+time.Nanosecond).d != time.Second {
		t.Fatalf("excessive duration did not use fallback: %v", waiter.d)
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
