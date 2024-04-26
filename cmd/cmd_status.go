package cmd

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v2"

	"github.com/vojkovic/teu/internal/db"
	"github.com/vojkovic/teu/internal/deployment"
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

	lastPull = ConvertUnixToHumanReadable(lastPull)

	lastCommit, err := db.GetRepositoryLastCommitFromDatabase()
	if err != nil {
		return err
	}

	lastCommit = ConvertUnixToHumanReadable(lastCommit)

	table.DefaultHeaderFormatter = func(format string, vals ...interface{}) string {
		return strings.ToUpper(fmt.Sprintf(format, vals...))
	}
	
	headerFmt := color.New(color.FgWhite, color.Underline).SprintfFunc()
  columnFmt := color.New(color.FgCyan).SprintfFunc()

	tbl := table.New("Information","").WithPadding(3).WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	
	tbl.AddRow("Node Name",deployment.GetHostname())
	tbl.AddRow("Repository", git_repository)
	tbl.AddRow("Last Pull",lastPull)
	tbl.AddRow("Last Commit",lastCommit)

	tbl.Print()

	return nil
}


// ConvertUnixToHumanReadable converts Unix timestamp to human-readable format
func ConvertUnixToHumanReadable(unixTimeStr string) string {
	unixTime, err := strconv.ParseInt(unixTimeStr, 10, 64)
	if err != nil {
		return "Invalid Unix timestamp"
	}

	// Convert Unix timestamp to time object
	timestamp := time.Unix(unixTime, 0)

	// Calculate time difference
	diff := time.Since(timestamp)

	// Convert time difference to human-readable format
	switch {
	case diff.Seconds() < 60:
		if int(diff.Seconds()) == 1 {
			return "1 second ago"
		}
		return fmt.Sprintf("%d seconds ago", int(diff.Seconds()))
	case diff.Minutes() < 60:
		if int(diff.Minutes()) == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", int(diff.Minutes()))
	case diff.Hours() < 24:
		if int(diff.Hours()) == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", int(diff.Hours()))
	default:
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}
}