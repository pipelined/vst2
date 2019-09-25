package vst2

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
		scanPaths = append(paths, envVstPath)
	}
}

// Open loads the plugin entry point into memory. It's DLL in windows.
func Open(path string) (entryPoint, error) {
	//Load plugin by path
	dll, err := syscall.LoadDLL(l.Path)
	if err != nil {
		return entryPoint{}, fmt.Errorf("failed to load VST from '%s': %v\n", l.Path, err)
	}
	l.library = unsafe.Pointer(dll)
	l.Name = strings.TrimSuffix(filepath.Base(dll.Name), filepath.Ext(dll.Name))

	//Get pointer to plugin's Main function
	entryPoint, err := syscall.GetProcAddress(dll.Handle, main)
	if err != nil {
		l.Close()
		return entryPoint{}, fmt.Errorf("failed to get entry point for plugin'%s': %v\n", l.Path, err)
	}
	return entryPoint{
		main: effectMain(entryPoint),
		handle: handle{
			dll: dll,
		},
	}, nil
}

// Close cleans up plugin refs.
func (h handle) close() error {
	if err := h.dll.Release(); err != nil {
		return fmt.Errorf("failed to release VST handle: %w", err)
	}
	h = nil
	return nil
}
