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
		mapping map[unsafe.Pointer]Plugin
	}{
		mapping: map[unsafe.Pointer]Plugin{},
	}
)

type (
	Plugin struct{}

	HostCallback struct {
		callbackFunc C.HostCallback
	}

	PluginAllocatorFunc func(HostCallback) Plugin
)
