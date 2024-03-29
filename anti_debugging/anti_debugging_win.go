package main

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"
)

const (
	GENERIC_READ      = 0x80000000
	OPEN_EXISTING     = 3
	MAX_PATH          = 260
	INVALID_HANDLE_VALUE = ^uintptr(0)
)

var (
	modkernel32 = syscall.NewLazyDLL("kernel32.dll")
	procGetModuleFileNameW = modkernel32.NewProc("GetModuleFileNameW")
	procCreateFileW = modkernel32.NewProc("CreateFileW")
	procCloseHandle = modkernel32.NewProc("CloseHandle")
	procIsDebuggerPresent = modkernel32.NewProc("IsDebuggerPresent")
	procOutputDebugStringA = modkernel32.NewProc("OutputDebugStringA")
)

func GetModuleFileNameW(hModule uintptr) (string, error) {
	var fileName [MAX_PATH]uint16
	_, _, err := procGetModuleFileNameW.Call(
		hModule,
		uintptr(unsafe.Pointer(&fileName[0])),
		MAX_PATH,
	)
	if err != nil {
		return "", err
	}
	return syscall.UTF16ToString(fileName[:]), nil
}

func CreateFileW(lpFileName string, dwDesiredAccess uint32, dwShareMode uint32, lpSecurityAttributes uintptr, dwCreationDisposition uint32, dwFlagsAndAttributes uint32, hTemplateFile uintptr) (uintptr, error) {
	fileNamePtr, err := syscall.UTF16PtrFromString(lpFileName)
	if err != nil {
		return 0, err
	}
	res, _, err := procCreateFileW.Call(
		uintptr(unsafe.Pointer(fileNamePtr)),
		uintptr(dwDesiredAccess),
		uintptr(dwShareMode),
		lpSecurityAttributes,
		uintptr(dwCreationDisposition),
		uintptr(dwFlagsAndAttributes),
		hTemplateFile,
	)
	if res == INVALID_HANDLE_VALUE {
		return 0, err
	}
	return res, nil
}

func CloseHandle(hObject uintptr) (bool, error) {
	ret, _, err := procCloseHandle.Call(hObject)
	if ret == 0 {
		return false, err
	}
	return true, nil
}

func IsDebuggerPresent() (bool, error) {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	procIsDebuggerPresent := kernel32.NewProc("IsDebuggerPresent")

	isDebuggerPresent, _, err := procIsDebuggerPresent.Call()
	if isDebuggerPresent != 0 {
		return true, nil
	}
	if err.(syscall.Errno) != 0 {
		return false, err
	}
	return false, nil
}

func getProcessFileHandle() (bool, error) {
	// Attempting to open its own executable
	fileName, err := GetModuleFileNameW(0)
	if err != nil && !strings.Contains(err.Error(), "successfully") {
		return false, err
	}
	res, _ := CreateFileW(fileName, GENERIC_READ, 0, 0, OPEN_EXISTING, 0, 0)
	if res == 0 {
		return false, nil
	}
	defer CloseHandle(res)
	return res == INVALID_HANDLE_VALUE, nil
}

func OutputDebugStringA(lpOutputString *byte) (bool, error) {
	// Outputs string to debugger, then tries to access an invalid memory address using inline assembly
	ret, _, err := procOutputDebugStringA.Call(uintptr(unsafe.Pointer(lpOutputString)))
	if ret != 0 {
		return true, nil
	}
	if err.(syscall.Errno) != 0 {
		return false, err
	}
	return false, nil
}

func isDebugged() (bool, error) {
	peb, err := IsDebuggerPresent()
	if err != nil {
		return false, err
	}
	
	fileHandle, err := getProcessFileHandle()
	if !fileHandle {
		return false, err
	}

	lpOutputString := []byte("aaaaa\x00")
	debuggerHandledException, err := OutputDebugStringA(&lpOutputString[0])
	if err != nil {
		return false, err
	}

	return  peb || fileHandle || debuggerHandledException, nil
}

func main() {
	debugged, err := isDebugged()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Is debugged:", debugged)
}
