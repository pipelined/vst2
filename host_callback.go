//go:build !plugin
// +build !plugin

package vst2

/*
#include <stdlib.h>
#include "include/host/host.c"
*/
import "C"

import (
	"fmt"
	"sync"
	"unsafe"

	"pipelined.dev/signal"
)

// global state for callbacks.
var callbacks = struct {
	sync.RWMutex
	mapping map[unsafe.Pointer]HostCallbackFunc
}{
	mapping: map[unsafe.Pointer]HostCallbackFunc{},
}

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
	Plugin struct {
		p *C.CPlugin
	}

	// pluginMain is a reference to VST main function.
	// wrapper on C entry point.
	pluginMain C.EntryPoint
)

type (
	// HostCallbackFunc used as callback function called by plugin. Use
	// closure wrapping technique to add more types to callback.
	HostCallbackFunc func(op HostOpcode, index int32, value int64, ptr unsafe.Pointer, opt float32) int64
)

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
		case HostGetBufferSize:
			if h.GetBufferSize != nil {
				return int64(h.GetBufferSize())
			}
		case HostGetTime:
			if h.GetTimeInfo != nil {
				return int64(uintptr(unsafe.Pointer(h.GetTimeInfo(TimeInfoFlag(value)))))
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
	callbacks.Lock()
	callbacks.mapping[unsafe.Pointer(p)] = c
	callbacks.Unlock()

	return &Plugin{p: p}
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
		in.cArray(),
		out.cArray(),
		C.int32_t(in.Frames),
	)
}

// ProcessFloat audio with VST plugin.
func (p *Plugin) ProcessFloat(in, out FloatBuffer) {
	C.processFloatHostBridge(
		(*C.CPlugin)(p.p),
		in.cArray(),
		out.cArray(),
		C.int32_t(in.Frames),
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

// Start executes the PlugOpen opcode.
func (p *Plugin) Start() {
	p.Dispatch(plugOpen, 0, 0, nil, 0.0)
}

// Close stops the plugin and cleans up C refs for plugin.
func (p *Plugin) Close() {
	p.Dispatch(plugClose, 0, 0, nil, 0.0)
	callbacks.Lock()
	delete(callbacks.mapping, unsafe.Pointer(p))
	callbacks.Unlock()
}

// Resume the plugin processing. It must be called before processing is
// done.
func (p *Plugin) Resume() {
	p.Dispatch(plugStateChanged, 0, 1, nil, 0)
}

// Suspend the plugin processing. It must be called after processing is
// done and no new signal is expected at this moment.
func (p *Plugin) Suspend() {
	p.Dispatch(plugStateChanged, 0, 0, nil, 0)
}

// SetBufferSize sets a buffer size per channel.
func (p *Plugin) SetBufferSize(bufferSize int) {
	p.Dispatch(plugSetBufferSize, 0, int64(bufferSize), nil, 0)
}

// SetSampleRate sets a sample rate for plugin.
func (p *Plugin) SetSampleRate(sampleRate signal.Frequency) {
	p.Dispatch(plugSetSampleRate, 0, 0, nil, float32(sampleRate))
}

// SetSpeakerArrangement creates and passes SpeakerArrangement structures to plugin
func (p *Plugin) SetSpeakerArrangement(in, out *SpeakerArrangement) {
	p.Dispatch(plugSetSpeakerArrangement, 0, int64(uintptr(unsafe.Pointer(in))), unsafe.Pointer(out), 0)
}

// ParamName returns the parameter label: "Release", "Gain", etc.
func (p *Plugin) ParamName(index int) string {
	var s ascii8
	p.Dispatch(plugGetParamName, int32(index), 0, unsafe.Pointer(&s), 0)
	return s.String()
}

// ParamValueName returns the parameter value label: "0.5", "HALL", etc.
func (p *Plugin) ParamValueName(index int) string {
	var s ascii8
	p.Dispatch(plugGetParamDisplay, int32(index), 0, unsafe.Pointer(&s), 0)
	return s.String()
}

// ParamUnitName returns the parameter unit label: "db", "ms", etc.
func (p *Plugin) ParamUnitName(index int) string {
	var s ascii8
	p.Dispatch(plugGetParamLabel, int32(index), 0, unsafe.Pointer(&s), 0)
	return s.String()
}

// CurrentProgramName returns current program name.
func (p *Plugin) CurrentProgramName() string {
	var s ascii24
	p.Dispatch(plugGetProgramName, 0, 0, unsafe.Pointer(&s), 0)
	return s.String()
}

// ProgramName returns program name for provided program index.
func (p *Plugin) ProgramName(index int) string {
	var s ascii24
	p.Dispatch(plugGetProgramNameIndexed, int32(index), 0, unsafe.Pointer(&s), 0)
	return s.String()
}

// SetCurrentProgramName sets new name to the current program. It will use
// up to 24 ASCII characters. Non-ASCII characters are ignored.
func (p *Plugin) SetCurrentProgramName(s string) {
	var ps ascii24
	copyASCII(ps[:], s)
	p.Dispatch(plugSetProgramName, 0, 0, unsafe.Pointer(&ps), 0)
}

// Program returns current program number.
func (p *Plugin) Program() int {
	return int(p.Dispatch(plugGetProgram, 0, 0, nil, 0))
}

// SetProgram changes current program index.
func (p *Plugin) SetProgram(index int) {
	p.Dispatch(plugSetProgram, 0, int64(index), nil, 0)
}

// ParamProperties returns parameter properties for provided parameter
// index. If opcode is not supported, boolean result is false.
func (p *Plugin) ParamProperties(index int) (*ParameterProperties, bool) {
	var props ParameterProperties
	r := p.Dispatch(plugGetParameterProperties, int32(index), 0, unsafe.Pointer(&props), 0)
	if r > 0 {
		return &props, true
	}
	return nil, false
}

// GetProgramData returns current preset data. Plugin allocates required
// memory, then this function allocates new byte slice of required length
// where data is copied.
func (p *Plugin) GetProgramData() []byte {
	var ptr unsafe.Pointer
	length := C.int(p.Dispatch(plugGetChunk, 1, 0, unsafe.Pointer(&ptr), 0))
	return C.GoBytes(ptr, length)
}

// SetProgramData sets preset data to the plugin. Data is the full preset
// including chunk header.
func (p *Plugin) SetProgramData(data []byte) {
	p.Dispatch(plugSetChunk, 1, int64(len(data)), unsafe.Pointer(&data[0]), 0)
}

// GetBankData returns current bank data. Plugin allocates required
// memory, then this function allocates new byte slice of required length
// where data is copied.
func (p *Plugin) GetBankData() []byte {
	var ptr unsafe.Pointer
	length := C.int(p.Dispatch(plugGetChunk, 0, 0, unsafe.Pointer(&ptr), 0))
	return C.GoBytes(ptr, length)
}

// SetBankData sets preset data to the plugin. Data is the full preset
// including chunk header.
func (p *Plugin) SetBankData(data []byte) {
	ptr := C.CBytes(data)
	p.Dispatch(plugSetChunk, 0, int64(len(data)), unsafe.Pointer(ptr), 0)
	C.free(ptr)
}

// Convert golang string to C string without allocations. Result string is
// not null-terminated.
func stringToCString(s string) *C.char {
	return (*C.char)(unsafe.Pointer(&[]byte(s)[0]))
}
