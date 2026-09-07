package marionette_client

import "errors"

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
		v, e := f.FindElement(by, value)
		if e != nil {
			var driverError *DriverError
			if errors.As(e, &driverError) && driverError.ErrorType == "no such element" {
				return true, nil, nil
			}
			return false, nil, e
		}

		return false, v, nil
	}
}
