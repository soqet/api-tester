package cli

import (
	"os"

	"github.com/urfave/cli/v2"
)

func initCommands() *cli.App {
	return &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "test",
				Aliases: []string{"t"},
				Usage:   "test json task",
				Action:  handleTest,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "file",
						Aliases: []string{"f"},
						Value:   "",
						Usage:   "file name",
					},
				},
			},
		},
	}
}

func Run() {
	app := initCommands()
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
