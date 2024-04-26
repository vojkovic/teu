package deployment

import (
	"fmt"

	"github.com/fatih/color"
)

func PrintStarting(appName string, i, total int) {
	fmt.Println(
		color.New(color.FgGreen).Sprint("•"),
		color.New(color.FgWhite).Sprint("Deploying"),
		color.New(color.FgGreen, color.Bold).Sprint(appName),
		color.New(color.FgBlack).Sprint("(" + fmt.Sprintf("%d", i) + "/" + fmt.Sprintf("%d", total) + ")"),
	)
}

func PrintReDeploying(appName string) {
	fmt.Println(
		color.New(color.FgGreen).Sprint("•"),
		color.New(color.FgWhite).Sprint("Re-Deploying"),
		color.New(color.FgGreen, color.Bold).Sprint(appName),
	)
}

func PrintAlreadyRunning(appName string) {
	fmt.Println(
		color.New(color.FgGreen).Sprint("•"),
		color.New(color.FgWhite).Sprint("Already running"),
		color.New(color.FgGreen, color.Bold).Sprint(appName),
	)
}

func PrintNoApplicationsToStop() {
	fmt.Println(
		color.New(color.FgYellow).Sprint("•"),
		color.New(color.FgWhite).Sprint("No applications to stop"),
	)
}

func PrintStartFailure(appName string) {
	fmt.Println(
		color.New(color.FgRed).Sprint("•"),
		color.New(color.FgWhite).Sprint("Failed to start"),
		color.New(color.FgGreen, color.Bold).Sprint(appName),
	)
}

func PrintStopFailure(appName string) {
	fmt.Println(
		color.New(color.FgRed).Sprint("•"),
		color.New(color.FgWhite).Sprint("Failed to stop"),
		color.New(color.FgGreen, color.Bold).Sprint(appName),
	)
}

func PrintStopping(appName string, i, total int) {
	fmt.Println(
		color.New(color.FgGreen).Sprint("•"),
		color.New(color.FgWhite).Sprint("Stopping"),
		color.New(color.FgGreen, color.Bold).Sprint(appName),
		color.New(color.FgBlack).Sprint("(" + fmt.Sprintf("%d", i) + "/" + fmt.Sprintf("%d", total) + ")"),
	)
}

func PrintSecretDecrypting(enc_secret_name, secret_name string) {
	fmt.Println(
		color.New(color.FgBlue).Sprint("•"),
		color.New(color.FgWhite).Sprint("Decrypting"),
		color.New(color.FgCyan).Sprint(enc_secret_name), 
		color.New(color.FgWhite).Sprint("→"),
		color.New(color.FgCyan).Sprint(secret_name),
	)
}

func PrintSecretFailedToDecrypt(enc_secret_name string) {
	fmt.Println(
		color.New(color.FgRed).Sprint("•"),
		color.New(color.FgWhite).Sprint("Failed to decrypt"),
		color.New(color.FgCyan).Sprint(enc_secret_name),
	)
}
