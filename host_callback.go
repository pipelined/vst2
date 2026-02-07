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
	// VST main function name.
	main = "VSTPluginMain"
	// VST API version.
	version = 2400
)

type (
	// VST is a reference to VST main function.
	// It also keeps reference to VST handle to clean up on Close.
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
		case HostAutomate:
			if h.Automate != nil {
				h.Automate(index, opt)
			}
		case HostIdle:
			if h.Idle != nil {
				h.Idle()
			}
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
		case HostProcessEvents:
			if h.ProcessEvents != nil {
				h.ProcessEvents((*EventsPtr)(ptr))
				return 1
			}
		case HostIOChanged:
			if h.IOChanged != nil {
				if h.IOChanged() {
					return 1
				}
			}
		case HostSizeWindow:
			if h.SizeWindow != nil {
				if h.SizeWindow(index, int32(value)) {
					return 1
				}
			}
		case HostGetInputLatency:
			if h.GetInputLatency != nil {
				return int64(h.GetInputLatency())
			}
		case HostGetOutputLatency:
			if h.GetOutputLatency != nil {
				return int64(h.GetOutputLatency())
			}
		case HostGetAutomationState:
			if h.GetAutomationState != nil {
				return int64(h.GetAutomationState())
			}
		case HostGetVendorString:
			if h.GetVendorString != nil {
				s := (*ascii64)(ptr)
				copyASCII(s[:], h.GetVendorString())
				return 1
			}
		case HostGetProductString:
			if h.GetProductString != nil {
				s := (*ascii64)(ptr)
				copyASCII(s[:], h.GetProductString())
				return 1
			}
		case HostGetVendorVersion:
			if h.GetVendorVersion != nil {
				return int64(h.GetVendorVersion())
			}
		case HostCanDo:
			if h.CanDo != nil {
				s := HostCanDoString(C.GoString((*C.char)(ptr)))
				return int64(h.CanDo(s))
			}
		case HostGetLanguage:
			if h.GetLanguage != nil {
				return int64(h.GetLanguage())
			}
		case HostGetDirectory:
			if h.GetDirectory != nil {
				dir := h.GetDirectory()
				return int64(uintptr(unsafe.Pointer(C.CString(dir))))
			}
		case HostUpdateDisplay:
			if h.UpdateDisplay != nil {
				if h.UpdateDisplay() {
					return 1
				}
			}
		case HostBeginEdit:
			if h.BeginEdit != nil {
				if h.BeginEdit(index) {
					return 1
				}
			}
		case HostEndEdit:
			if h.EndEdit != nil {
				if h.EndEdit(index) {
					return 1
				}
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
// This function also calls dispatch with EffOpen opcode.
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

// SendEvents sends MIDI events to the hosted plugin. The caller creates
// events via the Events() constructor and is responsible for calling
// Free() afterward.
func (p *Plugin) SendEvents(events *EventsPtr) {
	p.Dispatch(PlugProcessEvents, 0, 0, unsafe.Pointer(events), 0)
}

// EditorGetRect returns the editor rectangle. Returns nil if plugin has
// no editor.
func (p *Plugin) EditorGetRect() *EditorRectangle {
	var rect *EditorRectangle
	r := p.Dispatch(PlugEditGetRect, 0, 0, unsafe.Pointer(&rect), 0)
	if r == 0 && rect == nil {
		return nil
	}
	return rect
}

// EditorOpen opens the plugin editor in the provided native window handle
// (HWND on Windows, NSView* on macOS, X11 Window on Linux).
func (p *Plugin) EditorOpen(window unsafe.Pointer) {
	p.Dispatch(PlugEditOpen, 0, 0, window, 0)
}

// EditorClose closes the plugin editor window.
func (p *Plugin) EditorClose() {
	p.Dispatch(PlugEditClose, 0, 0, nil, 0)
}

// EditorIdle notifies the plugin editor to do idle processing.
func (p *Plugin) EditorIdle() {
	p.Dispatch(PlugEditIdle, 0, 0, nil, 0)
}

// CanDo queries the plugin about its capabilities. Returns YesCanDo,
// NoCanDo, or MaybeCanDo.
func (p *Plugin) CanDo(s PluginCanDoString) CanDoResponse {
	cs := C.CString(string(s))
	defer C.free(unsafe.Pointer(cs))
	return CanDoResponse(p.Dispatch(PlugCanDo, 0, 0, unsafe.Pointer(cs), 0))
}

// GetTailSize returns the tail size in samples (e.g. reverb time).
// Returns 0 as default, 1 means no tail.
func (p *Plugin) GetTailSize() int {
	return int(p.Dispatch(PlugGetTailSize, 0, 0, nil, 0))
}

// GetInputProperties returns pin properties for the given input index.
// If the opcode is not supported, the boolean result is false.
func (p *Plugin) GetInputProperties(index int) (*PinProperties, bool) {
	var props PinProperties
	r := p.Dispatch(PlugGetInputProperties, int32(index), 0, unsafe.Pointer(&props), 0)
	if r > 0 {
		return &props, true
	}
	return nil, false
}

// GetOutputProperties returns pin properties for the given output index.
// If the opcode is not supported, the boolean result is false.
func (p *Plugin) GetOutputProperties(index int) (*PinProperties, bool) {
	var props PinProperties
	r := p.Dispatch(PlugGetOutputProperties, int32(index), 0, unsafe.Pointer(&props), 0)
	if r > 0 {
		return &props, true
	}
	return nil, false
}

// GetVSTVersion returns the VST version supported by the plugin (e.g. 2400).
func (p *Plugin) GetVSTVersion() int {
	return int(p.Dispatch(PlugGetVstVersion, 0, 0, nil, 0))
}

// GetPluginName returns the name of the plugin (up to 32 chars).
func (p *Plugin) GetPluginName() string {
	var s ascii32
	p.Dispatch(PlugGetPluginName, 0, 0, unsafe.Pointer(&s), 0)
	return s.String()
}

// GetVendorString returns the plugin vendor string (up to 64 chars).
func (p *Plugin) GetVendorString() string {
	var s ascii64
	p.Dispatch(PlugGetVendorString, 0, 0, unsafe.Pointer(&s), 0)
	return s.String()
}

// GetProductString returns the plugin product string (up to 64 chars).
func (p *Plugin) GetProductString() string {
	var s ascii64
	p.Dispatch(PlugGetProductString, 0, 0, unsafe.Pointer(&s), 0)
	return s.String()
}

// GetVendorVersion returns the plugin vendor-specific version.
func (p *Plugin) GetVendorVersion() int {
	return int(p.Dispatch(PlugGetVendorVersion, 0, 0, nil, 0))
}

// GetCategory returns the plugin category.
func (p *Plugin) GetCategory() PluginCategory {
	return PluginCategory(p.Dispatch(PlugGetPlugCategory, 0, 0, nil, 0))
}

// NumInputs returns the number of audio input channels.
func (p *Plugin) NumInputs() int {
	return int(p.p.numInputs)
}

// NumOutputs returns the number of audio output channels.
func (p *Plugin) NumOutputs() int {
	return int(p.p.numOutputs)
}

// UniqueID returns the plugin's unique identifier.
func (p *Plugin) UniqueID() int32 {
	return int32(p.p.uniqueID)
}

// InitialDelay returns the plugin's latency in samples.
func (p *Plugin) InitialDelay() int {
	return int(p.p.initialDelay)
}
