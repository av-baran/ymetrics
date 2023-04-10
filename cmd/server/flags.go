package main

import (
	"flag"
	"os"
)

var flagServerAddress string

func parseFlags() {
	flagSet := flag.NewFlagSet("main", flag.ContinueOnError)
	flagSet.StringVar(&flagServerAddress, "a", "localhost:8080", "server address and port to listen")

	flagSet.Parse(os.Args)
}
