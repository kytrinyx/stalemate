package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"github.com/rjz/githubhook"
)

func hello(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "This is stalemate.")
}

func processPayload(rw http.ResponseWriter, req *http.Request) {
	hook, err := githubhook.Parse([]byte(os.Getenv("STALEMATE_SECRET_TOKEN")), req)
	if err != nil {
		log.Println(err)
		return
	}

	switch hook.Event {
	case "integration_installation":
		event := github.IntegrationInstallationEvent{}
		if err := json.Unmarshal(hook.Payload, &event); err != nil {
			log.Println(err)
			return
		}
		// Echo back the installation part of the payload.
		fmt.Fprintf(rw, event.Installation.String())

	default:
		log.Printf("not handling %s events yet", hook.Event)
	}

}

func main() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/events", processPayload)
	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Fatal(err)
	}
}
