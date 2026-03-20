// Package main is the entry point for Purgomatic.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/huangsam/purgomatic/internal/engine"
	"github.com/urfave/cli/v3"
	_ "modernc.org/sqlite"
)

func main() {
	app := &cli.Command{
		Name:  "purgomatic",
		Usage: "File indexing and migration planner",
		Commands: []*cli.Command{
			{
				Name:  "init",
				Usage: "Initialize the purgomatic SQLite database",
				Action: func(_ context.Context, _ *cli.Command) error {
					dbPath := engine.GetDBPath()
					db, err := sql.Open("sqlite", dbPath)
					if err != nil {
						return err
					}
					defer func() { _ = db.Close() }()
					if _, err := db.Exec(engine.DBInit); err != nil {
						return err
					}
					fmt.Printf("Initialized %s with Multi-Host support.\n", dbPath)
					return nil
				},
			},
			{
				Name:  "audit",
				Usage: "Scan all targets in scans.json & generate global library report",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "file", Aliases: []string{"f"}, Value: "scans.json", Usage: "Batch scan config (JSON)"},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					return engine.HandleAudit(cmd.String("file"))
				},
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Printf("Fatal error: %v\n", err)
		os.Exit(1)
	}
}
