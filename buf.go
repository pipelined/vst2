package vst2

// #include <stdlib.h>
import "C"
import (
	"unsafe"

	"github.com/pipelined/signal"
)

// DoubleBuffer is a samples buffer for VST Process Double function.
// C requires all buffer channels to be coallocated. This differs from
// Go slices.
type DoubleBuffer struct {
	numChannels int
	size        int
	data        []*C.double
}

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

// CopyFloat64 into DoubleBuffer. If dimensions differ - the lesser used.
func CopyFloat64(s signal.Float64, b DoubleBuffer) {
	// determine the size of data by picking up a lesser dimensions.
	numChannels := min(s.NumChannels(), b.numChannels)
	bufferSize := min(s.Size(), s.Size())

	// copy data.
	for i := 0; i < numChannels; i++ {
		row := (*[1 << 30]C.double)(unsafe.Pointer(b.data[i]))
		for j := 0; j < bufferSize; j++ {
			s[i][j] = float64(row[j])
		}
	}
}

// CopyDouble buffer into signal.Float64. If dimensions differ - the lesser used.
func CopyDouble(b DoubleBuffer, s signal.Float64) {
	// determine the size of data by picking up a lesser dimensions.
	numChannels := min(s.NumChannels(), b.numChannels)
	bufferSize := min(s.Size(), s.Size())

	// copy data.
	for i := 0; i < numChannels; i++ {
		row := (*[1 << 30]C.double)(unsafe.Pointer(b.data[i]))
		for j := 0; j < bufferSize; j++ {
			(*row)[j] = C.double(s[i][j])
		}
	}
}

// Free the allocated memory.
func (b DoubleBuffer) Free() {
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
