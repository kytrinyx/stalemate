package main

import "github.com/google/go-github/github"

func SyncTeam(event github.TeamEvent) error {
	if event.Repo == nil {
		return nil
	}
	if event.Repo.Permissions == nil {
		return nil
	}
	if (*event.Repo.Permissions)["push"] {

		team, err := FindTeam(*event.Org.Login, *event.Repo.Name, *event.Team.ID)
		if err != nil {
			return err
		}
		if team != nil {
			return nil
		}

		team = &Team{
			InstallationID: *event.Installation.ID,
			Account:        *event.Org.Login,
			Repository:     *event.Repo.Name,
			GitHubTeamID:   *event.Team.ID,
			Name:           *event.Team.Name,
		}
		if err := team.Create(); err != nil {
			return err
		}
	}

	return nil
}
