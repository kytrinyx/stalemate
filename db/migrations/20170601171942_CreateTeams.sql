
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE teams (
	    id serial PRIMARY KEY,
	    installation_id integer,
	    account character varying(255),
	    repository character varying(255),
	    github_team_id integer,
	    name character varying(255),
	    created_at timestamp without time zone NOT NULL,
	    updated_at timestamp without time zone NOT NULL
);

CREATE INDEX teams_installation ON teams (installation_id);
CREATE UNIQUE INDEX teams_repo ON teams (github_team_id, account, repository);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE teams;
