package vst2

/*
#cgo CFLAGS: -std=gnu99 -I${SRCDIR}/vst2sdk/

#include "vst2.h"
#include "aeffectx.h"
*/
import "C"

import (
	"log"
	"unsafe"
)

//Plugin type provides interface
type Plugin struct {
	entryPoint uintptr
	effect     *C.AEffect
}

//HostCallbackFunc used as callback from plugin
type HostCallbackFunc func(*Plugin, MasterOpcode, int64, int64, unsafe.Pointer, float64) int

//pluginOpcode used to wrap C opcodes values
type pluginOpcode uint64

//Constants for audio master callback opcodes
//Host -> Plugin
const (
	EffEditIdle      = pluginOpcode(C.effEditIdle)
	EffMainsChanged  = pluginOpcode(C.effMainsChanged)
	EffOpen          = pluginOpcode(C.effOpen)
	EffSetSampleRate = pluginOpcode(C.effSetSampleRate)
	EffSetBlockSize  = pluginOpcode(C.effSetBlockSize)
)

const (
	vstMain string = "VSTPluginMain"
)

//MasterOpcode used to wrap C opcodes values
type MasterOpcode uint64

//Constants for audio master callback opcodes
//Plugin -> Host
const (
	AudioMasterAutomate  = MasterOpcode(C.audioMasterAutomate)
	AudioMasterVersion   = MasterOpcode(C.audioMasterVersion)
	AudioMasterCurrentID = MasterOpcode(C.audioMasterCurrentId)
	AudioMasterIdle      = MasterOpcode(C.audioMasterIdle)
	AudioMasterGetTime   = MasterOpcode(C.audioMasterGetTime)
)

var (
	callback HostCallbackFunc = HostCallback
)

//NewPlugin loads the plugin into memory and stores callback func
//TODO: catch panic
func NewPlugin(path string) (*Plugin, error) {
	//Get pointer to plugin's Main function
	mainEntryPoint, err := getEntryPoint(path)
	if err != nil {
		log.Printf("Failed to obtain VST entry point '%s': %v\n", path, err)
		return nil, err
	}

	return &Plugin{entryPoint: mainEntryPoint}, nil
}

//Dispatch wraps-up C method to dispatch calls to plugin
func (plugin *Plugin) Dispatch(opcode pluginOpcode, index int64, value int64, ptr unsafe.Pointer, opt float64) {
	if plugin.effect != nil {
		C.dispatch(plugin.effect, C.int(opcode), C.int(index), C.int(value), ptr, C.float(opt))
	}
}

//Resume the plugin
func (plugin *Plugin) resume() {
	plugin.Dispatch(EffMainsChanged, 0, 1, nil, 0.0)
}

//Suspend the plugin
func (plugin *Plugin) suspend() {
	plugin.Dispatch(EffMainsChanged, 0, 0, nil, 0.0)
}

//Start the plugin
func (plugin *Plugin) Start() {
	//Convert to C++ pointer type
	vstEntryPoint := (C.vstPluginFuncPtr)(unsafe.Pointer(plugin.entryPoint))
	plugin.effect = C.loadEffect(vstEntryPoint)

	plugin.Dispatch(EffOpen, 0, 0, nil, 0.0)

	// Set default sample rate and block size
	sampleRate := 44100.0
	plugin.Dispatch(EffSetSampleRate, 0, 0, nil, sampleRate)

	blocksize := int64(4096)
	plugin.Dispatch(EffSetBlockSize, 0, blocksize, nil, 0.0)
}

//Process audio with VST plugin
func (plugin *Plugin) Process(samples [][]float64) (processed [][]float64) {
	//convert Samples to C type
	inSamples := (**C.double)(unsafe.Pointer(&samples[0][0]))
	blocksize := C.int(len(samples[0]))
	numChannels := C.int(len(samples))
	//call plugin and convert result to slice of slices
	outSamples := (*[1 << 30]*C.double)(unsafe.Pointer(C.processAudio(plugin.effect, numChannels, blocksize, inSamples)))[:numChannels]
	//convert slices to [][]float64
	processed = make([][]float64, numChannels)
	for channel, data := range outSamples {
		processed[channel] = (*[1 << 30]float64)(unsafe.Pointer(data))[:blocksize]
	}
	return processed
}

//SetHostCallback allows to override default host callback with custom implementation
func SetHostCallback(newCallback HostCallbackFunc) {
	if newCallback != nil {
		callback = newCallback
	}
}

//export hostCallback
//calls real callback
func hostCallback(effect *C.AEffect, opcode int64, index int64, value int64, ptr unsafe.Pointer, opt float64) int {
	if callback == nil {
		panic("Host callback is not defined!")
	}

	return callback(&Plugin{effect: effect}, MasterOpcode(opcode), index, value, ptr, opt)
}

//HostCallback is a default callback, can be overriden with SetHostCallback
func HostCallback(plugin *Plugin, opcode MasterOpcode, index int64, value int64, ptr unsafe.Pointer, opt float64) int {
	switch opcode {
	case AudioMasterVersion:
		return 2400
	case AudioMasterIdle:
		plugin.Dispatch(EffEditIdle, 0, 0, nil, 0)
	default:
		log.Printf("Plugin requested value of opcode %v\n", opcode)
		break
	}
	return 0
}
