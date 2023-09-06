package main

import (
	"github.com/oupson/univcpt/internal/app/univcpt"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	icalUrl, hasUrl := os.LookupEnv("ICAL_URL")
	if !hasUrl {
		logger.Error("missing ical url")
		os.Exit(-1)
	}

	app := univcpt.NewApp(logger, icalUrl)
	if err := app.Run(); err != nil {
		logger.Error("failed to run server", "err", err)
		os.Exit(-1)
	}
}
