package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

func hello(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "This is stalemate.")
}

func processPayload(rw http.ResponseWriter, req *http.Request) {
	dump, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(dump))
	fmt.Fprintf(rw, "{\"status\": \"ok\"}\n")
}

func main() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/events", processPayload)
	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Fatal(err)
	}
}
