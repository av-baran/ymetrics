package main

import (
	"flag"
	"log"
	"os"
)

var flagServerAddress string

func parseFlags() {
	flagSet := flag.NewFlagSet("main", flag.ContinueOnError)
	flagSet.StringVar(&flagServerAddress, "a", "localhost:8080", "server address and port to listen")

	if envServerAddress := os.Getenv("ADDRESS"); envServerAddress != "" {
		flagServerAddress = envServerAddress
		log.Printf("Set ADDRESS=%v from ENV", envServerAddress)
	}

	flagSet.Parse(os.Args)
}
