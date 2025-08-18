package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
)

func encodeForGopher(raw string, doubleEncode bool) string {
	encoded := url.PathEscape(raw)

	if doubleEncode {
		encoded = url.PathEscape(encoded)
	}

	return encoded
}

func main() {
	filePath := flag.String("file", "", "Path to raw HTTP request file")
	doubleEncode := flag.Bool("double-encode", false, "Apply double URL encoding")
	noSuffix := flag.Bool("no-suffix", false, "Remove trailing newline")
	host := flag.String("host", "127.0.0.1", "Target host for gopher payload")
	port := flag.String("port", "80", "Target port for gopher payload")
	flag.Parse()

	if *filePath == "" {
		fmt.Println("Usage: go run main.go -file request.txt [-double-encode] [-host 127.0.0.1] [-port 80]")
		os.Exit(1)
	}

	file, err := os.Open(*filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var requestLines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		requestLines = append(requestLines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	rawRequest := strings.Join(requestLines, "\r\n")

	if !strings.HasSuffix(rawRequest, "\r\n\r\n") && !(*noSuffix) {
		rawRequest += "\r\n\r\n"
	}

	encoded := encodeForGopher(rawRequest, *doubleEncode)

	delimiter := ":"
	if *doubleEncode {
		delimiter = "%3a"
	}

	gopherURL := fmt.Sprintf("gopher%s//%s%s%s/_%s", delimiter, *host, delimiter, *port, encoded)

	fmt.Println(gopherURL)
}
