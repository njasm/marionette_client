package marionette_client

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

type fixedMessageIDTransport struct{ id int }

func (t fixedMessageIDTransport) MessageID() int                    { return t.id }
func (fixedMessageIDTransport) Connect(string, int) error           { return nil }
func (fixedMessageIDTransport) Close() error                        { return nil }
func (fixedMessageIDTransport) Send(string, any) (*Response, error) { return nil, nil }
func (fixedMessageIDTransport) Receive() ([]byte, error)            { return nil, nil }

// test proto.go
func TestProto(t *testing.T) {
	t.Run("NewDecoderErrorTest", NewDecoderErrorTest)
	t.Run("DecodeErrorTest", DecodeErrorTest)
}

func NewDecoderErrorTest(t *testing.T) {
	var encID int32 = -1 //non existing protocol version
	_, err := NewDecoderEncoder(encID)
	if err == nil {
		t.Fatalf("%v", "Expected Error for non existing Marionette Protocol Version")
	}
}

func DecodeErrorTest(t *testing.T) {
	rv, err := NewDecoderEncoder(MarionetteProtocolV3)
	if err != nil {
		t.Fatalf("%v", err)
	}

	var errorJson = "12/{}ABC"
	var b = []byte(errorJson)
	if err = rv.Decode(b, nil); err == nil {
		t.Fatalf("%v", "Decoder call should error here.")
	}

	t.Logf("Expected error: %v", err)
}

func TestProtoV3Encode(t *testing.T) {
	encoder := ProtoV3DecoderEncoder{}
	encoded, err := encoder.Encode(fixedMessageIDTransport{id: 7}, "WebDriver:Command", map[string]any{"enabled": true})
	if err != nil {
		t.Fatal(err)
	}
	expected := `42:[0,7,"WebDriver:Command",{"enabled":true}]`
	if string(encoded) != expected {
		t.Fatalf("unexpected encoding %q", encoded)
	}

	if _, err = encoder.Encode(fixedMessageIDTransport{}, "bad", make(chan int)); err == nil {
		t.Fatal("expected JSON encoding error")
	}
}

func TestProtoV3DecodeValues(t *testing.T) {
	tests := []struct {
		name     string
		payload  string
		expected string
	}{
		{name: "object", payload: `[1,4,null,{"answer":42}]`, expected: `{"answer":42}`},
		{name: "array", payload: `[1,5,null,["a",2]]`, expected: `["a",2]`},
		{name: "scalar", payload: `[1,6,null,"ignored"]`, expected: ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := &Response{}
			if err := (ProtoV3DecoderEncoder{}).Decode([]byte(test.payload), response); err != nil {
				t.Fatal(err)
			}
			if response.MessageID == 0 || response.Size != int32(len(test.payload)) || response.Value != test.expected {
				t.Fatalf("unexpected response: %#v", response)
			}
		})
	}
}

func TestProtoV3DecodeDriverError(t *testing.T) {
	payload := `[1,8,{"error":"no such element","message":"missing","stacktrace":"trace"},null]`
	response := &Response{}
	err := (ProtoV3DecoderEncoder{}).Decode([]byte(payload), response)
	var driverError *DriverError
	if !errors.As(err, &driverError) {
		t.Fatalf("expected DriverError, got %T: %v", err, err)
	}
	if driverError.ErrorType != "no such element" || driverError.Message != "missing" || driverError.Stacktrace == nil || *driverError.Stacktrace != "trace" {
		t.Fatalf("unexpected driver error: %#v", driverError)
	}

	withoutTrace := `[1,9,{"error":"unknown error","message":"failed","stacktrace":null},null]`
	err = (ProtoV3DecoderEncoder{}).Decode([]byte(withoutTrace), &Response{})
	if !errors.As(err, &driverError) || driverError.Stacktrace != nil {
		t.Fatalf("unexpected nil-stacktrace error: %#v", err)
	}
}

func TestProtoV3EncodedMessageShape(t *testing.T) {
	encoded, err := (ProtoV3DecoderEncoder{}).Encode(fixedMessageIDTransport{id: 2}, "command", []string{"a"})
	if err != nil {
		t.Fatal(err)
	}
	var actual []any
	separator := 0
	for encoded[separator] != ':' {
		separator++
	}
	if err = json.Unmarshal(encoded[separator+1:], &actual); err != nil {
		t.Fatal(err)
	}
	expected := []any{float64(0), float64(2), "command", []any{"a"}}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("unexpected message: %#v", actual)
	}
}
