package sdk

// #include <stdlib.h>
import "C"
import (
	"unsafe"

	"pipelined.dev/signal"
)

type (
	// DoubleBuffer is a samples buffer for VST ProcessDouble function.
	// C requires all buffer channels to be coallocated. This differs from
	// Go slices.
	DoubleBuffer struct {
		numChannels int
		size        int
		data        []*C.double
	}

	// FloatBuffer is a samples buffer for VST Process function.
	// It should be used only if plugin doesn't support EffFlagsCanDoubleReplacing.
	FloatBuffer struct {
		numChannels int
		size        int
		data        []*C.float
	}
)

// NewDoubleBuffer allocates new memory for C-compatible buffer.
func NewDoubleBuffer(numChannels, bufferSize int) DoubleBuffer {
	b := make([]*C.double, numChannels)
	for i := 0; i < numChannels; i++ {
		b[i] = (*C.double)(C.malloc(C.size_t(C.sizeof_double * bufferSize)))
	}
	return DoubleBuffer{
		data:        b,
		size:        bufferSize,
		numChannels: numChannels,
	}
}

// CopyTo copies values to signal.Floating buffer. If dimensions differ - the lesser used.
func (b DoubleBuffer) CopyTo(s signal.Floating) {
	mustSameChannels(s.Channels(), b.numChannels)
	// determine the size of data by picking up a lesser dimensions.
	bufferSize := min(s.Length(), b.size)

	// copy data.
	for c := 0; c < s.Channels(); c++ {
		row := (*[1 << 30]C.double)(unsafe.Pointer(b.data[c]))
		for i := 0; i < bufferSize; i++ {
			s.SetSample(s.BufferIndex(c, i), float64(row[i]))
		}
	}
}

// CopyFrom copies values from signal.Float64. If dimensions differ - the lesser used.
func (b DoubleBuffer) CopyFrom(s signal.Floating) {
	mustSameChannels(s.Channels(), b.numChannels)
	// determine the size of data by picking up a lesser dimensions.
	bufferSize := min(s.Length(), b.size)

	// copy data.
	for c := 0; c < s.Channels(); c++ {
		row := (*[1 << 30]C.double)(unsafe.Pointer(b.data[c]))
		for i := 0; i < bufferSize; i++ {
			(*row)[i] = C.double(s.Sample(s.BufferIndex(c, i)))
		}
	}
}

// Free the allocated memory.
func (b DoubleBuffer) Free() {
	for i := range b.data {
		C.free(unsafe.Pointer(b.data[i]))
	}
}

// NewFloatBuffer allocates new memory for C-compatible buffer.
func NewFloatBuffer(numChannels, bufferSize int) FloatBuffer {
	b := make([]*C.float, numChannels)
	for i := 0; i < numChannels; i++ {
		b[i] = (*C.float)(C.malloc(C.size_t(C.sizeof_float * bufferSize)))
	}
	return FloatBuffer{
		data:        b,
		size:        bufferSize,
		numChannels: numChannels,
	}
}

// CopyTo copies values to signal.Float64 buffer. If dimensions differ - the lesser used.
func (b FloatBuffer) CopyTo(s signal.Floating) {
	mustSameChannels(s.Channels(), b.numChannels)
	// determine the size of data by picking up a lesser dimensions.
	bufferSize := min(s.Length(), b.size)

	// copy data.
	for c := 0; c < s.Channels(); c++ {
		row := (*[1 << 30]C.float)(unsafe.Pointer(b.data[c]))
		for i := 0; i < bufferSize; i++ {
			s.SetSample(s.BufferIndex(c, i), float64(row[i]))
		}
	}
}

// CopyFrom copies values from signal.Float64. If dimensions differ - the lesser used.
func (b FloatBuffer) CopyFrom(s signal.Floating) {
	mustSameChannels(s.Channels(), b.numChannels)
	// determine the size of data by picking up a lesser dimensions.
	bufferSize := min(s.Length(), b.size)

	// copy data.
	for c := 0; c < s.Channels(); c++ {
		row := (*[1 << 30]C.float)(unsafe.Pointer(b.data[c]))
		for i := 0; i < bufferSize; i++ {
			(*row)[i] = C.float(s.Sample(s.BufferIndex(c, i)))
		}
	}
}

// Free the allocated memory.
func (b FloatBuffer) Free() {
	for _, c := range b.data {
		C.free(unsafe.Pointer(c))
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func mustSameChannels(c1, c2 int) {
	if c1 != c2 {
		panic("different number of channels")
	}
}
