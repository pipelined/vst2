package vst2

/*
#cgo CFLAGS: -I${SRCDIR}
#include "vst.h"
*/
import "C"
import (
	"fmt"
	"sync"
	"unsafe"
)

var (
	// global state for callbacks.
	plugins = struct {
		sync.RWMutex
		mapping map[unsafe.Pointer]Plugin
	}{
		mapping: map[unsafe.Pointer]Plugin{},
	}
)

// Plugin is an instance of loaded VST plugin.
type Plugin struct {
	p  *C.CPlugin
	hc *C.HostCallback
}

func registerPlugin(p *C.CPlugin) {}

//export dispatchPluginBridge
// global dispatch, calls real plugin dispatch.
func dispatchPluginBridge(p *C.CPlugin, opcode int32, index int32, value int64, ptr unsafe.Pointer, opt float32) int64 {
	fmt.Printf("got dispatch opcode: %d\n", opcode)
	return 0
}

//export processDoublePluginBridge
// global processDouble, calls real plugin processDouble.
func processDoublePluginBridge(p *C.CPlugin, in, out **float64, sampleFrames int32) {
	return
}

//export processFloatPluginBridge
// global processFloat, calls real plugin processFloat.
func processFloatPluginBridge(p *C.CPlugin, in, out **float32, sampleFrames int32) {
	return
}

//export getParameterPluginBridge
// global getParameter, calls real plugin getParameter.
func getParameterPluginBridge(p *C.CPlugin, index int32) float32 {
	return 0
}

//export setParameterPluginBridge
// global setParameter, calls real plugin setParameter.
func setParameterPluginBridge(p *C.CPlugin, index int32, value float32) {
	return
}
