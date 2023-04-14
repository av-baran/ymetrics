package main

import (
	"flag"
	"log"
	"os"
)

type ArgFlags struct {
	ServerAddress  string
	ReportInterval int
	PollInterval   int
}

func parseFlags() *ArgFlags {
	flags := &ArgFlags{}
	flagSet := flag.NewFlagSet("main", flag.ContinueOnError)

	flagSet.StringVar(&flags.ServerAddress, "a", "localhost:8080", "server address and port to listen")
	flagSet.IntVar(&flags.ReportInterval, "r", 10, "report interval in seconds")
	flagSet.IntVar(&flags.PollInterval, "p", 2, "poll interval in seconds")

	if err := flagSet.Parse(os.Args[1:]); err != nil {
		log.Print(err.Error())
	}

	return flags
}
