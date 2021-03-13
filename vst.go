package vst2

/*
#cgo CFLAGS: -std=gnu99 -I${SRCDIR}
#include "sdk.h"
*/
import "C"
import (
	"unsafe"
)

const (
	// VST main function namp.
	main = "VSTPluginMain"
	// VST API version.
	version = 2400
)

type (
	// VST is a reference to VST main function.
	// It also keeps reference to VST handle to clean up on Closp.
	VST struct {
		main   pluginMain
		handle uintptr
		Name   string
	}

	// Plugin is an instance of loaded VST plugin.
	Plugin C.Plugin

	// pluginMain is a reference to VST main function.
	// wrapper on C entry point.
	pluginMain C.EntryPoint
)

// Plugin new instance of VST plugin with provided callback.
// This function also calls dispatch with EffOpen opcodp.
func (v *VST) Plugin(c HostCallbackFunc) *Plugin {
	if v.main == nil || c == nil {
		return nil
	}
	p := (*Plugin)(C.loadPluginBridge(v.main))
	mutex.Lock()
	callbacks[p] = c
	mutex.Unlock()

	return p
}

// Dispatch wraps-up C method to dispatch calls to plugin
func (p *Plugin) Dispatch(opcode PluginOpcode, index int32, value int64, ptr unsafe.Pointer, opt float32) uintptr {
	return uintptr(C.dispatchHostBridge((*C.Plugin)(p), C.int32_t(opcode), C.int32_t(index), C.int64_t(value), unsafe.Pointer(ptr), C.float(opt)))
}

// ScanPaths returns a slice of default vst2 locations.
// Locations are OS-specific.
func ScanPaths() []string {
	paths := make([]string, 0, len(scanPaths))
	paths = append(paths, scanPaths...)
	return paths
}

// NumParams returns the number of parameters.
func (p *Plugin) NumParams() int {
	return int(p.numParams)
}

// NumPrograms returns the number of programs.
func (p *Plugin) NumPrograms() int {
	return int(p.numPrograms)
}

// Flags returns the plugin flags.
func (p *Plugin) Flags() PluginFlag {
	return PluginFlag(p.flags)
}

// ProcessDouble audio with VST plugin.
func (p *Plugin) ProcessDouble(in, out DoubleBuffer) {
	C.processDoubleHostBridge(
		(*C.Plugin)(p),
		&in.data[0],
		&out.data[0],
		C.int32_t(in.size),
	)
}

// ProcessFloat audio with VST plugin.
func (p *Plugin) ProcessFloat(in, out FloatBuffer) {
	C.processFloatHostBridge(
		(*C.Plugin)(p),
		&in.data[0],
		&out.data[0],
		C.int32_t(in.size),
	)
}

// ParamValue returns the value of parameter.
func (p *Plugin) ParamValue(index int) float32 {
	return float32(C.getParameterHostBridge((*C.Plugin)(p), C.int32_t(index)))
}

// SetParamValue sets new value for parameter.
func (p *Plugin) SetParamValue(index int, value float32) {
	C.setParameterHostBridge((*C.Plugin)(p), C.int32_t(index), C.float(value))
}

// CanProcessFloat32 checks if plugin can process float32.
func (p *Plugin) CanProcessFloat32() bool {
	return PluginFlag(p.flags)&PluginFloatProcessing == PluginFloatProcessing
}

// CanProcessFloat64 checks if plugin can process float64.
func (p *Plugin) CanProcessFloat64() bool {
	return PluginFlag(p.flags)&PluginDoubleProcessing == PluginDoubleProcessing
}

// Convert golang string to C string without allocations. Result string is
// not null-terminated.
func stringToCString(s string) *C.char {
	return (*C.char)(unsafe.Pointer(&[]byte(s)[0]))
}
