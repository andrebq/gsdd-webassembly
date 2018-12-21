package main

import (
	"flag"
	"net/http"

	"github.com/rs/cors"

	hchan_server "github.com/andrebq/gsdd-webassembly/hchan/server"
	"github.com/sirupsen/logrus"
)

func main() {
	h := flag.Bool("h", false, "Show help")
	bind := flag.String("bind", ":8081", "Bind address to use")

	flag.Parse()

	if *h {
		flag.Usage()
	}

	handler := hchan_server.MakeRWHandler()
	handler = cors.Default().Handler(handler)

	logrus.WithField("bind", *bind).Info("starting server...")
	if err := http.ListenAndServe(*bind, handler); err != nil {
		logrus.WithError(err).Error("error starting server")
	}
}
