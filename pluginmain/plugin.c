#include <stdlib.h>
#include "plugin.h"

//Go dispatch prototype
int64_t dispatch(CPlugin *plugin, int32_t opcode, int32_t index, int64_t value, void *ptr, float opt);

//Go processDouble prototype
void processDouble(CPlugin *plugin, double ** inputs, double ** outputs, int32_t sampleFrames);

//Go processFloat prototype
void processFloat(CPlugin *plugin, float **inputs, float **outputs, int32_t sampleFrames);

//Go getParameter prototype
float getParameter(CPlugin *plugin, int32_t paramIndex);

//Go setParameter prototype
void setParameter(CPlugin *plugin, int32_t paramIndex, float value);

CPlugin* VSTPluginMain(HostCallback c) {
    // TODO: init values from go plugin
    CPlugin *p = malloc(sizeof(CPlugin));
    p->dispatcher = dispatch;
    p->getParameter = getParameter;
    p->setParameter = setParameter;
    p->processDouble = processDouble;
    p->processFloat = processFloat;
    return p;
}
