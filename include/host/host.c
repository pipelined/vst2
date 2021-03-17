#include <stdlib.h>
#include "include/vst.h"

//Go callback prototype
int64_t hostCallbackBridge(CPlugin *plugin, int32_t opcode, int32_t index, int64_t value, void *ptr, float opt);

//Bridge function to call entry point on Effect
CPlugin * loadPluginHostBridge(EntryPoint load){
	return load((HostCallback)hostCallbackBridge);
}

// Bridge to call dispatch function of loaded plugin
int64_t dispatchHostBridge(CPlugin *plugin, int32_t opcode, int32_t index, int64_t value, void *ptr, float opt){
	return plugin->dispatcher(plugin, opcode, index, value, ptr, opt);
}

// Bridge to call process replacing function of loaded plugin
void processDoubleHostBridge(CPlugin *plugin, double ** inputs, double ** outputs, int32_t sampleFrames){
	plugin -> processDouble(plugin, inputs, outputs, sampleFrames);
}

// Bridge to call process replacing function of loaded plugin
void processFloatHostBridge(CPlugin *plugin, float **inputs, float **outputs, int32_t sampleFrames){
	plugin -> processFloat(plugin, inputs, outputs, sampleFrames);
}

// Bridge to call get parameter fucntion of loaded plugin
float getParameterHostBridge(CPlugin *plugin, int32_t paramIndex) {
	return plugin->getParameter(plugin, paramIndex);
}

// Bridge to call set parameter fucntion of loaded plugin
void setParameterHostBridge(CPlugin *plugin, int32_t paramIndex, float value) {
	plugin->setParameter(plugin, paramIndex, value);
}