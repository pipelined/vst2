package vst2

/*
#cgo CFLAGS: -I${SRCDIR}
#include <stdlib.h>
#include <stdint.h>
#include "vst.h"
*/
import "C"
import (
	"strings"
	"unsafe"
)

// Plugin is VST2 plugin instance.
type Plugin struct {
	*effect
	Name string
	Path string
}

// ParamString used to get parameter string values: name, unit name and
// value name.
type ParamString [maxParamStrLen]byte

func (s ParamString) String() string {
	return trimNull(string(s[:]))
}

// ProgramString used to get and set program name.
type ProgramString [maxProgNameLen]byte

func (s ProgramString) String() string {
	return trimNull(string(s[:]))
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

// NumPrograms returns the number of programs.
func (p *Plugin) NumPrograms() int {
	return int(p.effect.numPrograms)
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

// ParamProperties returns parameter properties for provided parameter
// index. If opcode is not supported, boolean result is false.
func (p *Plugin) ParamProperties(index int) (*ParameterProperties, bool) {
	var props ParameterProperties
	r := p.Dispatch(EffGetParameterProperties, Index(index), 0, Ptr(&props), 0)
	if r > 0 {
		return &props, true
	}
	return nil, false
}

// ParamValue returns the value of parameter.
func (p *Plugin) ParamValue(index int) float32 {
	return float32(C.getParameter((*C.Effect)(p.effect), C.int32_t(index)))
}

// SetParamValue sets new value for parameter.
func (p *Plugin) SetParamValue(index int, value float32) {
	C.setParameter((*C.Effect)(p.effect), C.int32_t(index), C.float(value))
}

// ParamName returns the parameter label: "Release", "Gain", etc.
func (p *Plugin) ParamName(index int) string {
	var s ParamString
	p.Dispatch(EffGetParamName, Index(index), 0, Ptr(&s), 0)
	return s.String()
}

// ParamValueName returns the parameter value label: "0.5", "HALL", etc.
func (p *Plugin) ParamValueName(index int) string {
	var s ParamString
	p.Dispatch(EffGetParamDisplay, Index(index), 0, Ptr(&s), 0)
	return s.String()
}

// ParamUnitName returns the parameter unit label: "db", "ms", etc.
func (p *Plugin) ParamUnitName(index int) string {
	var s ParamString
	p.Dispatch(EffGetParamLabel, Index(index), 0, Ptr(&s), 0)
	return s.String()
}

// Program returns current program number.
func (p *Plugin) Program() int {
	return int(p.Dispatch(EffGetProgram, 0, 0, nil, 0))
}

// SetProgram changes current program index.
func (p *Plugin) SetProgram(index int) {
	p.Dispatch(EffSetProgram, 0, Value(index), nil, 0)
}

// CurrentProgramName returns current program name.
func (p *Plugin) CurrentProgramName() string {
	var s ProgramString
	p.Dispatch(EffGetProgramName, 0, 0, Ptr(&s), 0)
	return s.String()
}

// ProgramName returns program name for provided program index.
func (p *Plugin) ProgramName(index int) string {
	var s ProgramString
	p.Dispatch(EffGetProgramNameIndexed, Index(index), 0, Ptr(&s), 0)
	return s.String()
}

// SetProgramName sets new name to the current program.
func (p *Plugin) SetProgramName(name string) {
	var s ProgramString
	copy(s[:], []byte(name))
	p.Dispatch(EffSetProgramName, 0, 0, Ptr(&s), 0)
}

// GetProgramData returns current preset data. Plugin allocates required
// memory, then this function allocates new byte slice of required length
// where data is copied.
func (p *Plugin) GetProgramData() []byte {
	var ptr unsafe.Pointer
	length := C.int(p.Dispatch(EffGetChunk, 1, 0, Ptr(&ptr), 0))
	return C.GoBytes(ptr, length)
}

// SetProgramData sets preset data to the plugin. Data is the full preset
// including chunk header.
func (p *Plugin) SetProgramData(data []byte) {
	ptr := C.CBytes(data)
	p.Dispatch(EffSetChunk, 1, Value(len(data)), Ptr(ptr), 0)
	C.free(ptr)
}

// GetBankData returns current bank data. Plugin allocates required
// memory, then this function allocates new byte slice of required length
// where data is copied.
func (p *Plugin) GetBankData() []byte {
	var ptr unsafe.Pointer
	length := C.int(p.Dispatch(EffGetChunk, 0, 0, Ptr(&ptr), 0))
	return C.GoBytes(ptr, length)
}

// SetBankData sets preset data to the plugin. Data is the full preset
// including chunk header.
func (p *Plugin) SetBankData(data []byte) {
	ptr := C.CBytes(data)
	p.Dispatch(EffSetChunk, 0, Value(len(data)), Ptr(ptr), 0)
	C.free(ptr)
}

func trimNull(s string) string {
	return strings.Trim(s, "\x00")
}
