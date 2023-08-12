package main

import (
	"github.com/alecthomas/kong"

	"showman/cmd"
	"showman/domain"
)

var cli struct {
	SourcePath      string `required:"" short:"s" help:"location of source files"`
	DestinationPath string `required:"" short:"d" help:"location to destination files"`
	TransferredPath string `required:"" short:"t" help:"location to transferred files"`
	ApiKey          string `required:"" short:"k" help:"provider api key"`

	Boot cmd.Boot `cmd:"" default:"1" help:"start processing"`
}

func main() {
	ctx := kong.Parse(&cli)
	err := ctx.Run(&domain.Context{
		SourcePath:      cli.SourcePath,
		DestinationPath: cli.DestinationPath,
		TransferredPath: cli.TransferredPath,
		ApiKey:          cli.ApiKey,
	})
	ctx.FatalIfErrorf(err)
}
