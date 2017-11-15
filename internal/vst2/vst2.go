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
AEffect * createEffectInstance(AEffect * (*load)(audioMasterCallback), audioMasterCallback host){
	return load(host);
}

//Bridge to call dispatch function of loaded plugin
VstIntPtr dispatch(AEffect *effect, VstInt32 opCode, int index, int value, void *ptr, float opt){
	return effect->dispatcher(effect, opCode, index, value, ptr, opt);
}

//Bridge to call process replacing function of loaded plugin
double** processAudio(AEffect *effect, int numChannels, int blocksize, double *leftInputs, double *rightInputs){

	printf("in params: numChannels[%d] blocksize[%d]\n", numChannels, blocksize);

	double* inputs[numChannels];
	inputs[0] = &leftInputs[0];
	inputs[1] = &rightInputs[0];
	printf("inputs initialized!\n");

	double** outputs = (double**)malloc(sizeof(double**) * numChannels);
	for(int channel = 0; channel < numChannels; channel++) {
    	outputs[channel] = (double*)malloc(sizeof(double*) * blocksize);
  	}
  	printf("outputs initialized!\n");

	for (int j = 0; j < 10; j++) {
			printf("[%.6f][%.6f]\n", leftInputs[j], rightInputs[j]);
    }

	effect -> processDoubleReplacing(effect, &inputs, &outputs, blocksize);
	return outputs;
}
*/
import "C"

import (
	"fmt"
	"github.com/youpy/go-wav"
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
	modulePtr, err := syscall.LoadDLL(path)
	if err != nil {
		log.Printf("Failed trying to load VST from '%s'. %v\n", path, err)
		return nil, err
	}

	//Get pointer to plugin's Main function
	mainEntryPoint, err := syscall.GetProcAddress(modulePtr.Handle, "VSTPluginMain")
	if err != nil {
		log.Printf("Failed trying to obtain VST entry point '%s'. %v\n", path, err)
		return nil, err
	}

	//Convert to C++ pointer type
	vstEntryPoint := (C.vstPluginFuncPtr)(unsafe.Pointer(mainEntryPoint))

	//Convert callback function to C++ type
	callback := (C.audioMasterCallback)(C.hostCallback)

	return &Plugin{Effect: C.createEffectInstance(vstEntryPoint, callback)}, nil
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

	blocksize := C.int(512)
	C.dispatch(plugin.Effect, C.effSetBlockSize, 0, blocksize, nil, 0.0)
}

//TODO: catch panic
//process audio
func (plugin *Plugin) processAudio(wavSamples []wav.Sample) {

	samples := ConvertWavSamplesToFloat64(wavSamples)
	//convert Samples to float **

	rightCSamples := (*C.double)(unsafe.Pointer(&samples[0][0]))
	leftCSamples := (*C.double)(unsafe.Pointer(&samples[1][0]))

	//fmt.Printf("\nprocess samples type: %T\n", inSamples)
	blocksize := C.int(len(samples[0]))
	numChannels := C.int(len(samples))
	//fmt.Printf("\nblocksize: %v numchannels: %v\n", blocksize, numChannels)

	C.processAudio(plugin.Effect, numChannels, blocksize, leftCSamples, rightCSamples)
}

func ConvertWavSamplesToCDouble(wavSamples []wav.Sample) (samples [][]C.double) {
	samples = make([][]C.double, 2)
	//convert samples to float slice
	samples[0] = make([]C.double, len(wavSamples))
	samples[1] = make([]C.double, len(wavSamples))

	for i, sample := range wavSamples {
		samples[0][i] = C.double(float64(sample.Values[0]) / 32768)
		samples[1][i] = C.double(float64(sample.Values[1]) / 32768)
	}
	return samples
}

func ConvertWavSamplesToFloat64(wavSamples []wav.Sample) (samples [][]float64) {
	fmt.Println("Call from convert to float64")
	samples = make([][]float64, 2)
	//convert samples to float slice
	samples[0] = make([]float64, len(wavSamples))
	samples[1] = make([]float64, len(wavSamples))

	for i, sample := range wavSamples {
		samples[0][i] = float64(sample.Values[0]) / 32768
		samples[1][i] = float64(sample.Values[1]) / 32768
	}
	return samples
}
