// The offset is used to determine how many bytes are needed to overwrite
// the buffer and how much space we have around our shellcode. This can
// help us determine the exact number of bytes to reach the EIP

package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func generatePattern() string {
	var pattern strings.Builder

	for first := byte('A'); first <= 'D'; first++ {
		for second := byte('a'); second <= 'z'; second++ {
			for num := 0; num <= 9; num++ {
				pattern.WriteByte(first)
				pattern.WriteByte(second)
				pattern.WriteString(strconv.Itoa(num))
			}
		}
	}

	return pattern.String()
}

func findOffset(hexValue string) (int, error) {
	hexValue = strings.TrimPrefix(hexValue, "0x")

	bytes, err := hex.DecodeString(hexValue)
	if err != nil {
		return -1, fmt.Errorf("invalid hex value: %v", err)
	}

	// Reverse for little-endian
	var reversed strings.Builder
	for i := len(bytes) - 1; i >= 0; i-- {
		reversed.WriteByte(bytes[i])
	}
	searchString := reversed.String()

	pattern := generatePattern()
	offset := strings.Index(pattern, searchString)

	if offset == -1 {
		return -1, fmt.Errorf("pattern not found in the sequence")
	}

	return offset, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  go run script.go pattern                    - Generate pattern")
		fmt.Println("  go run script.go offset <hex_value>        - Find offset for hex value")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  go run script.go pattern")
		fmt.Println("  go run script.go offset 0x6F43376F")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "pattern":
		pattern := generatePattern()
		fmt.Println(pattern)

	case "offset":
		if len(os.Args) < 3 {
			fmt.Println("Error: Please provide a hex value for offset calculation")
			fmt.Println("Usage: go run script.go offset <hex_value>")
			os.Exit(1)
		}

		hexValue := os.Args[2]
		offset, err := findOffset(hexValue)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Offset: %d\n", offset)

	default:
		fmt.Printf("Error: Unknown command '%s'\n", command)
		fmt.Println("Valid commands: pattern, offset")
		os.Exit(1)
	}
}
