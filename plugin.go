// +build plugin

package vst2

//#include "include/vst.h"
import "C"
import (
	"sync"
	"unsafe"
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
		ProcessDoubleFunc
		Parameters []Parameter
	}

	HostCallback struct {
		callbackFunc C.HostCallback
	}

	DispatchFunc func(op PluginOpcode, index int32, value int64, ptr unsafe.Pointer, opt float32) int64

	ProcessDoubleFunc func(in, out DoubleBuffer)

	PluginAllocatorFunc func(HostCallback) Plugin
)

func getPlugin(cp *C.CPlugin) *Plugin {
	plugins.RLock()
	defer plugins.RUnlock()
	return plugins.mapping[unsafe.Pointer(cp)]
}
