package vst2

/*
#include "aeffectx.h"

typedef AEffect * (*vstPluginFuncPtr)(audioMasterCallback host);

AEffect * loadEffect(AEffect * (*load)(vstPluginFuncPtr));

double** processAudio(AEffect *effect, int numChannels, int blocksize, double ** goInputs);

int dispatch(AEffect *effect, int opCode, int index, int value, void *ptr, float opt);
*/
import "C"

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"
)

//Plugin type provides interface
type Plugin struct {
	entryPoint uintptr
	effect     *C.AEffect
	callback   HostCallbackFunc
}

//HostCallbackFunc used as callback from plugin
type HostCallbackFunc func(*Plugin, int64, int64, int64, unsafe.Pointer, float64) int

//init package
func init() {
	opCode := new(C.enum_AudioMasterOpcodes)
	fmt.Printf("OpCode: %v, type: %t\n", *opCode, *opCode)
}

//NewPlugin loads the plugin into memory and stores callback func
//TODO: catch panic
func NewPlugin(path string) (*Plugin, error) {
	//Load plugin by path
	moduleHandle, err := syscall.LoadLibrary(path)
	if err != nil {
		log.Printf("Failed to load VST from '%s': %v\n", path, err)
		return nil, err
	}

	//Get pointer to plugin's Main function
	mainEntryPoint, err := syscall.GetProcAddress(moduleHandle, "VSTPluginMain")
	if err != nil {
		log.Printf("Failed to obtain VST entry point '%s': %v\n", path, err)
		return nil, err
	}

	return &Plugin{entryPoint: mainEntryPoint}, nil
}

//Resumes the plugin
func (plugin *Plugin) resume() {
	C.dispatch(plugin.effect, C.effMainsChanged, 0, 1, nil, 0.0)
}

//Suspends the plugin
func (plugin *Plugin) suspend() {
	C.dispatch(plugin.effect, C.effMainsChanged, 0, 0, nil, 0.0)
}

//Starts the plugin
func (plugin *Plugin) start() {
	//Convert to C++ pointer type
	vstEntryPoint := (C.vstPluginFuncPtr)(unsafe.Pointer(plugin.entryPoint))
	plugin.effect = C.loadEffect(vstEntryPoint)

	C.dispatch(plugin.effect, C.effOpen, 0, 0, nil, 0.0)

	// Set default sample rate and block size
	sampleRate := C.float(44100.0)
	C.dispatch(plugin.effect, C.effSetSampleRate, 0, 0, nil, sampleRate)

	blocksize := C.int(4096)
	C.dispatch(plugin.effect, C.effSetBlockSize, 0, blocksize, nil, 0.0)
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

//export hostCallback
func hostCallback(effect *C.AEffect, opcode int64, index int64, value int64, ptr unsafe.Pointer, opt float64) int {
	switch opcode {
	case C.audioMasterVersion:
		return 2400
	case C.audioMasterIdle:
		C.dispatch(effect, C.effEditIdle, 0, 0, nil, 0)
	default:
		log.Printf("Plugin requested value of opcode %v\n", opcode)
		break
	}
	return 0
}
