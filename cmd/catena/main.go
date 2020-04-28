package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/bbengfort/catena"
	"github.com/bbengfort/catena/migrations"
	"github.com/joho/godotenv"
	"gopkg.in/urfave/cli.v1"

	// use postgres for the test database
	_ "github.com/lib/pq"
)

func main() {
	// Load the .env file if it exists
	godotenv.Load()

	// Instantiate the CLI application
	app := cli.NewApp()
	app.Name = "catena"
	app.Version = catena.PackageVersion
	app.Usage = "catena server and server utilities"
	app.Commands = []cli.Command{
		{
			Name:     "serve",
			Usage:    "run the catena API server",
			Action:   serve,
			Category: "server",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "a, addr",
					Usage:  "the address to bind the catena server on",
					EnvVar: "CATENA_ENDPOINT",
				},
				cli.StringFlag{
					Name:   "D, db",
					Usage:  "the database uri of the catena postgres database",
					EnvVar: "DATABASE_URL",
				},
			},
		},
		{
			Name:     "db:revision",
			Usage:    "print the current migration status of the database",
			Action:   revision,
			Category: "database",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "D, db",
					Usage:  "the database uri of the catena postgres database",
					EnvVar: "DATABASE_URL",
				},
				cli.Int64Flag{
					Name:  "r, revision",
					Usage: "specify a revision to get the status for",
					Value: -1,
				},
			},
		},
		{
			Name:     "db:migrate",
			Usage:    "migrate (or rollback) the database to the latest revision",
			Action:   migrate,
			Category: "database",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "D, db",
					Usage:  "the database uri of the catena postgres database",
					EnvVar: "DATABASE_URL",
				},
				cli.Int64Flag{
					Name:  "r, revision",
					Usage: "specify a revision to migrate (or rollback) the database to",
					Value: -1,
				},
			},
		},
	}

	// Run the program, it should not error
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

//===========================================================================
// Server Commands
//===========================================================================

func serve(c *cli.Context) (err error) {
	return cli.NewExitError("server not implemented yet", 42)
}

//===========================================================================
// Database Commands
//===========================================================================

func revision(c *cli.Context) (err error) {
	var db *sql.DB
	if uri := c.String("db"); uri != "" {
		if db, err = sql.Open("postgres", uri); err != nil {
			return cli.NewExitError(fmt.Errorf("could not connect to database: %s", err), 1)
		}
	} else {
		return cli.NewExitError("could not connect: no database url specified", 1)
	}

	var m migrations.Migration

	if r := c.Int64("revision"); r >= 0 {
		if m, err = migrations.Revision(r, db); err != nil {
			return cli.NewExitError(err, 1)
		}
	} else {
		if m, err = migrations.Current(db); err != nil {
			return cli.NewExitError(err, 1)
		}
	}

	fmt.Println(m.String())

	return nil
}

func migrate(c *cli.Context) (err error) {
	var db *sql.DB
	if uri := c.String("db"); uri != "" {
		if db, err = sql.Open("postgres", uri); err != nil {
			return cli.NewExitError(fmt.Errorf("could not connect to database: %s", err), 1)
		}
	} else {
		return cli.NewExitError("could not connect: no database url specified", 1)
	}

	var n int
	if n, err = migrations.Migrate(c.Int64("revision"), db); err != nil {
		return cli.NewExitError(err, 1)
	}
	fmt.Printf("%d migrations affected\n", n)

	var m migrations.Migration
	if m, err = migrations.Current(db); err != nil {
		return cli.NewExitError(err, 1)
	}

	fmt.Printf("\ncurrent migration:\n%s\n", m.String())
	return nil
}
