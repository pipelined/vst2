#include <stdlib.h>
#include "plugin.h"

//Go dispatch prototype
int64_t dispatch(Plugin *plugin, int32_t opcode, int32_t index, int64_t value, void *ptr, float opt);

//Go processDouble prototype
void processDouble(Plugin *plugin, double ** inputs, double ** outputs, int32_t sampleFrames);

//Go processFloat prototype
void processFloat(Plugin *plugin, float **inputs, float **outputs, int32_t sampleFrames);

//Go getParameter prototype
float getParameter(Plugin *plugin, int32_t paramIndex);

//Go setParameter prototype
void setParameter(Plugin *plugin, int32_t paramIndex, float value);

Plugin* VSTPluginMain(HostCallback c) {
    // TODO: init values from go plugin
    Plugin *p = malloc(sizeof(Plugin));
    p->dispatcher = dispatch;
    p->getParameter = getParameter;
    p->setParameter = setParameter;
    p->processDouble = processDouble;
    p->processFloat = processFloat;
    return p;
}

// // Bridge to call dispatch function of loaded plugin
// int64_t dispatchPluginBridge(Plugin *plugin, int32_t opcode, int32_t index, int64_t value, void *ptr, float opt){
// 	return dispatch(plugin, opcode, index, value, ptr, opt);
// }

// // Bridge to call process replacing function of loaded plugin
// void processDoublePluginBridge(Plugin *plugin, double ** inputs, double ** outputs, int32_t sampleFrames){
// 	processDouble(plugin, inputs, outputs, sampleFrames);
// }


// // Bridge to call process replacing function of loaded plugin
// void processFloatPluginBridge(Plugin *plugin, float **inputs, float **outputs, int32_t sampleFrames){
// 	processFloat(plugin, inputs, outputs, sampleFrames);
// }


// // Bridge to call get parameter fucntion of loaded plugin
// float getParameterPluginBridge(Plugin *plugin, int32_t paramIndex) {
// 	getParameter(plugin, paramIndex);
// }


// // Bridge to call set parameter fucntion of loaded plugin
// void setParameterPluginBridge(Plugin *plugin, int32_t paramIndex, float value) {
// 	setParameter(plugin, paramIndex, value);
// }