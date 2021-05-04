#include <stdlib.h>
#include "include/vst.h"

#ifdef _WIN32
#define VSTAPI __declspec(dllexport)
#else
#define VSTAPI __attribute__((visibility("default")))
#endif

void newGoPlugin(CPlugin *plugin, HostCallback c);

//Go dispatch prototype
int64_t dispatchPluginBridge(CPlugin *plugin, int32_t opcode, int32_t index, int64_t value, void *ptr, float opt);

//Go processDouble prototype
void processDoublePluginBridge(CPlugin *plugin, double ** inputs, double ** outputs, int32_t sampleFrames);

//Go processFloat prototype
void processFloatPluginBridge(CPlugin *plugin, float **inputs, float **outputs, int32_t sampleFrames);

//Go getParameter prototype
float getParameterPluginBridge(CPlugin *plugin, int32_t paramIndex);

//Go setParameter prototype
void setParameterPluginBridge(CPlugin *plugin, int32_t paramIndex, float value);

VSTAPI CPlugin* VSTPluginMain(HostCallback c) {
    CPlugin *p = malloc(sizeof(CPlugin));
    p->dispatcher = dispatchPluginBridge;
    p->getParameter = getParameterPluginBridge;
    p->setParameter = setParameterPluginBridge;
    p->processDouble = processDoublePluginBridge;
    p->processFloat = processFloatPluginBridge;
    newGoPlugin(p, c);
    return p;
}

int64_t callbackHost(HostCallback c, CPlugin* plugin, int32_t opcode, int32_t index, int64_t value, void* ptr, float opt) {
    return c(plugin, opcode, index, value, ptr, opt);
};
