#include <stdlib.h>
#include <stdio.h>
//#include <stdarg.h>
#include <stdint.h>
#include "effect.h"

//Go callback prototype
int hostCallback(Effect *effect, int opcode, int index, int64_t value, void *ptr, float opt);

//Bridge function to call entry point on Effect
Effect * loadEffect(EntryPoint load){
	return load((HostCallback)hostCallback);
}

// Bridge to call dispatch function of loaded plugin
int64_t dispatch(Effect *effect, int opcode, int index, int64_t value, void *ptr, float opt){
	return effect->dispatcher(effect, opcode, index, value, ptr, opt);
}

// Bridge to call process replacing function of loaded plugin
void processDouble(Effect *effect, int numChannels, int blocksize, double ** inputs, double ** outputs){
	effect -> processDoubleReplacing(effect, inputs, outputs, blocksize);
}

// Bridge to call process replacing function of loaded plugin
void processFloat(Effect *effect, int numChannels, int blocksize, float **inputs, float **outputs){
	effect -> processReplacing(effect, inputs, outputs, blocksize);
}