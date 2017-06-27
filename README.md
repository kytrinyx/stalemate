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

## Run the App

go install ./... && stalemate

The `install` command installs it to `$GOPATH/bin/stalemate` by default. That will need to be on your path.

### Localtunnel

Install per instructions on https://localtunnel.me

Start localtunnel with the same port that the app is running on:

The configured subdomain should be some random string of characters so that there are
no collisions. It's also whatever you put in your settings when you created the GitHub App.

```
$ lt --port 9090 --subdomain=$CONFIGURED_SUBDOMAIN
```
