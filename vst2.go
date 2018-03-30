package vst2

/*
#cgo CFLAGS: -std=gnu99 -I${SRCDIR}/vst2sdk/ -I${SRCDIR}

#include <stdlib.h>
#include <stdint.h>
#include "vst2.h"
#include "aeffect.h"
#include "aeffectx.h"
*/
import "C"

import (
	"log"
	"path/filepath"
	"sync"
	"unsafe"
)

// Plugin type provides interface
type Plugin struct {
	effect   *C.AEffect
	Name     string
	Path     string
	callback HostCallbackFunc
	timeInfo *vstTimeInfo
}

// speakerArrangement is a wrapper over vst2 VstSpeakerArrangement structure
type speakerArrangement C.struct_VstSpeakerArrangement

type vstTimeInfo C.struct_VstTimeInfo

// HostCallbackFunc used as callback from plugin
type HostCallbackFunc func(*Plugin, MasterOpcode, int64, int64, unsafe.Pointer, float64) int

const (
	vstMain string = "VSTPluginMain"
)

var (
	m           sync.RWMutex
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
func (l *Library) Open() (*Plugin, error) {
	plugin := &Plugin{
		Path:     l.Path,
		Name:     l.Name,
		callback: DefaultHostCallback,
		timeInfo: &vstTimeInfo{},
	}
	plugin.effect = C.loadEffect(C.vstPluginFuncPtr(l.entryPoint))
	m.Lock()
	plugins[plugin.effect] = plugin
	m.Unlock()
	return plugin, nil
}

// Close cleans up C refs for plugin
func (p *Plugin) Close() error {
	p.Dispatch(EffClose, 0, 0, nil, 0.0)
	m.Lock()
	delete(plugins, p.effect)
	m.Unlock()
	p.timeInfo = nil
	p.effect = nil
	return nil
}

// Dispatch wraps-up C method to dispatch calls to plugin
func (p *Plugin) Dispatch(opcode PluginOpcode, index int64, value int64, ptr unsafe.Pointer, opt float64) {
	if p.effect != nil {
		C.dispatch(p.effect, C.int(opcode), C.int(index), C.VstIntPtr(value), ptr, C.float(opt))
	}
}

// CanProcessFloat32 checks if plugin can process float32
func (p *Plugin) CanProcessFloat32() bool {
	if p.effect == nil {
		return false
	}
	return p.effect.flags&C.effFlagsCanReplacing == C.effFlagsCanReplacing
}

// CanProcessFloat64 checks if plugin can process float64
func (p *Plugin) CanProcessFloat64() bool {
	if p.effect == nil {
		return false
	}
	return p.effect.flags&C.effFlagsCanDoubleReplacing == C.effFlagsCanDoubleReplacing
}

// ProcessFloat64 audio with VST plugin
func (p *Plugin) ProcessFloat64(in [][]float64) [][]float64 {
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

	C.processDouble(p.effect, C.int(numChannels), C.int(blocksize), &input[0], &output[0])

	//convert []*C.double slices to [][]float64
	out := make([][]float64, numChannels)
	for i, data := range output {
		// copy data from C array to slice
		pa := (*[1 << 30]C.float)(unsafe.Pointer(data))
		out[i] = make([]float64, blocksize)
		for j := range out[i] {
			out[i][j] = float64(pa[j])
		}
	}
	return out
}

// ProcessFloat32 audio with VST plugin
func (p *Plugin) ProcessFloat32(in [][]float32) (out [][]float32) {
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

	C.processFloat(p.effect, C.int(numChannels), C.int(blocksize), &input[0], &output[0])

	//convert []*C.float slices to [][]float32
	out = make([][]float32, numChannels)
	for i, data := range output {
		// copy data from C array to slice
		pa := (*[1 << 30]C.float)(unsafe.Pointer(data))
		out[i] = make([]float32, blocksize)
		for j := range out[i] {
			out[i][j] = float32(pa[j])
		}
	}
	return out
}

// SetCallback overrides plugin's callback
func (p *Plugin) SetCallback(c HostCallbackFunc) {
	if c != nil {
		p.callback = c
	}
}

// SetSpeakerArrangement craetes and passes VstSpeakerArrangement structures to plugin
func (p *Plugin) SetSpeakerArrangement(numChannels int) {
	in := newSpeakerArrangement(numChannels)
	out := newSpeakerArrangement(numChannels)
	inp := uintptr(unsafe.Pointer(in))
	p.Dispatch(EffSetSpeakerArrangement, 0, int64(inp), unsafe.Pointer(out), 0.0)
}

func newSpeakerArrangement(numChannels int) *speakerArrangement {
	sa := speakerArrangement{}
	sa.numChannels = C.int(numChannels)
	switch numChannels {
	case 0:
		sa._type = C.kSpeakerArrEmpty
	case 1:
		sa._type = C.kSpeakerArrMono
	case 2:
		sa._type = C.kSpeakerArrStereo
	case 3:
		sa._type = C.kSpeakerArr30Music
	case 4:
		sa._type = C.kSpeakerArr40Music
	case 5:
		sa._type = C.kSpeakerArr50
	case 6:
		sa._type = C.kSpeakerArr60Music
	case 7:
		sa._type = C.kSpeakerArr70Music
	case 8:
		sa._type = C.kSpeakerArr80Music
	default:
		sa._type = C.kSpeakerArrUserDefined
	}

	for i := 0; i < numChannels; i++ {
		sa.speakers[i].azimuth = 0.0
		sa.speakers[i].elevation = 0.0
		sa.speakers[i].radius = 0.0
		sa.speakers[i].reserved = 0.0
		sa.speakers[i].name[0] = C.char(0)
		sa.speakers[i]._type = C.kSpeakerUndefined
	}
	return &sa
}

// SetTimeInfo sets new time info and returns pointer to it
func (p *Plugin) SetTimeInfo(sampleRate int, samplePos int64, tempo int, timeSigNum int, timeSigDenom int, nanoSeconds int64, ppqPos float64, barPos float64) int64 {
	// sample position
	p.timeInfo.sampleRate = C.double(sampleRate)
	p.timeInfo.samplePos = C.double(samplePos)
	p.timeInfo.flags |= C.kVstTransportPlaying
	p.timeInfo.flags |= C.kVstTransportChanged

	// nanoseconds
	p.timeInfo.nanoSeconds = C.double(nanoSeconds)
	p.timeInfo.flags |= C.kVstNanosValid

	// tempo
	p.timeInfo.tempo = C.double(tempo)
	p.timeInfo.flags |= C.kVstTempoValid

	// time signature
	p.timeInfo.timeSigNumerator = C.int(timeSigNum)
	p.timeInfo.timeSigDenominator = C.int(timeSigDenom)
	p.timeInfo.flags |= C.kVstTimeSigValid

	// ppq
	p.timeInfo.ppqPos = C.double(ppqPos)
	p.timeInfo.flags |= C.kVstPpqPosValid

	// bar start
	p.timeInfo.barStartPos = C.double(barPos)
	p.timeInfo.flags |= C.kVstBarsValid

	return int64(uintptr(unsafe.Pointer(p.timeInfo)))
}

//export hostCallback
// calls real callback
func hostCallback(effect *C.AEffect, opcode int64, index int64, value int64, ptr unsafe.Pointer, opt float64) int {
	// AudioMasterVersion is requested when plugin is created
	// It's never in map
	if MasterOpcode(opcode) == AudioMasterVersion {
		return vst2version
	}
	m.RLock()
	plugin, ok := plugins[effect]
	m.RUnlock()
	if !ok {
		panic("Plugin was closed")
	}

	if plugin.callback == nil {
		panic("Host callback is not defined!")
	}
	return plugin.callback(plugin, MasterOpcode(opcode), index, value, ptr, opt)
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
