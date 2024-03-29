package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <process name filter>")
		return
	}

	processNameFilter := os.Args[1]

	cmd := exec.Command("tasklist", "/fo", "csv", "/nh")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error running tasklist command: %v\n", err)
		return
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Split(line, ",")
		if len(fields) < 5 {
			continue
		}
		pid := strings.TrimSpace(fields[1])
		session := strings.TrimSpace(fields[2])
		name := strings.TrimSpace(fields[0])

		if strings.Contains(name, processNameFilter) {
			fmt.Printf("PID: %s, Session: %s, Name: %s\n", pid, session, name)
		}
	}
}
