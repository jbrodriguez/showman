package services

import (
	"os"
	"path/filepath"

	"github.com/fatih/color"

	"showman/domain"
)

func (c *Core) move(ctx *domain.Context, episodes []*domain.Episode) {
	color.Blue("\nmoving ...\n")

	for _, episode := range episodes {
		series := episode.Name
		season := episode.Season

		destination := filepath.Join(ctx.DestinationPath, series, season)

		color.Green("%s to %s", filepath.Base(episode.Location), destination)

		err := os.MkdirAll(destination, 0755)
		if err != nil {
			color.Red("unable to create destination %s: %s", destination, err)
			continue
		}

		err = os.Rename(episode.Location, filepath.Join(destination, filepath.Base(episode.Location)))
		if err != nil {
			color.Red("unable to move %s to %s: %s", episode.Location, destination, err)
			continue
		}
	}
}
