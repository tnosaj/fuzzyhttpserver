package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tnosaj/fuzzyhttpserver/server"
)

func evaluateInputs() (server.Settings, error) {
	var s server.Settings

	flag.BoolVar(&s.Debug, "v", false, "Enable verbose debugging output")

	flag.StringVar(&s.Port, "p", "8081", "Starts server on this port")

	flag.IntVar(&s.Timeout, "t", 2, "Timeout in seconds for a backend answer")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s: [flags] command [command args…]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	return s, nil
}
