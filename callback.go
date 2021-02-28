package vst2

/*
#cgo CFLAGS: -std=gnu99 -I${SRCDIR}
#include "sdk.h"
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
	callbacks = make(map[*Plugin]HostCallbackFunc)
)

//export hostCallback
// global hostCallback, calls real callback.
func hostCallback(p *Plugin, opcode int64, index int32, value int64, ptr unsafe.Pointer, opt float32) uintptr {
	// HostVersion is requested when plugin is created
	// It's never in map
	if HostOpcode(opcode) == HostVersion {
		return version
	}
	mutex.RLock()
	c, ok := callbacks[p]
	mutex.RUnlock()
	if !ok {
		panic("plugin was closed")
	}

	if c == nil {
		panic("host callback is undefined")
	}
	return c(HostOpcode(opcode), index, value, ptr, opt)
}

type (
	// HostCallbackFunc used as callback function called by plugin. Use closure
	// wrapping technique to add more types to callback.
	HostCallbackFunc func(op HostOpcode, index int32, value int64, ptr unsafe.Pointer, opt float32) uintptr

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

// Callback returns HostCallbackFunc that handles all vst types casts
// and allows to write handlers without usage of unsafe package.
func (h Host) Callback() HostCallbackFunc {
	return func(op HostOpcode, index int32, value int64, ptr unsafe.Pointer, opt float32) uintptr {
		switch op {
		case HostGetCurrentProcessLevel:
			if h.GetProcessLevel != nil {
				return uintptr(h.GetProcessLevel())
			}
		case HostGetSampleRate:
			if h.GetSampleRate != nil {
				return uintptr(h.GetSampleRate())
			}
		case HostGetBlockSize:
			if h.GetBufferSize != nil {
				return uintptr(h.GetBufferSize())
			}
		case HostGetTime:
			if h.GetTimeInfo != nil {
				return uintptr(unsafe.Pointer(h.GetTimeInfo()))
			}
		}
		return 0
	}
}

// NoopHostCallback returns dummy host callback that just prints received
// opcodes.
func NoopHostCallback() HostCallbackFunc {
	return func(op HostOpcode, index int32, value int64, ptr unsafe.Pointer, opt float32) uintptr {
		fmt.Printf("host received opcode: %v\n", op)
		return 0
	}
}
