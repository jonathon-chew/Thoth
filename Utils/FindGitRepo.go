package utils

import (
	"os"
	"path/filepath"
)

func FindGitRepos(root string) []string {
	var repos []string
	filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() && d.Name() == ".git" {
			repos = append(repos, filepath.Dir(path))
			return filepath.SkipDir // stop traversing this subdir
		}
		return nil
	})
	return repos
}
