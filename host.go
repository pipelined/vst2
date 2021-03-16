package vst2

/*
#cgo CFLAGS: -I${SRCDIR}
#include "vst.h"
*/
import "C"
import (
	"fmt"
	"sync"
	"unsafe"

	"pipelined.dev/signal"
)

// global state for callbacks.
var (
	mutex     sync.RWMutex
	callbacks = make(map[unsafe.Pointer]HostCallbackFunc)
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

	// pluginMain is a reference to VST main function.
	// wrapper on C entry point.
	pluginMain C.EntryPoint
)

type (
	// HostCallbackFunc used as callback function called by plugin. Use closure
	// wrapping technique to add more types to callback.
	HostCallbackFunc func(op HostOpcode, index int32, value int64, ptr unsafe.Pointer, opt float32) int64

	// Host handles all callbacks from plugin.
	Host struct {
		ProgressProcessed HostProgressProcessed
		GetSampleRate     HostGetSampleRateFunc
		GetBufferSize     HostGetBufferSizeFunc
		GetProcessLevel   HostGetProcessLevel
		GetTimeInfo       HostGetTimeInfo
	}

	// HostProgressProcessed is executed after every process call.
	HostProgressProcessed func(int)
	// HostGetSampleRateFunc returns host sample rate.
	HostGetSampleRateFunc func() signal.Frequency
	// HostGetBufferSizeFunc returns host buffer size.
	HostGetBufferSizeFunc func() int
	// HostGetProcessLevel returns the context of execution.
	HostGetProcessLevel func() ProcessLevel
	// HostGetTimeInfo returns current time info.
	HostGetTimeInfo func() *TimeInfo
)

// DispatchFunc called by host.
type DispatchFunc func(op PluginOpcode, index int32, value int64, ptr unsafe.Pointer, opt float32) int64

// Callback returns HostCallbackFunc that handles all vst types casts
// and allows to write handlers without usage of unsafe package.
func (h Host) Callback() HostCallbackFunc {
	return func(op HostOpcode, index int32, value int64, ptr unsafe.Pointer, opt float32) int64 {
		switch op {
		case HostGetCurrentProcessLevel:
			if h.GetProcessLevel != nil {
				return int64(h.GetProcessLevel())
			}
		case HostGetSampleRate:
			if h.GetSampleRate != nil {
				return int64(h.GetSampleRate())
			}
		case HostGetBlockSize:
			if h.GetBufferSize != nil {
				return int64(h.GetBufferSize())
			}
		case HostGetTime:
			if h.GetTimeInfo != nil {
				return int64(uintptr(unsafe.Pointer(h.GetTimeInfo())))
			}
		}
		return 0
	}
}

// NoopHostCallback returns dummy host callback that just prints received
// opcodes.
func NoopHostCallback() HostCallbackFunc {
	return func(op HostOpcode, index int32, value int64, ptr unsafe.Pointer, opt float32) int64 {
		fmt.Printf("host received opcode: %v\n", op)
		return 0
	}
}

// Plugin new instance of VST plugin with provided callback.
// This function also calls dispatch with EffOpen opcodp.
func (v *VST) Plugin(c HostCallbackFunc) *Plugin {
	if v.main == nil || c == nil {
		return nil
	}
	p := (*C.CPlugin)(C.loadPluginHostBridge(v.main))
	mutex.Lock()
	callbacks[unsafe.Pointer(p)] = c
	mutex.Unlock()

	return &Plugin{p: p}
}

//export hostCallback
// global hostCallback, calls real callback.
func hostCallback(p *C.CPlugin, opcode int32, index int32, value int64, ptr unsafe.Pointer, opt float32) int64 {
	// HostVersion is requested when plugin is created
	// It's never in map
	if HostOpcode(opcode) == HostVersion {
		return version
	}
	mutex.RLock()
	c, ok := callbacks[unsafe.Pointer(p)]
	mutex.RUnlock()
	if !ok {
		panic("plugin was closed")
	}

	if c == nil {
		panic("host callback is undefined")
	}
	return c(HostOpcode(opcode), index, value, ptr, opt)
}

// Dispatch wraps-up C method to dispatch calls to plugin
func (p *Plugin) Dispatch(opcode PluginOpcode, index int32, value int64, ptr unsafe.Pointer, opt float32) uintptr {
	return uintptr(C.dispatchHostBridge((*C.CPlugin)(p.p), C.int32_t(opcode), C.int32_t(index), C.int64_t(value), ptr, C.float(opt)))
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
	return int(p.p.numParams)
}

// NumPrograms returns the number of programs.
func (p *Plugin) NumPrograms() int {
	return int(p.p.numPrograms)
}

// Flags returns the plugin flags.
func (p *Plugin) Flags() PluginFlag {
	return PluginFlag(p.p.flags)
}

// ProcessDouble audio with VST plugin.
func (p *Plugin) ProcessDouble(in, out DoubleBuffer) {
	C.processDoubleHostBridge(
		(*C.CPlugin)(p.p),
		&in.data[0],
		&out.data[0],
		C.int32_t(in.size),
	)
}

// ProcessFloat audio with VST plugin.
func (p *Plugin) ProcessFloat(in, out FloatBuffer) {
	C.processFloatHostBridge(
		(*C.CPlugin)(p.p),
		&in.data[0],
		&out.data[0],
		C.int32_t(in.size),
	)
}

// ParamValue returns the value of parameter.
func (p *Plugin) ParamValue(index int) float32 {
	return float32(C.getParameterHostBridge((*C.CPlugin)(p.p), C.int32_t(index)))
}

// SetParamValue sets new value for parameter.
func (p *Plugin) SetParamValue(index int, value float32) {
	C.setParameterHostBridge((*C.CPlugin)(p.p), C.int32_t(index), C.float(value))
}

// CanProcessFloat32 checks if plugin can process float32.
func (p *Plugin) CanProcessFloat32() bool {
	return PluginFlag(p.p.flags)&PluginFloatProcessing == PluginFloatProcessing
}

// CanProcessFloat64 checks if plugin can process float64.
func (p *Plugin) CanProcessFloat64() bool {
	return PluginFlag(p.p.flags)&PluginDoubleProcessing == PluginDoubleProcessing
}

// Convert golang string to C string without allocations. Result string is
// not null-terminated.
func stringToCString(s string) *C.char {
	return (*C.char)(unsafe.Pointer(&[]byte(s)[0]))
}
