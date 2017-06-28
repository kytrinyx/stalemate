package main

import (
	"context"
	"net/http"
	"os"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/github"
)

func SyncTeam(event github.TeamEvent) error {
	team, err := SyncTeamEvent(event)
	if err != nil {
		return err
	}
	if team == nil {
		return nil
	}

	return SyncTeamMembers(team)
}

func SyncTeamEvent(event github.TeamEvent) (*Team, error) {
	if event.Repo == nil {
		return nil, nil
	}
	if event.Repo.Permissions == nil {
		return nil, nil
	}
	if !(*event.Repo.Permissions)["push"] {
		return nil, nil
	}

	team, err := FindTeam(*event.Org.Login, *event.Repo.Name, *event.Team.ID)
	if err != nil {
		return nil, err
	}
	if team != nil {
		return nil, nil
	}

	team = &Team{
		InstallationID: *event.Installation.ID,
		Account:        *event.Org.Login,
		Repository:     *event.Repo.Name,
		GitHubTeamID:   *event.Team.ID,
		Name:           *event.Team.Name,
	}
	if err := team.Create(); err != nil {
		return nil, err
	}
	return team, nil
}

func SyncTeamMembers(team *Team) error {
	client, err := cli(team.InstallationID)
	if err != nil {
		return err
	}

	opt := &github.OrganizationListTeamMembersOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var members []*github.User
	for {

		users, resp, err := client.Organizations.ListTeamMembers(context.Background(), team.GitHubTeamID, opt)
		if err != nil {
			return err
		}
		members = append(members, users...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	for _, member := range members {
		maintainer := &Maintainer{
			InstallationID: team.InstallationID,
			Account:        team.Account,
			Repository:     team.Repository,
			GitHubTeamID:   team.GitHubTeamID,
			GitHubUserID:   *member.ID,
			Username:       *member.Login,
		}
		if err := maintainer.Create(); err != nil {
			return err
		}
	}
	return nil
}

func cli(installationID int) (*github.Client, error) {
	itr, err := ghinstallation.New(http.DefaultTransport, AppID, installationID, []byte(os.Getenv("STALEMATE_PRIVATE_KEY")))
	if err != nil {
		return nil, err
	}
	return github.NewClient(&http.Client{Transport: itr}), nil
}
