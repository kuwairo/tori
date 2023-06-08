package core

import (
	"os"
	"path/filepath"
)

func Use(version string) error {
	home := getHome()
	target := filepath.Join(home, "versions", version, "go", "bin")

	if _, err := os.Stat(target); err != nil {
		return err
	}

	if err := symlink(version, home); err != nil {
		return err
	}

	return nil
}
