package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/bbengfort/catena"
	"github.com/bbengfort/catena/config"
	"github.com/bbengfort/catena/migrations"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"

	// use postgres for the test database
	_ "github.com/lib/pq"
)

func main() {
	// Load the .env file if it exists
	godotenv.Load()

	// Instantiate the CLI application
	app := cli.NewApp()
	app.Name = "catena"
	app.Version = catena.Version
	app.Usage = "catena server and server utilities"
	app.Before = makeConfig
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "c, conf",
			Usage:  "specify a path to a configuration file for loading",
			EnvVar: "CATENA_CONF_PATH",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:     "serve",
			Usage:    "run the catena API server",
			Action:   serve,
			Before:   updateConfig,
			Category: "server",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "a, addr",
					Usage:  "the address to bind the catena server on",
					EnvVar: "CATENA_BIND_ADDR",
				},
				cli.UintFlag{
					Name:   "p, port",
					Usage:  "the port the catena server listens on",
					EnvVar: "CATENA_PORT",
				},
				cli.BoolFlag{
					Name:   "S, no-tls",
					Usage:  "do not run the server with TLS security",
					EnvVar: "CATENA_NO_TLS",
				},
				cli.StringFlag{
					Name:   "D, db",
					Usage:  "the database uri of the catena postgres database",
					EnvVar: "DATABASE_URL",
				},
			},
		},
		{
			Name:     "configure",
			Usage:    "read or edit the current catena configuration",
			Action:   configure,
			Category: "server",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "p, paths",
					Usage: "list the discovered configuration paths and exit",
				},
			},
		},
		{
			Name:      "db:revision",
			Usage:     "print the current migration status of the database",
			ArgsUsage: "[-n revision name]",
			Action:    revision,
			Before:    updateConfig,
			Category:  "database",
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
				cli.BoolFlag{
					Name:  "n, new",
					Usage: "create a new revision with the name as args",
				},
			},
		},
		{
			Name:     "db:migrate",
			Usage:    "migrate (or rollback) the database to the latest revision",
			Action:   migrate,
			Category: "database",
			Before:   updateConfig,
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
// Setup and Teardown
//===========================================================================

// global configuration for all CLI commands
var conf config.Config

func makeConfig(c *cli.Context) (err error) {
	// Load the configuration from the environment
	if conf, err = config.New(); err != nil {
		return cli.NewExitError(err, 1)
	}

	// If a path to the configuration is provided, loaded it
	if confpath := c.String("conf"); confpath != "" {
		if conf, err = conf.LoadFile(confpath); err != nil {
			return cli.NewExitError(err, 1)
		}
	} else {
		// Otherwise, attempt to load the configuration from the search path
		if conf, err = conf.LoadSystem(); err != nil {
			return cli.NewExitError(err, 1)
		}
	}

	return nil
}

func updateConfig(c *cli.Context) (err error) {
	// Update the config from the context manually for CLI flags
	// TODO: should we do this with reflection as well?
	if addr := c.String("addr"); addr != "" {
		conf.Addr = addr
	}

	if port := c.Uint("port"); port != 0 {
		conf.Port = uint16(port)
	}

	if noTLS := c.Bool("no-tls"); noTLS {
		conf.NoTLS = true
	}

	if db := c.String("db"); db != "" {
		conf.DBURL = db
	}

	return nil
}

//===========================================================================
// Server Commands
//===========================================================================

func serve(c *cli.Context) (err error) {
	var api *catena.Catena
	if api, err = catena.New(conf); err != nil {
		return cli.NewExitError(err, 1)
	}

	if err = api.Serve(); err != nil {
		return cli.NewExitError(err, 1)
	}

	return nil
}

func configure(c *cli.Context) (err error) {
	if c.Bool("paths") {
		if path := c.GlobalString("conf"); path != "" {
			fmt.Printf("*%s\n", path)
			return nil
		}

		paths := config.Find()
		if len(paths) == 0 {
			fmt.Println("no system configuration files")
		} else {
			for i, path := range paths {
				if i == 0 {
					fmt.Printf("*%s\n", path)
				} else {
					fmt.Println(path)
				}
			}
		}
		return nil
	}

	// TODO: specify output format
	// TODO: allow user to edit with vim
	data, err := yaml.Marshal(conf)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	fmt.Println(string(data))
	return nil
}

//===========================================================================
// Database Commands
//===========================================================================

func revision(c *cli.Context) (err error) {
	if c.Bool("new") {
		var path string
		if path, err = migrations.New(strings.Join(c.Args(), "_"), "migrations"); err != nil {
			return cli.NewExitError(err, 1)
		}

		fmt.Printf("blank migration created at %s\n", path)
		fmt.Println("remember to run go generate ./... to add the migration")
		return
	}

	if conf.DBURL == "" {
		return cli.NewExitError("could not connect: no database url specified", 1)
	}

	var db *sql.DB
	if db, err = sql.Open("postgres", conf.DBURL); err != nil {
		return cli.NewExitError(fmt.Errorf("could not connect to database: %s", err), 1)
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
	if conf.DBURL == "" {
		return cli.NewExitError("could not connect: no database url specified", 1)
	}

	var db *sql.DB
	if db, err = sql.Open("postgres", conf.DBURL); err != nil {
		return cli.NewExitError(fmt.Errorf("could not connect to database: %s", err), 1)
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
