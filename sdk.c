#include <stdlib.h>
#include <stdio.h>
#include <stdint.h>
#include "sdk.h"

//Go callback prototype
int hostCallback(Effect *effect, int32_t opcode, int32_t index, int64_t value, void *ptr, float opt);

//Bridge function to call entry point on Effect
Effect * loadEffect(EntryPoint load){
	return load((HostCallback)hostCallback);
}

// Bridge to call dispatch function of loaded plugin
int64_t dispatch(Effect *effect, int32_t opcode, int32_t index, int64_t value, void *ptr, float opt){
	return effect->dispatcher(effect, opcode, index, value, ptr, opt);
}

// Bridge to call process replacing function of loaded plugin
void processDouble(Effect *effect, int32_t numChannels, int32_t blocksize, double ** inputs, double ** outputs){
	effect -> processDoubleReplacing(effect, inputs, outputs, blocksize);
}

// Bridge to call process replacing function of loaded plugin
void processFloat(Effect *effect, int32_t numChannels, int32_t blocksize, float **inputs, float **outputs){
	effect -> processReplacing(effect, inputs, outputs, blocksize);
}

// Bridge to call get parameter fucntion of loaded plugin
float getParameter(Effect *effect, int32_t paramIndex) {
	return effect->getParameter(effect, paramIndex);
}

// Bridge to call set parameter fucntion of loaded plugin
void setParameter(Effect *effect, int32_t paramIndex, float value) {
	effect->setParameter(effect, paramIndex, value);
}