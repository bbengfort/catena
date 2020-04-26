package main

import (
	"os"

	"github.com/bbengfort/catena"
	"github.com/joho/godotenv"
	"gopkg.in/urfave/cli.v1"
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
			},
		},
	}

	// Run the program, it should not error
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func serve(c *cli.Context) (err error) {
	return nil
}
