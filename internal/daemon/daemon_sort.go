package daemon

import (
	"fmt"
	"path/filepath"

	"github.com/vojkovic/teu/config"
	"github.com/vojkovic/teu/internal/db"
	"github.com/vojkovic/teu/internal/deployment"
)

// GetSameDeploy finds apps that are present in both old and new lists.
func GetSameDeploy(oldApps, newApps []string) []string {
	sameDeploy := make([]string, 0)
	for _, app := range newApps {
		for _, oldApp := range oldApps {
			if app == oldApp {
				sameDeploy = append(sameDeploy, app)
				break
			}
		}
	}
	return sameDeploy
}

// DeleteOldApps finds apps in the old list not present in the new list.
func DeleteOldApps(oldApps, newApps []string) []string {
	oldDeploy := make([]string, 0)
	for _, app := range oldApps {
		found := false
		for _, newApp := range newApps {
			if app == newApp {
				found = true
				break
			}
		}
		if !found {
			oldDeploy = append(oldDeploy, app)
		}
	}
	return oldDeploy
}

// NewDeployApps finds apps in the new list not present in the old list.
func NewDeployApps(oldApps, newApps []string) []string {
	newDeploy := make([]string, 0)
	for _, app := range newApps {
		found := false
		for _, oldApp := range oldApps {
			if app == oldApp {
				found = true
				break
			}
		}
		if !found {
			newDeploy = append(newDeploy, app)
		}
	}
	return newDeploy
}

// GetDeployments compares old and new app lists, returns apps to deploy and delete.
func GetDeployments(teuRepoPath string, c *config.TeuConfig) ([]string, []string, error) {
	var newApps []string
	var oldApps []string

	for _, app := range c.Applications {
		newApps = append(newApps, app.Name)
	}

	// Retrieve old apps from the database.
	oldApps, err := db.GetAllApplicationsFromDatabase()
	if err != nil {
		return nil, nil, err
	}

	// Find apps to deploy and delete.
	newDeploy := NewDeployApps(oldApps, newApps)
	deleteOld := DeleteOldApps(oldApps, newApps)

	// Check for changes in same apps.
	sameDeploy := GetSameDeploy(oldApps, newApps)
	for _, app := range sameDeploy {
		hash, err := db.GetApplicationHashFromDatabase(app)
		if err != nil {
			return nil, nil, err
		}

		// Find app deploy location.
		var deployPath string
		for _, newApp := range c.Applications {
			if newApp.Name == app {
				deployPath = newApp.Deploy
				break
			}
		}

		if deployPath == "" {
			return nil, nil, fmt.Errorf("deploy path for %s not found in teu.yml", app)
		}

		newHash, err := deployment.GetApplicationHash(filepath.Join(teuRepoPath, deployPath))
		if err != nil {
			return nil, nil, err
		}

		// Update deployment lists if hash differs.
		if hash != newHash {
			newDeploy = append(newDeploy, app)
			deleteOld = append(deleteOld, app)
		}
	}

	return newDeploy, deleteOld, nil
}
