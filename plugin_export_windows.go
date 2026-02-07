//go:build plugin
// +build plugin

package vst2

import (
	"reflect"
	"syscall"
	"unsafe"
)

// fix vst host crash in windows, when plugin gets unloaded by pinning dll
// see https://docs.microsoft.com/en-gb/windows/win32/api/libloaderapi/nf-libloaderapi-getmodulehandleexa?redirectedfrom=MSDN
// for information about GET_MODULE_HANDLE_EX_FLAG_PIN
// workaround for issue: https://github.com/golang/go/issues/11100
func loadHook() {
	const (
		GET_MODULE_HANDLE_EX_FLAG_FROM_ADDRESS = 4
		GET_MODULE_HANDLE_EX_FLAG_PIN          = 1
	)
	var (
		kernel32, _          = syscall.LoadLibrary("kernel32.dll")
		getModuleHandleEx, _ = syscall.GetProcAddress(kernel32, "GetModuleHandleExW")
		handle               uintptr
	)
	defer func(handle syscall.Handle) {
		err := syscall.FreeLibrary(handle)
		if err != nil {
			panic("cant unload kernel32 lib")
		}
	}(kernel32)
	if _, _, callErr := syscall.Syscall(uintptr(getModuleHandleEx), 3, GET_MODULE_HANDLE_EX_FLAG_FROM_ADDRESS|GET_MODULE_HANDLE_EX_FLAG_PIN, reflect.ValueOf(loadHook).Pointer(), uintptr(unsafe.Pointer(&handle))); callErr != 0 {
		panic("cant pin dll")
	}
}
