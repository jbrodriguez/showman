package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jbrodriguez/mlog"
)

const allowed = ".mp4;.mkv;.avi;.srt"

// var re = regexp.MustCompile(`(.*)\.S(\d\d)E(\d\d)\.`)
var re = regexp.MustCompile(`(.*)[\.\s](s(\d\d)e(\d\d))[\.\s]`)

// Scan -
func Scan(settings *Settings) (Shows, error) {
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

		// it can be a single episode
		// or it can be full seasons
		// let's test for a single episode
		// regex to obtain episode showname, seasonnumber, episodenumber
		// lost.s01e01.scene
		matches := re.FindStringSubmatch(strings.ToLower(entry.Name()))
		if matches == nil {
			// it's a full season folder
			handleFullSeason(root, entry, shows)

			// mlog.Warning("Unable to get show data from filename %s", entry.Name())
			continue
		}

		handleSingleEpisode(root, entry, matches, shows)
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

func handleSingleEpisode(root string, entry os.FileInfo, matches []string, shows Shows) {
	location := filepath.Join(root, entry.Name())

	items, err := ioutil.ReadDir(location)
	if err != nil {
		mlog.Warning("Unable to scan %s: %s", location, err)
		return
	}

	if len(items) == 0 {
		mlog.Warning("Folder %s has no items. Skipping this folder ...", location)
		return
	}

	showName := strings.ToLower(strings.Replace(matches[1], ".", " ", -1))
	seasonNumber := matches[3]
	episodeNumber := matches[4]

	// lookup this show in the shows map, if it doesn't exist, create the new entry
	// shows["lost"] exists ?
	var show *Show
	var ok bool
	if show, ok = shows[showName]; !ok {
		show = &Show{
			Seasons: make(map[string][]*Episode),
			Multi:   false,
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

func handleFullSeason(root string, entry os.FileInfo, shows Shows) {
	// location is
	// <path>/Locked Up S01 <scene>
	location := filepath.Join(root, entry.Name())

	items, err := ioutil.ReadDir(location)
	if err != nil {
		mlog.Warning("Unable to scan %s: %s", location, err)
		return
	}

	if len(items) == 0 {
		mlog.Warning("Folder %s has no items. Skipping this folder ...", location)
		return
	}

	// temp structure to gather all files related to an episode
	episodes := make(map[string]*Episode)

	// "globals" that define what show and which season we're dealing with
	var show *Show
	var season string

	// items is the list of single episodes
	for _, item := range items {
		// ignore .jpg, .nfo, files without extension, etc
		if filepath.Ext(item.Name()) == "" || !strings.Contains(allowed, filepath.Ext(item.Name())) {
			continue
		}

		// item is
		// Locked Up S01E02 <scene>.{mkv|srt}
		matches := re.FindStringSubmatch(strings.ToLower(item.Name()))
		if matches == nil {
			mlog.Warning("Unable to get show data from filename %s", item.Name())
			continue
		}

		// Locked Up = showName = matches[1]
		// S01E02 = episodeId = matches[2]
		// 01 = seasonNumber = matches[3]
		// 02 = episodeNumber = matches[4]
		showName := strings.ToLower(strings.Replace(matches[1], ".", " ", -1))
		seasonNumber := matches[3]
		episodeNumber := matches[4]
		episodeID := matches[2]

		// this holds the season number to be used outside the loop
		season = seasonNumber

		// lookup this show in the shows map, if it doesn't exist, create the new entry
		// shows["locked up"] exists ?
		var ok bool
		if show, ok = shows[showName]; !ok {
			show = &Show{
				Seasons: make(map[string][]*Episode),
				Multi:   true,
			}

			// this adds the show to the "global" shows instance
			shows[showName] = show
		}

		// lookup this episode in the episodes map, it it doesn't exist, create the entry
		// episodes["S01E02"] exists ?
		var episode *Episode
		if episode, ok = episodes[episodeID]; !ok {
			// create an episode with the info we have so far
			episode = &Episode{
				Season:   seasonNumber,
				Episode:  episodeNumber,
				Location: location,
				Files:    make([]string, 0),
			}

			episodes[episodeID] = episode
		}

		episode.Files = append(episode.Files, item.Name())
	}

	// shows["Locked Up"].seasons["01"] = []
	show.Seasons[season] = make([]*Episode, 0)

	// add episodes to the show
	for _, episode := range episodes {
		show.Seasons[season] = append(show.Seasons[season], episode)
	}
}
