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
	version := cCtx.Args().First()
	makeDefault := cCtx.Bool("use")

	if err := core.Install(version, makeDefault, true); err != nil {
		fmt.Printf("An error occurred during execution: %v\n", err)
	}

	return nil
}

func handleList(cCtx *cli.Context) error {
	online := cCtx.Bool("available")

	if err := core.List(online); err != nil {
		fmt.Printf("An error occurred during execution: %v\n", err)
	}

	return nil
}

// TODO: consider accepting several versions
func handleRemove(cCtx *cli.Context) error {
	version := cCtx.Args().First()

	if err := core.Remove(version); err != nil {
		fmt.Printf("An error occurred during execution: %v\n", err)
	}

	return nil
}

func handleUse(cCtx *cli.Context) error {
	version := cCtx.Args().First()

	if err := core.Use(version); err != nil {
		fmt.Printf("An error occurred during execution: %v\n", err)
	}

	return nil
}
