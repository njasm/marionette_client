package marionette_client

type DriverError struct {
	ErrorType  string `json:="Error"`
	Message    string
	Stacktrace *string
}

func (e DriverError) Error() string {
	return e.Message
}

func (e DriverError) String() string {
	return e.Error()
}
