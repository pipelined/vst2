package vst2

/*
#cgo CFLAGS: -std=gnu99 -I${SRCDIR}
#include "sdk.h"
*/
import "C"
import (
	"sync"
	"time"
	"unsafe"
)

type (
	// HostCallbackFunc used as callback function called by plugin. Use closure
	// wrapping technique to add more types to callback.
	HostCallbackFunc func(HostOpcode, Index, Value, unsafe.Pointer, Opt) Return

	// Index is index in plugin dispatch/host callback.
	Index int64
	// Value is value in plugin dispatch/host callback.
	Value int64
	// Opt is opt in plugin dispatch/host callback.
	Opt float64
	// Return is returned value for dispatch/host callback.
	Return int64
)

// global state for callbacks.
var (
	mutex     sync.RWMutex
	callbacks = make(map[*Plugin]HostCallbackFunc)
)

//export hostCallback
// global hostCallback, calls real callback.
func hostCallback(p *Plugin, opcode int64, index int64, value int64, ptr unsafe.Pointer, opt float64) Return {
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

// DefaultHostCallback returns default vst2 host callback.
func DefaultHostCallback(props *HostProperties) HostCallbackFunc {
	return func(opcode HostOpcode, index Index, value Value, ptr unsafe.Pointer, opt Opt) Return {
		switch opcode {
		case HostGetCurrentProcessLevel:
			return Return(ProcessLevelRealtime)
		case HostGetSampleRate:
			return Return(props.SampleRate)
		case HostGetBlockSize:
			return Return(props.SampleRate)
		case HostGetTime:
			ti := &TimeInfo{
				SampleRate:         float64(props.SampleRate),
				SamplePos:          float64(props.CurrentPosition),
				NanoSeconds:        float64(time.Now().UnixNano()),
				TimeSigNumerator:   4,
				TimeSigDenominator: 4,
			}
			return ti.Return()
		default:
			break
		}
		return 0
	}
}
