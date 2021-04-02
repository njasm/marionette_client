package marionette_client

import (
	"encoding/json"
	"fmt"
)

type Point struct {
	X float32
	Y float32
}

type Size struct {
	Width  float64
	Height float64
}

type WindowRect struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

type ElementRect struct {
	Point
	Size
}

type WebElement struct {
	id string //`json:"element-6066-11e4-a52e-4f735466cecf"`
	c  *Client
}

func (e *WebElement) Id() string {
	return e.id
}

func (e *WebElement) GetActiveElement() (*WebElement, error) {
	return getActiveElement(e.c)
}

func (e *WebElement) FindElement(by By, value string) (*WebElement, error) {
	return findElement(e.c, by, value, &e.id)
}

func (e *WebElement) FindElements(by By, value string) ([]*WebElement, error) {
	return findElements(e.c, by, value, &e.id)
}

func (e *WebElement) Enabled() bool {
	return isElementEnabled(e.c, e.id)
}

func (e *WebElement) Selected() bool {
	return isElementSelected(e.c, e.id)
}

func (e *WebElement) Displayed() bool {
	return isElementDisplayed(e.c, e.id)
}

func (e *WebElement) TagName() string {
	return getElementTagName(e.c, e.id)
}

func (e *WebElement) Text() string {
	return getElementText(e.c, e.id)
}

func (e *WebElement) Attribute(name string) string {
	return getElementAttribute(e.c, e.id, name)
}

func (e *WebElement) CssValue(property string) string {
	return getElementCssPropertyValue(e.c, e.id, property)
}

func (e *WebElement) Rect() (*ElementRect, error) {
	return getElementRect(e.c, e.id)
}

func (e *WebElement) Click() {
	clickElement(e.c, e.id)
}

func (e *WebElement) SendKeys(keys string) error {
	return sendKeysToElement(e.c, e.id, keys)
}

func (e *WebElement) Clear() {
	clearElement(e.c, e.id)
}

func (e *WebElement) Location() (*Point, error) {
	r, err := getElementRect(e.c, e.id)
	if err != nil {
		return nil, err
	}

	return &r.Point, nil
}

func (e *WebElement) Size() (*Size, error) {
	r, err := getElementRect(e.c, e.id)
	if err != nil {
		return nil, err
	}

	return &r.Size, nil
}

func (e *WebElement) Screenshot() (string, error) {
	id := e.Id()
	return takeScreenshot(e.c, &id)
}

func (e *WebElement) UnmarshalJSON(data []byte) error {
	var d map[string]map[string]string
	err := json.Unmarshal(data, &d)
	if err != nil {
		return err
	}

	if newId, ok := d["value"][WEBDRIVER_ELEMENT_KEY]; ok {
		e.id = newId
		return nil
	}

	return DriverError{
		ErrorType: "WebDriverElementKey",
		Message: fmt.Sprintf("key %v expected in response but not found", WEBDRIVER_ELEMENT_KEY),
		Stacktrace: nil,
	}
}
