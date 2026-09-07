package marionette_client

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
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
	id      string
	actions []PointerAction
}

func MouseActions(id string, actions ...PointerAction) ActionSequence {
	return ActionSequence{id: id, actions: actions}
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
	return json.Marshal(map[string]any{
		"type":       "pointer",
		"id":         s.id,
		"parameters": map[string]string{"pointerType": "mouse"},
		"actions":    s.actions,
	})
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
		if len(sequence.actions) == 0 {
			return nil, fmt.Errorf("action sequence %q has no actions", sequence.id)
		}
		for actionIndex, action := range sequence.actions {
			if err := validatePointerAction(action); err != nil {
				return nil, fmt.Errorf("action sequence %q action %d: %w", sequence.id, actionIndex, err)
			}
		}
	}

	return c.transport.Send("WebDriver:PerformActions", map[string]any{"actions": actions})
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
