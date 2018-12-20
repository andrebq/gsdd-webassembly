package main

import (
	"fmt"
	"time"
)

func main() {
	Render(Div().Add(Text("Starting.....")))
	go func() {
		for {
			time.Sleep(time.Second)

			Render(
				Div().Add(
					Text("Time now:"),
					Text(fmt.Sprintf(time.Now().Format(time.RFC3339))),
				),
			)
		}
	}()
	select {}
}
