package services

import (
	"os"
	"path/filepath"

	"github.com/fatih/color"

	"showman/domain"
)

func (c *Core) cleanup(ctx *domain.Context) error {
	color.Blue("\ncleaning up ...\n")

	err := filepath.WalkDir(ctx.SourcePath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == ctx.SourcePath {
			return nil
		}

		// color.Green("removing %s", path)
		parentDir := filepath.Dir(path)
		if parentDir == ctx.SourcePath {
			// Print direct children of the root directory
			color.Green("transferring %s", path)

			e := os.Rename(path, filepath.Join(ctx.TransferredPath, d.Name()))
			if e != nil {
				return e
			}

			if d.IsDir() {
				return filepath.SkipDir // Skip traversal of subdirectories
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
