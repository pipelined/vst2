package vst2

/* 
#cgo CFLAGS: -std=gnu99 -I${SRCDIR}/../vendor/vst2/
#cgo LDFLAGS: -Wl,--allow-multiple-definition

#include <stdlib.h>
#include <stdio.h>
#include <stdarg.h>
#include "aeffectx.h"

// Plugin's entry point
typedef AEffect * (*vstPluginFuncPtr)(audioMasterCallback host);

//Go callback prototype
int hostCallback(AEffect *effect, int opcode, int index, int value, void *ptr, float opt);

//Bridge function to call entry point on AEffect
AEffect * loadEffect(AEffect * (*load)(audioMasterCallback)){
	return load((audioMasterCallback)hostCallback);
}

//Bridge to call dispatch function of loaded plugin
VstIntPtr dispatch(AEffect *effect, int opCode, int index, int value, void *ptr, float opt){
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

//Go callback prototype
int hostCallback(AEffect *effect, int opcode, int index, int value, void *ptr, float opt);
*/
import "C"