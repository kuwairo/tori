package core

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	gversion "github.com/hashicorp/go-version"
)

func Remove(version string) error {
	if _, err := gversion.NewVersion(version); err != nil {
		return fmt.Errorf("%q is not a valid version", version)
	}

	home := getHome()
	link := filepath.Join(home, "bin")
	target := filepath.Join(home, "versions", version)

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
