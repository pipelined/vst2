// +build plugin

package vst2

//#include "include/plugin/plugin.c"
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
		PluginCanDoFunc
		ProcessEventsFunc
	}

	// Dispatcher handles plugin dispatch calls from the host.
	Dispatcher struct {
		SetBufferSizeFunc func(size int)
	}

	// ProcessDoubleFunc defines logic for double signal processing.
	ProcessDoubleFunc func(in, out DoubleBuffer)

	// ProcessFloatFunc defines logic for float signal processing.
	ProcessFloatFunc func(in, out FloatBuffer)

	ProcessEventsFunc func(*EventsPtr)

	PluginCanDoFunc func(PluginCanDoString) CanDoResponse

	callbackHandler struct {
		callback C.HostCallback
	}

	dispatchFunc func(op PluginOpcode, index int32, value int64, ptr unsafe.Pointer, opt float32) int64
)

func (d Dispatcher) dispatchFunc(p Plugin) dispatchFunc {
	return func(op PluginOpcode, index int32, value int64, ptr unsafe.Pointer, opt float32) int64 {
		switch op {
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
			if p.PluginCanDoFunc != nil {
				s := PluginCanDoString(C.GoString((*C.char)(ptr)))
				return int64(p.PluginCanDoFunc(s))
			}
			return 0
		case PlugProcessEvents:
			var e *EventsPtr = (*EventsPtr)(ptr)
			if p.ProcessEventsFunc != nil {
				p.ProcessEventsFunc(e)
			}
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
		GetTimeInfo: func() *TimeInfo {
			return (*TimeInfo)(unsafe.Pointer(uintptr(C.callbackHost(h.callback, cp, C.int(HostGetTime), 0, 0, nil, 0))))
		},
	}
}

func getPlugin(cp *C.CPlugin) *Plugin {
	plugins.RLock()
	defer plugins.RUnlock()
	return plugins.mapping[uintptr(unsafe.Pointer(cp))]
}
