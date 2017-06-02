package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"github.com/rjz/githubhook"
)

// Hard-coding this for now.
var psqlInfo = "host=localhost port=5432 user=stalemate dbname=stalemate_development sslmode=disable"
var db *sql.DB

func InitDB(dbInfo string) {
	if db != nil {
		return
	}

	var err error
	db, err = sql.Open("postgres", dbInfo)
	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
}

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
	case "team_add":
		// Ignore this one.
	case "team":
		event := github.TeamEvent{}
		if err := json.Unmarshal(hook.Payload, &event); err != nil {
			log.Println(err)
		}
		if err := SyncTeam(event); err != nil {
			log.Println(err)
		}
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
