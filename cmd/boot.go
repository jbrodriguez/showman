package cmd

import (
	"showman/domain"
	"showman/services"
)

type Boot struct {
}

func (b *Boot) Run(ctx *domain.Context) error {
	core := &services.Core{}
	return core.Run(ctx)
}
