// This script sets up a cron job on Linux to execute `whoami` every minute.
// The cron job can be used to ensure that a command or script runs regularly.

package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func main() {
	err := createCronJob()
	if err != nil {
		fmt.Println("Failed to create cron job:", err)
	} else {
		fmt.Println("Cron job created successfully.")
	}
}

// createCronJob creates a cron job to run the `whoami` command every minute.
func createCronJob() error {
	cronCommand := "* * * * * whoami\n"
	crontabCmd := exec.Command("crontab", "-l")
	currentCron, _ := crontabCmd.Output()

	// Avoid duplicate cron entry if it already exists
	if strings.Contains(string(currentCron), cronCommand) {
		return nil
	}

	cmd := exec.Command("bash", "-c", fmt.Sprintf(`(crontab -l 2>/dev/null; echo "%s") | crontab -`, cronCommand))
	return cmd.Run()
}
