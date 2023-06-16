package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kuwairo/tori/core"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:                   "tori",
		Usage:                  "manage the Go toolchain versions",
		UseShortOptionHandling: true,

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
					&cli.IntFlag{
						Name:    "limit",
						Aliases: []string{"l"},
						Usage:   "limit the number of displayed versions to the specified value",
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

// TODO: consider adding '-q' flag for the third argument
func handleInstall(cCtx *cli.Context) error {
	version := cCtx.Args().First()
	makeDefault := cCtx.Bool("use")
	return handleError(core.Install(version, makeDefault, true))
}

func handleList(cCtx *cli.Context) error {
	online := cCtx.Bool("available")
	limit := cCtx.Int("limit")
	return handleError(core.List(online, limit))
}

func handleRemove(cCtx *cli.Context) error {
	version := cCtx.Args().First()
	return handleError(core.Remove(version))
}

func handleUse(cCtx *cli.Context) error {
	version := cCtx.Args().First()
	return handleError(core.Use(version))
}

func handleError(err error) error {
	if err != nil {
		fmt.Fprintf(os.Stderr, "An error occurred during execution: %s\n", err)
		os.Exit(1)
	}

	return nil
}
