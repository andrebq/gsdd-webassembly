//+build !js

package hchan

import (
	"net/url"
)

func readChan(target *url.URL) (string, error) {
	panic("not implemented - bad developer")
}

func writeChan(target *url.URL, data string) error {
	panic("not implemented - bad developer")
}
