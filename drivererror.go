package marionette_client

import "fmt"

var _ fmt.Stringer = (*DriverError)(nil)
var _ error = (*DriverError)(nil)

type DriverError struct {
	ErrorType  string `json:"Error"`
	Message    string
	Stacktrace *string
}

func (e *DriverError) Error() string {
	return e.Message
}

func (e *DriverError) String() string {
	return e.Error()
}
