package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mbsabath/promsd/api"
)

var help bool
var hostname string
var port string

func main() {
	flag.BoolVar(&help, "help", false, "Print Usage")
	flag.StringVar(&hostname, "host", "localhost", "Host to listen on (usually localhost or 0.0.0.0)")
	flag.StringVar(&port, "port", "5000", "Port to listen on")

	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(0)
	}

	handler := api.NewSdHandler()
	err := http.ListenAndServe(fmt.Sprintf("%s:%s", hostname, port), handler)
	log.Fatal(err)
}
