package core

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
)

const tmpPostfixSize = 16

func getHome() string {
	home := os.Getenv("TORI_HOME")
	if home == "" {
		home = filepath.Join(os.Getenv("HOME"), ".tori")
	}

	return home
}

func symlink(version, path string) error {
	b := make([]byte, tmpPostfixSize)
	if _, err := rand.Read(b); err != nil {
		return err
	}

	link := filepath.Join(path, "bin")
	tmpLink := link + hex.EncodeToString(b)

	target := filepath.Join("versions", version, "go", "bin")
	if err := os.Symlink(target, tmpLink); err != nil {
		return err
	}

	if err := os.Rename(tmpLink, link); err != nil {
		return err
	}

	return nil
}
