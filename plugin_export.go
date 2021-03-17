// +build plugin

package vst2

//#include "include/vst.h"
import "C"
import (
	"fmt"
	"unsafe"
)

//export newGoPlugin
// instantiate go plugin
func newGoPlugin(cp *C.CPlugin, c C.HostCallback) {
	p := PluginAllocator(HostCallback{c})
	plugins.Lock()
	plugins.mapping[unsafe.Pointer(cp)] = p
	plugins.Unlock()
}

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
