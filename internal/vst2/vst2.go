package vst2

/*
#cgo CFLAGS: -std=gnu99
#cgo CPPFLAGS: -I${SRCDIR}/../../vendor/vst2/

#include <stdlib.h>
#include <stdio.h>
#include <stdarg.h>
#include "aeffectx.h"

// Plugin's entry point
typedef AEffect * (*vstPluginFuncPtr)(audioMasterCallback host);

// Main host callback
VstIntPtr VSTCALLBACK hostCallback(AEffect *effect, VstInt32 opcode, VstInt32 index, VstInt32 value, void *ptr, float opt){
  switch(opcode) {
    case audioMasterVersion:
      return 2400;
    case audioMasterIdle:
      effect->dispatcher(effect, effEditIdle, 0, 0, 0, 0);
    // Handle other opcodes here... there will be lots of them
    default:
      //printf("Plugin requested value of opcode %d\n", opcode);
      break;
  }
}

//Bridge function to call entry point on AEffect
AEffect * loadEffect(AEffect * (*load)(audioMasterCallback), audioMasterCallback hostCallback){
	return load(hostCallback);
}

//Bridge to call dispatch function of loaded plugin
VstIntPtr dispatch(AEffect *effect, VstInt32 opCode, int index, int value, void *ptr, float opt){
	return effect->dispatcher(effect, opCode, index, value, ptr, opt);
}

//Bridge to call process replacing function of loaded plugin
double** processAudio(AEffect *effect, int numChannels, int blocksize, double ** goInputs){

	double** inputs = (double**)malloc(sizeof(double**) * numChannels);
	for(int channel = 0; channel < numChannels; channel++) {
    	inputs[channel] = (double*)&goInputs[channel];
  	}

	double** outputs = (double**)malloc(sizeof(double**) * numChannels);
	for(int channel = 0; channel < numChannels; channel++) {
    	outputs[channel] = (double*)malloc(sizeof(double*) * blocksize);
  	}

	effect -> processDoubleReplacing(effect, inputs, outputs, blocksize);

	return outputs;
}
*/
import "C"

import (
	"log"
	"syscall"
	"unsafe"
)

type Plugin struct {
	Effect *C.AEffect
}

//TODO: catch panic
//Loads the plugin into memory and invokes entry point
func LoadPlugin(path string) (*Plugin, error) {

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

	//Convert to C++ pointer type
	vstEntryPoint := (C.vstPluginFuncPtr)(unsafe.Pointer(mainEntryPoint))
	//Convert callback function to C++ type
	callback := (C.audioMasterCallback)(C.hostCallback)

	return &Plugin{Effect: C.loadEffect(vstEntryPoint, callback)}, nil
}

//Resumes the plugin
func (plugin *Plugin) resume() {
	C.dispatch(plugin.Effect, C.effMainsChanged, 0, 1, nil, 0.0)
}

//Suspends the plugin
func (plugin *Plugin) suspend() {
	C.dispatch(plugin.Effect, C.effMainsChanged, 0, 0, nil, 0.0)
}

//Starts the plugin
func (plugin *Plugin) start() {
	C.dispatch(plugin.Effect, C.effOpen, 0, 0, nil, 0.0)

	// Set default sample rate and block size
	sampleRate := C.float(44100.0)
	C.dispatch(plugin.Effect, C.effSetSampleRate, 0, 0, nil, sampleRate)

	blocksize := C.int(4096)
	C.dispatch(plugin.Effect, C.effSetBlockSize, 0, blocksize, nil, 0.0)
}

//Process audio with VST plugin
func (plugin *Plugin) Process(samples [][]float64) (processed [][]float64) {
	//convert Samples to float **
	inSamples := (**C.double)(unsafe.Pointer(&samples[0][0]))
	blocksize := C.int(len(samples[0]))
	numChannels := C.int(len(samples))
	//call plugin and convert result to slice of slices
	outSamples := (*[1 << 30]*C.double)(unsafe.Pointer(C.processAudio(plugin.Effect, numChannels, blocksize, inSamples)))[:numChannels]
	//convert slices to [][]float64
	processed = make([][]float64, numChannels)
	for channel, data := range outSamples {
		processed[channel] = (*[1 << 30]float64)(unsafe.Pointer(data))[:blocksize]
	}
	return processed
}
