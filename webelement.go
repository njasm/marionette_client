package marionette_client

import (
	"encoding/json"
)

type WebElement struct {
	id string //`json:"element-6066-11e4-a52e-4f735466cecf"`
	c  *Client
}

func (e *WebElement) Id() string {
	return e.id
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

func (e *WebElement) Rect() map[string]interface{} {
	return getElementRect(e.c, e.id)
}

func (e *WebElement) Click() {
	clickElement(e.c, e.id)
}

func (e *WebElement) SendKeys(keys string) {
	sendKeysToElement(e.c, e.id, keys)
}

func (e *WebElement) Clear() {
	clearElement(e.c, e.id)
}

func (e *WebElement) UnmarshalJSON(data []byte) error {
	var d map[string]map[string]string
	err := json.Unmarshal([]byte(data), &d)
	if err != nil {
		return err
	}

	e.id = d["value"][WEBDRIVER_ELEMENT_KEY]

	return nil
}
