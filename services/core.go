package services

import (
	"time"

	"github.com/fatih/color"

	"showman/domain"
)

type Core struct {
}

func (c *Core) Run(ctx *domain.Context) error {
	now := time.Now()
	color.Blue("starting ...\n")

	// r := rand.New(rand.NewSource(seed))

	// walk over the source path and collect all the nzbs
	episodes, err := c.scan(ctx)
	if err != nil {
		return err
	}

	c.scrape(ctx, episodes)

	c.move(ctx, episodes)

	err = c.cleanup(ctx)
	if err != nil {
		return err
	}

	color.Blue("\ncompleted (elapsed: %s)", time.Since(now))

	return nil
}
