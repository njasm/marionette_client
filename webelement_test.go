package marionette_client

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

func TestWebElementUnmarshalJSON(t *testing.T) {
	t.Run("InvalidJSONFormat", func(t *testing.T) {
		invalidJSON := []byte(`{"invalid": "json"}`)

		var element WebElement
		err := json.Unmarshal(invalidJSON, &element)
		if err == nil {
			t.Fatal("Expected error for invalid JSON format, got nil")
		}
	})

	t.Run("MissingWebDriverElementKey", func(t *testing.T) {
		missingKeyJSON := []byte(`{"value": {"someKey": "someValue"}}`)

		var element WebElement
		err := json.Unmarshal(missingKeyJSON, &element)
		if err == nil {
			t.Fatal("Expected error for missing WebDriverElementKey, got nil")
		}

		expectedErrorType := "WebDriverElementKey"
		expectedMessage := "key element-6066-11e4-a52e-4f735466cecf expected in response but not found"
		if err.Error() != expectedMessage {
			t.Errorf("Expected error message %q, but got %q", expectedMessage, err.Error())
		}

		driverError, ok := err.(*DriverError)
		if !ok {
			t.Fatalf("Expected error to be of type *DriverError, got %T", err)
		}

		if driverError.ErrorType != expectedErrorType {
			t.Errorf("Expected error type %q, but got %q", expectedErrorType, driverError.ErrorType)
		}
	})
}

func TestWebElementValueWrappers(t *testing.T) {
	transport := &recordingTransport{response: &Response{Value: `{"value":"result"}`}}
	client := NewClient()
	client.Transport(transport)
	element := &WebElement{id: "element-id", c: client}

	tests := []struct {
		name    string
		command string
		values  map[string]any
		call    func() string
	}{
		{name: "tag name", command: "WebDriver:GetElementTagName", values: map[string]any{"id": "element-id"}, call: element.TagName},
		{name: "text", command: "WebDriver:GetElementText", values: map[string]any{"id": "element-id"}, call: element.Text},
		{name: "attribute", command: "WebDriver:GetElementAttribute", values: map[string]any{"id": "element-id", "name": "title"}, call: func() string { return element.Attribute("title") }},
		{name: "property", command: "WebDriver:GetElementProperty", values: map[string]any{"id": "element-id", "name": "value"}, call: func() string { return element.Property("value") }},
		{name: "css", command: "WebDriver:GetElementCSSValue", values: map[string]any{"id": "element-id", "propertyName": "color"}, call: func() string { return element.CssValue("color") }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if value := test.call(); value != "result" {
				t.Fatalf("unexpected value %q", value)
			}
			if transport.command != test.command || !reflect.DeepEqual(transport.values, test.values) {
				t.Fatalf("unexpected call: command=%q values=%#v", transport.command, transport.values)
			}
		})
	}

	transport.err = errors.New("transport failed")
	if element.TagName() != "" || element.Text() != "" || element.Attribute("title") != "" || element.Property("value") != "" || element.CssValue("color") != "" {
		t.Fatal("value wrappers must return an empty string after transport errors")
	}
}

func TestWebElementBooleanWrappers(t *testing.T) {
	transport := &recordingTransport{response: &Response{Value: `{"value":true}`}}
	client := NewClient()
	client.Transport(transport)
	element := &WebElement{id: "element-id", c: client}

	tests := []struct {
		name    string
		command string
		call    func() bool
	}{
		{name: "enabled", command: "WebDriver:IsElementEnabled", call: element.Enabled},
		{name: "selected", command: "WebDriver:IsElementSelected", call: element.Selected},
		{name: "displayed", command: "WebDriver:IsElementDisplayed", call: element.Displayed},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if !test.call() || transport.command != test.command {
				t.Fatalf("unexpected result for %s using %q", test.name, transport.command)
			}
			if !reflect.DeepEqual(transport.values, map[string]any{"id": "element-id"}) {
				t.Fatalf("unexpected payload %#v", transport.values)
			}
		})
	}

	transport.response = &Response{Value: `{"value":false}`}
	if element.Enabled() || element.Selected() || element.Displayed() {
		t.Fatal("false responses must remain false")
	}
	transport.err = errors.New("transport failed")
	if element.Enabled() || element.Selected() || element.Displayed() {
		t.Fatal("transport errors must return false")
	}
}

func TestWebElementRectLocationAndSize(t *testing.T) {
	transport := &recordingTransport{response: &Response{Value: `{"x":1.5,"y":2.5,"width":30,"height":40}`}}
	client := NewClient()
	client.Transport(transport)
	element := &WebElement{id: "element-id", c: client}

	rect, err := element.Rect()
	if err != nil || rect.X != 1.5 || rect.Y != 2.5 || rect.Width != 30 || rect.Height != 40 {
		t.Fatalf("unexpected rectangle: %#v, %v", rect, err)
	}
	location, err := element.Location()
	if err != nil || location.X != 1.5 || location.Y != 2.5 {
		t.Fatalf("unexpected location: %#v, %v", location, err)
	}
	size, err := element.Size()
	if err != nil || size.Width != 30 || size.Height != 40 {
		t.Fatalf("unexpected size: %#v, %v", size, err)
	}
	if transport.command != "WebDriver:GetElementRect" || !reflect.DeepEqual(transport.values, map[string]any{"id": "element-id"}) {
		t.Fatalf("unexpected rectangle call: %q %#v", transport.command, transport.values)
	}

	transport.response = &Response{Value: "invalid"}
	if _, err = element.Rect(); err == nil {
		t.Fatal("expected malformed rectangle error")
	}
	transport.err = errors.New("transport failed")
	if _, err = element.Location(); err == nil {
		t.Fatal("expected location transport error")
	}
	if _, err = element.Size(); err == nil {
		t.Fatal("expected size transport error")
	}
}

func TestWebElementUnmarshalJSONSuccess(t *testing.T) {
	var element WebElement
	payload := []byte(`{"value":{"element-6066-11e4-a52e-4f735466cecf":"element-id"}}`)
	if err := json.Unmarshal(payload, &element); err != nil {
		t.Fatal(err)
	}
	if element.Id() != "element-id" {
		t.Fatalf("unexpected id %q", element.Id())
	}
}
