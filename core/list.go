package core

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/hashicorp/go-version"
)

const (
	defaultRefs    = "https://go.googlesource.com/go/+refs"
	tagExpression  = `tags/go\d+\.\d+(?:beta\d+|rc\d+|\.\d+)?`
	minimumVersion = "1.11" // preliminary support for modules
)

func getVersionsOffline() ([]*version.Version, error) {
	home := getHome()
	target := filepath.Join(home, "versions")

	entries, err := os.ReadDir(target)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	versions := make([]*version.Version, len(entries))
	for i, entry := range entries {
		v, err := version.NewVersion(entry.Name())
		if err != nil {
			return nil, err
		}
		versions[i] = v
	}

	sort.Sort(sort.Reverse(version.Collection(versions)))
	return versions, nil
}

func getVersionsOnline() ([]*version.Version, error) {
	res, err := http.Get(defaultRefs)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	refs, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	tags := regexp.MustCompile(tagExpression).FindAll(refs, -1)

	versions := make([]*version.Version, 0, len(tags))
	min := version.Must(version.NewVersion(minimumVersion))
	for _, tag := range tags {
		v, err := version.NewVersion(string(tag[7:]))
		if err != nil {
			return nil, err
		}

		if !v.LessThan(min) {
			versions = append(versions, v)
		}
	}

	sort.Sort(sort.Reverse(version.Collection(versions)))
	return versions, nil
}

func list(versions []*version.Version, limit int) {
	if l := len(versions); limit < 1 || l < limit {
		limit = l
	}

	for _, version := range versions[:limit] {
		fmt.Println(version.Original())
	}
}

func List(online bool, limit int) (err error) {
	var versions []*version.Version

	if !online {
		versions, err = getVersionsOffline()
	} else {
		versions, err = getVersionsOnline()
	}
	if err != nil {
		return
	}

	list(versions, limit)
	return
}
