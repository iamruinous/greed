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

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func main() {
	app := &cli.App{
		Name:  "feedbin-cli",
		Usage: "A CLI client for Feedbin",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "username",
				Aliases: []string{"u"},
				EnvVars: []string{"FEEDBIN_USERNAME"},
				Usage:   "Feedbin username",
			},
			&cli.StringFlag{
				Name:    "password",
				Aliases: []string{"P"},
				EnvVars: []string{"FEEDBIN_PASSWORD"},
				Usage:   "Feedbin password",
			},
			&cli.IntFlag{
				Name:    "fetch-limit",
				Aliases: []string{"l"},
				EnvVars: []string{"GREED_FETCH_LIMIT"},
				Usage:   "Number of entries to fetch",
				Value:   20,
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
			&cli.BoolFlag{
				Name:    "random",
				Aliases: []string{"r"},
				EnvVars: []string{"GREED_RANDOM"},
				Usage:   "Display random entries",
				Value:   false,
			},
			&cli.BoolFlag{
				Name:    "show-progress",
				Aliases: []string{"p"},
				EnvVars: []string{"GREED_SHOW_PROGRESS"},
				Usage:   "Show progress spinner",
				Value:   false,
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "fetch",
				Usage:  "Fetch latest entries and update cache",
				Action: fetchAndUpdateCache,
			},
			{
				Name:   "list",
				Usage:  "List latest entries",
				Action: run,
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "display-limit",
						Aliases: []string{"d"},
						EnvVars: []string{"GREED_DISPLAY_LIMIT"},
						Usage:   "Number of entries to display",
						Value:   5,
					},
					&cli.BoolFlag{
						Name:    "ignore-cache",
						Aliases: []string{"i"},
						EnvVars: []string{"GREED_IGNORE_CACHE"},
						Usage:   "Ignore cache and fetch latest entries",
						Value:   false,
					},
				},
			},
			{
				Name:   "version",
				Usage:  "Print the version information",
				Action: printVersion,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	username := c.String("username")
	password := c.String("password")

	if username == "" || password == "" {
		return fmt.Errorf("username and password are required")
	}

	fetchLimit := c.Int("fetch-limit")
	displayLimit := c.Int("display-limit")
	cacheDir := c.String("cache-dir")
	showProgress := c.Bool("show-progress")
	ignoreCache := c.Bool("ignore-cache")

	client := feedbin.NewClient(username, password, fetchLimit)

	cache, err := cache.New(cacheDir)
	if err != nil {
		return fmt.Errorf("failed to create cache: %w", err)
	}

	entries, err := ui.FetchEntries(client, cache, showProgress, ignoreCache)
	if err != nil {
		return err
	}

	randomEntries := c.Bool("random")
	ui.DisplayEntries(entries, displayLimit, randomEntries)
	return nil
}

func fetchAndUpdateCache(c *cli.Context) error {
	username := c.String("username")
	password := c.String("password")

	if username == "" || password == "" {
		return fmt.Errorf("username and password are required")
	}

	fetchLimit := c.Int("fetch-limit")
	cacheDir := c.String("cache-dir")
	showProgress := c.Bool("show-progress")

	client := feedbin.NewClient(username, password, fetchLimit)

	cache, err := cache.New(cacheDir)
	if err != nil {
		return fmt.Errorf("failed to create cache: %w", err)
	}

	entries, err := ui.FetchEntries(client, cache, showProgress, true) // always ignoreCache
	if err != nil {
		return err
	}

	cache.Set(entries)
	return nil
}

func printVersion(c *cli.Context) error {
	fmt.Printf("Greed version %s\n", Version)
	fmt.Printf("Commit: %s\n", Commit)
	fmt.Printf("Built: %s\n", Date)
	return nil
}
