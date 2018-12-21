package main

import (
	"fmt"
	"sync"
	"syscall/js"
	"time"

	"github.com/andrebq/gsdd-webassembly/hchan"
)

type (
	applicationState struct {
		sync.Mutex
		Time     time.Time
		LastPing time.Time
		LastPong time.Time
	}

	stateMutator func(*applicationState)
)

var (
	globalState applicationState
)

func (as *applicationState) Update(mutation stateMutator) {
	as.Lock()
	defer as.Unlock()
	mutation(as)

	renderState(as)
}

func renderState(state *applicationState) {
	rootDiv := Div().Add(
		Div().Add(
			Text("Time now: "),
			Text(fmt.Sprintf(state.Time.Format(time.RFC3339))),
		),
	)

	if !state.LastPing.IsZero() {
		rootDiv.Add(P(Text("Last Ping: "), Text(state.LastPing.Format(time.RFC3339))))
	}
	if !state.LastPong.IsZero() {
		rootDiv.Add(P(Text("Last Pong: "), Text(state.LastPong.Format(time.RFC3339))))
	}

	Render(rootDiv)
}

func main() {
	go func() {
		updateTime := func(s *applicationState) {
			s.Time = time.Now()
		}

		for {
			globalState.Update(updateTime)
			time.Sleep(time.Second)
		}
	}()

	type pingPong struct {
		ID string
	}

	go func() {
		makeUpdatePing := func(t time.Time) stateMutator {
			return func(s *applicationState) {
				s.LastPing = t
			}
		}
		for {
			ping := pingPong{ID: "blah"}
			err := hchan.Write("http://localhost:8082/write/ping", ping)
			if err != nil {
				consoleErr("ping", err)
				continue
			}
			globalState.Update(makeUpdatePing(time.Now()))
			time.Sleep(time.Second * 2)
		}
	}()

	go func() {
		makeUpdatePong := func(t time.Time) stateMutator {
			return func(s *applicationState) {
				s.LastPong = t
			}
		}
		for {
			var pong pingPong
			err := hchan.Read(&pong, "http://localhost:8082/read/ping")
			if err != nil {
				consoleErr("pong", err)
				continue
			}
			globalState.Update(makeUpdatePong(time.Now()))
		}
	}()
	select {}
}

func consoleErr(val ...interface{}) {
	if len(val) == 0 {
		return
	}
	msg := make([]interface{}, 0, len(val))
	for _, v := range val {
		msg = append(msg, fmt.Sprintf("%v", v))
	}
	js.Global().Get("console").Call("error", msg...)
}
