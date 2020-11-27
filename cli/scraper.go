package main

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pioz/tvdb"
	"github.com/ryanbradynd05/go-tmdb"

	"github.com/jbrodriguez/mlog"
)

// Scrape -
func Scrape(settings *Settings, shows Shows) (Shows, error) {
	if settings.Provider == "tmdb" {
		return ScrapeTMDB(settings, shows)
	}

	return ScrapeTVDB(settings, shows)
}

// ScrapeTVDB -
func ScrapeTVDB(settings *Settings, shows Shows) (Shows, error) {
	mlog.Info("")
	mlog.Info("Started scraping shows ...")

	c := tvdb.Client{Apikey: settings.APIKey, Username: settings.UserName, Userkey: settings.UserKey, Language: "en"}
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

// ScrapeTMDB -
func ScrapeTMDB(settings *Settings, shows Shows) (Shows, error) {
	mlog.Info("")
	mlog.Info("Started scraping shows ...")

	config := tmdb.Config{
		APIKey:   settings.APIKey,
		Proxies:  nil,
		UseProxy: false,
	}

	c := tmdb.Init(config)

	mlog.Info("Connected to TMDB ...")

	options := map[string]string{}

	// shows["lost"].seasons["01"][{season: "01", episode: "01", files: ["1.mkv", "1.srt"]}]
	for name, show := range shows {
		mlog.Info("Looking up %s ...", name)

		// var series tvdb.Series
		// var err error

		// if strings.HasPrefix(name, "tvdbid-") {
		// 	series.ID, _ = strconv.Atoi(name[7:])
		// 	err = c.GetSeries(&series)
		// } else {
		// 	series, err = c.BestSearch(name)
		// }

		results, err := c.SearchTv(name, options)
		if err != nil {
			mlog.Warning("Unable to locate %s: %s", name, err)
			time.Sleep(2 * time.Second)
			continue
		}

		if results.TotalResults == 0 {
			mlog.Warning("Couldn't find a show by this name: %s", name)
			time.Sleep(2 * time.Second)
			continue
		}

		series := results.Results[0]

		show.Scraped = true
		show.Name = series.Name

		for season, episodes := range show.Seasons {
			if len(episodes) > 1 {
				// multiple episodes, let's do a season series call
				mlog.Info("Performing multi-episodes lookup ...")

				seasonNumber, _ := strconv.Atoi(season)

				tvSeason, err := c.GetTvSeasonInfo(series.ID, seasonNumber, options)
				if err != nil {
					mlog.Warning("Unable to get multi-episodes lookup: %s", err)
					continue
				}

				for _, episode := range episodes {
					episodeNumber, _ := strconv.Atoi(episode.Episode)

					tvEpisode := getEpisode(tvSeason, seasonNumber, episodeNumber)
					episode.Name = tvEpisode.Name

					mlog.Info("Found [%s - S%sE%s - %s]", show.Name, season, episode.Episode, episode.Name)
				}

			} else {
				// single episode, let's do an episode series call
				mlog.Info("Performing single-episode lookup ...")

				seasonNumber, _ := strconv.Atoi(season)
				episodeNumber, _ := strconv.Atoi(episodes[0].Episode)

				tvEpisode, err := c.GetTvEpisodeInfo(series.ID, seasonNumber, episodeNumber, options)
				if err != nil {
					mlog.Warning("Unable to get single-episode lookup: %s", err)
					continue
				}

				episodes[0].Name = tvEpisode.Name

				mlog.Info("Found [%s - S%sE%s - %s]", show.Name, season, episodes[0].Episode, episodes[0].Name)
			}
		}

		// avoid rate limiting from the api
		time.Sleep(2 * time.Second)
	}

	mlog.Info("Finished scraping shows ...")

	return shows, nil
}

func getEpisode(season *tmdb.TvSeason, seasonNumber, episodeNumber int) *tmdb.TvEpisode {
	for i := range season.Episodes {
		if season.Episodes[i].SeasonNumber == seasonNumber && season.Episodes[i].EpisodeNumber == episodeNumber {
			return &season.Episodes[i]
		}
	}

	return nil
}
