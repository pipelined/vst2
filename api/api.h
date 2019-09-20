#include <stdint.h>

typedef struct Effect Effect;

typedef	int64_t (*audioMasterCallback) (Effect* effect, int32_t opcode, int32_t index, int64_t value, void* ptr, float opt);
typedef int64_t (*EffectDispatcherProc) (Effect* effect, int32_t opcode, int32_t index, int64_t value, void* ptr, float opt);
typedef void (*EffectProcessProc) (Effect* effect, float** inputs, float** outputs, int32_t sampleFrames);
typedef void (*EffectProcessDoubleProc) (Effect* effect, double** inputs, double** outputs, int32_t sampleFrames);
typedef void (*EffectSetParameterProc) (Effect* effect, int32_t index, float parameter);
typedef float (*EffectGetParameterProc) (Effect* effect, int32_t index);

struct Effect
{
	int32_t magic;			///< must be #kEffectMagic ('VstP')
	EffectDispatcherProc dispatcher;

	// Deprecated.
	EffectProcessProc process;

	EffectSetParameterProc setParameter;
	EffectGetParameterProc getParameter;

	int32_t numPrograms;   ///< number of programs
	int32_t numParams;		///< all programs are assumed to have numParams parameters
	int32_t numInputs;		///< number of audio inputs
	int32_t numOutputs;	///< number of audio outputs

	int32_t flags;			///< @see VstEffectFlags

	int64_t resvd1;		///< reserved for Host, must be 0
	int64_t resvd2;		///< reserved for Host, must be 0

	int32_t initialDelay;	///< for algorithms which need input in the first place (Group delay or latency in Samples). This value should be initialized in a resume state.

	// Deprecated.
	int32_t realQualities;	///< \deprecated unused member
	// Deprecated.
	int32_t offQualities;		///< \deprecated unused member
	// Deprecated.
	float ioRatio;			///< \deprecated unused member

	void* object;			///< #AudioEffect class pointer
	void* user;				///< user-defined pointer

	int32_t uniqueID;		///< registered unique identifier (register it at Steinberg 3rd party support Web). This is used to identify a plug-in during save+load of preset and project.
	int32_t version;		///< plug-in version (example 1100 for version 1.1.0.0)


	EffectProcessProc processReplacing;
	EffectProcessDoubleProc processDoubleReplacing;

	char future[56];		///< reserved for future use (please zero)
};

// Plugin's entry point
typedef Effect * (*vstPluginFuncPtr)(audioMasterCallback host);

// Bridge function to call entry point on Effect
Effect * loadEffect(vstPluginFuncPtr load);

// Bridge to call dispatch function of loaded plugin
int64_t dispatch(Effect *effect, int opcode, int index, int64_t value, void *ptr, float opt);

// Bridge to call process replacing function of loaded plugin
void processDouble(Effect *effect, int numChannels, int blocksize, double **inputs, double **outputs);

// Bridge to call process replacing function of loaded plugin
void processFloat(Effect *effect, int numChannels, int blocksize, float **inputs, float **outputs);