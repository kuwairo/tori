package core

import (
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func getHome() string {
	home := os.Getenv("TORI_HOME")
	if home == "" {
		home = filepath.Join(os.Getenv("HOME"), ".tori")
	}

	return home
}

func symlink(version, path string) error {
	link := filepath.Join(path, "bin")
	tmpLink := link + uuid.New().String()

	target := filepath.Join("versions", version, "go", "bin")
	if err := os.Symlink(target, tmpLink); err != nil {
		return err
	}

	if err := os.Rename(tmpLink, link); err != nil {
		return err
	}

	return nil
}
