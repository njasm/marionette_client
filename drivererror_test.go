package marionette_client

import (
	"testing"
)

func TestDriverError(t *testing.T) {
	t.Run("ErrorMethod", func(t *testing.T) {
		errorType := "SomeErrorType"
		message := "This is an error message"
		stacktrace := "stacktrace details"

		err := DriverError{
			ErrorType:  errorType,
			Message:    message,
			Stacktrace: &stacktrace,
		}

		expected := message
		actual := err.Error()

		if actual != expected {
			t.Errorf("Expected error message %q, but got %q", expected, actual)
		}
	})

	t.Run("StringMethod", func(t *testing.T) {
		errorType := "SomeErrorType"
		message := "This is an error message"
		stacktrace := "stacktrace details"

		err := DriverError{
			ErrorType:  errorType,
			Message:    message,
			Stacktrace: &stacktrace,
		}

		expected := message
		actual := err.String()

		if actual != expected {
			t.Errorf("Expected string representation %q, but got %q", expected, actual)
		}
	})
}
