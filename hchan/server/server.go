package server

import (
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"sync"
	"time"
)

type (
	chanReg struct {
		sync.Mutex
		items map[string]chan string
	}
)

func (c *chanReg) getOrCreate(chName string) chan string {
	c.Lock()
	defer c.Unlock()
	if c.items == nil {
		c.items = map[string]chan string{}
	}
	val, ok := c.items[chName]
	if !ok {
		val = make(chan string, 100)
		c.items[chName] = val
	}
	return val
}

// MakeRWHandler creates a new read/write handler pair which are capable of reading/writing data
// using a shared channel.
//
// TODO: find a way to handle this, but I'm too lazy to think right now
// Once created a channel will live forever (ie, no close is called)
func MakeRWHandler() http.Handler {
	reg := chanReg{}
	read := http.StripPrefix("/read", makeReadHandler(&reg))
	write := http.StripPrefix("/write", makeWriteHandler(&reg))

	mux := http.NewServeMux()
	mux.Handle("/read/", read)
	mux.Handle("/write/", write)
	return mux
}

func makeReadHandler(r *chanReg) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			http.Error(w, "Bad method", http.StatusMethodNotAllowed)
			return
		}
		chName := path.Clean(req.URL.Path)
		if !validPath(chName) {
			http.Error(w, "Bad path", http.StatusBadRequest)
			return
		}

		ch := r.getOrCreate(chName)
		select {
		case str := <-ch:
			hdr := w.Header()
			hdr.Set("Content-Type", "application/json; encoding=utf-8")
			hdr.Set("Content-Length", strconv.Itoa(len(str)))
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, str)
		case <-time.After(time.Minute):
			http.Error(w, "no data received from the other end", http.StatusGatewayTimeout)
		}
	}
}

func makeWriteHandler(r *chanReg) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			http.Error(w, "Bad method", http.StatusMethodNotAllowed)
			return
		}
		chName := path.Clean(req.URL.Path)
		if !validPath(chName) {
			http.Error(w, "Bad path", http.StatusBadRequest)
			return
		}

		body, err := ioutil.ReadAll(req.Body)
		req.Body.Close()
		if err != nil {
			// only possible cause is something internal
			http.Error(w, "Something right is not wrong", http.StatusInternalServerError)
			return
		}

		ch := r.getOrCreate(chName)
		select {
		case ch <- string(body):
			hdr := w.Header()
			hdr.Set("Content-Type", "application/json; encoding=utf-8")
			w.WriteHeader(http.StatusCreated)
		case <-time.After(time.Minute):
			http.Error(w, "no data received from the other end", http.StatusGatewayTimeout)
		}
	}
}

func validPath(p string) bool {
	return p != "/" && p != "."
}
