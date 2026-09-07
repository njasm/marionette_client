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
	t.Run("ElementIsPresentSuccessTest", ElementIsPresentSuccessTest)
	t.Run("ElementIsNotPresentTest", ElementIsNotPresentTest)
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

func ElementIsPresentSuccessTest(t *testing.T) {
	ok, element, err := ElementIsPresent(Id, "element")(fakeFinder{})
	if err != nil || !ok || element == nil {
		t.Fatalf("unexpected present result: ok=%v element=%#v err=%v", ok, element, err)
	}
}

func ElementIsNotPresentTest(t *testing.T) {
	t.Run("no such element", func(t *testing.T) {
		finder := errorFinder{err: &DriverError{ErrorType: "no such element", Message: "missing"}}
		ok, element, err := ElementIsNotPresent(Id, "missing")(finder)
		if err != nil || !ok || element != nil {
			t.Fatalf("unexpected absent result: ok=%v element=%#v err=%v", ok, element, err)
		}
	})

	t.Run("unrelated error", func(t *testing.T) {
		expected := errors.New("transport failed")
		ok, _, err := ElementIsNotPresent(Id, "missing")(errorFinder{err: expected})
		if ok || !errors.Is(err, expected) {
			t.Fatalf("unexpected error result: ok=%v err=%v", ok, err)
		}
	})

	t.Run("element remains present", func(t *testing.T) {
		ok, element, err := ElementIsNotPresent(Id, "element")(fakeFinder{})
		if err != nil || ok || element == nil {
			t.Fatalf("unexpected present result: ok=%v element=%#v err=%v", ok, element, err)
		}
	})
}

type errorFinder struct{ err error }

func (f errorFinder) FindElement(By, string) (*WebElement, error) { return nil, f.err }
func (f errorFinder) FindElements(By, string) ([]*WebElement, error) {
	return nil, f.err
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
