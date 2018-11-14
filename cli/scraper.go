package main

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pioz/tvdb"

	"github.com/jbrodriguez/mlog"
)

func Scrape(settings *Settings, shows Shows) (Shows, error) {
	mlog.Info("")
	mlog.Info("Started scraping shows ...")

	c := tvdb.Client{Apikey: settings.ApiKey, Username: settings.UserName, Userkey: settings.UserKey, Language: "en"}
	err := c.Login()
	if err != nil {
		mlog.Warning("Unable to connect to TVDBv2: %s", err)
		return shows, err
	}

	mlog.Info("Connected to TVDBv2 ...")

	// shows["lost"].seasons["01"][{season: "01", episode: "01", files: ["1.mkv", "1.srt"]}]
	for name, show := range shows {
		mlog.Info("Looking up %s ...", name)

		var series tvdb.Series
		var err error

		if strings.HasPrefix(name, "tvdbid-") {
			series.ID, _ = strconv.Atoi(name[7:])
			err = c.GetSeries(&series)
		} else {
			series, err = c.BestSearch(name)
		}

		if err != nil {
			mlog.Warning("Unable to locate %s: %s", name, err)
			time.Sleep(2 * time.Second)
			continue
		}

		show.Scraped = true
		show.Name = series.SeriesName

		for season, episodes := range show.Seasons {
			if len(episodes) > 1 {
				// multiple episodes, let's do a season series call
				mlog.Info("Performing multi-episodes lookup ...")

				err := c.GetSeriesEpisodes(&series, url.Values{"airedSeason": {season}})
				if err != nil {
					mlog.Warning("Unable to get multi-episodes lookup: %s", err)
					continue
				}

				for _, episode := range episodes {
					seasonNumber, _ := strconv.Atoi(season)
					episodeNumber, _ := strconv.Atoi(episode.Episode)

					episode.Name = series.GetEpisode(seasonNumber, episodeNumber).EpisodeName
					mlog.Info("Found [%s - S%sE%s - %s]", show.Name, season, episode.Episode, episode.Name)
				}

			} else {
				// single episode, let's do an episode series call
				mlog.Info("Performing single-episode lookup ...")

				err := c.GetSeriesEpisodes(&series, url.Values{"airedSeason": {season}, "airedEpisode": {episodes[0].Episode}})
				if err != nil {
					mlog.Warning("Unable to get single-episode lookup: %s", err)
					continue
				}

				seasonNumber, _ := strconv.Atoi(season)
				episodeNumber, _ := strconv.Atoi(episodes[0].Episode)

				episodes[0].Name = series.GetEpisode(seasonNumber, episodeNumber).EpisodeName

				mlog.Info("Found [%s - S%sE%s - %s]", show.Name, season, episodes[0].Episode, episodes[0].Name)
			}
		}

		// avoid rate limiting from the api
		time.Sleep(2 * time.Second)
	}

	mlog.Info("Finished scraping shows ...")

	return shows, nil
}
