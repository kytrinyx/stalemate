package main

import "testing"

func init() {
	InitDB(psqlInfo)
}

func TestTeamCRUD(t *testing.T) {
	// It should not blow up if it's not found.
	if _, err := FindTeam("bogus-org", "bogus-repo", 2); err != nil {
		t.Error(err)
	}

	team1 := &Team{
		InstallationID: 1,
		Account:        "bogus-org",
		Repository:     "bogus-repo",
		GitHubTeamID:   2,
		Name:           "bogus-team",
	}
	defer team1.Delete()

	// It gets created.
	if err := team1.Create(); err != nil {
		t.Fatal(err)
	}
	if team1.ID == 0 {
		t.Fatal("Team has no ID. Did it get created?")
	}

	// It can't create duplicates.
	if err := team1.Create(); err == nil {
		t.Errorf("It didn't blow up when saving the same team twice.")
	}

	team2, err := FindTeam(team1.Account, team1.Repository, team1.GitHubTeamID)
	if err != nil {
		t.Error(err)
	}
	if team2.Name != team1.Name {
		t.Errorf("Team did not save all the data.")
	}

	team3 := &Team{
		Account:      team1.Account,
		Repository:   team1.Repository,
		GitHubTeamID: team1.GitHubTeamID,
	}
	if err := team3.Delete(); err != nil {
		t.Fatal(err)
	}

	team4, err := FindTeam(team3.Account, team3.Repository, team3.GitHubTeamID)
	if err != nil {
		t.Error(err)
	}
	if team4 != nil {
		t.Error("Team did not get deleted.")
	}
}
