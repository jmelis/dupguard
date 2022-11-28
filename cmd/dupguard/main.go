package main

import (
	"log"
	"os"

	"github.com/jmelis/dupguard/internal/indexer"
	"github.com/jmelis/dupguard/internal/reporter"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "index",
				Usage: "index a path",
				Action: func(cCtx *cli.Context) error {
					indexer.Index(cCtx.Args().Slice())
					return nil
				},
			},
			{
				Name:  "check",
				Usage: "check a path",
				Action: func(cCtx *cli.Context) error {
					reporter.Check(cCtx.Args().Slice())
					return nil
				},
			},
			{
				Name:  "dupes",
				Usage: "report duplicated files",
				Action: func(cCtx *cli.Context) error {
					reporter.Dupes()
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
