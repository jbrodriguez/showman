package main

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jbrodriguez/mlog"
)

const allowed = ".mkv;.srt"

func Scan(settings *Settings) (Shows, error) {
	re := regexp.MustCompile(`(.*)\.S(\d\d)E(\d\d)\.`)

	root := settings.SourceDir

	mlog.Info("")
	mlog.Info("Started scanning %s ...", root)

	entries, err := ioutil.ReadDir(root)
	if err != nil {
		return nil, err
	}

	shows := make(Shows)

	for _, entry := range entries {
		if !entry.IsDir() {
			// TODO: handle file case
			mlog.Warning("Leaving %s unprocessed", entry.Name())
			continue
		}

		// it's a folder-episode
		location := filepath.Join(root, entry.Name())

		items, err := ioutil.ReadDir(location)
		if err != nil {
			mlog.Warning("Unable to scan %s: %s", location, err)
			continue
		}

		if len(items) == 0 {
			mlog.Warning("Folder %s has no items. Skipping this folder ...", location)
			continue
		}

		// regex to obtain episode showname, seasonnumber, episodenumber
		// lost.s01e01.scene
		matches := re.FindStringSubmatch(entry.Name())
		if matches == nil {
			mlog.Warning("Unable to get show data from filename %s", entry.Name())
			continue
		}

		showName := strings.ToLower(strings.Replace(matches[1], ".", " ", -1))
		seasonNumber := matches[2]
		episodeNumber := matches[3]

		// lookup this show in the shows map, if it doesn't exist, create the new entry
		// shows["lost"] exists ?
		var show *Show
		var ok bool
		if show, ok = shows[showName]; !ok {
			show = &Show{
				Seasons: make(map[string][]*Episode),
			}
		}

		// create an episode with the info we have so far
		episode := &Episode{
			Season:   seasonNumber,
			Episode:  episodeNumber,
			Location: location,
			Files:    make([]string, 0),
		}

		// add all files in the folder for this episode
		for _, item := range items {
			// ignore .jpg, .nfo, files without extension, etc
			if filepath.Ext(item.Name()) == "" || !strings.Contains(allowed, filepath.Ext(item.Name())) {
				continue
			}

			episode.Files = append(episode.Files, item.Name())
		}

		// if the seasons haven't been initialized yet, do it now, so we can append to it
		// shows["lost"].seasons["01"] = []
		if _, ok := show.Seasons[seasonNumber]; !ok {
			show.Seasons[seasonNumber] = make([]*Episode, 0)
		}

		show.Seasons[seasonNumber] = append(show.Seasons[seasonNumber], episode)

		shows[showName] = show
	}

	for name, show := range shows {
		mlog.Info("Found show(%s):", name)
		for season, episodes := range show.Seasons {
			for _, episode := range episodes {
				mlog.Info("s%se%s - %s", season, episode.Episode, episode.Files)
				// for _, file := range episode.Files {
				// 	mlog.Info("%s", file)
				// }
				// mlog.Info("")
			}
			mlog.Info("")
		}
	}

	mlog.Info("Finished scanning ...")

	return shows, nil
}
