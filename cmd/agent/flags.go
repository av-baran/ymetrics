package main

import (
	"flag"
	"os"
)

var (
	flagServerAddress  string
	flagReportInterval int
	flagPollInterval   int
)

func parseFlags() {
	flagSet := flag.NewFlagSet("main", flag.ContinueOnError)

	flagSet.StringVar(&flagServerAddress, "a", "localhost:8080", "server address and port to listen")
	flagSet.IntVar(&flagReportInterval, "r", 10, "report interval in seconds")
	flagSet.IntVar(&flagPollInterval, "p", 2, "poll interval in seconds")

	flagSet.Parse(os.Args)
}
