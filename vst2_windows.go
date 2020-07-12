package vst2

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const (
	// Extension of Vst2 files
	Extension = ".dll"
)

var (
	// ScanPaths of Vst2 files
	scanPaths = []string{
		"C:\\Program Files (x86)\\Steinberg\\VSTPlugins",
		"C:\\Program Files\\Steinberg\\VSTPlugins ",
	}
)

// handle keeps a DLL reference to clean up on close.
type handle struct {
	dll *syscall.DLL
}

func init() {
	envVstPath := os.Getenv("VST_PATH")
	if len(envVstPath) > 0 {
		scanPaths = append(scanPaths, envVstPath)
	}
}

// open loads the plugin entry point into memory. It's DLL in windows.
func open(path string) (effectMain, handle, error) {
	//Load plugin by path
	dll, err := syscall.LoadDLL(path)
	if err != nil {
		return nil, handle{}, fmt.Errorf("failed to load VST from '%s': %v\n", path, err)
	}

	//Get pointer to plugin's Main function
	m, err := syscall.GetProcAddress(dll.Handle, main)
	if err == nil {
		return effectMain(unsafe.Pointer(m)), handle{dll: dll}, nil
	}
	err = fmt.Errorf("failed to get entry point for plugin '%s': %w\n", path, err)
	if errRelease := dll.Release(); errRelease != nil {
		return nil, handle{}, fmt.Errorf("failed to release DLL '%s': %w after: %v", path, errRelease, err)
	}
	return nil, handle{}, err
}

// Close cleans up plugin refs.
func (h handle) close() error {
	if err := h.dll.Release(); err != nil {
		return fmt.Errorf("failed to release VST handle: %w", err)
	}
	return nil
}
