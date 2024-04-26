package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v2"

	"github.com/vojkovic/teu/internal/db"
	"github.com/vojkovic/teu/pkg/common"
)

func createStatus() *cli.Command {
	return &cli.Command{
		Name:  "status",
		Usage: "Show the status of the deployment of the node.",
		Action: status,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "json",
				Usage: "Output status in JSON format",
			},
		},
	}
}

func status(ctx *cli.Context) error {
	if ctx.Bool("json") {
		log.Println("TODO: Outputting status in JSON format")
		return nil
	}

	repo, err := db.GetRepositoryFromDatabase()
	if err != nil {
		return err
	}

	if repo == "" {
		return fmt.Errorf("No repository is set. Run 'teu join' to join a cluster.")
	}

	err = printNodeInfoTable()
	if err != nil {
		return err
	}

	println("\n")

	return nil
}

func printNodeInfoTable() error {

	git_repository, err := db.GetRepositoryFromDatabase()
	if err != nil {
		return err
	}
	
	lastPull, err := db.GetRepositoryLastPullFromDatabase()
	if err != nil {
		return err
	}

	lastPull = common.ConvertUnixToHumanReadable(lastPull)

	lastCommit, err := db.GetRepositoryLastCommitFromDatabase()
	if err != nil {
		return err
	}

	lastCommit = common.ConvertUnixToHumanReadable(lastCommit)

	table.DefaultHeaderFormatter = func(format string, vals ...interface{}) string {
		return strings.ToUpper(fmt.Sprintf(format, vals...))
	}
	
	headerFmt := color.New(color.FgWhite, color.Underline).SprintfFunc()
  columnFmt := color.New(color.FgCyan).SprintfFunc()

	tbl := table.New("Information","").WithPadding(3).WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	
	tbl.AddRow("Node Name",common.GetHostname())
	tbl.AddRow("Repository", git_repository)
	tbl.AddRow("Last Pull",lastPull)
	tbl.AddRow("Last Commit",lastCommit)

	tbl.Print()

	return nil
}
