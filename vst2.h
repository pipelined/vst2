#include "aeffectx.h"

// Plugin's entry point
typedef AEffect * (*vstPluginFuncPtr)(audioMasterCallback host);

//Bridge function to call entry point on AEffect
AEffect * loadEffect(AEffect * (*load)(audioMasterCallback));

//Bridge to call dispatch function of loaded plugin
VstIntPtr dispatch(AEffect *effect, int opcode, int index, int value, void *ptr, float opt);

//Bridge to call process replacing function of loaded plugin
double** processAudio(AEffect *effect, int numChannels, int blocksize, double ** goInputs);