# Catena

**The links that associate us.**

[![GoDoc](https://godoc.org/github.com/bbengfort/catena?status.svg)](https://godoc.org/github.com/bbengfort/catena)
[![Go Report Card](https://goreportcard.com/badge/github.com/bbengfort/catena)](https://goreportcard.com/report/github.com/bbengfort/catena)
[![Build Status](https://travis-ci.com/bbengfort/catena.svg?branch=master)](https://travis-ci.com/bbengfort/catena)
[![codecov](https://codecov.io/gh/bbengfort/catena/branch/master/graph/badge.svg)](https://codecov.io/gh/bbengfort/catena)


## Database Migrations

The schema of the database is managed through migration files that can be applied or rolled back to ensure the database version matches the expected version of the server.

> **NOTE**: Currently only PostgreSQL is tested with this method.

To create a migration, create a new SQL file in the `migrations` folder with the format `XXXX_my_migration.sql` where the `XXXX` should be the next migration number in sequence and the text describes the migration. You can use the catena command to do this as long as you're in the project root:

```
$ go run ./cmd/catena --new my revision name
```

In the SQL file you should have the following two comments:

```sql
-- migrate: up
-- insert up migration sql here

-- migrate: down
-- insert down migration sql here
```

All of the SQL under `-- migrate: up` will be run when the migration is applied and all of the SQL under `-- migrate: down` will be run if the migration is rolled back. Ensure that you separate multiple SQL commands with `;` because they will be executed all at once and with newlines removed.

Next generate the migration by running `go generate` in the project root:

```
$ go generate ./...
```

This will generate the migrations code from the SQL files and allow you to apply it with the `catena migrate` command.