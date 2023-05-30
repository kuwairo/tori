package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

func download(url, path string) error {
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

func main() {
	version := "1.20.4"
	goos, goarch := runtime.GOOS, runtime.GOARCH
	url := fmt.Sprintf("https://go.dev/dl/go%s.%s-%s.tar.gz", version, goos, goarch)

	tarball, err := os.CreateTemp("", "tori-")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Downloading go%s.%s-%s...\n", version, goos, goarch)
	if err := download(url, tarball.Name()); err != nil {
		log.Fatal(err)
	}

	home, err := filepath.Abs(os.Getenv("TORI_HOME"))
	if err != nil {
		log.Fatal(err)
	}

	target := filepath.Join(home, "versions", version)
	if err := os.MkdirAll(target, 0755); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Extracting go%s.%s-%s...\n", version, goos, goarch)
	if err := extract(tarball, target); err != nil {
		log.Fatal(err)
	}
	os.Remove(tarball.Name())

	symlink := filepath.Join(home, "bin")
	if _, err := os.Lstat(symlink); err == nil {
		if err := os.Remove(symlink); err != nil {
			log.Fatal(err)
		}
	}
	os.Symlink(filepath.Join(target, "go", "bin"), symlink)
}
