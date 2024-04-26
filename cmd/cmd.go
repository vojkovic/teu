package cmd

import "github.com/urfave/cli/v2"

func CreateCommands() []*cli.Command {
	return []*cli.Command{
		createJoin(),
		createLeave(),
		createStatus(),
		createDaemon(),
	}
}