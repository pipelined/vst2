package vst2

/*
#cgo CFLAGS: -std=gnu99 -I${SRCDIR}/vst2sdk/ -I${SRCDIR}

#include <stdlib.h>
#include <stdint.h>
#include "vst2.h"
#include "aeffectx.h"
*/
import "C"

import (
	"log"
	"path/filepath"
	"unsafe"
)

// Library used to instantiate new instances of plugin
type Library struct {
	entryPoint unsafe.Pointer
	library    uintptr
	Name       string
	Path       string
}

// Plugin type provides interface
type Plugin struct {
	effect   *C.AEffect
	Name     string
	Path     string
	callback HostCallbackFunc
}

// HostCallbackFunc used as callback from plugin
type HostCallbackFunc func(*Plugin, MasterOpcode, int64, int64, unsafe.Pointer, float64) int

const (
	vstMain string = "VSTPluginMain"
)

var (
	plugins     = make(map[*C.AEffect]*Plugin)
	vst2version = 2400
)

// Open loads the library into memory and stores entry point func
//TODO: catch panic
func Open(path string) (*Library, error) {
	fullPath, err := filepath.Abs(path)
	if err != nil {
		log.Printf("Failed to obtain absolute path for '%s': %v\n", path, err)
		return nil, err
	}
	library := &Library{
		Path: fullPath,
	}
	//Get pointer to plugin's Main function
	err = library.load()
	if err != nil {
		log.Printf("Failed to load VST library '%s': %v\n", path, err)
		return nil, err
	}

	return library, nil
}

// Open creates new instance of plugin
func (library *Library) Open() (*Plugin, error) {
	plugin := &Plugin{
		Path:     library.Path,
		Name:     library.Name,
		callback: DefaultHostCallback,
	}
	plugin.effect = C.loadEffect(C.vstPluginFuncPtr(library.entryPoint))
	plugins[plugin.effect] = plugin
	return plugin, nil
}

// Close cleans up C refs for plugin
func (plugin *Plugin) Close() error {
	plugin.Dispatch(EffClose, 0, 0, nil, 0.0)
	delete(plugins, plugin.effect)
	plugin.effect = nil
	return nil
}

// Dispatch wraps-up C method to dispatch calls to plugin
func (plugin *Plugin) Dispatch(opcode PluginOpcode, index int64, value int64, ptr unsafe.Pointer, opt float64) {
	if plugin.effect != nil {
		C.dispatch(plugin.effect, C.int(opcode), C.int(index), C.int64_t(value), ptr, C.float(opt))
	}
}

// Process audio with VST plugin
func (plugin *Plugin) Process(samples [][]float64) (processed [][]float64) {
	//convert Samples to C type
	inSamples := (**C.double)(unsafe.Pointer(&samples[0][0]))
	blocksize := C.int(len(samples[0]))
	numChannels := C.int(len(samples))
	//call plugin and convert result to slice of slices
	outSamples := (*[1 << 30]*C.double)(unsafe.Pointer(C.processDouble(plugin.effect, numChannels, blocksize, inSamples)))[:numChannels]
	//convert slices to [][]float64
	processed = make([][]float64, numChannels)
	for channel, data := range outSamples {
		processed[channel] = (*[1 << 30]float64)(unsafe.Pointer(data))[:blocksize]
	}
	return processed
}

// ProcessFloat audio with VST plugin
func (plugin *Plugin) ProcessFloat(samples [][]float32) (processed [][]float32) {
	//convert Samples to C type
	inSamples := (**C.float)(unsafe.Pointer(&samples[0][0]))
	blocksize := C.int(len(samples[0]))
	numChannels := C.int(len(samples))
	plugin.Dispatch(EffSetBlockSize, 0, int64(len(samples[0])), nil, 0.0)
	//call plugin and convert result to slice of slices
	outSamples := (*[1 << 30]*C.float)(unsafe.Pointer(C.processFloat(plugin.effect, numChannels, blocksize, inSamples)))[:numChannels]
	//convert slices to [][]float64

	processed = make([][]float32, numChannels)
	for channel, data := range outSamples {
		processed[channel] = (*[1 << 30]float32)(unsafe.Pointer(data))[:blocksize]
	}
	return processed
}

//export hostCallback
// calls real callback
func hostCallback(effect *C.AEffect, opcode int64, index int64, value int64, ptr unsafe.Pointer, opt float64) int {
	// AudioMasterVersion is requested when plugin is created
	// It's never in map
	if MasterOpcode(opcode) == AudioMasterVersion {
		return vst2version
	}

	plugin, ok := plugins[effect]
	if !ok {
		panic("Plugin was closed")
	}

	if plugin.callback == nil {
		panic("Host callback is not defined!")
	}
	return plugin.callback(plugin, MasterOpcode(opcode), index, value, ptr, opt)
}

// SetCallback overrides plugin's callback
func (plugin *Plugin) SetCallback(c HostCallbackFunc) {
	if c != nil {
		plugin.callback = c
	}
}

// DefaultHostCallback is a default callback, just prints incoming opcodes should be overriden with SetHostCallback
func DefaultHostCallback(plugin *Plugin, opcode MasterOpcode, index int64, value int64, ptr unsafe.Pointer, opt float64) int {
	switch opcode {
	case AudioMasterVersion:
		log.Printf("AudioMasterVersion")
		return 2400
	case AudioMasterIdle:
		log.Printf("AudioMasterIdle")
		plugin.Dispatch(EffEditIdle, 0, 0, nil, 0)

	case AudioMasterGetCurrentProcessLevel:
		log.Printf("AudioMasterGetCurrentProcessLevel")
		return C.kVstProcessLevelUnknown

	default:
		log.Printf("Plugin requested value of opcode %v\n", opcode)
		break
	}
	return 0
}
