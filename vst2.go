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

// CanProcessFloat32 checks if plugin can process float32
func (plugin *Plugin) CanProcessFloat32() bool {
	if plugin.effect == nil {
		return false
	}
	return plugin.effect.flags&C.effFlagsCanReplacing == C.effFlagsCanReplacing
}

// CanProcessFloat64 checks if plugin can process float64
func (plugin *Plugin) CanProcessFloat64() bool {
	if plugin.effect == nil {
		return false
	}
	return plugin.effect.flags&C.effFlagsCanDoubleReplacing == C.effFlagsCanDoubleReplacing
}

// ProcessFloat64 audio with VST plugin
func (plugin *Plugin) ProcessFloat64(in [][]float64) (out [][]float64) {
	numChannels := len(in)
	blocksize := len(in[0])

	// convert [][]float64 to []*C.double
	input := make([]*C.double, numChannels)
	output := make([]*C.double, numChannels)
	for i, row := range in {
		// allocate input memory for C layout
		inp := (*C.double)(C.malloc(C.size_t(C.sizeof_double * blocksize)))
		input[i] = inp
		defer C.free(unsafe.Pointer(inp))

		// copy data from slice to C array
		pa := (*[1 << 30]C.double)(unsafe.Pointer(inp))
		for j, v := range row {
			(*pa)[j] = C.double(v)
		}

		// allocate output memory for C layout
		outp := (*C.double)(C.malloc(C.size_t(C.sizeof_double * blocksize)))
		output[i] = outp
		defer C.free(unsafe.Pointer(outp))
	}

	C.processDouble(plugin.effect, C.int(numChannels), C.int(blocksize), &input[0], &output[0])

	//convert []*C.double slices to [][]float64
	out = make([][]float64, numChannels)
	for c, data := range output {
		out[c] = (*[1 << 30]float64)(unsafe.Pointer(data))[:blocksize]
	}
	return out
}

// ProcessFloat32 audio with VST plugin
func (plugin *Plugin) ProcessFloat32(in [][]float32) (out [][]float32) {
	numChannels := len(in)
	blocksize := len(in[0])

	// convert [][]float32 to []*C.float
	input := make([]*C.float, numChannels)
	output := make([]*C.float, numChannels)
	for i, row := range in {
		// allocate input memory for C layout
		inp := (*C.float)(C.malloc(C.size_t(C.sizeof_float * blocksize)))
		input[i] = inp
		defer C.free(unsafe.Pointer(inp))

		// copy data from slice to C array
		pa := (*[1 << 30]C.float)(unsafe.Pointer(inp))
		for j, v := range row {
			(*pa)[j] = C.float(v)
		}

		// allocate output memory for C layout
		outp := (*C.float)(C.malloc(C.size_t(C.sizeof_float * blocksize)))
		output[i] = outp
		defer C.free(unsafe.Pointer(outp))
	}

	C.processFloat(plugin.effect, C.int(numChannels), C.int(blocksize), &input[0], &output[0])

	//convert []*C.float slices to [][]float32
	out = make([][]float32, numChannels)
	for c, data := range output {
		out[c] = (*[1 << 30]float32)(unsafe.Pointer(data))[:blocksize]
	}
	return out
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
		//	log.Printf("AudioMasterVersion")
		return 2400
	case AudioMasterIdle:
		// log.Printf("AudioMasterIdle")
		plugin.Dispatch(EffEditIdle, 0, 0, nil, 0)

	case AudioMasterGetCurrentProcessLevel:
		// log.Printf("AudioMasterGetCurrentProcessLevel")
		return C.kVstProcessLevelUnknown

	default:
		// log.Printf("Plugin requested value of opcode %v\n", opcode)
		break
	}
	return 0
}
