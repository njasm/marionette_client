package marionette_client

import (
	"testing"
)

func TestNewSession(t *testing.T) {
	client := NewClient()
    err := client.Connect("", 0)
    if err != nil {
        t.Error(err)
    }

	_, err = client.NewSession("", nil)
	if err != nil {
		t.Error(err)
	}

    client.Close()
}
