package core

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func Remove(version string) error {
	home := getHome()
	link := filepath.Join(home, "bin")
	versions := filepath.Join(home, "versions")

	// TODO: fix possible out-of-home RemoveAll
	target := filepath.Join(versions, version)
	if target == versions {
		return errors.New("no valid version is specified for the command")
	}

	linked, err := os.Readlink(link)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	if linked == filepath.Join("versions", version, "go", "bin") {
		if err := os.Remove(link); err != nil {
			return err
		}
	}

	if _, err := os.Stat(target); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("unable to locate version %q: %w", version, err)
		}
		return err
	}

	if err := os.RemoveAll(target); err != nil {
		return err
	}

	return nil
}
