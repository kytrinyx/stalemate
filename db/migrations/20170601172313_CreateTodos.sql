
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE todos (
	    id serial PRIMARY KEY,
	    installation_id integer,
	    account character varying(255),
	    repository character varying(255),
	    github_user_id integer, -- author
	    username character varying(255), -- author
	    url character varying(255), -- of issue or pr
	    title character varying(255), -- of issue or pr
	    waiting_since timestamp without time zone NOT NULL, -- earliest unhandled activity by author
	    created_at timestamp without time zone NOT NULL,
	    updated_at timestamp without time zone NOT NULL
);

CREATE INDEX todos_installation ON todos (installation_id);
CREATE UNIQUE INDEX todos_url ON todos (installation_id);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE todos;
