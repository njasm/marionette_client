package marionette_client

import (
	"errors"
	"testing"
	"time"
)

type fakeFinder struct {
	ReturnError bool
}

func (f fakeFinder) FindElement(by By, value string) (*WebElement, error) {
	if f.ReturnError {
		return nil, errors.New("ReturnError set.")
	}

	return new(WebElement), nil
}

func (f fakeFinder) FindElements(by By, value string) ([]*WebElement, error) {
	if f.ReturnError {
		return nil, errors.New("ReturnError set.")
	}

	var e []*WebElement
	return e, nil
}

func TestExpected(t *testing.T) {
	t.Run("ElementIsPresentFalseTest", ElementIsPresentFalseTest)
}

func ElementIsPresentFalseTest(t *testing.T) {
	fake := new(fakeFinder)
	fake.ReturnError = true

	fun := ElementIsPresent(By(ID), "")

	r, _, _ := fun(fake)
	if r {
		t.Fatalf("%v", "Result should be false")
	}
}

// required test in sequential main client test: client_test.go
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
