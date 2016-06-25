package expected

import (
	"errors"
)

func ElementIsPresent(by By, value string, startNodeID *string) func(f *Finder) (bool, *WebElement, error) {
	return func (f Finder) (*WebElement, error) {
		if f == nil {
			return nil, errors.New("Client is nil.")
		}

		if value == "" {
			nil, errors.New("Value is empty.")
		}

		v, e := f.FindElement(by, value)
		if e != nil {
			false, nil, e
		}

		return true, v, e
	}
}

func ElementIsNotPresent(by By, value string, startNodeID *string) func(f *Finder) (bool, *WebElement, error) {
	return func (f Finder) (*WebElement, error) {
		if f == nil {
			return nil, errors.New("Client is nil.")
		}

		if value == "" {
			nil, errors.New("Value is empty.")
		}

		v, e := f.FindElement(by, value)
		if e != nil {
			false, nil, e
		}

		return true, v, e
	}
}


