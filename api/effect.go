package api

/*
#cgo CFLAGS: -std=gnu99 -I${SRCDIR}
#include <stdlib.h>
#include <stdint.h>
#include "effect.h"
*/
import "C"
import (
	"sync"
	"unsafe"
)

// global state for callbacks.
var (
	mutex     sync.RWMutex
	callbacks = make(map[*Effect]HostCallbackFunc)
)

//export hostCallback
// calls real callback
func hostCallback(e *Effect, opcode int64, index int64, value int64, ptr unsafe.Pointer, opt float64) Return {
	// AudioMasterVersion is requested when plugin is created
	// It's never in map
	if HostOpcode(opcode) == HostVersion {
		return version
	}
	mutex.RLock()
	c, ok := callbacks[e]
	mutex.RUnlock()
	if !ok {
		panic("plugin was closed")
	}

	if c == nil {
		panic("host callback is undefined")
	}
	return c(e, HostOpcode(opcode), Index(index), Value(value), Ptr(ptr), Opt(opt))
}

const (
	// VST main function name.
	main = "VSTPluginMain"
	// VST API version.
	version = 2400
)

type (
	// EntryPoint is a reference to VST main function. It also keeps
	// reference to VST handle to clean up on Close.
	EntryPoint struct {
		main effectMain
		handle
	}

	// Effect is an alias on C effect type.
	Effect C.Effect

	// HostCallbackFunc used as callback function called by plugin.
	HostCallbackFunc func(*Effect, HostOpcode, Index, Value, Ptr, Opt) Return

	// Index is index in plugin dispatch/host callback.
	Index int64
	// Value is value in plugin dispatch/host callback.
	Value int64
	// Ptr is ptr in plugin dispatch/host callback.
	Ptr unsafe.Pointer
	// Opt is opt in plugin dispatch/host callback.
	Opt float64
	// Return is returned value for dispatch/host callback.
	Return int64

	effectMain C.EntryPoint
)

// Close cleans up VST handle.
func (e EntryPoint) Close() error {
	if e.main == nil {
		return nil
	}
	e.main = nil
	return e.handle.close()
}

// Load new instance of VST plugin with provided callback.
// This function also calls dispatch with EffOpen opcode.
func (e EntryPoint) Load(c HostCallbackFunc) *Effect {
	ef := (*Effect)(C.loadEffect(e.main))
	mutex.Lock()
	callbacks[ef] = c
	mutex.Unlock()
	ef.Dispatch(EffOpen, 0, 0, nil, 0.0)
	return ef
}

// Dispatch wraps-up C method to dispatch calls to plugin
func (e *Effect) Dispatch(opcode EffectOpcode, index Index, value Value, ptr Ptr, opt Opt) Return {
	return Return(C.dispatch((*C.Effect)(e), C.int(opcode), C.int(index), C.int64_t(value), unsafe.Pointer(ptr), C.float(opt)))
}

// CanProcessFloat32 checks if plugin can process float32.
func (e *Effect) CanProcessFloat32() bool {
	if e == nil {
		return false
	}
	return EffectFlags(e.flags)&EffFlagsCanReplacing == EffFlagsCanReplacing
}

// CanProcessFloat64 checks if plugin can process float64.
func (e *Effect) CanProcessFloat64() bool {
	if e == nil {
		return false
	}
	return EffectFlags(e.flags)&EffFlagsCanDoubleReplacing == EffFlagsCanDoubleReplacing
}

// ProcessFloat64 audio with VST plugin.
// TODO: add c buffer parameter.
func (e *Effect) ProcessFloat64(in [][]float64) [][]float64 {
	numChannels := len(in)
	blocksize := len(in[0])

	// convert [][]float64 to []*C.double
	input := make([]*C.double, numChannels)
	output := make([]*C.double, numChannels)
	for i, row := range in {
		// allocate input memory for C layout
		inp := (*C.double)(C.malloc(C.size_t(C.sizeof_double * blocksize)))
		input[i] = inp
		defer C.free(unsafe.Pointer(inp))

		// copy data from slice to C array
		pa := (*[1 << 30]C.double)(unsafe.Pointer(inp))
		for j, v := range row {
			(*pa)[j] = C.double(v)
		}

		// allocate output memory for C layout
		outp := (*C.double)(C.malloc(C.size_t(C.sizeof_double * blocksize)))
		output[i] = outp
		defer C.free(unsafe.Pointer(outp))
	}

	C.processDouble((*C.Effect)(e), C.int(numChannels), C.int(blocksize), &input[0], &output[0])

	//convert []*C.double slices to [][]float64
	out := make([][]float64, numChannels)
	for i, data := range output {
		// copy data from C array to slice
		pa := (*[1 << 30]C.float)(unsafe.Pointer(data))
		out[i] = make([]float64, blocksize)
		for j := range out[i] {
			out[i][j] = float64(pa[j])
		}
	}
	return out
}

// ProcessFloat32 audio with VST plugin.
// TODO: add c buffer parameter.
func (e *Effect) ProcessFloat32(in [][]float32) (out [][]float32) {
	numChannels := len(in)
	blocksize := len(in[0])

	// convert [][]float32 to []*C.float
	input := make([]*C.float, numChannels)
	output := make([]*C.float, numChannels)
	for i, row := range in {
		// allocate input memory for C layout
		inp := (*C.float)(C.malloc(C.size_t(C.sizeof_float * blocksize)))
		input[i] = inp
		defer C.free(unsafe.Pointer(inp))

		// copy data from slice to C array
		pa := (*[1 << 30]C.float)(unsafe.Pointer(inp))
		for j, v := range row {
			(*pa)[j] = C.float(v)
		}

		// allocate output memory for C layout
		outp := (*C.float)(C.malloc(C.size_t(C.sizeof_float * blocksize)))
		output[i] = outp
		defer C.free(unsafe.Pointer(outp))
	}

	C.processFloat((*C.Effect)(e), C.int(numChannels), C.int(blocksize), &input[0], &output[0])

	//convert []*C.float slices to [][]float32
	out = make([][]float32, numChannels)
	for i, data := range output {
		// copy data from C array to slice
		pa := (*[1 << 30]C.float)(unsafe.Pointer(data))
		out[i] = make([]float32, blocksize)
		for j := range out[i] {
			out[i][j] = float32(pa[j])
		}
	}
	return out
}

// Start the plugin.
func (e *Effect) Start() {
	e.Dispatch(EffStateChanged, 0, 1, nil, 0.0)
}

// Stop the plugin.
func (e *Effect) Stop() {
	e.Dispatch(EffStateChanged, 0, 0, nil, 0.0)
}

// SetBufferSize sets a buffer size
func (e *Effect) SetBufferSize(bufferSize int) {
	e.Dispatch(EffSetBufferSize, 0, Value(bufferSize), nil, 0.0)
}

// SetSampleRate sets a sample rate for plugin
func (e *Effect) SetSampleRate(sampleRate int) {
	e.Dispatch(EffSetSampleRate, 0, 0, nil, Opt(sampleRate))
}

// SetSpeakerArrangement craetes and passes SpeakerArrangement structures to plugin
func (e *Effect) SetSpeakerArrangement(in, out *SpeakerArrangement) {
	e.Dispatch(EffSetSpeakerArrangement, 0, in.Value(), out.Ptr(), 0.0)
}
