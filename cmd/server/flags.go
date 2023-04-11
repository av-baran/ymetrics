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
	if err := flagSet.Parse(os.Args[1:]); err != nil {
		log.Print(err.Error())
	}
	log.Printf("flagServerAddress: %v, %v", flagServerAddress, &flagServerAddress)

	if envServerAddress := os.Getenv("ADDRESS"); envServerAddress != "" {
		flagServerAddress = envServerAddress
		log.Printf("Set ADDRESS=%v from ENV", envServerAddress)
	}

}
