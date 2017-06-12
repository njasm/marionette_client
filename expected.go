package marionette_client

func ElementIsPresent(by By, value string) func(f Finder) (bool, *WebElement, error) {
	return func(f Finder) (bool, *WebElement, error) {
		result := true
		v, e := f.FindElement(by, value)
		if e != nil || v == nil {
			result = false
		}

		return result, v, e
	}
}

func ElementIsNotPresent(by By, value string) func(f Finder) (bool, *WebElement, error) {
	return func(f Finder) (bool, *WebElement, error) {
		result := false
		v, e := f.FindElement(by, value)
		if e != nil {
			result = true
		}

		return result, v, e
	}
}
