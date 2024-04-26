package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/vojkovic/teu/internal/db"
	"github.com/vojkovic/teu/internal/deployment"
	"github.com/vojkovic/teu/internal/git"
)

func createJoin() *cli.Command {
	return &cli.Command{
		Name:  "join",
		Usage: "Join this node to the cluster",
		Action: join,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "repo",
				Usage: "[Required] The repository where teu.yml is stored.",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "token",
				Usage: "The Git Authentication Token to use to pull from the repository.",
			},
		},
	}
}

func join(ctx *cli.Context) error {

	repo := ctx.String("repo")
	token := ctx.String("token")

	if token == "" {
		log.Println("No token provided, using anonymous access")
	}

	err := db.DatabaseInit()
	if err != nil {
		return err
	}

	err = deployment.DeploymentDown()
	if err != nil {
		return err
	}

	err = git.Update(repo, token)
	if err != nil {
		return err
	}

	err = db.SetRepositoryInDatabase(repo)
	if err != nil {
		return err
	}

	if token != "" {
		err = db.SetRepositoryTokenInDatabase(token)
		if err != nil {
			return err
		}
	}
	
	currentTimeInUnixSeconds := fmt.Sprintf("%d", time.Now().Unix())
	
	err = db.SetRepositoryLastCommitInDatabase(currentTimeInUnixSeconds)
	if err != nil {
		return err
	}

	err = db.SetRepositoryLastPullInDatabase(currentTimeInUnixSeconds)
	if err != nil {
		return err
	}

	path := filepath.Join(os.Getenv("HOME"), ".teu", "repo")
	err = deployment.DeploymentUp(path)
	if err != nil {
		return err
	}

	return nil
}