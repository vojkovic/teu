package common

import (
	"fmt"
	"os"
	"path/filepath"
)

// checks that teu.yml exists from a given path and that it's a file and not a directory.
func TeuFileExists (path string) error {
	info, err := os.Stat(filepath.Clean(path) + "/teu.yml")
	if err != nil {
		return fmt.Errorf("could not find teu.yml in %s: %s", path, err)
	}

	if info.IsDir() {
		return fmt.Errorf("teu.yml %s is a directory, expected file", path)
	}

	return nil
}

// checks that the repository folder exists from a given path and that it's a directory and not a file.
func RepoFolderExists (path string) error {
	info, err := os.Stat(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("could not find the repo in %s: %s", path, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("repo %s is a file, expected direcotry", path)
	}

	return nil
}

func IsTeuRepo (path string) error {
	err := TeuFileExists(path)
	if err != nil {
		return err
	}
	
	
	return nil
}