# stalemate
GitHub app to catch unhandled issues and PRs before they go stale.

## Migrations

We're using https://bitbucket.org/liamstask/goose for database migrations.

Run migrations:

```
$ go get bitbucket.org/liamstask/goose/cmd/goose
$ createuser stalemate
$ createdb stalemate_development
$ goose up
```

Create a new migration with:

```
$ goose create NameOfMigration sql
```
