package steps

import (
	"context"
	"os"
	"path/filepath"
)

type CleanDirectoryStep struct {
	Directory string
}

func (c *CleanDirectoryStep) Do(ctx context.Context) error {
	// remove all files in the directory
	dirEntries, err := os.ReadDir(c.Directory)
	if err != nil {
		return err
	}

	for _, entry := range dirEntries {
		entryPath := filepath.Join(c.Directory, entry.Name())
		err := os.RemoveAll(entryPath)
		if err != nil {
			return err
		}
	}

	return nil
}
