package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

type Payload struct {
	Action string `json:"action"`
}

func TestFixture(t *testing.T) {
	raw, err := ioutil.ReadFile("./fixtures/installation-created.json")
	if err != nil {
		t.Fatal(err)
	}
	var p Payload
	if err := json.Unmarshal(raw, &p); err != nil {
		t.Fatal(err)
	}
	if p.Action != "created" {
		t.Errorf("Failed to read the fixture. This is a silly test, but whatevs.")
	}
}
