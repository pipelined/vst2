//go:build plugin
// +build plugin

package vst2

/*
#include <stdlib.h>
#include "include/plugin/plugin.c"
*/
import "C"

import (
	"sync"
	"unsafe"

	"pipelined.dev/signal"
)

var (
	PluginAllocator PluginAllocatorFunc

	// global state for callbacks.
	plugins = struct {
		sync.RWMutex
		mapping map[uintptr]*Plugin
	}{
		mapping: map[uintptr]*Plugin{},
	}
)

type (
	// PluginAllocatorFunc allocates new plugin instance and its
	// dispatcher.
	PluginAllocatorFunc func(Host) (Plugin, Dispatcher)

	// Plugin is a VST2 effect that processes float/double signal buffers.
	Plugin struct {
		UniqueID       [4]byte
		Version        int32
		Name           string
		Category       PluginCategory
		Vendor         string
		InputChannels  int
		OutputChannels int
		Flags          PluginFlag
		inputDouble    DoubleBuffer
		outputDouble   DoubleBuffer
		inputFloat     FloatBuffer
		outputFloat    FloatBuffer
		ProcessDoubleFunc
		ProcessFloatFunc
		Parameters []*Parameter
		dispatchFunc
	}

	// Dispatcher handles plugin dispatch calls from the host.
	Dispatcher struct {
		SetBufferSizeFunc func(size int)
		CanDoFunc         func(PluginCanDoString) CanDoResponse // called by host to query the plugin about its capabilities
		CloseFunc         func()                                // called by host right before deleting the plugin, use to free up resources
		ProcessEventsFunc func(*EventsPtr)                      // called by host to pass events (e.g. MIDI events) along with their time stamps (frames) within the next processing block
		GetChunkFunc      func(isPreset bool) []byte            // called by host to get the current state of the plugin. You should define GetChunkFunc & SetChunkFunc in pairs; defining both sets the PluginProgramChunks flag advertizing the capability to the host.
		SetChunkFunc      func(data []byte, isPreset bool)      // called by host to set the current state of the plugin
	}

	// ProcessDoubleFunc defines logic for double signal processing.
	ProcessDoubleFunc func(in, out DoubleBuffer)

	// ProcessFloatFunc defines logic for float signal processing.
	ProcessFloatFunc func(in, out FloatBuffer)

	callbackHandler struct {
		callback C.HostCallback
	}

	dispatchFunc func(op PluginOpcode, index int32, value int64, ptr unsafe.Pointer, opt float32) int64
)

func (d Dispatcher) dispatchFunc(p Plugin) dispatchFunc {
	return func(op PluginOpcode, index int32, value int64, ptr unsafe.Pointer, opt float32) int64 {
		switch op {
		case plugClose:
			if d.CloseFunc != nil {
				d.CloseFunc()
			}
			return 0
		case plugGetParamName:
			s := (*ascii8)(ptr)
			copyASCII(s[:], p.Parameters[index].Name)
		case plugGetParamDisplay:
			s := (*ascii8)(ptr)
			copyASCII(s[:], p.Parameters[index].GetValueLabel())
		case plugGetParamLabel:
			s := (*ascii8)(ptr)
			copyASCII(s[:], p.Parameters[index].Unit)
		case PlugCanBeAutomated:
			if p.Parameters[index].NotAutomated {
				return 0
			}
			return 1
		case plugSetBufferSize:
			if d.SetBufferSizeFunc == nil {
				return 0
			}
			d.SetBufferSizeFunc(int(value))
		case PlugGetPluginName:
			s := (*ascii32)(ptr)
			copyASCII(s[:], p.Name)
		case PlugGetProductString:
			s := (*ascii64)(ptr)
			copyASCII(s[:], p.Name)
		case PlugGetVendorString:
			s := (*ascii64)(ptr)
			copyASCII(s[:], p.Vendor)
		case PlugGetPlugCategory:
			return int64(p.Category)
		case PlugCanDo:
			if d.CanDoFunc == nil {
				return 0
			}
			s := PluginCanDoString(C.GoString((*C.char)(ptr)))
			return int64(d.CanDoFunc(s))
		case PlugProcessEvents:
			if d.ProcessEventsFunc == nil {
				return 0
			}
			var e *EventsPtr = (*EventsPtr)(ptr)
			d.ProcessEventsFunc(e)
		case plugGetChunk:
			if d.GetChunkFunc == nil {
				return 0
			}
			chunk := d.GetChunkFunc(index > 0)
			if len(chunk) == 0 {
				return 0
			}
			*(*unsafe.Pointer)(ptr) = unsafe.Pointer(C.CBytes(chunk))
			return int64(len(chunk))
		case plugSetChunk:
			if d.SetChunkFunc == nil {
				return 0
			}
			bytes := C.GoBytes(ptr, C.int(value))
			d.SetChunkFunc(bytes, index > 0)
			return 0
		default:
			return 0
		}
		return 1
	}
}

func (h callbackHandler) host(cp *C.CPlugin) Host {
	return Host{
		GetSampleRate: func() signal.Frequency {
			return signal.Frequency(C.callbackHost(h.callback, cp, C.int(HostGetSampleRate), 0, 0, nil, 0))
		},
		GetBufferSize: func() int {
			return int(C.callbackHost(h.callback, cp, C.int(HostGetBufferSize), 0, 0, nil, 0))
		},
		GetProcessLevel: func() ProcessLevel {
			return ProcessLevel(C.callbackHost(h.callback, cp, C.int(HostGetCurrentProcessLevel), 0, 0, nil, 0))
		},
		GetTimeInfo: func(flags TimeInfoFlag) *TimeInfo {
			return (*TimeInfo)(unsafe.Pointer(uintptr(C.callbackHost(h.callback, cp, C.int(HostGetTime), 0, C.int64_t(flags), nil, 0))))
		},
		UpdateDisplay: func() bool {
			return C.callbackHost(h.callback, cp, C.int(HostUpdateDisplay), 0, 0, nil, 0) > 0
		},
		Automate: func(index int32, value float32) {
			C.callbackHost(h.callback, cp, C.int(HostAutomate), C.int32_t(index), 0, nil, C.float(value))
		},
		Idle: func() {
			C.callbackHost(h.callback, cp, C.int(HostIdle), 0, 0, nil, 0)
		},
		ProcessEvents: func(events *EventsPtr) {
			C.callbackHost(h.callback, cp, C.int(HostProcessEvents), 0, 0, unsafe.Pointer(events), 0)
		},
		IOChanged: func() bool {
			return C.callbackHost(h.callback, cp, C.int(HostIOChanged), 0, 0, nil, 0) > 0
		},
		SizeWindow: func(width, height int32) bool {
			return C.callbackHost(h.callback, cp, C.int(HostSizeWindow), C.int32_t(width), C.int64_t(height), nil, 0) > 0
		},
		GetInputLatency: func() int32 {
			return int32(C.callbackHost(h.callback, cp, C.int(HostGetInputLatency), 0, 0, nil, 0))
		},
		GetOutputLatency: func() int32 {
			return int32(C.callbackHost(h.callback, cp, C.int(HostGetOutputLatency), 0, 0, nil, 0))
		},
		GetAutomationState: func() AutomationState {
			return AutomationState(C.callbackHost(h.callback, cp, C.int(HostGetAutomationState), 0, 0, nil, 0))
		},
		GetVendorString: func() string {
			var s ascii64
			C.callbackHost(h.callback, cp, C.int(HostGetVendorString), 0, 0, unsafe.Pointer(&s), 0)
			return s.String()
		},
		GetProductString: func() string {
			var s ascii64
			C.callbackHost(h.callback, cp, C.int(HostGetProductString), 0, 0, unsafe.Pointer(&s), 0)
			return s.String()
		},
		GetVendorVersion: func() int32 {
			return int32(C.callbackHost(h.callback, cp, C.int(HostGetVendorVersion), 0, 0, nil, 0))
		},
		CanDo: func(s HostCanDoString) CanDoResponse {
			cs := C.CString(string(s))
			defer C.free(unsafe.Pointer(cs))
			return CanDoResponse(C.callbackHost(h.callback, cp, C.int(HostCanDo), 0, 0, unsafe.Pointer(cs), 0))
		},
		GetLanguage: func() HostLanguage {
			return HostLanguage(C.callbackHost(h.callback, cp, C.int(HostGetLanguage), 0, 0, nil, 0))
		},
		GetDirectory: func() string {
			ptr := unsafe.Pointer(uintptr(C.callbackHost(h.callback, cp, C.int(HostGetDirectory), 0, 0, nil, 0)))
			if ptr == nil {
				return ""
			}
			return C.GoString((*C.char)(ptr))
		},
		BeginEdit: func(index int32) bool {
			return C.callbackHost(h.callback, cp, C.int(HostBeginEdit), C.int32_t(index), 0, nil, 0) > 0
		},
		EndEdit: func(index int32) bool {
			return C.callbackHost(h.callback, cp, C.int(HostEndEdit), C.int32_t(index), 0, nil, 0) > 0
		},
	}
}

func getPlugin(cp *C.CPlugin) *Plugin {
	plugins.RLock()
	defer plugins.RUnlock()
	return plugins.mapping[uintptr(unsafe.Pointer(cp))]
}
