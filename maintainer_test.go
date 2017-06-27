package main

import "testing"

func init() {
	InitDB(psqlInfo)
}

func TestMaintainerCRUD(t *testing.T) {
	// It should not blow up if it's not found.
	if _, err := FindMaintainer("bogus-org", "bogus-repo", 2); err != nil {
		t.Error(err)
	}

	maintainer1 := &Maintainer{
		InstallationID: 1,
		Account:        "bogus-org",
		Repository:     "bogus-repo",
		GitHubTeamID:   100,
		GitHubUserID:   2,
		Username:       "bogus-user",
	}
	defer maintainer1.Delete()

	// It gets created.
	if err := maintainer1.Create(); err != nil {
		t.Fatal(err)
	}
	if maintainer1.ID == 0 {
		t.Fatal("Maintainer 1 has no ID. Did it get created?")
	}

	// It can't create duplicates.
	if err := maintainer1.Create(); err == nil {
		t.Errorf("It didn't blow up when saving the same maintainer twice. It should have.")
	}

	maintainer2, err := FindMaintainer(maintainer1.Account, maintainer1.Repository, maintainer1.GitHubUserID)
	if err != nil {
		t.Error(err)
	}
	if maintainer2.Username != maintainer1.Username {
		t.Errorf("Maintainer did not save all the data.")
	}

	maintainer3 := &Maintainer{
		Account:      maintainer1.Account,
		Repository:   maintainer1.Repository,
		GitHubTeamID: maintainer1.GitHubTeamID,
		GitHubUserID: maintainer1.GitHubUserID,
	}
	if err := maintainer3.Delete(); err != nil {
		t.Fatal(err)
	}

	maintainer4, err := FindMaintainer(maintainer3.Account, maintainer3.Repository, maintainer3.GitHubUserID)
	if err != nil {
		t.Error(err)
	}
	if maintainer4 != nil {
		t.Error("Maintainer did not get deleted.")
	}

	// It can, however, create the same user in a different team, on the same repo.
	// This is because I'm lazy, and I am not going to normalize the data.
	maintainer5 := &Maintainer{
		InstallationID: 1,
		Account:        "bogus-org",
		Repository:     "bogus-repo",
		GitHubTeamID:   200,
		GitHubUserID:   2,
		Username:       "bogus-user",
	}
	defer maintainer5.Delete()

	if err := maintainer5.Create(); err != nil {
		t.Fatal(err)
	}
	if maintainer5.ID == 0 {
		t.Fatal("Maintainer 5 has no ID. Did it get created?")
	}
}
