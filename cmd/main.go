package main

import (
	"log"
	"net/http"

	"github.com/mbsabath/promsd/api"
)

func main() {
	handler := api.NewSdHandler()
	err := http.ListenAndServe("localhost:6969", handler)
	log.Fatal(err)
}
