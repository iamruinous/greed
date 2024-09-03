package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/iamruinous/greed/internal/cache"
	"github.com/iamruinous/greed/internal/feedbin"
	"github.com/iamruinous/greed/internal/ui"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "feedbin-cli",
		Usage: "A CLI client for Feedbin",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "username",
				Aliases:  []string{"u"},
				EnvVars:  []string{"FEEDBIN_USERNAME"},
				Usage:    "Feedbin username",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "password",
				Aliases:  []string{"p"},
				EnvVars:  []string{"FEEDBIN_PASSWORD"},
				Usage:    "Feedbin password",
				Required: true,
			},
			&cli.IntFlag{
				Name:    "fetch-limit",
				Aliases: []string{"l"},
				EnvVars: []string{"FEEDBIN_ENTRY_LIMIT"},
				Usage:   "Number of entries to fetch",
				Value:   5,
			},
			&cli.IntFlag{
				Name:    "display-limit",
				Aliases: []string{"d"},
				EnvVars: []string{"FEEDBIN_DISPLAY_LIMIT"},
				Usage:   "Number of entries to display",
				Value:   5,
			},
			&cli.StringFlag{
				Name:    "cache-dir",
				Aliases: []string{"c"},
				EnvVars: []string{"GREED_CACHE_DIR"},
				Usage:   "Cache directory",
				Value: func() string {
					cacheDir, err := os.UserCacheDir()
					if err != nil {
						log.Fatal("Failed to get user cache directory:", err)
					}
					return filepath.Join(cacheDir, "greed")
				}(),
			},
		},
		Action: run,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	username := c.String("username")
	password := c.String("password")
	fetchLimit := c.Int("fetch-limit")
	displayLimit := c.Int("display-limit")
	cacheDir := c.String("cache-dir")

	client := feedbin.NewClient(username, password, fetchLimit)

	cache, err := cache.New(cacheDir)
	if err != nil {
		return fmt.Errorf("failed to create cache: %w", err)
	}

	entries, err := ui.FetchEntries(client, cache)
	if err != nil {
		return err
	}

	ui.DisplayEntries(entries, displayLimit)
	return nil
}
