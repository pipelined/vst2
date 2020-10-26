package vst2

/*
#cgo CFLAGS: -I${SRCDIR}
#include <stdlib.h>
#include <stdint.h>
#include "vst.h"
*/
import "C"
import (
	"unsafe"
)

// Plugin is VST2 plugin instance.
type Plugin struct {
	*effect
	Name string
	Path string
}

// Close cleans up C refs for plugin
func (p *Plugin) Close() error {
	if p.effect == nil {
		return nil
	}
	p.Dispatch(EffClose, 0, 0, nil, 0.0)
	mutex.Lock()
	delete(callbacks, p.effect)
	mutex.Unlock()
	p.effect = nil
	return nil
}

// NumParams returns the number of parameters.
func (p *Plugin) NumParams() int {
	return int(p.effect.numParams)
}

// Dispatch wraps-up C method to dispatch calls to plugin
func (p *Plugin) Dispatch(opcode EffectOpcode, index Index, value Value, ptr Ptr, opt Opt) Return {
	return Return(C.dispatch((*C.Effect)(p.effect), C.int32_t(opcode), C.int32_t(index), C.int64_t(value), unsafe.Pointer(ptr), C.float(opt)))
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
		C.int32_t(in.numChannels),
		C.int32_t(in.size),
		&in.data[0],
		&out.data[0],
	)
}

// ProcessFloat audio with VST plugin.
func (p *Plugin) ProcessFloat(in, out FloatBuffer) {
	C.processFloat(
		(*C.Effect)(p.effect),
		C.int32_t(in.numChannels),
		C.int32_t(in.size),
		&in.data[0],
		&out.data[0],
	)
}

// Start the plugin.
func (p *Plugin) Start() {
	p.Dispatch(EffStateChanged, 0, 1, nil, 0)
}

// Stop the plugin.
func (p *Plugin) Stop() {
	p.Dispatch(EffStateChanged, 0, 0, nil, 0)
}

// SetBufferSize sets a buffer size per channel.
func (p *Plugin) SetBufferSize(bufferSize int) {
	p.Dispatch(EffSetBufferSize, 0, Value(bufferSize), nil, 0)
}

// SetSampleRate sets a sample rate for plugin.
func (p *Plugin) SetSampleRate(sampleRate int) {
	p.Dispatch(EffSetSampleRate, 0, 0, nil, Opt(sampleRate))
}

// SetSpeakerArrangement creates and passes SpeakerArrangement structures to plugin
func (p *Plugin) SetSpeakerArrangement(in, out *SpeakerArrangement) {
	p.Dispatch(EffSetSpeakerArrangement, 0, in.Value(), out.Ptr(), 0)
}

// ParameterProperties returns parameter properties for provided parameter
// index. If opcode is not supported, boolean result is false.
func (p *Plugin) ParameterProperties(index int) (*ParameterProperties, bool) {
	var props ParameterProperties
	r := p.Dispatch(EffGetParameterProperties, Index(index), 0, Ptr(&props), 0)
	if r > 0 {
		return &props, true
	}
	return nil, false
}

// GetParameterValue returns the value of parameter.
func (p *Plugin) GetParameterValue(index int) float32 {
	return float32(C.getParameter((*C.Effect)(p.effect), C.int32_t(index)))
}

// SetParameterValue sets new value for parameter.
func (p *Plugin) SetParameterValue(index int, value float32) {
	C.setParameter((*C.Effect)(p.effect), C.int32_t(index), C.float(value))
}

// GetParamName returns the parameter label: "Release", "Gain", etc.
func (p *Plugin) GetParamName(index int) string {
	var val [maxParamStrLen]byte
	p.Dispatch(EffGetParamName, Index(index), 0, Ptr(&val), 0)
	return string(val[:])
}

// GetParamValueName returns the parameter value label: "0.5", "HALL", etc.
func (p *Plugin) GetParamValueName(index int) string {
	var val [maxParamStrLen]byte
	p.Dispatch(EffGetParamDisplay, Index(index), 0, Ptr(&val), 0)
	return string(val[:])
}

// GetParamUnitName returns the parameter unit label: "db", "ms", etc.
func (p *Plugin) GetParamUnitName(index int) string {
	var val [maxParamStrLen]byte
	p.Dispatch(EffGetParamLabel, Index(index), 0, Ptr(&val), 0)
	return string(val[:])
}
