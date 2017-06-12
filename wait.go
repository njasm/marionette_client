package marionette_client

import (
	"errors"
	"time"
)

type Waiter struct {
	f Finder
	d time.Duration // Milliseconds
}

type Finder interface {
	FindElement(by By, value string) (*WebElement, error)
	FindElements(by By, value string) ([]*WebElement, error)
}

func Wait(f Finder) *Waiter {
	return &Waiter{f: f, d: time.Duration(1)}
}

func (w *Waiter) For(d time.Duration) *Waiter {
	if d < 0 || d > time.Minute*10 {
		w.d = time.Duration(time.Second)
		return w
	}

	w.d = d
	return w
}

func (w *Waiter) Until(f func(c Finder) (bool, *WebElement, error)) (bool, *WebElement, error) {
	firstRun := true
	delta := time.Now()
	for time.Since(delta) < w.d || firstRun {
		firstRun = false

		ok, value, err := f(w.f)
		if err != nil {
			_, de := err.(DriverError)
			if de {
				return false, nil, err
			}
		}

		if ok {
			return true, value, err
		}

		time.Sleep(time.Second)
	}

	return false, nil, errors.New("Condition never occurred.")
}
