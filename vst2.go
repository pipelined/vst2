// Package vst2 provides interface to VST2 plugins.
package vst2

/*
#cgo CFLAGS: -std=gnu99 -I${SRCDIR}
#include <stdlib.h>
#include <stdint.h>
#include "vst.h"
*/
import "C"
import (
	"fmt"
	"path/filepath"
	"sync"
	"unsafe"
)

// global state for callbacks.
var (
	mutex     sync.RWMutex
	callbacks = make(map[*effect]HostCallbackFunc)
)

//export hostCallback
// global hostCallback, calls real callback.
func hostCallback(e *effect, opcode int64, index int64, value int64, ptr unsafe.Pointer, opt float64) Return {
	// AudioMasterVersion is requested when plugin is created
	// It's never in map
	if HostOpcode(opcode) == HostVersion {
		return version
	}
	mutex.RLock()
	c, ok := callbacks[e]
	mutex.RUnlock()
	if !ok {
		panic("plugin was closed")
	}

	if c == nil {
		panic("host callback is undefined")
	}
	return c(HostOpcode(opcode), Index(index), Value(value), Ptr(ptr), Opt(opt))
}

const (
	// VST main function name.
	main = "VSTPluginMain"
	// VST API version.
	version = 2400
)

type (
	// HostCallbackFunc used as callback function called by plugin. Use closure
	// wrapping technic to add more types to callback.
	HostCallbackFunc func(HostOpcode, Index, Value, Ptr, Opt) Return

	// Index is index in plugin dispatch/host callback.
	Index int64
	// Value is value in plugin dispatch/host callback.
	Value int64
	// Ptr is ptr in plugin dispatch/host callback.
	Ptr unsafe.Pointer
	// Opt is opt in plugin dispatch/host callback.
	Opt float64
	// Return is returned value for dispatch/host callback.
	Return int64
)

type (
	// Effect is an alias on C effect type.
	effect C.Effect

	// effectMain is a reference to VST main function.
	// wrapper on C entry point.
	effectMain C.EntryPoint

	// VST used to create new instances of plugin.
	// It also keeps reference to VST handle to clean up on Close.
	VST struct {
		main effectMain
		// handle is OS-specific.
		handle
		Name string
		Path string
	}

	// Plugin is VST2 plugin instance.
	Plugin struct {
		*effect
		Name string
		Path string
	}
)

// Open loads the VST into memory and stores entry point func.
func Open(path string) (VST, error) {
	p, err := filepath.Abs(path)
	if err != nil {
		return VST{}, fmt.Errorf("failed to get absolute path: %w", err)
	}

	m, h, err := open(p)
	if err != nil {
		return VST{}, fmt.Errorf("failed to load VST '%s': %w", path, err)
	}

	return VST{
		Path:   p,
		main:   m,
		handle: h,
	}, nil
}

// Close cleans up VST resoures.
func (v VST) Close() error {
	if v.main == nil {
		return nil
	}
	v.main = nil
	if err := v.handle.close(); err != nil {
		return fmt.Errorf("failed close VST %s: %w", v.Name, err)
	}
	return nil
}

// Load new instance of VST plugin with provided callback.
// This function also calls dispatch with EffOpen opcode.
func (v VST) Load(c HostCallbackFunc) *Plugin {
	if v.main == nil || c == nil {
		return nil
	}
	e := (*effect)(C.loadEffect(v.main))
	mutex.Lock()
	callbacks[e] = c
	mutex.Unlock()

	p := &Plugin{
		effect: e,
		Path:   v.Path,
		Name:   v.Name,
	}
	p.Dispatch(EffOpen, 0, 0, nil, 0.0)
	return p
}

// Close cleans up C refs for plugin
func (p *Plugin) Close() error {
	if p.effect == nil {
		return nil
	}
	p.Dispatch(EffClose, 0, 0, nil, 0.0)
	p.effect = nil
	return nil
}

// Dispatch wraps-up C method to dispatch calls to plugin
func (p *Plugin) Dispatch(opcode EffectOpcode, index Index, value Value, ptr Ptr, opt Opt) Return {
	return Return(C.dispatch((*C.Effect)(p.effect), C.int(opcode), C.int(index), C.int64_t(value), unsafe.Pointer(ptr), C.float(opt)))
}

// CanProcessFloat32 checks if plugin can process float32.
func (p *Plugin) CanProcessFloat32() bool {
	if p == nil {
		return false
	}
	return EffectFlags(p.effect.flags)&EffFlagsCanReplacing == EffFlagsCanReplacing
}

// CanProcessFloat64 checks if plugin can process float64.
func (p *Plugin) CanProcessFloat64() bool {
	if p == nil {
		return false
	}
	return EffectFlags(p.effect.flags)&EffFlagsCanDoubleReplacing == EffFlagsCanDoubleReplacing
}

// ProcessDouble audio with VST plugin.
func (p *Plugin) ProcessDouble(in, out DoubleBuffer) {
	C.processDouble(
		(*C.Effect)(p.effect),
		C.int(in.numChannels),
		C.int(in.size),
		&in.data[0],
		&out.data[0],
	)
}

// ProcessFloat audio with VST plugin.
func (p *Plugin) ProcessFloat(in, out FloatBuffer) {
	C.processFloat(
		(*C.Effect)(p.effect),
		C.int(in.numChannels),
		C.int(in.size),
		&in.data[0],
		&out.data[0],
	)
}

// Start the plugin.
func (p *Plugin) Start() {
	p.Dispatch(EffStateChanged, 0, 1, nil, 0.0)
}

// Stop the plugin.
func (p *Plugin) Stop() {
	p.Dispatch(EffStateChanged, 0, 0, nil, 0.0)
}

// SetBufferSize sets a buffer size
func (p *Plugin) SetBufferSize(bufferSize int) {
	p.Dispatch(EffSetBufferSize, 0, Value(bufferSize), nil, 0.0)
}

// SetSampleRate sets a sample rate for plugin
func (p *Plugin) SetSampleRate(sampleRate int) {
	p.Dispatch(EffSetSampleRate, 0, 0, nil, Opt(sampleRate))
}

// SetSpeakerArrangement craetes and passes SpeakerArrangement structures to plugin
func (p *Plugin) SetSpeakerArrangement(in, out *SpeakerArrangement) {
	p.Dispatch(EffSetSpeakerArrangement, 0, in.Value(), out.Ptr(), 0.0)
}

// ScanPaths returns a slice of default vst2 locations.
// Locations are OS-specific.
func ScanPaths() (paths []string) {
	return append([]string{}, scanPaths...)
}

func newSpeakerArrangement(numChannels int) *SpeakerArrangement {
	sa := SpeakerArrangement{}
	sa.NumChannels = int32(numChannels)
	switch numChannels {
	case 0:
		sa.Type = SpeakerArrEmpty
	case 1:
		sa.Type = SpeakerArrMono
	case 2:
		sa.Type = SpeakerArrStereo
	case 3:
		sa.Type = SpeakerArr30Music
	case 4:
		sa.Type = SpeakerArr40Music
	case 5:
		sa.Type = SpeakerArr50
	case 6:
		sa.Type = SpeakerArr60Music
	case 7:
		sa.Type = SpeakerArr70Music
	case 8:
		sa.Type = SpeakerArr80Music
	default:
		sa.Type = SpeakerArrUserDefined
	}

	for i := 0; i < int(numChannels); i++ {
		sa.Speakers[i].Type = SpeakerUndefined
	}
	return &sa
}
