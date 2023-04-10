package main

import (
	"flag"
)

var flagServerAddress string

func parseFlags() {
	flag.StringVar(&flagServerAddress, "a", "localhost:8080", "server address and port to listen")

	flag.Parse()
}
