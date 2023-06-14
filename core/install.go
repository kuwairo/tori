package core

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

const defaultSource = "https://dl.google.com/go"

func buildURL(source, version string) string {
	goarch := runtime.GOARCH
	if goarch == "arm" {
		goarch = "armv6l"
	}

	goos := runtime.GOOS
	return fmt.Sprintf("%s/go%s.%s-%s.tar.gz", source, version, goos, goarch)
}

func fetch(url, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// TODO: add a descriptive error message for 404
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", res.Status)
	}

	if _, err := io.Copy(file, res.Body); err != nil {
		return err
	}

	return nil
}

func extract(archive io.Reader, path string) error {
	zr, err := gzip.NewReader(archive)
	if err != nil {
		return err
	}

	tr, mkdir := tar.NewReader(zr), map[string]bool{}
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		hdrMode := hdr.FileInfo().Mode()
		hdrPath := filepath.Join(path, filepath.FromSlash(hdr.Name))

		switch {
		case hdrMode.IsDir():
			if err := os.MkdirAll(hdrPath, 0755); err != nil {
				return err
			}
			mkdir[hdrPath] = true
		case hdrMode.IsRegular():
			if dir := filepath.Dir(hdrPath); !mkdir[dir] {
				if err := os.MkdirAll(dir, 0755); err != nil {
					return err
				}
				mkdir[dir] = true
			}

			file, err := os.OpenFile(hdrPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, hdrMode.Perm())
			if err != nil {
				return err
			}

			if _, err := io.Copy(file, tr); err != nil {
				return err
			}
			file.Close()
		}
	}

	return nil
}

func Install(version string, makeDefault, verbose bool) error {
	url := buildURL(defaultSource, version)
	archive, err := os.CreateTemp("", "tori-")
	if err != nil {
		return err
	}

	if verbose {
		fmt.Printf("Fetching %s...\n", version)
	}
	if err := fetch(url, archive.Name()); err != nil {
		return err
	}

	home := getHome()
	target := filepath.Join(home, "versions", version)
	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	if verbose {
		fmt.Printf("Extracting %s...\n", version)
	}
	if err := extract(archive, target); err != nil {
		return err
	}
	os.Remove(archive.Name())

	if makeDefault {
		if err := symlink(version, home); err != nil {
			return err
		}
	}

	return nil
}
