package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/health", WithHeaders(HandleHealth))
	http.HandleFunc("/", WithHeaders(HandleResources))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
