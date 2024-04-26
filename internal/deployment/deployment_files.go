package deployment

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Copy a folder from the deploy location into ~/.teu/deployments/<app_name>
// This is used to copy the deployment folder into the deployments directory
// This is so we can keep track of what has been deployed already.
// Function takes a string as a path and an application name
func CopyDeployment(path string, app_name string) (string, error) {
	
	// Set the deploy location to ~/.teu/deployments/<app_name>
	home_dir, _ := os.UserHomeDir()
	deploy_location := filepath.Join(home_dir, ".teu", "deployments", strings.ToLower(app_name))

	// copy the deployment folder into the deploy location
	err := copyDir(path, deploy_location)
	if err != nil {
		return "", err
	}

	return deploy_location, nil
}


func copyDir(src, dest string) error {
	// Open the source directory
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !srcInfo.IsDir() {
		return nil // Or return an error if src is not a directory
	}

	// Create the destination directory if it doesn't exist
	if err := os.MkdirAll(dest, srcInfo.Mode()); err != nil {
		return err
	}

	// Get a list of all files and directories in the source directory
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// Copy each file/directory to the destination directory
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			// Recursively copy subdirectories
			if err := copyDir(srcPath, destPath); err != nil {
				return err
			}
		} else {
			// Copy files
			if err := copyFile(srcPath, destPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
			return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
			return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
			return err
	}

	return destFile.Sync()
}

// Remove a folder
func DeleteDeployment(app_name string) (error) {
	deployLocation := filepath.Join(os.Getenv("HOME"), ".teu", "deployments", strings.ToLower(app_name))

	err := os.RemoveAll(deployLocation)
	if err != nil {
		return err
	}

	return nil
}