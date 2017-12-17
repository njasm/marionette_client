package marionette_client

import "testing"

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
	rv, err := NewDecoderEncoder(MARIONETTE_PROTOCOL_V3)
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
