#include <stdint.h>

typedef struct Effect Effect;

typedef	int64_t (*HostCallback) (Effect* effect, int32_t opcode, int32_t index, int64_t value, void* ptr, float opt);
typedef int64_t (*DispatchProc) (Effect* effect, int32_t opcode, int32_t index, int64_t value, void* ptr, float opt);
typedef void (*EffectProcessProc) (Effect* effect, float** inputs, float** outputs, int32_t sampleFrames);
typedef void (*EffectProcessDoubleProc) (Effect* effect, double** inputs, double** outputs, int32_t sampleFrames);
typedef void (*SetParameterProc) (Effect* effect, int32_t index, float parameter);
typedef float (*GetParameterProc) (Effect* effect, int32_t index);

struct Effect
{
	// EffectMagic value.
	int32_t magic;
	// Host to plugin dispatcher function.
	DispatchProc dispatcher;

	// Deprecated.
	EffectProcessProc process;

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
	EffectProcessProc processReplacing;
	// Process audio samples in replacing mode with double precision.
	EffectProcessDoubleProc processDoubleReplacing;

	// Reserved for extension.
	char future[56];
};

// Plugin's entry point
typedef Effect* (*EntryPoint)(HostCallback host);

// Bridge function to call entry point on Effect
Effect* loadEffect(EntryPoint load);

// Bridge to call dispatch function of loaded plugin
int64_t dispatch(Effect *effect, int opcode, int index, int64_t value, void *ptr, float opt);

// Bridge to call process replacing function of loaded plugin
void processDouble(Effect *effect, int numChannels, int blocksize, double **inputs, double **outputs);

// Bridge to call process replacing function of loaded plugin
void processFloat(Effect *effect, int numChannels, int blocksize, float **inputs, float **outputs);

// Bridge to call get parameter fucntion of loaded plugin
float getParameter(Effect *effect, int32_t paramIndex);

// Bridge to call set parameter fucntion of loaded plugin
void setParameter(Effect *effect, int32_t paramIndex, float value);