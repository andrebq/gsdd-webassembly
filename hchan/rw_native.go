//+build !js

package hchan

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	httpCli = http.DefaultClient
)

func readChan(target *url.URL) (string, error) {
	return post(target, "")
}

func writeChan(target *url.URL, data string) error {
	_, err := post(target, data)
	return err
}

func post(target *url.URL, data string) (string, error) {
	res, err := httpCli.Post(target.String(), "application/json; encoding=utf-8", bytes.NewBufferString(data))
	if err != nil {
		return "", err
	}
	body, err := readAndCloseBody(res)
	if err != nil {
		return "", err
	}

	if !isSuccess(res) {
		return "", fmt.Errorf("unexpected status code %v: body: %s", res.StatusCode, body)
	}
	return body, nil
}

func isSuccess(res *http.Response) bool {
	return res.StatusCode >= 200 && res.StatusCode < 300
}

func readAndCloseBody(res *http.Response) (string, error) {
	defer res.Body.Close()
	buf, err := ioutil.ReadAll(res.Body)
	return string(buf), err
}
