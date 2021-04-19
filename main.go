package main

import (
	"fmt"
	"os"

	"github.com/sheb-gregor/angiograph/analysis"
	"github.com/urfave/cli"
)

func initApp() *cli.App {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name:      "repl",
			ShortName: "i",
			Usage:     "",
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "private"},
				cli.BoolFlag{Name: "deps"},
				cli.BoolFlag{Name: "json"},
			},
			Action: func(c *cli.Context) error {
				a := analysis.AnalyseOpts{
					Deps:      false,
					Private:   false,
					PrintJSON: false,
					Path:      c.Args().First(),
					Patterns:  c.Args().Tail(),
				}

				return a.Run()
			},
		},
	}
	app.Action = func(c *cli.Context) error {
		a := analysis.AnalyseOpts{
			Deps:      true,
			Private:   false,
			PrintJSON: false,
			Path:      "/box/projects/sheb/hermes",
			// Patterns:  c.Args().Tail(),
		}
		_, err := a.ImportsTree()
		// fmt.Println(res)
		return err
	}
	return app
}

func main() {
	app := initApp()

	if err := app.Run(os.Args); err != nil {
		fmt.Println("FATAL: ", err)
		os.Exit(1)
	}
}
