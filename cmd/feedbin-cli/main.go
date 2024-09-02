package main

import (
	"log"
	"os"

	"github.com/iamruinous/greed/internal/cache"
	"github.com/iamruinous/greed/internal/feedbin"
	"github.com/iamruinous/greed/internal/ui"
)

func main() {
	client := feedbin.NewClient(
		os.Getenv("FEEDBIN_USERNAME"),
		os.Getenv("FEEDBIN_PASSWORD"),
	)

	cache := cache.New()
	tui := ui.NewTUI(client, cache)

	if _, err := tui.Run(); err != nil {
		log.Fatal(err)
	}
}
