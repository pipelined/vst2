package main

/*
#cgo CFLAGS: -I${SRCDIR}/..
#include "plugin.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// (Plugin* plugin, int32_t opcode, int32_t index, int64_t value, void* ptr, float opt);

//export dispatch
// global dispatch, calls real plugin dispatch.
func dispatch(_ *C.Plugin, opcode int32, index int32, value int64, ptr unsafe.Pointer, opt float32) int64 {
	fmt.Printf("got dispatch opcode: %d\n", opcode)
	return 0
}

//export processDouble
// global processDouble, calls real plugin processDouble.
func processDouble(_ *C.Plugin, in, out **float64, sampleFrames int32) {
	return
}

//export processFloat
// global processFloat, calls real plugin processFloat.
func processFloat(_ *C.Plugin, in, out **float32, sampleFrames int32) {
	return
}

//export getParameter
// global getParameter, calls real plugin getParameter.
func getParameter(_ *C.Plugin, index int32) float32 {
	return 0
}

//export setParameter
// global setParameter, calls real plugin setParameter.
func setParameter(_ *C.Plugin, index int32, value float32) {
	return
}
