package deployment

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/vojkovic/teu/internal/db"
)

func DeploymentDown() error {
	applications, err := db.GetAllApplicationsFromDatabase()
	if err != nil {
		return err
	}

	if len(applications) == 0 {
		PrintNoApplicationsToStop()
		return nil
	}

	for i, app := range applications {
		PrintStopping(app, i+1, len(applications))

		compose_file_path := filepath.Join(filepath.Join(os.Getenv("HOME"), ".teu", "deployments", app, "docker-compose.yml"))

		err = DeploymentComposeDown(compose_file_path, strings.ToLower(app))
		if err != nil {
			PrintStopFailure(app)
			return err
		}

		err := db.RemoveApplicationFromDatabase(app)
		if err != nil {
			PrintStartFailure(app)
			return err
		}

		if err := DeleteDeployment(app); err != nil {
			return err
		}

		// print a newline if there are more applications to remove
		if i < len(applications) - 1 {
			fmt.Println()
		}
	}

	err = db.SetRepositoryInDatabase("")
	if err != nil {
		return err
	}

	err = db.SetRepositoryTokenInDatabase("")
	if err != nil {
		return err
	}

	err = db.SetRepositoryLastCommitInDatabase("")
	if err != nil {
		return err
	}

	err = db.SetRepositoryLastPullInDatabase("")
	if err != nil {
		return err
	}

	err = os.RemoveAll(filepath.Join(os.Getenv("HOME"), ".teu", "repo"))
	if err != nil {
		return err
	}
	return nil;
}