package core

import (
	"os"
	"path/filepath"
)

// TODO: add a descriptive error message for non-existent versions
func Remove(version string) error {
	home := getHome()
	link := filepath.Join(home, "bin")
	target := filepath.Join(home, "versions", version)

	linked, err := os.Readlink(link)
	if err != nil {
		return err
	}

	if linked == filepath.Join("versions", version, "go", "bin") {
		if err := os.Remove(link); err != nil {
			return err
		}
	}

	if err := os.RemoveAll(target); err != nil {
		return err
	}

	return nil
}
