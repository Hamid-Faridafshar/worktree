package main

import (
	"log"
	"os"
	"worktree/internal/layout"

	"github.com/rivo/tview"
	"github.com/urfave/cli"
)

var (
	entryPoint string
)

func main() {
	app := &cli.App{
		Name:   "worktree cli",
		Action: cmd,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "entry-point",
				EnvVar:      "WT_ENTRY_POINT",
				Value:       ".",
				Usage:       "Directory that has all the repositories",
				Destination: &entryPoint,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("Failed to run app: %v", err)
	}
}

func cmd(c *cli.Context) error {
	app := tview.NewApplication()

	layout := layout.NewLayout(app, entryPoint)
	layout.SetupLayoutContentMenus(entryPoint)

	if err := layout.SetRoot(); err != nil {
		layout.Log("Failed to run app: %v", err)
	}

	return nil
}
