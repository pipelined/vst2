package vst2

//#include "vst2.h"
//#include <stdlib.h>
import "C"
import (
	"fmt"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

func (library *Library) load() error {
	//Load plugin by path
	vstDLL, err := syscall.LoadDLL(library.Path)
	if err != nil {
		return fmt.Errorf("Failed to load VST from '%s': %v\n", library.Path, err)
	}
	library.library = uintptr(unsafe.Pointer(vstDLL))
	library.Name = strings.TrimSuffix(filepath.Base(vstDLL.Name), filepath.Ext(vstDLL.Name))

	//Get pointer to plugin's Main function
	entryPoint, err := syscall.GetProcAddress(vstDLL.Handle, vstMain)
	if err != nil {
		library.Close()
		return fmt.Errorf("Failed to get entry point for plugin'%s': %v\n", library.Path, err)
	}
	library.entryPoint = unsafe.Pointer(entryPoint)
	return nil
}

//Close cleans up plugin refs
func (library *Library) Close() {
	vstDLL := (*syscall.DLL)(unsafe.Pointer((library.library)))
	vstDLL.Release()
}
