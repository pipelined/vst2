#ifndef VST_H
#define VST_H
#include <stdint.h>

typedef struct CPlugin CPlugin;
typedef struct Events Events;

typedef	int64_t (*HostCallback) (CPlugin* plugin, int32_t opcode, int32_t index, int64_t value, void* ptr, float opt);
typedef int64_t (*DispatchProc) (CPlugin* plugin, int32_t opcode, int32_t index, int64_t value, void* ptr, float opt);
typedef void (*ProcessFloatFunc) (CPlugin* plugin, float** inputs, float** outputs, int32_t sampleFrames);
typedef void (*ProcessDoubleFunc) (CPlugin* plugin, double** inputs, double** outputs, int32_t sampleFrames);
typedef void (*SetParameterProc) (CPlugin* plugin, int32_t index, float parameter);
typedef float (*GetParameterProc) (CPlugin* plugin, int32_t index);

struct CPlugin
{
	// EffectMagic value.
	int32_t magic;
	// Host to plugin dispatcher function.
	DispatchProc dispatcher;

	// Deprecated.
	ProcessFloatFunc process;

	// Set new value of automatable parameter.
	SetParameterProc setParameter;
	// Returns current value of automatable parameter.
	GetParameterProc getParameter;

	// Number of presets.
	int32_t numPrograms;
	// Number of parameters per preset.
	int32_t numParams;
	// Number of audio inputs.
	int32_t numInputs;
	// Number of audio outputs.
	int32_t numOutputs;

	// EffectFlags values.
	int32_t flags;

	// Reserved for Host, must be 0.
	int64_t resvd1;
	// Reserved for Host, must be 0.
	int64_t resvd2;

	// InitialDelay is for algorithms which need input in the first place.
	// This value should be initialized in a resume state.
	int32_t initialDelay;

	// Deprecated.
	int32_t realQualities;
	// Deprecated.
	int32_t offQualities;
	// Deprecated.
	float ioRatio;

	// Internal class pointer.
	void* object;
	// User-defined pointer.
	void* user;

	// Registered unique identifier (register it at Steinberg 3rd party support Web).
	// This is used to identify a plug-in during save+load of preset and project.
	int32_t uniqueID;
	// Version of plugin (example 1100 for version 1.1.0.0).
	int32_t version;

	// Process audio samples in replacing mode with single precision.
	ProcessFloatFunc processFloat;
	// Process audio samples in replacing mode with double precision.
	ProcessDoubleFunc processDouble;

	// Reserved for extension.
	char future[56];
};

// CPlugin's entry point
typedef CPlugin* (*EntryPoint)(HostCallback host);

struct Events
{
	// Number of Events in array.
	int32_t numEvents;
	// Not used.
	int64_t reserved;
	// Event pointer array, variable size.
	void** events;
};

// Bridge to allocate events structure.
Events* newEvents(int32_t numEvents);

// sets event into events array. This function is needed because there is
// no way to assign values to void** from Go.
void setEvent(Events *events, void *event, int32_t pos);

// gets the event from events container. This function is needed because there is
// no way to assign values to void** from Go.
void *getEvent(Events *events, int32_t pos);
#endif