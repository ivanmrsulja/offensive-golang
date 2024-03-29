package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <driver name filter>")
		return
	}
	driverNameFilter := os.Args[1]

	cmd := exec.Command("sc", "query", "type=", "driver")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error running sc query command: %v\n", err)
		return
	}

	lines := strings.Split(string(output), "\n")
	skipDriver := false
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			skipDriver = false
			continue
		}

		if skipDriver {
			continue
		}

		if strings.HasPrefix(fields[0], "SERVICE_NAME") {
			if !strings.Contains(strings.ToLower(fields[1]), strings.ToLower(driverNameFilter)) {
				skipDriver = true
				continue
			}
			driverName := strings.TrimSpace(fields[1])
			fmt.Printf("Driver Name: %s\n", driverName)
		} else if strings.HasPrefix(fields[0], "TYPE") {
			driverType := strings.TrimSpace(strings.Join(fields[2:], " "))
			fmt.Printf("Type: %s\n", driverType)
		} else if strings.HasPrefix(fields[0], "STATE") {
			driverState := strings.TrimSpace(strings.Join(fields[2:], " "))
			fmt.Printf("State: %s\n", driverState)
		} else if strings.HasPrefix(fields[0], "WIN32_EXIT_CODE") {
			exitCode := strings.TrimSpace(strings.Join(fields[2:], " "))
			fmt.Printf("Win32 Exit Code: %s\n", exitCode)
		} else if strings.HasPrefix(fields[0], "SERVICE_EXIT_CODE") {
			serviceExitCode := strings.TrimSpace(strings.Join(fields[2:], " "))
			fmt.Printf("Service Exit Code: %s\n", serviceExitCode)
		} else if strings.HasPrefix(fields[0], "CHECKPOINT") {
			checkpoint := strings.TrimSpace(strings.Join(fields[2:], " "))
			fmt.Printf("Checkpoint: %s\n", checkpoint)
		} else if strings.HasPrefix(fields[0], "WAIT_HINT") {
			waitHint := strings.TrimSpace(strings.Join(fields[2:], " "))
			fmt.Printf("Wait Hint: %s\n\n", waitHint)
		}
	}
}
