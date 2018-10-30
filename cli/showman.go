package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jbrodriguez/mlog"
)

// Version - app version
var Version string

func main() {
	settings, err := setup(Version)
	if err != nil {
		log.Printf("Unable to start the app: %s", err)
		os.Exit(1)
	}

	run(settings)
}

func setup(version string) (*Settings, error) {
	home := os.Getenv("HOME")

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// look for showman.conf at the following places
	// $HOME/.config/showman/showman.conf
	// $HOME/.showman/showman.conf
	// <current dir>/showman.conf
	locations := []string{
		filepath.Join(home, ".config", "showman"),
		filepath.Join(home, ".showman"),
		cwd,
	}

	settings, err := NewSettings("showman.conf", version, home, locations)

	return settings, err
}

func run(settings *Settings) {
	mlog.DefaultFlags = mlog.DefaultFlags &^ (log.Ldate | log.Ltime | log.Lshortfile)
	if settings.LogToFile {
		mlog.Start(mlog.LevelInfo, filepath.Join(settings.DataDir, "logs", "showman.log"))
	} else {
		mlog.Start(mlog.LevelInfo, "")
	}

	mlog.Info("showman v%s starting ...", settings.Version)

	var msg string
	if settings.Location == "" {
		msg = "No config file specified. Using app defaults ..."
	} else {
		msg = fmt.Sprintf("Using config file located at %s ...", settings.Location)
	}
	mlog.Info(msg)

	shows, err := Scan(settings)
	if err != nil {
		mlog.Warning("Unable to scan for shows: %s", err)
	}

	shows, err = Scrape(settings, shows)
	if err != nil {
		mlog.Warning("Unable to scrape shows: %s", err)
	}

	Move(settings, shows)

	mlog.Stop()
}
