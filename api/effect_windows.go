package api

import (
	"fmt"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

type handle struct {
	dll *syscall.DLL
}

func Open(path string) (EntryPoint, error) {
	//Load plugin by path
	dll, err := syscall.LoadDLL(l.Path)
	if err != nil {
		return EntryPoint{}, fmt.Errorf("failed to load VST from '%s': %v\n", l.Path, err)
	}
	l.library = unsafe.Pointer(dll)
	l.Name = strings.TrimSuffix(filepath.Base(dll.Name), filepath.Ext(dll.Name))

	//Get pointer to plugin's Main function
	entryPoint, err := syscall.GetProcAddress(dll.Handle, main)
	if err != nil {
		l.Close()
		return EntryPoint{}, fmt.Errorf("failed to get entry point for plugin'%s': %v\n", l.Path, err)
	}
	return EntryPoint{
		main: effectMain(entryPoint),
		handle: handle{
			dll: dll,
		},
	}, nil
}

//Close cleans up plugin refs
func (h handle) close() error {
	if err := h.dll.Release(); err != nil {
		return fmt.Errorf("failed to release VST handle: %w", err)
	}
	h = nil
	return nil
}
