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

	fun := ElementIsPresent(Id, "")

	r, _, _ := fun(fake)
	if r {
		t.Fatalf("%v", "Result should be false")
	}
}

// required test in sequential main client test: client_test.go
func NotPresentTest(t *testing.T) {
	timeout := time.Duration(5) * time.Second
	condition := ElementIsNotPresent(Id, "non-existing-element")
	ok, _, _ := Wait(client).For(timeout).Until(condition)

	if !ok {
		t.Fatal("Element Was Found in ElementIsNotPresent test.")
	}
}
