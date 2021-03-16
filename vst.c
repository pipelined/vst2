#include <stdlib.h>
#include <stdio.h>
#include <stdint.h>
#include "vst.h"

//Go callback prototype
int64_t hostCallback(CPlugin *plugin, int32_t opcode, int32_t index, int64_t value, void *ptr, float opt);

//Bridge function to call entry point on Effect
CPlugin * loadPluginBridge(EntryPoint load){
	return load((HostCallback)hostCallback);
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

typedef struct Event
{
	int32_t type;
	int32_t byteSize;
	int32_t deltaFrames;	///< sample frames related to the current block start sample position
	int32_t flags;			///< none defined yet (should be zero)
	int32_t dumpBytes;		///< byte size of sysexDump
	int64_t resvd1;		///< zero (Reserved for future use)
	char* sysexDump;		///< sysex dump
	int64_t resvd2;		///< zero (Reserved for future use)
} Event;

Events* newEvents(int32_t numEvents) {
	Events *e = malloc(sizeof(*e));
	e->numEvents = numEvents;
	e->events = malloc(sizeof(e->events) * numEvents);
	return e;
}

void setEvent(Events *events, void *event, int32_t pos) {
	events->events[pos] = event;
}

void* getEvent(Events *events, int32_t pos) {
	return events->events[pos];
}