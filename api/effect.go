package api

/*
#cgo CFLAGS: -std=gnu99 -I${SRCDIR}
#include <stdlib.h>
#include <stdint.h>
#include "api.h"
*/
import "C"
import (
	"sync"
	"unsafe"
)

const (
	main    = "VSTPluginMain"
	version = 2400
)

var (
	mutex     sync.RWMutex
	callbacks = make(map[*Effect]HostCallbackFunc)
)

type (
	EntryPoint struct {
		main effectMain
		handle
	}

	Effect C.Effect

	// HostCallbackFunc used as callback from plugin
	HostCallbackFunc func(*Effect, HostOpcode, int64, int64, unsafe.Pointer, float64) int

	effectMain C.vstPluginFuncPtr
)

func (e EntryPoint) Close() error {
	return e.handle.close()
}

func Load(m EntryPoint, c HostCallbackFunc) *Effect {
	e := (*Effect)(C.loadEffect(m.main))
	mutex.Lock()
	callbacks[e] = c
	mutex.Unlock()
	return e
}

//export hostCallback
// calls real callback
func hostCallback(e *Effect, opcode int64, index int64, value int64, ptr unsafe.Pointer, opt float64) int {
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
	return c(e, HostOpcode(opcode), index, value, ptr, opt)
}

// Dispatch wraps-up C method to dispatch calls to plugin
func (e *Effect) Dispatch(opcode EffectOpcode, index int64, value int64, ptr unsafe.Pointer, opt float64) {
	C.dispatch((*C.Effect)(e), C.int(opcode), C.int(index), C.int64_t(value), ptr, C.float(opt))
}

// CanProcessFloat32 checks if plugin can process float32
func (e *Effect) CanProcessFloat32() bool {
	if e == nil {
		return false
	}
	return EffectFlags(e.flags)&EffFlagsCanReplacing == EffFlagsCanReplacing
}

// CanProcessFloat64 checks if plugin can process float64
func (e *Effect) CanProcessFloat64() bool {
	if e == nil {
		return false
	}
	return EffectFlags(e.flags)&EffFlagsCanDoubleReplacing == EffFlagsCanDoubleReplacing
}

// ProcessFloat64 audio with VST plugin
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

// ProcessFloat32 audio with VST plugin
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
