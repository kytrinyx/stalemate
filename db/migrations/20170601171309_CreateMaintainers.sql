
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE maintainers (
	    id serial PRIMARY KEY,
	    installation_id integer,
	    account character varying(255),
	    repository character varying(255),
	    github_team_id integer,
	    github_user_id integer,
	    username character varying(255),
	    created_at timestamp without time zone NOT NULL,
	    updated_at timestamp without time zone NOT NULL
);

CREATE INDEX maintainers_installation ON maintainers (installation_id);
CREATE UNIQUE INDEX maintainer_team_repo ON maintainers (installation_id, github_team_id, github_user_id, account, repository);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE maintainers;
