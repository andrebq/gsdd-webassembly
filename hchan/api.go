package hchan

import (
	"encoding/json"
	"net/url"
)

// Read one entry for the given url
func Read(urlStr string) (string, error) {
	target, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	return readChan(target)
}

// Write one entry to the given url
func Write(urlStr string, data interface{}) error {
	target, err := url.Parse(urlStr)
	if err != nil {
		return err
	}
	buf, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return writeChan(target, string(buf))
}
