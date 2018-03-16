#include <stdlib.h>
#include <stdio.h>
//#include <stdarg.h>
#include <stdint.h>
#include "aeffectx.h"

// Plugin's entry point
typedef AEffect * (*vstPluginFuncPtr)(audioMasterCallback host);

//Go callback prototype
int hostCallback(AEffect *effect, int opcode, int index, int64_t value, void *ptr, float opt);

//Bridge function to call entry point on AEffect
AEffect * loadEffect(vstPluginFuncPtr load){
	return load((audioMasterCallback)hostCallback);
}

// Bridge to call dispatch function of loaded plugin
VstIntPtr dispatch(AEffect *effect, int opcode, int index, VstIntPtr value, void *ptr, float opt){
	return effect->dispatcher(effect, opcode, index, value, ptr, opt);
}

// Bridge to call process replacing function of loaded plugin
void processDouble(AEffect *effect, int numChannels, int blocksize, double ** inputs, double ** outputs){
	effect -> processDoubleReplacing(effect, inputs, outputs, blocksize);
}

// Bridge to call process replacing function of loaded plugin
void processFloat(AEffect *effect, int numChannels, int blocksize, float **inputs, float **outputs){
	effect -> processReplacing(effect, inputs, outputs, blocksize);
}