// +build plugin

package vst2

//#include "include/vst.h"
import "C"
import (
	"unsafe"
)

// instantiate go plugin
//export newGoPlugin
func newGoPlugin(cp *C.CPlugin, c C.HostCallback) {
	p, d := PluginAllocator(callbackHandler{c}.host(cp))
	p.dispatchFunc = d.dispatchFunc(p.Parameters)
	cp.magic = C.int(EffectMagic)
	cp.numInputs = C.int(p.InputChannels)
	cp.numOutputs = C.int(p.OutputChannels)
	cp.numParams = C.int(len(p.Parameters))
	cp.version = C.int(p.Version)
	cp.uniqueID = C.int(p.UniqueID[0]<<24 | p.UniqueID[1]<<16 | p.UniqueID[2]<<8 | p.UniqueID[3]<<0)
	if p.ProcessDoubleFunc != nil {
		cp.flags = cp.flags | C.int(PluginDoubleProcessing)
		p.inputDouble = DoubleBuffer{data: make([]*C.double, p.InputChannels)}
		p.outputDouble = DoubleBuffer{data: make([]*C.double, p.OutputChannels)}
	}
	if p.ProcessFloatFunc != nil {
		cp.flags = cp.flags | C.int(PluginFloatProcessing)
		p.inputFloat = FloatBuffer{data: make([]*C.float, p.InputChannels)}
		p.outputFloat = FloatBuffer{data: make([]*C.float, p.OutputChannels)}
	}
	plugins.Lock()
	plugins.mapping[uintptr(unsafe.Pointer(cp))] = &p
	plugins.Unlock()
}

//export dispatchPluginBridge
// global dispatch, calls real plugin dispatch.
func dispatchPluginBridge(cp *C.CPlugin, opcode int32, index int32, value int64, ptr unsafe.Pointer, opt float32) int64 {
	p := getPlugin(cp)
	return p.dispatchFunc(PluginOpcode(opcode), index, value, ptr, opt)
}

//export processDoublePluginBridge
// global processDouble, calls real plugin processDouble.
func processDoublePluginBridge(cp *C.CPlugin, in, out **C.double, sampleFrames int32) {
	p := getPlugin(cp)
	for i := range p.inputDouble.data {
		p.inputDouble.data[i] = getDoubleChannel(in, i)
	}
	for i := range p.outputDouble.data {
		p.outputDouble.data[i] = getDoubleChannel(out, i)
	}
	p.inputDouble.Frames = int(sampleFrames)
	p.outputDouble.Frames = int(sampleFrames)
	p.ProcessDoubleFunc(p.inputDouble, p.outputDouble)
	return
}

//export processFloatPluginBridge
// global processFloat, calls real plugin processFloat.
func processFloatPluginBridge(cp *C.CPlugin, in, out **C.float, sampleFrames int32) {
	p := getPlugin(cp)
	for i := range p.inputFloat.data {
		p.inputFloat.data[i] = getFloatChannel(in, i)
	}
	for i := range p.outputFloat.data {
		p.outputFloat.data[i] = getFloatChannel(out, i)
	}
	p.inputFloat.Frames = int(sampleFrames)
	p.outputFloat.Frames = int(sampleFrames)
	p.ProcessFloatFunc(p.inputFloat, p.outputFloat)
	return
}

//export getParameterPluginBridge
// global getParameter, calls real plugin getParameter.
func getParameterPluginBridge(cp *C.CPlugin, index int32) float32 {
	return getPlugin(cp).Parameters[index].Value
}

//export setParameterPluginBridge
// global setParameter, calls real plugin setParameter.
func setParameterPluginBridge(cp *C.CPlugin, index int32, value float32) {
	getPlugin(cp).Parameters[index].Value = value
}
