package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/iamruinous/greed/internal/cache"
	"github.com/iamruinous/greed/internal/feedbin"
	"github.com/iamruinous/greed/internal/ui"
	"github.com/urfave/cli/v3"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func main() {
	app := &cli.Command{
		Name:                  "feedbin-cli",
		Usage:                 "A CLI client for Feedbin",
		EnableShellCompletion: true,
		Suggest:               true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "username",
				Aliases: []string{"u"},
				Sources: cli.EnvVars("FEEDBIN_USERNAME"),
				Usage:   "Feedbin username",
			},
			&cli.StringFlag{
				Name:    "password",
				Aliases: []string{"P"},
				Sources: cli.EnvVars("FEEDBIN_PASSWORD"),
				Usage:   "Feedbin password",
			},
			&cli.IntFlag{
				Name:    "fetch-limit",
				Aliases: []string{"l"},
				Sources: cli.EnvVars("GREED_FETCH_LIMIT"),
				Usage:   "Number of entries to fetch",
				Value:   20,
			},
			&cli.StringFlag{
				Name:    "cache-dir",
				Aliases: []string{"c"},
				Sources: cli.EnvVars("GREED_CACHE_DIR"),
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
				Name:    "show-progress",
				Aliases: []string{"p"},
				Sources: cli.EnvVars("GREED_SHOW_PROGRESS"),
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
				Action: listEntries,
				Flags: []cli.Flag{
					&cli.DurationFlag{
						Name:    "cache-expires-after",
						Aliases: []string{"e"},
						Sources: cli.EnvVars("GREED_CACHE_EXPIRES_AFTER"),
						Usage:   "Cache expiration duration",
						Value:   time.Minute * 5,
					},
					&cli.IntFlag{
						Name:    "display-limit",
						Aliases: []string{"d"},
						Sources: cli.EnvVars("GREED_DISPLAY_LIMIT"),
						Usage:   "Number of entries to display",
						Value:   5,
					},
					&cli.BoolFlag{
						Name:    "ignore-cache",
						Aliases: []string{"i"},
						Sources: cli.EnvVars("GREED_IGNORE_CACHE"),
						Usage:   "Ignore cache and fetch latest entries",
						Value:   false,
					},
					&cli.BoolFlag{
						Name:    "random",
						Aliases: []string{"r"},
						Sources: cli.EnvVars("GREED_RANDOM"),
						Usage:   "Display random entries",
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

	err := app.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func listEntries(ctx context.Context, cmd *cli.Command) error {
	username := cmd.String("username")
	password := cmd.String("password")

	if username == "" || password == "" {
		return fmt.Errorf("username and password are required")
	}

	fetchLimit := cmd.Int("fetch-limit")
	if fetchLimit <= 0 {
		return fmt.Errorf("fetch-limit must be greater than 0")
	}

	displayLimit := cmd.Int("display-limit")
	if displayLimit <= 0 {
		return fmt.Errorf("display-limit must be greater than 0")
	}

	cacheDir := cmd.String("cache-dir")
	showProgress := cmd.Bool("show-progress")
	ignoreCache := cmd.Bool("ignore-cache")
	cacheDuration := cmd.Duration("cache-expires-after")

	client := feedbin.NewClient(username, password, fetchLimit)

	cache, err := cache.New(cacheDir, cacheDuration)
	if err != nil {
		return fmt.Errorf("failed to create cache: %w", err)
	}

	entries, err := ui.FetchEntries(client, cache, showProgress, ignoreCache)
	if err != nil {
		return err
	}

	randomEntries := cmd.Bool("random")
	ui.DisplayEntries(entries, displayLimit, randomEntries)
	return nil
}

func fetchAndUpdateCache(ctx context.Context, cmd *cli.Command) error {
	username := cmd.String("username")
	password := cmd.String("password")

	if username == "" || password == "" {
		return fmt.Errorf("username and password are required")
	}

	fetchLimit := cmd.Int("fetch-limit")
	if fetchLimit <= 0 {
		return fmt.Errorf("fetch-limit must be greater than 0")
	}

	cacheDir := cmd.String("cache-dir")
	showProgress := cmd.Bool("show-progress")
	cacheDuration := cmd.Duration("cache-expires-after")

	client := feedbin.NewClient(username, password, fetchLimit)

	cache, err := cache.New(cacheDir, cacheDuration)
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

func printVersion(ctx context.Context, cmd *cli.Command) error {
	fmt.Printf("Greed version %s\n", Version)
	fmt.Printf("Commit: %s\n", Commit)
	fmt.Printf("Built: %s\n", Date)
	return nil
}
