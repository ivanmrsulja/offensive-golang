package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <PID>")
		return
	}
	pid := os.Args[1]

	cmd := exec.Command("tasklist", "/m", "/fi", fmt.Sprintf("PID eq %s", pid))
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error running tasklist command: %v\n", err)
		return
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 1 {
			dllName := fields[1]
			if !strings.HasSuffix(dllName, ".dll,") {
				continue
			}
			fmt.Printf("DLL Name: %s\n", dllName)
		}
	}
}
