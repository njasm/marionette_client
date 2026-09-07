package marionette_client

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"
)

type recordingTransport struct {
	command string
	values  any
	err     error
}

func (t *recordingTransport) MessageID() int            { return 0 }
func (t *recordingTransport) Connect(string, int) error { return nil }
func (t *recordingTransport) Close() error              { return nil }
func (t *recordingTransport) Receive() ([]byte, error)  { return nil, nil }
func (t *recordingTransport) Send(command string, values any) (*Response, error) {
	t.command = command
	t.values = values
	if t.err != nil {
		return nil, t.err
	}
	return &Response{Value: "{}"}, nil
}

func TestPerformActionsPayload(t *testing.T) {
	transport := &recordingTransport{}
	client := NewClient()
	client.Transport(transport)
	element := &WebElement{id: "element-id", c: client}

	_, err := client.PerformActions(
		MouseActions("mouse",
			PointerMove(10, -5, 1500*time.Millisecond, ViewportOrigin()),
			PointerMove(1, 2, 0, PointerOriginCurrent()),
			PointerMove(0, 0, 10*time.Millisecond, ElementOrigin(element)),
			PointerDown(0),
			Pause(25*time.Millisecond),
			PointerUp(0),
		),
		MouseActions("second", Pause(2*time.Millisecond)),
	)
	if err != nil {
		t.Fatal(err)
	}
	if transport.command != "WebDriver:PerformActions" {
		t.Fatalf("unexpected command %q", transport.command)
	}

	encoded, err := json.Marshal(transport.values)
	if err != nil {
		t.Fatal(err)
	}
	var actual any
	if err := json.Unmarshal(encoded, &actual); err != nil {
		t.Fatal(err)
	}
	var expected any
	expectedJSON := `{"actions":[{"actions":[{"duration":1500,"origin":"viewport","type":"pointerMove","x":10,"y":-5},{"duration":0,"origin":"pointer","type":"pointerMove","x":1,"y":2},{"duration":10,"origin":{"element-6066-11e4-a52e-4f735466cecf":"element-id"},"type":"pointerMove","x":0,"y":0},{"button":0,"type":"pointerDown"},{"duration":25,"type":"pause"},{"button":0,"type":"pointerUp"}],"id":"mouse","parameters":{"pointerType":"mouse"},"type":"pointer"},{"actions":[{"duration":2,"type":"pause"}],"id":"second","parameters":{"pointerType":"mouse"},"type":"pointer"}]}`
	if err := json.Unmarshal([]byte(expectedJSON), &expected); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("payload mismatch\nactual:   %s\nexpected: %s", encoded, expectedJSON)
	}
}

func TestPerformActionsValidation(t *testing.T) {
	client := NewClient()
	transport := &recordingTransport{}
	client.Transport(transport)

	tests := []struct {
		name      string
		sequences []ActionSequence
	}{
		{"no sequences", nil},
		{"empty id", []ActionSequence{MouseActions("", Pause(0))}},
		{"duplicate id", []ActionSequence{MouseActions("mouse", Pause(0)), MouseActions("mouse", Pause(0))}},
		{"empty actions", []ActionSequence{MouseActions("mouse")}},
		{"invalid action", []ActionSequence{MouseActions("mouse", PointerAction{})}},
		{"negative move duration", []ActionSequence{MouseActions("mouse", PointerMove(0, 0, -time.Millisecond, ViewportOrigin()))}},
		{"invalid origin", []ActionSequence{MouseActions("mouse", PointerMove(0, 0, 0, PointerOrigin{}))}},
		{"nil element", []ActionSequence{MouseActions("mouse", PointerMove(0, 0, 0, ElementOrigin(nil)))}},
		{"empty element id", []ActionSequence{MouseActions("mouse", PointerMove(0, 0, 0, ElementOrigin(&WebElement{})))}},
		{"invalid button", []ActionSequence{MouseActions("mouse", PointerDown(5))}},
		{"negative pause", []ActionSequence{MouseActions("mouse", Pause(-time.Millisecond))}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			transport.command = ""
			if _, err := client.PerformActions(test.sequences...); err == nil {
				t.Fatal("expected validation error")
			}
			if transport.command != "" {
				t.Fatal("transport was called after validation failed")
			}
		})
	}
}

func TestPerformActionsPropagatesTransportError(t *testing.T) {
	expected := errors.New("transport failed")
	client := NewClient()
	client.Transport(&recordingTransport{err: expected})

	_, err := client.PerformActions(MouseActions("mouse", Pause(0)))
	if !errors.Is(err, expected) {
		t.Fatalf("expected transport error, got %v", err)
	}
}
