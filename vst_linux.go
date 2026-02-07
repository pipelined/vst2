//go:build !plugin
// +build !plugin

package vst2

// #cgo LDFLAGS: -ldl
// #include <dlfcn.h>
// #include <stdlib.h>
import "C"

import (
	"fmt"
	"os"
	"path/filepath"
	"unsafe"
)

const (
	// FileExtension of VST2 files in Windows.
	FileExtension = ".so"
)

// ScanPaths of Vst2 files
var scanPaths = []string{
	"/usr/local/lib/vst",
	"/usr/lib/vst",
}

func init() {
	home := os.Getenv("HOME")
	scanPaths = append(scanPaths,
		filepath.Join(home, ".vst"),
		filepath.Join(home, "vst"),
	)
}

// Open loads the plugin entry point into memory. It's SO in linux.
func Open(path string) (*VST, error) {
	handle := C.dlopen(stringToCString(path), C.RTLD_LAZY)
	if handle == nil {
		return nil, fmt.Errorf("failed loading vst: %v", dlerror())
	}

	// clear previous errors as stated in the man.
	C.dlerror()

	m := C.dlsym(handle, stringToCString(main))
	if m == nil {
		return nil, fmt.Errorf("failed finding vst main: %v", dlerror())
	}

	return &VST{
		main:   pluginMain(m),
		handle: uintptr(handle),
		Name:   filepath.Base(path[:len(path)-len(filepath.Ext(path))]),
	}, nil
}

func dlerror() string {
	CError := C.dlerror()
	defer C.free(unsafe.Pointer(CError))
	return C.GoString(CError)
}

// Close frees plugin handle.
func (m *VST) Close() error {
	if C.dlclose(unsafe.Pointer(m.handle)) != 0 {
		return fmt.Errorf("error unloading vst: %v", dlerror())
	}
	return nil
}
