package main

import (
	"net/http"

	"github.com/av-baran/ymetrics/internal/router"
	"github.com/av-baran/ymetrics/internal/storage/memstor"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	repo := memstor.New()
	return http.ListenAndServe("0.0.0.0:8080", router.New(repo))
}
