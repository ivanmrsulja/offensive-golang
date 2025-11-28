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

func generatePattern(length int) string {
	var pattern strings.Builder
	chars := []byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
		'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

	for pattern.Len() < length {
		for _, first := range chars {
			for second := byte('a'); second <= 'z'; second++ {
				for num := 0; num <= 9; num++ {
					if pattern.Len() >= length {
						return pattern.String()[:length]
					}
					pattern.WriteByte(first)
					pattern.WriteByte(second)
					pattern.WriteString(strconv.Itoa(num))
				}
			}
		}
	}

	return pattern.String()[:length]
}

func findOffset(hexValue string, pattern string) (int, error) {
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

	offset := strings.Index(pattern, searchString)

	if offset == -1 {
		return -1, fmt.Errorf("pattern not found in the sequence")
	}

	return offset, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  go run script.go pattern [length]           - Generate pattern (default: 1000)")
		fmt.Println("  go run script.go offset <hex_value> [pattern_length] - Find offset for hex value")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  go run script.go pattern")
		fmt.Println("  go run script.go pattern 5000")
		fmt.Println("  go run script.go offset 0x6F43376F")
		fmt.Println("  go run script.go offset 0x6F43376F 8000")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "pattern":
		length := 1000 // default length
		if len(os.Args) >= 3 {
			var err error
			length, err = strconv.Atoi(os.Args[2])
			if err != nil {
				fmt.Printf("Error: Invalid length '%s'. Please provide a valid integer\n", os.Args[2])
				os.Exit(1)
			}
			if length <= 0 {
				fmt.Printf("Error: Length must be positive, got %d\n", length)
				os.Exit(1)
			}
		}
		pattern := generatePattern(length)
		fmt.Println(pattern)

	case "offset":
		if len(os.Args) < 3 {
			fmt.Println("Error: Please provide a hex value for offset calculation")
			fmt.Println("Usage: go run script.go offset <hex_value> [pattern_length]")
			os.Exit(1)
		}

		hexValue := os.Args[2]
		patternLength := 10000 // default pattern length

		if len(os.Args) >= 4 {
			var err error
			patternLength, err = strconv.Atoi(os.Args[3])
			if err != nil {
				fmt.Printf("Error: Invalid pattern length '%s'. Please provide a valid integer\n", os.Args[3])
				os.Exit(1)
			}
			if patternLength <= 0 {
				fmt.Printf("Error: Pattern length must be positive, got %d\n", patternLength)
				os.Exit(1)
			}
		}

		pattern := generatePattern(patternLength)
		offset, err := findOffset(hexValue, pattern)
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
