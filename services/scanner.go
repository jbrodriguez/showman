package services

import (
	"os"
	"path/filepath"
	"regexp"
	"showman/domain"
	"strings"

	"github.com/fatih/color"
)

const allowed = ".mp4;.mkv;.avi;.srt"

var re = regexp.MustCompile(`/([^/]+)[\.-]s(\d{2})e(\d{2})[\.-]`)
var reYear = regexp.MustCompile(`(\s\d{4})$`)

func (c *Core) scan(ctx *domain.Context) ([]*domain.Episode, error) {
	color.Blue("\nscanning ...\n")

	root := ctx.SourcePath
	episodes := make([]*domain.Episode, 0)

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == root {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(d.Name()) == "" || !strings.Contains(allowed, filepath.Ext(d.Name())) {
			return nil
		}

		matches := re.FindStringSubmatch(strings.ToLower(path))
		if matches == nil {
			// mlog.Warning("Unable to get show data from filename %s", entry.Name())
			// transferrable = append(transferrable, path)
			color.Yellow("unable to process: %s", d.Name())
			return nil
		}

		showname := strings.ToLower(strings.Replace(matches[1], ".", " ", -1))
		showname = strings.ReplaceAll(showname, ":", "")
		showname = reYear.ReplaceAllString(showname, "")
		seasonNumber := matches[2]
		episodeNumber := matches[3]

		episode := &domain.Episode{
			Series:   showname,
			Season:   seasonNumber,
			Episode:  episodeNumber,
			Location: path,
		}

		episodes = append(episodes, episode)

		// color.Green("path: %s", path)
		color.Green("%s/s.%s.e%s/%s", episode.Series, episode.Season, episode.Episode, d.Name())

		return nil
	})
	if err != nil {
		return nil, err
	}

	// color.Blue("finished scanning ...\n")
	return episodes, nil
}

// func handleSingleFileEpisode(root string, entry os.DirEntry, shows domain.Shows) {
// 	// it's a single file episode
// 	// regex to obtain episode showname, seasonnumber, episodenumber
// 	// lost.s01e01.scene
// 	matches := re.FindStringSubmatch(strings.ToLower(entry.Name()))
// 	if matches == nil {
// 		// mlog.Warning("Unable to get show data from filename %s", entry.Name())
// 		color.Yellow("unable to process: %s", entry.Name())
// 		return
// 	}

// 	showname := strings.ToLower(strings.Replace(matches[1], ".", " ", -1))
// 	showname = strings.ReplaceAll(showname, ":", "")
// 	season := matches[3]
// 	episode := matches[4]

// 	// lookup this show in the shows map, if it doesn't exist, create the new entry
// 	// shows["lost"] exists ?
// 	var show domain.Show
// 	var ok bool
// 	if show, ok = shows[showname]; !ok {
// 		show = domain.Show{
// 			Seasons: map[string][]domain.Episode{},
// 			Multi:   false,
// 		}
// 	}

// 	ep := domain.Episode{
// 		Name:    showname,
// 		Season:  season,
// 		Episode: episode,
// 		Files:   []string{entry.Name()},
// 	}

// 	if _, ok := show.Seasons[season]; !ok {
// 		show.Seasons[season] = []domain.Episode{}
// 	}

// 	show.Seasons[season] = append(show.Seasons[season], ep)

// 	shows[showname] = show
// }

// func handleSingleEpisode(root string, entry os.DirEntry, matches []string, shows domain.Shows) {
// 	location := filepath.Join(root, entry.Name())

// 	entries, err := os.ReadDir(location)
// 	if err != nil {
// 		color.Yellow("unable to read dir: %s", location)
// 		return
// 	}

// 	if len(entries) == 0 {
// 		color.Yellow("empty dir: %s", location)
// 		return
// 	}

// 	showname := strings.ToLower(strings.Replace(matches[1], ".", " ", -1))
// 	showname = strings.ReplaceAll(showname, ":", "")
// 	season := matches[3]
// 	episode := matches[4]

// 	// lookup this show in the shows map, if it doesn't exist, create the new entry
// 	// shows["lost"] exists ?
// 	var show domain.Show
// 	var ok bool
// 	if show, ok = shows[showname]; !ok {
// 		show = domain.Show{
// 			Seasons: map[string][]domain.Episode{},
// 			Multi:   false,
// 		}
// 	}

// 	ep := domain.Episode{
// 		Name:    showname,
// 		Season:  season,
// 		Episode: episode,
// 		Files:   []string{},
// 	}

// 	// add all files in the folder for this episode
// 	for _, entry := range entries {
// 		// ignore .jpg, .nfo, files without extension, etc
// 		if filepath.Ext(entry.Name()) == "" || !strings.Contains(allowed, filepath.Ext(entry.Name())) {
// 			continue
// 		}

// 		ep.Files = append(ep.Files, entry.Name())
// 	}

// 	// if the seasons haven't been initialized yet, do it now, so we can append to it
// 	// shows["lost"].seasons["01"] = []
// 	if _, ok := show.Seasons[season]; !ok {
// 		show.Seasons[season] = make([]domain.Episode, 0)
// 	}

// 	show.Seasons[season] = append(show.Seasons[season], ep)

// 	shows[showname] = show
// }
