package marionette_client

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestWebDriverKeys(t *testing.T) {
	keys := []string{
		KeyNull, KeyCancel, KeyHelp, KeyBackspace, KeyTab, KeyClear, KeyReturn, KeyEnter,
		KeyShift, KeyControl, KeyAlt, KeyPause, KeyEscape, KeySpace, KeyPageUp, KeyPageDown,
		KeyEnd, KeyHome, KeyLeft, KeyUp, KeyRight, KeyDown, KeyInsert, KeyDelete, KeySemicolon,
		KeyEquals, KeyNumpad0, KeyNumpad1, KeyNumpad2, KeyNumpad3, KeyNumpad4, KeyNumpad5,
		KeyNumpad6, KeyNumpad7, KeyNumpad8, KeyNumpad9, KeyMultiply, KeyAdd, KeySeparator,
		KeySubtract, KeyDecimal, KeyDivide, KeyF1, KeyF2, KeyF3, KeyF4, KeyF5, KeyF6, KeyF7,
		KeyF8, KeyF9, KeyF10, KeyF11, KeyF12, KeyMeta, KeyZenkakuHankaku, KeyRightShift,
		KeyRightControl, KeyRightAlt, KeyRightMeta,
	}
	expected := []string{
		"\ue000", "\ue001", "\ue002", "\ue003", "\ue004", "\ue005", "\ue006", "\ue007",
		"\ue008", "\ue009", "\ue00a", "\ue00b", "\ue00c", "\ue00d", "\ue00e", "\ue00f",
		"\ue010", "\ue011", "\ue012", "\ue013", "\ue014", "\ue015", "\ue016", "\ue017",
		"\ue018", "\ue019", "\ue01a", "\ue01b", "\ue01c", "\ue01d", "\ue01e", "\ue01f",
		"\ue020", "\ue021", "\ue022", "\ue023", "\ue024", "\ue025", "\ue026", "\ue027",
		"\ue028", "\ue029", "\ue031", "\ue032", "\ue033", "\ue034", "\ue035", "\ue036",
		"\ue037", "\ue038", "\ue039", "\ue03a", "\ue03b", "\ue03c", "\ue03d", "\ue040",
		"\ue050", "\ue051", "\ue052", "\ue053",
	}

	if !reflect.DeepEqual(keys, expected) {
		t.Fatalf("special key catalog mismatch\ngot:  %q\nwant: %q", keys, expected)
	}
	if KeyBackSpace != KeyBackspace || KeyLeftShift != KeyShift || KeyLeftControl != KeyControl ||
		KeyLeftAlt != KeyAlt || KeyArrowLeft != KeyLeft || KeyArrowUp != KeyUp ||
		KeyArrowRight != KeyRight || KeyArrowDown != KeyDown || KeyLeftMeta != KeyMeta ||
		KeyCommand != KeyMeta || KeyLeftCommand != KeyCommand || KeyLeftOption != KeyLeftAlt ||
		KeyRightOption != KeyRightAlt {
		t.Fatal("key aliases do not match their canonical values")
	}
}

func TestSendKeysPayload(t *testing.T) {
	transport := &recordingTransport{}
	client := NewClient()
	client.Transport(transport)
	element := &WebElement{id: "element-id", c: client}

	if err := element.SendKeys(KeyControl, "a", KeyNull); err != nil {
		t.Fatal(err)
	}
	if transport.command != "WebDriver:ElementSendKeys" {
		t.Fatalf("unexpected command %q", transport.command)
	}
	expected := map[string]any{"id": "element-id", "text": KeyControl + "a" + KeyNull}
	if !reflect.DeepEqual(transport.values, expected) {
		t.Fatalf("payload mismatch: got %#v, want %#v", transport.values, expected)
	}
}

func TestKeyboardActionsPayload(t *testing.T) {
	transport := &recordingTransport{}
	client := NewClient()
	client.Transport(transport)

	_, err := client.PerformActions(KeyboardActions("keyboard",
		KeyDownAction(KeyControl),
		KeyDownAction("a"),
		KeyUpAction("a"),
		KeyUpAction(KeyControl),
		KeyboardPause(25*time.Millisecond),
	))
	if err != nil {
		t.Fatal(err)
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
	expectedJSON := `{"actions":[{"actions":[{"type":"keyDown","value":"\ue009"},{"type":"keyDown","value":"a"},{"type":"keyUp","value":"a"},{"type":"keyUp","value":"\ue009"},{"duration":25,"type":"pause"}],"id":"keyboard","type":"key"}]}`
	if err := json.Unmarshal([]byte(expectedJSON), &expected); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("payload mismatch\nactual:   %s\nexpected: %s", encoded, expectedJSON)
	}
}

func TestKeyboardActionsValidation(t *testing.T) {
	client := NewClient()
	transport := &recordingTransport{}
	client.Transport(transport)

	tests := []ActionSequence{
		KeyboardActions("keyboard", KeyAction{}),
		KeyboardActions("keyboard", KeyDownAction("")),
		KeyboardActions("keyboard", KeyUpAction("ab")),
		KeyboardActions("keyboard", KeyboardPause(-time.Millisecond)),
	}
	for index, sequence := range tests {
		transport.command = ""
		if _, err := client.PerformActions(sequence); err == nil {
			t.Fatalf("case %d: expected validation error", index)
		}
		if transport.command != "" {
			t.Fatalf("case %d: transport was called", index)
		}
	}
}
