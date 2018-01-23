#include "aeffectx.h"

// Plugin's entry point
typedef AEffect * (*vstPluginFuncPtr)(audioMasterCallback host);

//Bridge function to call entry point on AEffect
AEffect * loadEffect(AEffect * (*load)(audioMasterCallback));

//Bridge to call dispatch function of loaded plugin
VstIntPtr dispatch(AEffect *effect, int opcode, int index, int64_t value, void *ptr, float opt);

//Bridge to call process replacing function of loaded plugin
double** processDouble(AEffect *effect, int numChannels, int blocksize, double ** goInputs);

//Bridge to call process replacing function of loaded plugin
float** processFloat(AEffect *effect, int numChannels, int blocksize, float ** goInputs);