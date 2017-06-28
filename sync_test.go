package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/google/go-github/github"
)

func init() {
	InitDB(psqlInfo)
}

func TestSyncTeamEvent(t *testing.T) {

	team := &Team{
		Account:      "tornerose",
		Repository:   "sleep100",
		GitHubTeamID: 2378504,
	}
	defer team.Delete()

	var event github.TeamEvent

	event = loadTeamEvent(t, "./fixtures/team-created.json")
	if _, err := SyncTeamEvent(event); err != nil {
		t.Error(err)
	}

	if countTeamsForAccount(t, team.Account) != 0 {
		t.Error("Should not have created the team (no repo yet).")
	}

	event = loadTeamEvent(t, "./fixtures/team-added-to-repo.json")
	if _, err := SyncTeamEvent(event); err != nil {
		t.Error(err)
	}

	if countTeamsForAccount(t, team.Account) != 0 {
		t.Error("Should not have created the team (permissions are read-only).")
	}

	event = loadTeamEvent(t, "./fixtures/team-edited-permissions-push.json")
	if _, err := SyncTeamEvent(event); err != nil {
		t.Error(err)
	}

	target, err := FindTeam(team.Account, team.Repository, team.GitHubTeamID)
	if err != nil {
		t.Fatal(err)
	}
	if target.Name != "sleepers" {
		t.Error("Expected team to be created.")
	}
}

func loadTeamEvent(t *testing.T, filename string) github.TeamEvent {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	var event github.TeamEvent
	if err := json.Unmarshal(raw, &event); err != nil {
		t.Fatal(err)
	}
	return event
}

func countTeamsForAccount(t *testing.T, name string) int {
	var count int
	if err := db.QueryRow("SELECT COUNT(id) FROM teams WHERE account=$1", name).Scan(&count); err != nil {
		t.Fatal(err)
	}
	return count
}
