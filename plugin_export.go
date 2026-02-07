//go:build plugin
// +build plugin

package vst2

//#include "include/vst.h"
import "C"

import (
	"unsafe"
)

// instantiate go plugin
//
//export newGoPlugin
func newGoPlugin(cp *C.CPlugin, c C.HostCallback) {
	loadHook()
	p, d := PluginAllocator(callbackHandler{c}.host(cp))
	cp.magic = C.int(EffectMagic)
	cp.numInputs = C.int(p.InputChannels)
	cp.numOutputs = C.int(p.OutputChannels)
	cp.numParams = C.int(len(p.Parameters))
	cp.version = C.int(p.Version)
	cp.uniqueID = C.int(uint(p.UniqueID[0])<<24 | uint(p.UniqueID[1])<<16 | uint(p.UniqueID[2])<<8 | uint(p.UniqueID[3])<<0)
	cp.flags = cp.flags | C.int(p.Flags)
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

	// GetChunk and SetChunk should be defined in pairs. If both are
	// defined, we advertise the capability to save and load plugin settings,
	// by settings flag PluginProgramChunks.
	if d.GetChunkFunc != nil && d.SetChunkFunc != nil {
		cp.flags = cp.flags | C.int(PluginProgramChunks)
	}
	p.dispatchFunc = d.dispatchFunc(p)
	plugins.Lock()
	plugins.mapping[uintptr(unsafe.Pointer(cp))] = &p
	plugins.Unlock()
}

// global dispatch, calls real plugin dispatch.
//
//export dispatchPluginBridge
func dispatchPluginBridge(cp *C.CPlugin, opcode int32, index int32, value int64, ptr unsafe.Pointer, opt float32) int64 {
	p := getPlugin(cp)
	pluginOpcode := PluginOpcode(opcode)
	ret := p.dispatchFunc(pluginOpcode, index, value, ptr, opt)
	if pluginOpcode == plugClose {
		plugins.RLock()
		defer plugins.RUnlock()
		delete(plugins.mapping, uintptr(unsafe.Pointer(cp)))
	}
	return ret
}

// global processDouble, calls real plugin processDouble.
//
//export processDoublePluginBridge
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

// global processFloat, calls real plugin processFloat.
//
//export processFloatPluginBridge
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

// global getParameter, calls real plugin getParameter.
//
//export getParameterPluginBridge
func getParameterPluginBridge(cp *C.CPlugin, index int32) float32 {
	return getPlugin(cp).Parameters[index].Value
}

// global setParameter, calls real plugin setParameter.
//
//export setParameterPluginBridge
func setParameterPluginBridge(cp *C.CPlugin, index int32, value float32) {
	getPlugin(cp).Parameters[index].Value = value
}
