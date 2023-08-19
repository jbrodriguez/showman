package services

import (
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/ryanbradynd05/go-tmdb"

	"showman/domain"
)

func (c *Core) scrape(ctx *domain.Context, episodes []*domain.Episode) {
	color.Blue("\nscraping ...\n")
	scrapeTMDB(ctx, episodes)
}

func scrapeTMDB(ctx *domain.Context, episodes []*domain.Episode) {
	r := strings.NewReplacer(":", "", "*", "", "?", "")

	config := tmdb.Config{
		APIKey:   ctx.ApiKey,
		Proxies:  nil,
		UseProxy: false,
	}

	t := tmdb.Init(config)
	options := map[string]string{}

	cache := map[string]string{}
	notFound := map[string]bool{}

	for _, episode := range episodes {
		if _, ok := notFound[episode.Series]; ok {
			continue
		}

		name := ""
		if _, ok := cache[episode.Series]; ok {
			name = cache[episode.Series]
		} else {
			results, err := t.SearchTv(episode.Series, options)
			if err != nil {
				color.Yellow("query failed for %s: %s", episode.Series, err)
				time.Sleep(2 * time.Second)
				continue
			}

			if results.TotalResults == 0 {
				color.Yellow("show not found: %s", episode.Series)
				time.Sleep(2 * time.Second)
				notFound[episode.Series] = true
				continue
			}

			series := results.Results[0]
			name = r.Replace(series.Name)
			cache[episode.Series] = name
			color.Green("%s -> %s", episode.Series, name)
			time.Sleep(1 * time.Second)
		}

		episode.Name = name
	}
}
