package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jbrodriguez/mlog"
)

func Move(settings *Settings, shows Shows) error {
	mlog.Info("Starting mover ...")

	r := strings.NewReplacer(":", "", "*", "", "/", "", "?", "")

	for _, show := range shows {
		if !show.Scraped {
			continue
		}

		for season, episodes := range show.Seasons {
			destination := filepath.Join(settings.DestDir, show.Name, season)

			mlog.Info(" -------- Creating destination %s ...", destination)

			err := os.MkdirAll(destination, 0755)
			if err != nil {
				mlog.Warning("Unable to create destination %s: %s", destination, err)
				continue
			}

			for _, episode := range episodes {
				mlog.Info("s%se%s - %s", episode.Season, episode.Episode, episode.Files)

				for _, file := range episode.Files {
					// move file from source to destination
					src := filepath.Join(episode.Location, file)
					dst := filepath.Join(destination, fmt.Sprintf("%s - S%sE%s - %s.%s", show.Name, episode.Season, episode.Episode, r.Replace(episode.Name), file[len(file)-3:]))

					mlog.Info("Moving [%s] -> [%s] ...", src, dst)
					err := os.Rename(src, dst)
					if err != nil {
						mlog.Warning("Unable to move (%s) to (%s): %s", src, dst, err)
						mlog.Info("")
						continue
					}
				}

				empty, err := IsDirEmpty(episode.Location)
				if err != nil {
					mlog.Warning("Unable to check for empty folder (%s): %s", episode.Location, err)
					continue
				}

				if !empty {
					continue
				}

				src := episode.Location
				dst := filepath.Join(settings.TransferDir, filepath.Base(episode.Location))

				mlog.Info("Transferring [%s] -> [%s] ...", src, dst)
				err = os.Rename(src, dst)
				if err != nil {
					mlog.Warning("Unable to move (%s) to (%s): %s", src, dst, err)
					mlog.Info("")
					continue
				}

				mlog.Info("")
			}
		}
	}

	mlog.Info("Finshed mover...")

	return nil
}
