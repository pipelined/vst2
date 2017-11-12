package vst2

//TODO: add exceptions handling

/*
#cgo CFLAGS: -std=gnu99
#cgo CPPFLAGS:  -I${SRCDIR}/../../vendor/vst2/
#include <stdlib.h>
#include <stdio.h>
#include "aeffectx.h"

// Plugin's entry point
typedef AEffect * (*vstPluginFuncPtr)(audioMasterCallback host);
// Plugin's dispatcher function
typedef VstIntPtr (*dispatcherFuncPtr)(AEffect *effect, VstInt32 opCode, VstInt32 index, VstInt32 value, void *ptr, float opt);
// Plugin's getParameter() method
typedef float (*getParameterFuncPtr)(AEffect *effect, VstInt32 index);
// Plugin's setParameter() method
typedef void (*setParameterFuncPtr)(AEffect *effect, VstInt32 index, float value);
// Plugin's processEvents() method
typedef void (*processFuncPtr)(AEffect *effect, float **inputs,  float **outputs, VstInt32 sampleFrames);

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
AEffect * createEffectInstance(AEffect * (*load)(audioMasterCallback), audioMasterCallback host){
	return load(host);
}

//Bridge to call dispatch function of loaded plugin
VstIntPtr dispatch(AEffect *effect, VstInt32 opCode, VstInt32 index, VstInt32 value, void *ptr, float opt){
	return effect->dispatcher(effect, opCode, index, value, ptr, opt);
}

//Bridge to call process replacing function of loaded plugin
double** processAudio(AEffect *effect, double *leftInputs, double *rightInputs, int blocksize){
	double** inputs;
	double** outputs;
	for(int channel = 0; channel < 2; channel++) {
    	outputs[channel] = (double*)malloc(sizeof(double*) * blocksize);
  	}

	effect -> processDoubleReplacing(effect, inputs, outputs, blocksize);
	return outputs;
}
*/
import "C"

import (
	"fmt"
	"github.com/youpy/go-wav"
	"syscall"
	"unsafe"
)

type Plugin struct {
	Effect *C.AEffect
}

//Loads the plugin into memory and invokes entry point
func LoadPlugin(path string) (plugin *Plugin) {

	//Load plugin by path
	modulePtr, err := syscall.LoadDLL(path)
	if err != nil {
		fmt.Printf("Failed trying to load VST from '%s', error %v\n", path, err.Error())
		return nil
	}

	//Get pointer to plugin's Main function
	mainEntryPoint, err := syscall.GetProcAddress(modulePtr.Handle, "VSTPluginMain")
	if err != nil {
		fmt.Printf("Failed trying to obtain VST entry point '%s', error %d\n", path, err.Error())
		return nil
	}

	//Convert to C++ pointer type
	vstEntryPoint := (C.vstPluginFuncPtr)(unsafe.Pointer(mainEntryPoint))

	//Convert callback function to C++ type
	callback := (C.audioMasterCallback)(C.hostCallback)
	return &Plugin{C.createEffectInstance(vstEntryPoint, callback)}
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

	blocksize := C.VstInt32(512)
	C.dispatch(plugin.Effect, C.effSetBlockSize, 0, blocksize, nil, 0.0)
}

func (plugin *Plugin) processAudio(samples []wav.Sample) {
	//convert Samples to float **
	var rightSamples []C.double
	var leftSamples []C.double

	for _, sample := range samples {
		rightSamples = append(rightSamples, C.double(float64(sample.Values[0])/32768))
		leftSamples = append(leftSamples, C.double(float64(sample.Values[1])/32768))
	}

	rightCSamples := (*C.double)(unsafe.Pointer(&rightSamples[0]))
	leftCSamples := (*C.double)(unsafe.Pointer(&rightSamples[0]))
	//fmt.Printf("\nprocess samples type: %T\n", inSamples)
	blocksize := C.int(len(samples))
	//fmt.Printf("\nblocksize type: %T value: %v\n", blocksize, blocksize)

	C.processAudio(plugin.Effect, leftCSamples, rightCSamples, blocksize)
}
