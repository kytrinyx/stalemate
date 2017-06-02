package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

var (
	stmtFindTeam   *sql.Stmt
	stmtCreateTeam *sql.Stmt
	stmtDeleteTeam *sql.Stmt
)

type Team struct {
	ID             int
	InstallationID int
	Account        string
	Repository     string
	GitHubTeamID   int
	Name           string
	CreatedAt      *time.Time
	UpdatedAt      *time.Time
}

func FindTeam(account, repository string, githubID int) (*Team, error) {
	stmt, err := findTeamStmt()
	if err != nil {
		return nil, err
	}

	t := &Team{}
	row := stmt.QueryRow(account, repository, githubID)
	err = row.Scan(&t.ID, &t.InstallationID, &t.Account, &t.Repository, &t.GitHubTeamID, &t.Name, &t.CreatedAt, &t.UpdatedAt)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	}
	return t, nil
}

func (t *Team) Create() error {
	stmt, err := createTeamStmt()
	if err != nil {
		return err
	}

	ts := time.Now()
	row := stmt.QueryRow(t.InstallationID, t.Account, t.Repository, t.GitHubTeamID, t.Name, ts, ts)
	return row.Scan(&t.ID)
}

func (t *Team) Delete() error {
	if t.Account == "" || t.Repository == "" || t.GitHubTeamID == 0 {
		s := "unable to delete team %s@%s (%d) for repo %s - insufficient identifying information"
		return fmt.Errorf(s, t.Account, t.Name, t.GitHubTeamID, t.Repository)
	}

	stmt, err := deleteTeamStmt()
	if err != nil {
		return err
	}

	_, err = stmt.Exec(t.Account, t.Repository, t.GitHubTeamID)
	return err
}

func findTeamStmt() (*sql.Stmt, error) {
	if stmtFindTeam != nil {
		return stmtFindTeam, nil
	}

	s := `
	SELECT id, installation_id, account, repository, github_team_id, name, created_at, updated_at
	FROM teams
	WHERE account=$1 AND repository=$2 AND github_team_id=$3
	`
	stmt, err := db.Prepare(s)
	if err != nil {
		return nil, err
	}

	stmtFindTeam = stmt
	return stmt, nil
}

func createTeamStmt() (*sql.Stmt, error) {
	if stmtCreateTeam != nil {
		return stmtCreateTeam, nil
	}

	s := `
	INSERT INTO teams (installation_id, account, repository, github_team_id, name, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id`
	stmt, err := db.Prepare(s)
	if err != nil {
		return nil, err
	}

	stmtCreateTeam = stmt
	return stmt, nil
}

func deleteTeamStmt() (*sql.Stmt, error) {
	if stmtDeleteTeam != nil {
		return stmtDeleteTeam, nil
	}

	stmt, err := db.Prepare(`DELETE FROM teams WHERE account=$1 AND repository=$2 AND github_team_id=$3`)
	if err != nil {
		return nil, err
	}

	stmtDeleteTeam = stmt
	return stmt, nil
}
