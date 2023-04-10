package main

import (
	"flag"
)

var (
	flagServerAddress  string
	flagReportInterval int
	flagPollInterval   int
)

func parseFlags() {
	flag.StringVar(&flagServerAddress, "a", "localhost:8080", "server address and port to listen")
	flag.IntVar(&flagReportInterval, "r", 10, "report interval in seconds")
	flag.IntVar(&flagPollInterval, "p", 2, "poll interval in seconds")

	flag.Parse()
}
