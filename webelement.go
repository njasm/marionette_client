package marionette_client

import (
	"encoding/json"
)

type webElement struct {
	id string //`json:"element-6066-11e4-a52e-4f735466cecf"`
	c  *Client
}

func (e *webElement) Id() string {
	return e.id
}

func (e *webElement) FindElement(by By, value string) (*webElement, error) {
	return e.c.FindElement(by, value, &e.id)
}

func (e *webElement) FindElements(by By, value string) ([]*webElement, error) {
	return findElements(e.c, by, value, &e.id)
}

func (e *webElement) IsEnabled() bool {
	return isElementEnabled(e.c, e.id)
}

func (e *webElement) IsSelected() bool {
	return isElementSelected(e.c, e.id)
}

func (e *webElement) IsDisplayed() bool {
	return isElementDisplayed(e.c, e.id)
}

func (e *webElement) TagName() string {
	return getElementTagName(e.c, e.id)
}

func (e *webElement) Text() string {
	return getElementText(e.c, e.id)
}

func (e *webElement) Attribute(name string) string {
	return getElementAttribute(e.c, e.id, name)
}

func (e *webElement) CssValue(property string) string {
	return getElementCssPropertyValue(e.c, e.id, property)
}

func (e *webElement) Rect() map[string]interface{} {
	return getElementRect(e.c, e.id)
}

func (e *webElement) Click() {
	clickElement(e.c, e.id)
}

func (e *webElement) SendKeys(keys string) {
	sendKeysToElement(e.c, e.id, keys)
}

func (e *webElement) Clear() {
	clearElement(e.c, e.id)
}

func (e *webElement) UnmarshalJSON(data []byte) error {
	var d map[string]map[string]string
	err := json.Unmarshal([]byte(data), &d)
	if err != nil {
		return err
	}

	e.id = d["value"][WEBDRIVER_ELEMENT_KEY]

	return nil
}
