// services_vulns.go performs scans on a specified target to check for vulnerabilities and Bitcoin service information

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var bitcoinServiceDetected bool = false

func main() {
	if len(os.Args) < 2 {
		fmt.Println("[-] Usage: go run services_vulns.go <TARGET_IP> [--skip-vuln]")
		os.Exit(1)
	}
	target := os.Args[1]
	skipVulnCheck := false

	if len(os.Args) > 2 && os.Args[2] == "--skip-vuln" {
		skipVulnCheck = true
	}

	if !skipVulnCheck {
		results := nmapScriptVulnScan(target)
		fmt.Println("[*] Vulnerability Scan Results:\n")
		fmt.Println(parseNmapVulnOutput(results))
	}

	if bitcoinServiceDetected || skipVulnCheck {
		fmt.Println("\n[*] Checking for bitcoin info and node addresses\n\n[*]Bitcoin-info")
		results := nmapScriptBitcoinInfoScan(target)
		fmt.Println(parseNmapBitcoinOutput(results) + "\n")

		fmt.Println("[*]Bitcoin-getaddr")
		results = nmapScriptBitcoinGetaddrScan(target)
		fmt.Println(parseNmapBitcoinOutput(results))
	}
}

func nmapScriptVulnScan(target string) string {
	args := []string{"-Pn", "--script", "vuln", target}

	output := runNmapScan(args)

	outputStr := string(output)
	return outputStr
}

func nmapScriptBitcoinInfoScan(target string) string {
	args := []string{"-p", "8333", "--script", "bitcoin-info", target}

	output := runNmapScan(args)

	outputStr := string(output)
	return outputStr
}

func nmapScriptBitcoinGetaddrScan(target string) string {
	args := []string{"-p", "8333", "--script", "bitcoin-getaddr", target}

	output := runNmapScan(args)

	outputStr := string(output)
	return outputStr
}

func runNmapScan(args []string) []byte {
	cmd := exec.Command("nmap", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[-] Nmap scan failed: %v\n", err)
		os.Exit(1)
	}

	return output
}

func parseNmapVulnOutput(output string) string {
	scanner := bufio.NewScanner(strings.NewReader(output))
	var services string
	capture := false

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "rDNS record for") {
			line = "\033[35m" + line + "\033[0m" // Magenta text
			services += line + "\n\n"
		}

		if strings.HasPrefix(line, "PORT") {
			capture = true
			services += line + "\n"
			continue
		}
		if capture {
			if line == "" {
				break
			}
			if strings.Contains(line, "8333/tcp") && strings.Contains(line, "bitcoin") {
				bitcoinServiceDetected = true
			}
			if strings.Contains(line, "VULNERABLE") {
				line = "\033[31m" + line + "\033[0m" // Red text
			}
			if strings.Contains(line, "CVE") {
				line = "\033[32m" + line + "\033[0m" // Green text
			}
			if strings.Contains(line, "Risk factor: Critical") {
				line = "\033[41;33m" + line + "\033[0m" // Yellow background, red text
			}
			if strings.Contains(line, "Risk factor: High") {
				line = "\033[31m" + line + "\033[0m" // Red text
			}
			if strings.Contains(line, "Risk factor: Medium") {
				line = "\033[33m" + line + "\033[0m" // Yellow text
			}
			if strings.Contains(line, "Risk factor: Low") {
				line = "\033[34m" + line + "\033[0m" // Blue text
			}
			services += line + "\n"
		}
	}

	return services
}

func parseNmapBitcoinOutput(output string) string {
	scanner := bufio.NewScanner(strings.NewReader(output))
	var services string
	capture := false

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "(not scanned)") {
			line = "\033[34m" + line + "\033[0m" // Blue text
			services += line + "\n\n"
		}

		if strings.Contains(line, "rDNS record for") {
			line = "\033[35m" + line + "\033[0m" // Magenta text
			services += line + "\n\n"
		}

		if strings.HasPrefix(line, "PORT") {
			capture = true
			services += line + "\n"
			continue
		}
		if capture {
			if line == "" {
				break
			}
			services += line + "\n"
		}
	}

	return services
}
