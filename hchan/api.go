package hchan

import (
	"encoding/json"
	"net/url"
)

// Read one entry for the given url
func Read(out interface{}, urlStr string) error {
	target, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	data, err := readChan(target)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), out)
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
