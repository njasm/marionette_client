package marionette_client

import (
	"testing"
	"time"
)

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

// working - if called before other tests all hell will break loose
func TestCloseWindow(t *testing.T) {
	r, err := client.CloseWindow()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}

// working
func TestQuitApplication(t *testing.T) {
	r, err := client.QuitApplication()
	if err != nil {
		t.Fatalf("%#v", err)
	}

	t.Log(r.Value)
}
