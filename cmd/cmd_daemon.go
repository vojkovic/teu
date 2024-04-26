package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/vojkovic/teu/internal/daemon"
)

func createDaemon() *cli.Command {
	return &cli.Command{
		Name:  "daemon",
		Usage: "Run the teu daemon",
		Hidden: true,
		Action: startDaemon,
	}
}

func startDaemon(ctx *cli.Context) error {
	err := daemon.StartDaemon()
	if err != nil {
		return err
	}

	return nil
}