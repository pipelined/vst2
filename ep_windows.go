package vst2

import (
	"log"
	"syscall"
	"unsafe"
)

//TODO: refactor to plugin.load method
func getEntryPoint(path string) (unsafe.Pointer, error) {
	//Load plugin by path
	moduleHandle, err := syscall.LoadLibrary(path)
	if err != nil {
		log.Printf("Failed to load VST from '%s': %v\n", path, err)
		return nil, err
	}

	//Get pointer to plugin's Main function
	result, err := syscall.GetProcAddress(moduleHandle, vstMain)
	if err != nil {
		log.Printf("Failed to get entry point for plugin'%s': %v\n", path, err)
		return nil, err
	}
	return unsafe.Pointer(result), nil
}
