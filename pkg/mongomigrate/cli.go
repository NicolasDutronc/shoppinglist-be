package mongomigrate

import (
	"github.com/urfave/cli/v2"
)

// GetCLI returns the mongomigrate CLI
func GetCLI(m *Mongomigrate) *cli.App {
	return &cli.App{
		Name:  "mongomigrate CLI",
		Usage: "Run migrations, rollback and check your database version",
		Commands: []*cli.Command{
			{
				Name: "up",
				Action: func(c *cli.Context) error {
					return upAction(c, m)
				},
			},
			{
				Name: "down",
				Action: func(c *cli.Context) error {
					return downAction(c, m)
				},
			},
			{
				Name: "version",
				Action: func(c *cli.Context) error {
					return getVersionAction(c, m)
				},
			},
			{
				Name: "seed",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Usage:   "Run seeders matching the name",
					},
				},
				Action: func(c *cli.Context) error {
					return seedAction(c, m)
				},
			},
			{
				Name: "search",
				Action: func(c *cli.Context) error {
					return search(c, m)
				},
			},
		},
	}
}
