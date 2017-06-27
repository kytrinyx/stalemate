package main

import (
	"database/sql"
	"fmt"
	"time"
)

var (
	stmtFindMaintainer   *sql.Stmt
	stmtCreateMaintainer *sql.Stmt
	stmtDeleteMaintainer *sql.Stmt
)

type Maintainer struct {
	ID             int
	InstallationID int
	Account        string
	Repository     string
	GitHubTeamID   int
	GitHubUserID   int
	Username       string
	CreatedAt      *time.Time
	UpdatedAt      *time.Time
}

func FindMaintainer(account, repository string, githubID int) (*Maintainer, error) {
	stmt, err := findMaintainerStmt()
	if err != nil {
		return nil, err
	}

	t := &Maintainer{}
	row := stmt.QueryRow(account, repository, githubID)
	err = row.Scan(&t.ID, &t.InstallationID, &t.Account, &t.Repository, &t.GitHubTeamID, &t.GitHubUserID, &t.Username, &t.CreatedAt, &t.UpdatedAt)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	}
	return t, nil
}

func (t *Maintainer) Create() error {
	stmt, err := createMaintainerStmt()
	if err != nil {
		return err
	}

	ts := time.Now()
	row := stmt.QueryRow(t.InstallationID, t.Account, t.Repository, t.GitHubTeamID, t.GitHubUserID, t.Username, ts, ts)
	return row.Scan(&t.ID)
}

func (t *Maintainer) Delete() error {
	if t.Account == "" || t.Repository == "" || t.GitHubTeamID == 0 || t.GitHubUserID == 0 {
		s := "unable to delete maintainer %s (id: %d, team: %d) from %s/%s - insufficient identifying information"
		return fmt.Errorf(s, t.Username, t.GitHubUserID, t.GitHubTeamID, t.Account, t.Repository)
	}

	stmt, err := deleteMaintainerStmt()
	if err != nil {
		return err
	}

	_, err = stmt.Exec(t.Account, t.Repository, t.GitHubTeamID, t.GitHubUserID)
	return err
}

func findMaintainerStmt() (*sql.Stmt, error) {
	if stmtFindMaintainer != nil {
		return stmtFindMaintainer, nil
	}

	s := `
	SELECT id, installation_id, account, repository, github_team_id, github_user_id, username, created_at, updated_at
	FROM maintainers
	WHERE account=$1 AND repository=$2 AND github_user_id=$3
	LIMIT 1
	`
	stmt, err := db.Prepare(s)
	if err != nil {
		return nil, err
	}

	stmtFindMaintainer = stmt
	return stmt, nil
}

func createMaintainerStmt() (*sql.Stmt, error) {
	if stmtCreateMaintainer != nil {
		return stmtCreateMaintainer, nil
	}

	s := `
	INSERT INTO maintainers (installation_id, account, repository, github_team_id, github_user_id, username, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id`
	stmt, err := db.Prepare(s)
	if err != nil {
		return nil, err
	}

	stmtCreateMaintainer = stmt
	return stmt, nil
}

func deleteMaintainerStmt() (*sql.Stmt, error) {
	if stmtDeleteMaintainer != nil {
		return stmtDeleteMaintainer, nil
	}

	stmt, err := db.Prepare(`DELETE FROM maintainers WHERE account=$1 AND repository=$2 AND github_team_id=$3 AND github_user_id=$4`)
	if err != nil {
		return nil, err
	}

	stmtDeleteMaintainer = stmt
	return stmt, nil
}
