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

	// find iframe switch to and click button
	frame, err := client.FindElement(By(ID), "iframeResult")
	if err != nil {
		t.Fatalf("%#v", err)
	}

	err = client.SwitchToFrame(By(WEBELEMENT_ID), frame.Id())
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
	//f, _ := client.GetActiveFrame()
	//t.Log(f.Attribute("id"))

	timeout := time.Duration(10) * time.Second
	condition := ElementIsNotPresent(By(ID), "non-existing-element")
	ok, _, _ := Wait(client).For(timeout).Until(condition)

	if !ok {
		t.Fatal("Element Was Found in ElementIsNotPresent test.")
	}
}
