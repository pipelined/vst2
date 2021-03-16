#include "vst.h"

CPlugin* VSTPluginMain(HostCallback c);

// Bridge to call dispatch function of loaded plugin
int64_t dispatchPluginBridge(CPlugin* plugin, int32_t opcode, int32_t index, int64_t value, void *ptr, float opt);

// Bridge to call process replacing function of loaded plugin
void processDoublePluginBridge(CPlugin* plugin, double **inputs, double **outputs, int32_t sampleFrames);

// Bridge to call process replacing function of loaded plugin
void processFloatPluginBridge(CPlugin* plugin, float **inputs, float **outputs, int32_t sampleFrames);

// Bridge to call get parameter fucntion of loaded plugin
float getParameterPluginBridge(CPlugin* plugin, int32_t paramIndex);

// Bridge to call set parameter fucntion of loaded plugin
void setParameterPluginBridge(CPlugin* plugin, int32_t paramIndex, float value);