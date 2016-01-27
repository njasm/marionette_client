package marionette_client

import (
	"testing"
    "fmt"
)

func TestNewSession(t *testing.T) {
	client := NewClient()
    err := client.Connect("", 0)
    if err != nil {
        t.Error(err)
    }

	r, err := client.NewSession("", nil)
	if err != nil {
		t.Error(err)
	}

    fmt.Println(r.Value)
    client.Close()
}
