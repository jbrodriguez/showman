package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

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

	settings := NewSettings("showman.conf", version, home, locations)

	return settings, nil
}

func run(settings *Settings) {
	mlog.DefaultFlags &^= (log.Ldate | log.Ltime | log.Lshortfile)

	if settings.LogDir != "" {
		mlog.Start(mlog.LevelInfo, filepath.Join(settings.LogDir, "showman.log"))
	} else {
		mlog.Start(mlog.LevelInfo, "")
	}

	mlog.Info("showman v%s starting [%s] ...", settings.Version, time.Now().Format(time.RFC3339))

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

	if len(shows) == 0 {
		mlog.Info("No new shows found. Nothing do, exiting now ...")
		return
	}

	shows, err = Scrape(settings, shows)
	if err != nil {
		mlog.Warning("Unable to scrape one or more shows: %s", err)
	}

	Move(settings, shows)

	mlog.Stop()
}
