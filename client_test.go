package marionette_client

import (
	"fmt"
	"testing"
)

var client *Client

func init() {
	client = NewClient()
}

func TestNewSession(t *testing.T) {

	err := client.Connect("", 0)
	if err != nil {
		t.Error(err)
	}

	r, err := client.NewSession("", nil)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("using marionette protocol: ", client.transport.MarionetteProtocol)
	fmt.Println(r.Value)
	client.Close()
}

func TestGetSessionID(t *testing.T) {
	if client.SessionId != client.GetSessionID() {
		fmt.Println("SessionId differ...")
		t.FailNow()
	}

	fmt.Println("session is : ", client.SessionId)
}

func TestConnectWithActiveConnection(t *testing.T) {
	err := client.Connect("", 0)
	if err == nil {
		fmt.Println(err)
	}

	fmt.Println("No Error..")
}

func TestGetMarionetteID(t *testing.T) {
	//not availiable command in protocol 2?
	if client.MarionetteProtocol == 2 {
		t.SkipNow()
	}

	r, err := client.GetMarionetteID()
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	fmt.Println(r.Value)
}

func TestGetSessionCapabilities(t *testing.T) {
	r, err := client.GetSessionCapabilities()
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	fmt.Println(r)
}

func TestLog(t *testing.T) {
    r, err := client.Log("message testing", "warning")
    if err != nil {
        fmt.Println(err)
        t.FailNow()
    }

    fmt.Println(r.Value)
}

func TestGetLogs(t *testing.T) {
	r, err := client.GetLogs()
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	fmt.Println(r)
}
