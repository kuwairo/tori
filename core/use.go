package core

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func Use(version string) error {
	home := getHome()
	target := filepath.Join(home, "versions", version, "go", "bin")

	if _, err := os.Stat(target); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("unable to locate version %q: %w", version, err)
		}
		return err
	}

	if err := symlink(version, home); err != nil {
		return err
	}

	return nil
}
