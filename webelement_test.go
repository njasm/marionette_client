package marionette_client

import (
	"encoding/json"
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
