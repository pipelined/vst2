package vst2

import (
	"log"
	"syscall"
)

func getEntryPoint(path string) (uintptr, error) {
	//Load plugin by path
	moduleHandle, err := syscall.LoadLibrary(path)
	if err != nil {
		log.Printf("Failed to load VST from '%s': %v\n", path, err)
		return 0, err
	}

	//Get pointer to plugin's Main function
	return syscall.GetProcAddress(moduleHandle, vstMain)
}
