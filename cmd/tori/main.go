package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "tori",
		Usage: "manage the Go toolchain versions",
		Commands: []*cli.Command{
			{
				Name:    "install",
				Aliases: []string{"i"},
				Usage:   "install the specified version",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "use",
						Aliases: []string{"u"},
						Usage:   "use the specified version as the default",
					},
				},
				Action: handleInstall,
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "list installed versions",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "available",
						Aliases: []string{"a"},
						Usage:   "list versions available for installation",
					},
				},
				Action: handleList,
			},
			{
				Name:    "remove",
				Aliases: []string{"r"},
				Usage:   "remove the specified version",
				Action:  handleRemove,
			},
			{
				Name:    "use",
				Aliases: []string{"u"},
				Usage:   "use the specified version as the default",
				Action:  handleUse,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func handleInstall(cCtx *cli.Context) error {
	return nil
}

func handleList(cCtx *cli.Context) error {
	return nil
}

func handleRemove(cCtx *cli.Context) error {
	return nil
}

func handleUse(cCtx *cli.Context) error {
	return nil
}
