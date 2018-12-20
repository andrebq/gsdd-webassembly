package main

import (
	"fmt"
	"time"

	"github.com/andrebq/gsdd-webassembly/hchan"
)

func main() {
	value, err := hchan.Read("http://example.com/some/random/channel/name")
	go func() {
		for {
			time.Sleep(time.Second)

			Render(
				Div().Add(
					Div().Add(
						Text("Time now:"),
						Text(fmt.Sprintf(time.Now().Format(time.RFC3339))),
					),
					Div().Add(Text(
						fmt.Sprintf("Value: %v / Err: %v", value, err),
					))))
		}
	}()
	select {}
}
