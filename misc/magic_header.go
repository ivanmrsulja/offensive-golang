// simple script to generate empty files with custom magic headers in order to bypass validation
// used to solve some of the pwn.college labs

package main

import (
	"os"
)

func main() {
	file, err := os.Create("valid.cimg")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Write the magic header
	_, err = file.Write([]byte("<MAG"))
	if err != nil {
		panic(err)
	}

	// Write 100 null bytes, or anything really
	padding := make([]byte, 100)
	_, err = file.Write(padding)
	if err != nil {
		panic(err)
	}
}
