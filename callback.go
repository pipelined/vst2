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
func hostCallback(p *Plugin, opcode int64, index int64, value int64, ptr unsafe.Pointer, opt float64) uintptr {
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
	return c(HostOpcode(opcode), Index(index), Value(value), ptr, Opt(opt))
}

type (
	// HostCallbackFunc used as callback function called by plugin. Use closure
	// wrapping technique to add more types to callback.
	HostCallbackFunc func(HostOpcode, Index, Value, unsafe.Pointer, Opt) uintptr

	// Index is index in plugin dispatch/host callback.
	Index int64
	// Value is value in plugin dispatch/host callback.
	Value int64
	// Opt is opt in plugin dispatch/host callback.
	Opt float64
	// // Return is returned value for dispatch/host callback.
	// Return uintptr

	// HostCallbackAllocator returns new host callback function that does
	// proper casting of plugin calls.
	HostCallbackAllocator struct {
		GetSampleRate   HostGetSampleRateFunc
		GetBufferSize   HostGetBufferSizeFunc
		GetProcessLevel HostGetProcessLevel
		GetTimeInfo     HostGetTimeInfo
	}

	// HostGetSampleRateFunc returns host sample rate.
	HostGetSampleRateFunc func() signal.Frequency
	// HostGetBufferSizeFunc returns host buffer size.
	HostGetBufferSizeFunc func() int
	// HostGetProcessLevel returns the context of execution.
	HostGetProcessLevel func() ProcessLevel
	// HostGetTimeInfo returns current time info.
	HostGetTimeInfo func() *TimeInfo
)

func HostCallback(h HostCallbackAllocator) HostCallbackFunc {
	return func(opcode HostOpcode, index Index, value Value, ptr unsafe.Pointer, opt Opt) uintptr {
		switch opcode {
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

// DefaultHostCallback returns default vst2 host callback.
func DefaultHostCallback() HostCallbackFunc {
	return func(opcode HostOpcode, index Index, value Value, ptr unsafe.Pointer, opt Opt) uintptr {
		fmt.Printf("host received opcode: %v\n", opcode)
		return 0
	}
}
