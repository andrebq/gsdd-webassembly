package hchan

import (
	"errors"
	"net/url"
	"syscall/js"
)

func readChan(target *url.URL) (string, error) {
	value := make(chan string)
	err := make(chan error)
	var cb js.Callback
	cb = js.NewCallback(func(args []js.Value) {
		defer close(err)
		defer close(value)
		defer cb.Release()

		if len(args) == 1 {
			// only success
			value <- args[0].String()
		}

		switch {
		case args[1] == js.Null():
			value <- args[0].String()
		default:
			err <- errors.New(args[1].String())
		}
	})
	js.Global().Get("httpChan").Call("read", target.String(), cb)

	select {
	case v := <-value:
		return v, nil
	case e := <-err:
		return "", e
	}
}

func writeChan(target *url.URL, data string) error {
	value := make(chan struct{})
	err := make(chan error)
	var cb js.Callback
	cb = js.NewCallback(func(args []js.Value) {
		defer close(err)
		defer close(value)
		defer cb.Release()

		if len(args) == 1 {
			// only success
			value <- struct{}{}
		}

		switch {
		case args[1] == js.Null():
			value <- struct{}{}
		default:
			err <- errors.New(args[1].String())
		}
	})
	js.Global().Get("httpChan").Call("write", target.String(), data, cb)

	select {
	case <-value:
		return nil
	case e := <-err:
		return e
	}
}
