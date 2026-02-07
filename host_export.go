//go:build !plugin
// +build !plugin

package vst2

//#include "include/vst.h"
import "C"
import "unsafe"

// global hostCallbackBridge, calls real callback.
//
//export hostCallbackBridge
func hostCallbackBridge(p *C.CPlugin, opcode int32, index int32, value int64, ptr unsafe.Pointer, opt float32) int64 {
	// HostVersion is requested when plugin is created
	// It's never in map
	if HostOpcode(opcode) == HostVersion {
		return version
	}
	callbacks.RLock()
	c, ok := callbacks.mapping[unsafe.Pointer(p)]
	callbacks.RUnlock()
	if !ok {
		panic("plugin was closed")
	}

	if c == nil {
		panic("host callback is undefined")
	}
	return c(HostOpcode(opcode), index, value, ptr, opt)
}
