package main

import (
	"context"
	"net/http"
	"os"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/github"
)

func cli(installationID int) (*github.Client, error) {
	itr, err := ghinstallation.New(http.DefaultTransport, AppID, installationID, []byte(os.Getenv("STALEMATE_PRIVATE_KEY")))
	if err != nil {
		return nil, err
	}
	return github.NewClient(&http.Client{Transport: itr}), nil
}

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
				InstallationID: *event.Installation.ID,
				Account:        *event.Org.Login,
				Repository:     *event.Repo.Name,
				GitHubTeamID:   *event.Team.ID,
				GitHubUserID:   *member.ID,
				Username:       *member.Login,
			}
			if err := maintainer.Create(); err != nil {
				return err
			}
		}
	}
	return nil
}
