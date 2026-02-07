//go:build !plugin
// +build !plugin

package vst2

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
)

const (
	// FileExtension of VST2 files in Windows.
	FileExtension = ".dll"
)

// ScanPaths of Vst2 files
var scanPaths = []string{
	"C:\\Program Files (x86)\\Steinberg\\VSTPlugins",
	"C:\\Program Files\\Steinberg\\VSTPlugins ",
}

func init() {
	envVstPath := os.Getenv("VST_PATH")
	if len(envVstPath) > 0 {
		scanPaths = append(scanPaths, envVstPath)
	}
}

// Open loads the plugin entry point into memory. It's DLL in windows.
func Open(path string) (*VST, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for '%s': %w", path, err)
	}
	// Load plugin by path
	dll, err := syscall.LoadDLL(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load VST from '%s': %w", path, err)
	}

	// Get pointer to plugin's Main function
	m, err := syscall.GetProcAddress(dll.Handle, main)
	if err == nil {
		return &VST{
			main:   pluginMain(unsafe.Pointer(m)),
			handle: uintptr(dll.Handle),
			Name:   filepath.Base(path[:len(path)-len(filepath.Ext(path))]),
		}, nil
	}

	err = fmt.Errorf("failed to get entry point for plugin '%s': %w", path, err)
	if errRelease := dll.Release(); errRelease != nil {
		return nil, fmt.Errorf("failed to release DLL '%s': %w after: %v", path, errRelease, err)
	}
	return nil, err
}

// Close frees plugin handle.
func (m *VST) Close() error {
	if err := syscall.FreeLibrary(syscall.Handle(m.handle)); err != nil {
		return fmt.Errorf("failed to release VST handle: %w", err)
	}
	return nil
}
