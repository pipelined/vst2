#include <stdlib.h>
#include <stdio.h>
#include <stdint.h>
#include "sdk.h"

//Go callback prototype
int hostCallback(Plugin *plugin, int32_t opcode, int32_t index, int64_t value, void *ptr, float opt);

//Bridge function to call entry point on Effect
Plugin * loadPlugin(EntryPoint load){
	return load((HostCallback)hostCallback);
}

// Bridge to call dispatch function of loaded plugin
int64_t dispatch(Plugin *plugin, int32_t opcode, int32_t index, int64_t value, void *ptr, float opt){
	return plugin->dispatcher(plugin, opcode, index, value, ptr, opt);
}

// Bridge to call process replacing function of loaded plugin
void processDouble(Plugin *plugin, int32_t numChannels, int32_t blocksize, double ** inputs, double ** outputs){
	plugin -> processDoubleReplacing(plugin, inputs, outputs, blocksize);
}

// Bridge to call process replacing function of loaded plugin
void processFloat(Plugin *plugin, int32_t numChannels, int32_t blocksize, float **inputs, float **outputs){
	plugin -> processReplacing(plugin, inputs, outputs, blocksize);
}

// Bridge to call get parameter fucntion of loaded plugin
float getParameter(Plugin *plugin, int32_t paramIndex) {
	return plugin->getParameter(plugin, paramIndex);
}

// Bridge to call set parameter fucntion of loaded plugin
void setParameter(Plugin *plugin, int32_t paramIndex, float value) {
	plugin->setParameter(plugin, paramIndex, value);
}