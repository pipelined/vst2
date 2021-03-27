package vst2

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
		Frames int
		data   []*C.double
	}

	// FloatBuffer is a samples buffer for VST Process function.
	// It should be used only if plugin doesn't support EffFlagsCanDoubleReplacing.
	FloatBuffer struct {
		Frames int
		data   []*C.float
	}
)

// NewDoubleBuffer allocates new memory for C-compatible buffer.
func NewDoubleBuffer(numChannels, bufferSize int) DoubleBuffer {
	b := make([]*C.double, numChannels)
	for i := 0; i < numChannels; i++ {
		b[i] = (*C.double)(C.malloc(C.size_t(C.sizeof_double * bufferSize)))
	}
	return DoubleBuffer{
		data:   b,
		Frames: bufferSize,
	}
}

// CopyTo copies values to signal.Floating buffer. If dimensions differ - the lesser used.
func (b DoubleBuffer) CopyTo(s signal.Floating) {
	mustSameChannels(s.Channels(), len(b.data))
	// determine the size of data by picking up a lesser dimensions.
	bufferSize := min(s.Length(), b.Frames)

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
	mustSameChannels(s.Channels(), len(b.data))
	// determine the size of data by picking up a lesser dimensions.
	bufferSize := min(s.Length(), b.Frames)

	// copy data.
	for c := 0; c < s.Channels(); c++ {
		row := (*[1 << 30]C.double)(unsafe.Pointer(b.data[c]))
		for i := 0; i < bufferSize; i++ {
			(*row)[i] = C.double(s.Sample(s.BufferIndex(c, i)))
		}
	}
}

// cArray returns C array that is used as storage for buffer.
func (b DoubleBuffer) cArray() **C.double {
	return (**C.double)(unsafe.Pointer(&b.data[0]))
}

// Channel returns slice that's backed by C array and stores samples from
// single channel.
func (b DoubleBuffer) Channel(i int) []float64 {
	return (*(*[1 << 30]float64)(unsafe.Pointer(b.data[i])))[:b.Frames:b.Frames]
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
		data:   b,
		Frames: bufferSize,
	}
}

// CopyTo copies values to signal.Float64 buffer. If dimensions differ - the lesser used.
func (b FloatBuffer) CopyTo(s signal.Floating) {
	mustSameChannels(s.Channels(), len(b.data))
	// determine the size of data by picking up a lesser dimensions.
	bufferSize := min(s.Length(), b.Frames)

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
	mustSameChannels(s.Channels(), len(b.data))
	// determine the size of data by picking up a lesser dimensions.
	bufferSize := min(s.Length(), b.Frames)

	// copy data.
	for c := 0; c < s.Channels(); c++ {
		row := (*[1 << 30]C.float)(unsafe.Pointer(b.data[c]))
		for i := 0; i < bufferSize; i++ {
			(*row)[i] = C.float(s.Sample(s.BufferIndex(c, i)))
		}
	}
}

// cArray returns C array that is used as storage for buffer.
func (b FloatBuffer) cArray() **C.float {
	return (**C.float)(unsafe.Pointer(&b.data[0]))
}

// Channel returns slice that's backed by C array and stores samples from
// single channel.
func (b FloatBuffer) Channel(i int) []float32 {
	return (*(*[1 << 30]float32)(unsafe.Pointer(b.data[i])))[:b.Frames:b.Frames]
}

// Free the allocated memory.
func (b FloatBuffer) Free() {
	for _, c := range b.data {
		C.free(unsafe.Pointer(c))
	}
}

// getDoubleChannel returns single channel of C buffer. This function
// refers C type, so it shouldn't be used by users of the package.
func getDoubleChannel(buf **C.double, i int) *C.double {
	ptrPtr := (**C.double)(unsafe.Pointer(uintptr(unsafe.Pointer(buf)) + uintptr(i)*unsafe.Sizeof(*buf)))
	return *ptrPtr
}

// getFloatChannel returns single channel of C buffer. This function
// refers C type, so it shouldn't be used by users of the package.
func getFloatChannel(buf **C.float, i int) *C.float {
	ptrPtr := (**C.float)(unsafe.Pointer(uintptr(unsafe.Pointer(buf)) + uintptr(i)*unsafe.Sizeof(*buf)))
	return *ptrPtr
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
