#include <stdlib.h>
#include <stdio.h>
//#include <stdarg.h>
#include <stdint.h>
#include "aeffectx.h"

// Plugin's entry point
typedef AEffect * (*vstPluginFuncPtr)(audioMasterCallback host);

//Go callback prototype
int hostCallback(AEffect *effect, int opcode, int index, int64_t value, void *ptr, float opt);

//Bridge function to call entry point on AEffect
AEffect * loadEffect(vstPluginFuncPtr load){
	return load((audioMasterCallback)hostCallback);
}

// struct VstSpeakerArrangement* getSpeakerArrangement(int channels) {
// 	struct VstSpeakerArrangement speakerArrangement;
// 	memset(&speakerArrangement, 0, sizeof(struct VstSpeakerArrangement));
//   	speakerArrangement.numChannels = channels;

//   	if (channels <= 8) {
//     	speakerArrangement.numChannels = channels;
//   	} else {
//     	printf("Number of channels = %d. Will only arrange 8 speakers.", channels);
//     	speakerArrangement.numChannels = 8;
//   	}

// 	switch (speakerArrangement.numChannels) {
// 	case 0:
// 	  	speakerArrangement.type = kSpeakerArrEmpty;
// 	  	break;

// 	case 1:
// 	  	speakerArrangement.type = kSpeakerArrMono;
// 	  	break;

// 	case 2:
// 	  	speakerArrangement.type = kSpeakerArrStereo;
// 	  	break;

// 	case 3:
// 	  	speakerArrangement.type = kSpeakerArr30Music;
// 	  	break;

// 	case 4:
// 	  	speakerArrangement.type = kSpeakerArr40Music;
// 	  	break;

// 	case 5:
// 	  	speakerArrangement.type = kSpeakerArr50;
// 	  	break;

// 	case 6:
// 	  	speakerArrangement.type = kSpeakerArr60Music;
// 	  	break;

// 	case 7:
// 	  	speakerArrangement.type = kSpeakerArr70Music;
// 	  	break;

// 	case 8:
// 	  	speakerArrangement.type = kSpeakerArr80Music;
// 	  	break;

// 	default:
// 	  	printf("Cannot arrange more than 8 speakers.");
// 	  	break;
// 	}

// 	for (int i = 0; i < speakerArrangement.numChannels; i++) {
//     	speakerArrangement.speakers[i].azimuth = 0.0f;
//     	speakerArrangement.speakers[i].elevation = 0.0f;
//     	speakerArrangement.speakers[i].radius = 0.0f;
//     	speakerArrangement.speakers[i].reserved = 0.0f;
//     	speakerArrangement.speakers[i].name[0] = '\0';
//     	speakerArrangement.speakers[i].type = kSpeakerUndefined;
//   	}
//   	return &speakerArrangement;
// }

//Bridge to call dispatch function of loaded plugin
VstIntPtr dispatch(AEffect *effect, int opcode, int index, int value, void *ptr, float opt){
	return effect->dispatcher(effect, opcode, index, value, ptr, opt);
}

//Bridge to call process replacing function of loaded plugin
double** processDouble(AEffect *effect, int numChannels, int blocksize, double ** goInputs){

	double** inputs = (double**)malloc(sizeof(double*) * numChannels);
	for(int channel = 0; channel < numChannels; channel++) {
    	inputs[channel] = (double*)&goInputs[channel];
  	}

	double** outputs = (double**)malloc(sizeof(double*) * numChannels);
	for(int channel = 0; channel < numChannels; channel++) {
    	outputs[channel] = (double*)malloc(sizeof(double) * blocksize);
  	}

	effect -> processDoubleReplacing(effect, inputs, outputs, blocksize);
	free(inputs);
	return outputs;
}

//Bridge to call process replacing function of loaded plugin
float** processFloat(AEffect *effect, int numChannels, int blocksize, float ** goInputs){
	float** inputs = (float**)malloc(sizeof(float*) * numChannels);
	for(int channel = 0; channel < numChannels; channel++) {
    	inputs[channel] = (float*)&goInputs[channel];
  	}
	
	float** outputs = (float**)malloc(sizeof(float*) * numChannels);
	for(int channel = 0; channel < numChannels; channel++) {
    	outputs[channel] = (float*)malloc(sizeof(float) * blocksize);
  	}

	effect -> processReplacing(effect, inputs, outputs, blocksize);
	free(inputs);
	return outputs;
}