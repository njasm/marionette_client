package marionette_client

import (
	"testing"
	"time"
)

func WaitForUntilIntegrationTest(t *testing.T) {
	client.SetContext(Context(CONTENT))
	client.Navigate("http://www.w3schools.com/xml/tryit.asp?filename=tryajax_get")

	timeout := time.Duration(10) * time.Second
	condition := ElementIsPresent(By(CSS_SELECTOR), "a.w3-button.w3-bar-item.topnav-icons.fa.fa-rotate")
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
		t.Fatalf("%#v", err)
	}

	e.Click()
}

func NotPresentTest(t *testing.T) {
	client.SwitchToParentFrame()
	r, err := client.ActiveFrame()

	if err != nil {
		t.Fatalf("Getting Active Frame Error: %#v", err)
	}

	t.Log(r.id)

	timeout := time.Duration(10) * time.Second
	condition := ElementIsNotPresent(By(ID), "non-existing-element")
	ok, _, _ := Wait(client).For(timeout).Until(condition)

	if !ok {
		t.Fatal("Element Was Found in ElementIsNotPresent test.")
	}
}