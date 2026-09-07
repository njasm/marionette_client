package marionette_client

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

type PointerOrigin struct {
	kind    string
	element *WebElement
}

func ViewportOrigin() PointerOrigin {
	return PointerOrigin{kind: "viewport"}
}

func PointerOriginCurrent() PointerOrigin {
	return PointerOrigin{kind: "pointer"}
}

func ElementOrigin(element *WebElement) PointerOrigin {
	return PointerOrigin{kind: "element", element: element}
}

type PointerAction struct {
	typeName string
	x        int
	y        int
	duration time.Duration
	origin   PointerOrigin
	button   int
}

func PointerMove(x, y int, duration time.Duration, origin PointerOrigin) PointerAction {
	return PointerAction{typeName: "pointerMove", x: x, y: y, duration: duration, origin: origin}
}

func PointerDown(button int) PointerAction {
	return PointerAction{typeName: "pointerDown", button: button}
}

func PointerUp(button int) PointerAction {
	return PointerAction{typeName: "pointerUp", button: button}
}

func Pause(duration time.Duration) PointerAction {
	return PointerAction{typeName: "pause", duration: duration}
}

type ActionSequence struct {
	id             string
	typeName       string
	pointerActions []PointerAction
	keyActions     []KeyAction
}

func MouseActions(id string, actions ...PointerAction) ActionSequence {
	return ActionSequence{id: id, typeName: "pointer", pointerActions: actions}
}

type KeyAction struct {
	typeName string
	value    string
	duration time.Duration
}

func KeyDownAction(value string) KeyAction {
	return KeyAction{typeName: "keyDown", value: value}
}

func KeyUpAction(value string) KeyAction {
	return KeyAction{typeName: "keyUp", value: value}
}

func KeyboardPause(duration time.Duration) KeyAction {
	return KeyAction{typeName: "pause", duration: duration}
}

func KeyboardActions(id string, actions ...KeyAction) ActionSequence {
	return ActionSequence{id: id, typeName: "key", keyActions: actions}
}

func (a PointerAction) MarshalJSON() ([]byte, error) {
	value := map[string]any{"type": a.typeName}
	switch a.typeName {
	case "pointerMove":
		value["x"] = a.x
		value["y"] = a.y
		value["duration"] = a.duration.Milliseconds()
		switch a.origin.kind {
		case "viewport", "pointer":
			value["origin"] = a.origin.kind
		case "element":
			if a.origin.element != nil {
				value["origin"] = map[string]string{WebdriverElementKey: a.origin.element.id}
			}
		}
	case "pointerDown", "pointerUp":
		value["button"] = a.button
	case "pause":
		value["duration"] = a.duration.Milliseconds()
	}
	return json.Marshal(value)
}

func (s ActionSequence) MarshalJSON() ([]byte, error) {
	value := map[string]any{"type": s.typeName, "id": s.id}
	if s.typeName == "pointer" {
		value["parameters"] = map[string]string{"pointerType": "mouse"}
		value["actions"] = s.pointerActions
	} else {
		value["actions"] = s.keyActions
	}
	return json.Marshal(value)
}

func (a KeyAction) MarshalJSON() ([]byte, error) {
	value := map[string]any{"type": a.typeName}
	if a.typeName == "pause" {
		value["duration"] = a.duration.Milliseconds()
	} else {
		value["value"] = a.value
	}
	return json.Marshal(value)
}

func (c *Client) PerformActions(actions ...ActionSequence) (*Response, error) {
	if len(actions) == 0 {
		return nil, errors.New("at least one action sequence is required")
	}

	ids := make(map[string]struct{}, len(actions))
	for sequenceIndex, sequence := range actions {
		if strings.TrimSpace(sequence.id) == "" {
			return nil, fmt.Errorf("action sequence %d has an empty id", sequenceIndex)
		}
		if _, exists := ids[sequence.id]; exists {
			return nil, fmt.Errorf("action sequence id %q is duplicated", sequence.id)
		}
		ids[sequence.id] = struct{}{}
		if sequence.actionCount() == 0 {
			return nil, fmt.Errorf("action sequence %q has no actions", sequence.id)
		}
		switch sequence.typeName {
		case "pointer":
			for actionIndex, action := range sequence.pointerActions {
				if err := validatePointerAction(action); err != nil {
					return nil, fmt.Errorf("action sequence %q action %d: %w", sequence.id, actionIndex, err)
				}
			}
		case "key":
			for actionIndex, action := range sequence.keyActions {
				if err := validateKeyAction(action); err != nil {
					return nil, fmt.Errorf("action sequence %q action %d: %w", sequence.id, actionIndex, err)
				}
			}
		default:
			return nil, fmt.Errorf("action sequence %q has an invalid type", sequence.id)
		}
	}

	return c.transport.Send("WebDriver:PerformActions", map[string]any{"actions": actions})
}

func (s ActionSequence) actionCount() int {
	if s.typeName == "pointer" {
		return len(s.pointerActions)
	}
	return len(s.keyActions)
}

// ReleaseActions releases all depressed input-source actions and clears their state.
func (c *Client) ReleaseActions() (*Response, error) {
	return c.transport.Send("WebDriver:ReleaseActions", nil)
}

func validatePointerAction(action PointerAction) error {
	switch action.typeName {
	case "pointerMove":
		if action.duration < 0 {
			return errors.New("duration cannot be negative")
		}
		switch action.origin.kind {
		case "viewport", "pointer":
			return nil
		case "element":
			if action.origin.element == nil {
				return errors.New("element origin cannot be nil")
			}
			if action.origin.element.id == "" {
				return errors.New("element origin has an empty id")
			}
			return nil
		default:
			return errors.New("pointer move origin is invalid")
		}
	case "pointerDown", "pointerUp":
		if action.button < 0 || action.button > 4 {
			return fmt.Errorf("mouse button %d is outside the range 0..4", action.button)
		}
		return nil
	case "pause":
		if action.duration < 0 {
			return errors.New("duration cannot be negative")
		}
		return nil
	default:
		return errors.New("pointer action type is invalid")
	}
}

func validateKeyAction(action KeyAction) error {
	switch action.typeName {
	case "keyDown", "keyUp":
		if !utf8.ValidString(action.value) || utf8.RuneCountInString(action.value) != 1 {
			return errors.New("key value must contain exactly one Unicode code point")
		}
		return nil
	case "pause":
		if action.duration < 0 {
			return errors.New("duration cannot be negative")
		}
		return nil
	default:
		return errors.New("key action type is invalid")
	}
}
