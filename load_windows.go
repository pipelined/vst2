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

//TODO: refactor to plugin.load method
func (plugin *Plugin) load() error {
	//Load plugin by path
	vstDLL, err := syscall.LoadDLL(plugin.Path)
	if err != nil {
		return fmt.Errorf("Failed to load VST from '%s': %v\n", plugin.Path, err)
	}
	plugin.library = unsafe.Pointer(vstDLL)
	plugin.Name = strings.TrimSuffix(filepath.Base(vstDLL.Name), filepath.Ext(vstDLL.Name))

	//Get pointer to plugin's Main function
	mainEntryPoint, err := syscall.GetProcAddress(vstDLL.Handle, vstMain)
	if err != nil {
		plugin.Unload()
		return fmt.Errorf("Failed to get entry point for plugin'%s': %v\n", plugin.Path, err)
	}

	plugin.effect = C.loadEffect(C.vstPluginFuncPtr(unsafe.Pointer(mainEntryPoint)))
	return nil
}

//Unload cleans up plugin refs
func (plugin *Plugin) Unload() {
	vstDLL := (*syscall.DLL)(plugin.library)
	vstDLL.Release()
	C.free(unsafe.Pointer(plugin.effect))
}
