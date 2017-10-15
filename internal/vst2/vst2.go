package vst2

//TODO: add exceptions handling

/*
#cgo CPPFLAGS:  -I${SRCDIR}/../../vendor/vst2/
#include "aeffectx.h"
#include "stdlib.h"

// Main host callback
VstIntPtr VSTCALLBACK hostCallback(AEffect *effect, VstInt32 opcode, VstInt32 index, VstInt32 value, void *ptr, float opt){
  switch(opcode) {
    case audioMasterVersion:
      return 2400;
    case audioMasterIdle:
      effect->dispatcher(effect, effEditIdle, 0, 0, 0, 0);
    // Handle other opcodes here... there will be lots of them
    default:
      //printf("Plugin requested value of opcode %d\n", opcode);
      break;
  }
}

// Plugin's entry point
typedef AEffect * (*vstPluginFuncPtr)(audioMasterCallback host);
// Plugin's dispatcher function
typedef VstIntPtr (*dispatcherFuncPtr)(AEffect *effect, VstInt32 opCode, VstInt32 index, VstInt32 value, void *ptr, float opt);
// Plugin's getParameter() method
typedef float (*getParameterFuncPtr)(AEffect *effect, VstInt32 index);
// Plugin's setParameter() method
typedef void (*setParameterFuncPtr)(AEffect *effect, VstInt32 index, float value);
// Plugin's processEvents() method
typedef void (*processFuncPtr)(AEffect *effect, float **inputs,  float **outputs, VstInt32 sampleFrames);

//Bridge function to call entry point on AEffect
AEffect * createEffectInstance(AEffect * (*load)(audioMasterCallback), audioMasterCallback host){
	return load(host);
}

//Bridge to call dispatch function of loaded plugin
VstIntPtr dispatch(AEffect *effect, VstInt32 opCode, VstInt32 index, VstInt32 value, void *ptr, float opt){
	return effect->dispatcher(effect, opCode, index, value, ptr, opt);
}

*/
import "C"

import (
	"fmt"
	"syscall"
	"unsafe"
)

type Plugin struct {
	Effect *C.AEffect
}

//Loads the plugin into memory and invokes entry point
func LoadPlugin(path string) (plugin *Plugin) {

	//Load plugin by path
	modulePtr, err := syscall.LoadLibrary(path)
	if err != nil {
		fmt.Printf("Failed trying to load VST from '%s', error %d\n", path, err.Error())
		return nil
	}

	//Get pointer to plugin's Main function
	mainEntryPoint, err := syscall.GetProcAddress(modulePtr, "VSTPluginMain")
	if err != nil {
		fmt.Printf("Failed trying to obtain VST entry point '%s', error %d\n", path, err.Error())
		return nil
	}

	//Convert to C++ pointer type
	vstEntryPoint := (C.vstPluginFuncPtr)(unsafe.Pointer(mainEntryPoint))

	//Convert callback function to C++ type
	callback := (C.audioMasterCallback)(C.hostCallback)
	return &Plugin{C.createEffectInstance(vstEntryPoint, callback)}
}

func configurePluginCallbacks(plugin *C.AEffect) {
	// Check plugin's magic number
	// If incorrect, then the file either was not loaded properly, is not a
	// real VST plugin, or is otherwise corrupt.
	if plugin.magic != C.kEffectMagic {
		fmt.Printf("Plugin's magic number is bad\n")
		// return -1
	}

	fmt.Printf("Plugin name is %v", plugin.uniqueID)

	// Set up plugin callback functions
	plugin.getParameter = plugin.getParameter
	plugin.setParameter = plugin.setParameter

	// return plugin
}

/*func resumePlugin(plugin *C.AEffect) {
	dispatcher(plugin, C.effMainsChanged, 0, 1, nil, 0.0)
}

func suspendPlugin(plugin *C.AEffect, dispatcher C.dispatcherFuncPtr) {
	dispatcher(plugin, C.effMainsChanged, 0, 0, nil, 0.0)
}*/

//Starts the plugin
func (plugin *Plugin) start() {
	C.dispatch(plugin.Effect, C.effOpen, 0, 0, nil, 0.0)

	// Set default sample rate and block size
	sampleRate := C.float(44100.0)
	C.dispatch(plugin.Effect, C.effSetSampleRate, 0, 0, nil, sampleRate)

	blocksize := C.VstInt32(512)
	C.dispatch(plugin.Effect, C.effSetBlockSize, 0, blocksize, nil, 0.0)
}
