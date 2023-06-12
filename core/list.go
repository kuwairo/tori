package core

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/hashicorp/go-version"
)

func listOffline() ([]*version.Version, error) {
	home := getHome()
	target := filepath.Join(home, "versions")

	entries, err := os.ReadDir(target)
	if err != nil {
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

	sort.Sort(version.Collection(versions))
	return versions, nil
}

func List(online bool) error {
	if !online {
		versions, err := listOffline()
		if err != nil {
			return err
		}

		for _, version := range versions {
			fmt.Println(version.Original())
		}
	}

	// TODO: listOnline()

	return nil
}
