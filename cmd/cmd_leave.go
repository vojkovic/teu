package cmd

import (
	"github.com/urfave/cli/v2"
	"github.com/vojkovic/teu/internal/deployment"
)

func createLeave() *cli.Command {
	return &cli.Command{
		Name:  "leave",
		Usage: "Stop all applications and remove them from the database. Stop syncing the repository.",
		Action: leave,
	}
}

func leave(ctx *cli.Context) error {
	err := deployment.DeploymentDown()
	if err != nil {
		return err
	}

	return nil
}