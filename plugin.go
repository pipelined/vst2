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
	Plugin struct {
		InputChannels  int
		OutputChannels int
		inputDouble    DoubleBuffer
		outputDouble   DoubleBuffer
		inputFloat     FloatBuffer
		outputFloat    FloatBuffer
		ProcessDoubleFunc
		ProcessFloatFunc
		Parameters []*Parameter
		dispatchFunc
	}

	Dispatcher struct {
		paramNameFunc      func(int) string
		paramValueNameFunc func(int) string
		paramUnitNameFunc  func(int) string
	}

	callbackHandler struct {
		callback C.HostCallback
	}

	dispatchFunc func(op PluginOpcode, index int32, value int64, ptr unsafe.Pointer, opt float32) int64

	ProcessDoubleFunc func(in, out DoubleBuffer)

	ProcessFloatFunc func(in, out FloatBuffer)

	PluginAllocatorFunc func(Host) (Plugin, Dispatcher)
)

func (d Dispatcher) dispatchFunc(params []*Parameter) dispatchFunc {
	return func(op PluginOpcode, index int32, value int64, ptr unsafe.Pointer, opt float32) int64 {
		return 0
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
