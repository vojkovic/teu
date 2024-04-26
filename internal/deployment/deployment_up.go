package deployment

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/vojkovic/teu/config"
	"github.com/vojkovic/teu/internal/db"
	"github.com/vojkovic/teu/pkg/common"
)

func DeploymentUp(path string) error {
	
	err := common.IsTeuRepo(path)
	if err != nil {
		return err
	}

	c, err := config.LoadConfig(path)
	if err != nil {
		return err
	}

	for i, app := range c.Applications {
		PrintStarting(app.Name, i+1, len(c.Applications))

		if db.IsApplicationInDatabase(app.Name) {
			hash, err := GetApplicationHash(filepath.Join(path, app.Deploy))
			if err != nil {
				PrintStartFailure(app.Name)
				return err
			}

			previous_hash, err := db.GetApplicationHashFromDatabase(app.Name)
			if err != nil {
				PrintStartFailure(app.Name)
				return err
			}

			if hash == previous_hash {
				PrintAlreadyRunning(app.Name)
				continue
			} else {
				// Remove the application from the database, stop application and remove the deployment
				err := db.RemoveApplicationFromDatabase(app.Name)
				if err != nil {
					PrintStartFailure(app.Name)
					return err
				}

				deploy_location, err := CopyDeployment(filepath.Join(path, app.Deploy), app.Name)
				if err != nil {
					PrintStartFailure(app.Name)
					return err
				}
		
				compose_file_path := filepath.Join(deploy_location, "docker-compose.yml")

				err = DeploymentComposeDown(compose_file_path, strings.ToLower(app.Name))
				if err != nil {
					PrintStartFailure(app.Name)
					return err
				}
			}
		}
		
		deploy_location, err := CopyDeployment(filepath.Join(path, app.Deploy), app.Name)
		if err != nil {
			PrintStartFailure(app.Name)
			return err
		}

		err = DecryptSecretsInApp(app.Secrets, c.Teu.AgeSecretKey, deploy_location, app.Name)
		if err != nil {
			return err
		}

		compose_file_path := filepath.Join(deploy_location, "docker-compose.yml")
		
		err = DeploymentComposeUp(compose_file_path, strings.ToLower(app.Name))
		if err != nil {
			os.RemoveAll(deploy_location)
			PrintStartFailure(app.Name)
			return err
		}

		// print a newline if there are more applications to deploy
		if i < len(c.Applications) - 1 {
			fmt.Println()
		}

		if !db.IsApplicationInDatabase(app.Name) {
			hash, err := GetApplicationHash(filepath.Join(path, app.Deploy))
			if err != nil {
				PrintStartFailure(app.Name)
				return err
			}

			err = db.AddApplicationToDatabase(app.Name, hash)
			if err != nil {
				PrintStartFailure(app.Name)
				return err
			}
		}
	}
	return nil;
}