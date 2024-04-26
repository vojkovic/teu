package daemon

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/vojkovic/teu/config"
	"github.com/vojkovic/teu/internal/db"
	"github.com/vojkovic/teu/internal/deployment"
	"github.com/vojkovic/teu/internal/git"
	"github.com/vojkovic/teu/pkg/common"
)

func reconcile() error {

	repo, token, err := GetRepositoryAndTokenFromDatabase()
	if err != nil {
		return err
	}

	err = git.Update(repo, token)
	if err != nil {
		return err
	}

	teu_repository_path := filepath.Join(os.Getenv("HOME"), ".teu", "repo")
	// read teu.yml
	// compare teu.yml with the database
	// if there are any differences, update the database
	
	// read teu.yml, list applications as a slice of strings
	err = common.IsTeuRepo(teu_repository_path)
	if err != nil {
		return err
	}

	c, err := config.LoadConfig(teu_repository_path)
	if err != nil {
		return err
	}

	new_deploy, delete_old, err := GetDeployments(teu_repository_path, c)
	if err != nil {
		return err
	}

	fmt.Println("New apps: ", new_deploy)
	fmt.Println("Delete apps: ", delete_old)
	
	// delete the old apps
	for _, app := range delete_old {
		fmt.Println("Deleting app: ", app)

		compose_file_path := filepath.Join(filepath.Join(os.Getenv("HOME"), ".teu", "deployments", app, "docker-compose.yml"))

		err = deployment.DeploymentComposeDown(compose_file_path, strings.ToLower(app))
		if err != nil {
			return err
		}

		err = db.RemoveApplicationFromDatabase(app)
		if err != nil {
			return err
		}
	}

	// add the new apps
	for _, app := range new_deploy {

		fmt.Println("Deploying app: ", app)
		var deploy_location string = ""
		var secrets[] string = make([]string, 0)
		var hash string = ""
		for _, new_app := range c.Applications {
			if new_app.Name == app {
				deploy_location, err = deployment.CopyDeployment(filepath.Join(teu_repository_path, new_app.Deploy), new_app.Name)
				if err != nil {
					return err
				}
				hash, err = deployment.GetApplicationHash(filepath.Join(teu_repository_path, new_app.Deploy))
				if err != nil {
					return err
				}
				secrets = new_app.Secrets
			}
		}

		if deploy_location == "" {
			return fmt.Errorf("deploy location for %s not found", app)
		}


		err = deployment.DecryptSecretsInApp(secrets, c.Teu.AgeSecretKey, deploy_location, app)
		if err != nil {
			return err
		}

		compose_file_path := filepath.Join(deploy_location, "docker-compose.yml")
		
		err = deployment.DeploymentComposeUp(compose_file_path, strings.ToLower(app))
		if err != nil {
			os.RemoveAll(deploy_location)
			return err
		}

		err = db.AddApplicationToDatabase(app, hash)
		if err != nil {
			return err
		}
	}


	return nil
}

func GetRepositoryAndTokenFromDatabase() (string, string, error) {
	// check if database contains a repository
	// if so, run reconcile for that repository
	repo, err := db.GetRepositoryFromDatabase()
	if err != nil {
		return "", "", err
	}

	if repo == "" {
		return "", "", fmt.Errorf("No repository found in database")
	}

	token, err := db.GetRepositoryTokenFromDatabase()
	if err != nil {
		return "", "", err
	}

	return repo, token, nil
}
