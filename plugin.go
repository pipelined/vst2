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
		mapping map[unsafe.Pointer]*Plugin
	}{
		mapping: map[unsafe.Pointer]*Plugin{},
	}
)

type (
	Plugin struct {
		InputChannels  int
		OutputChannels int
		DispatchFunc
		inputDouble  DoubleBuffer
		outputDouble DoubleBuffer
		inputFloat   FloatBuffer
		outputFloat  FloatBuffer
		ProcessDoubleFunc
		ProcessFloatFunc
		Parameters []*Parameter
	}

	callbackHandler struct {
		callback C.HostCallback
	}

	DispatchFunc func(op PluginOpcode, index int32, value int64, ptr unsafe.Pointer, opt float32) int64

	ProcessDoubleFunc func(in, out DoubleBuffer)

	ProcessFloatFunc func(in, out FloatBuffer)

	PluginAllocatorFunc func(Host) Plugin
)

func (h callbackHandler) host(cp *C.CPlugin) Host {
	return Host{
		GetSampleRate: func() signal.Frequency {
			return signal.Frequency(C.callbackHost(h.callback, cp, C.int(HostGetSampleRate), 0, 0, nil, 0))
		},
		GetBufferSize: func() int {
			return int(C.callbackHost(h.callback, cp, C.int(HostGetBufferSize), 0, 0, nil, 0))
		},
	}
}

func getPlugin(cp *C.CPlugin) *Plugin {
	plugins.RLock()
	defer plugins.RUnlock()
	return plugins.mapping[unsafe.Pointer(cp)]
}
