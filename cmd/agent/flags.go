package main

import (
	"flag"
	"log"
	"os"
	"strconv"
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

	if err := flagSet.Parse(os.Args[1:]); err != nil {
		log.Print(err.Error())
	}

	if envServerAddress := os.Getenv("ADDRESS"); envServerAddress != "" {
		flagServerAddress = envServerAddress
		log.Printf("Set ADDRESS=%v from ENV", envServerAddress)
	}
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		if r, err := strconv.Atoi(envReportInterval); err == nil {
			flagReportInterval = r
			log.Printf("Set REPORT_INTERVAL=%v from ENV", r)
		}
		log.Printf("Invalid value in ENV variable REPORT_INTERVAL=%v", envReportInterval)
	}
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		if r, err := strconv.Atoi(envPollInterval); err == nil {
			flagReportInterval = r
			log.Printf("Set POLL_INTERVAL=%v from ENV", r)
		}
		log.Printf("Invalid value in ENV variable POLL_INTERVAL=%v", envPollInterval)
	}
}
