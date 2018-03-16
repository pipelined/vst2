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

// Library used to instantiate new instances of plugin
type Library struct {
	entryPoint unsafe.Pointer
	library    unsafe.Pointer
	Name       string
	Path       string
}

func (l *Library) load() error {
	//Load plugin by path
	vstDLL, err := syscall.LoadDLL(l.Path)
	if err != nil {
		return fmt.Errorf("Failed to load VST from '%s': %v\n", l.Path, err)
	}
	l.library = unsafe.Pointer(vstDLL)
	l.Name = strings.TrimSuffix(filepath.Base(vstDLL.Name), filepath.Ext(vstDLL.Name))

	//Get pointer to plugin's Main function
	entryPoint, err := syscall.GetProcAddress(vstDLL.Handle, vstMain)
	if err != nil {
		l.Close()
		return fmt.Errorf("Failed to get entry point for plugin'%s': %v\n", l.Path, err)
	}
	l.entryPoint = unsafe.Pointer(entryPoint)
	return nil
}

//Close cleans up plugin refs
func (l *Library) Close() {
	vstDLL := (*syscall.DLL)(l.library)
	vstDLL.Release()
	l.library = nil
}
