package vst2

/*
#include <stdlib.h>
*/
import "C"
import (
	"strings"
	"unsafe"

	"pipelined.dev/audio/vst2/sdk"
)

// Plugin is VST2 plugin instance.
type Plugin struct {
	*sdk.Effect
}

// ParamString used to get parameter string values: name, unit name and
// value name.
type ParamString [sdk.MaxParamStrLen]byte

func (s ParamString) String() string {
	return trimNull(string(s[:]))
}

// ProgramString used to get and set program name.
type ProgramString [sdk.MaxProgNameLen]byte

func (s ProgramString) String() string {
	return trimNull(string(s[:]))
}

// Close cleans up C refs for plugin
func (p *Plugin) Close() error {
	p.Dispatch(sdk.EffClose, 0, 0, nil, 0.0)
	if err := p.Effect.Close(); err != nil {
		return err
	}
	p.Effect = nil
	return nil
}

// CanProcessFloat32 checks if plugin can process float32.
func (p *Plugin) CanProcessFloat32() bool {
	if p == nil {
		return false
	}
	return p.Flags()&sdk.EffFlagsCanReplacing == sdk.EffFlagsCanReplacing
}

// CanProcessFloat64 checks if plugin can process float64.
func (p *Plugin) CanProcessFloat64() bool {
	if p == nil {
		return false
	}
	return p.Flags()&sdk.EffFlagsCanDoubleReplacing == sdk.EffFlagsCanDoubleReplacing
}

// Start the plugin.
func (p *Plugin) Start() {
	p.Dispatch(sdk.EffStateChanged, 0, 1, nil, 0)
}

// Stop the plugin.
func (p *Plugin) Stop() {
	p.Dispatch(sdk.EffStateChanged, 0, 0, nil, 0)
}

// SetBufferSize sets a buffer size per channel.
func (p *Plugin) SetBufferSize(bufferSize int) {
	p.Dispatch(sdk.EffSetBufferSize, 0, sdk.Value(bufferSize), nil, 0)
}

// SetSampleRate sets a sample rate for plugin.
func (p *Plugin) SetSampleRate(sampleRate int) {
	p.Dispatch(sdk.EffSetSampleRate, 0, 0, nil, sdk.Opt(sampleRate))
}

// SetSpeakerArrangement creates and passes SpeakerArrangement structures to plugin
func (p *Plugin) SetSpeakerArrangement(in, out *sdk.SpeakerArrangement) {
	p.Dispatch(sdk.EffSetSpeakerArrangement, 0, in.Value(), out.Ptr(), 0)
}

// ParamProperties returns parameter properties for provided parameter
// index. If opcode is not supported, boolean result is false.
func (p *Plugin) ParamProperties(index int) (*sdk.ParameterProperties, bool) {
	var props sdk.ParameterProperties
	r := p.Dispatch(sdk.EffGetParameterProperties, sdk.Index(index), 0, sdk.Ptr(&props), 0)
	if r > 0 {
		return &props, true
	}
	return nil, false
}

// ParamName returns the parameter label: "Release", "Gain", etc.
func (p *Plugin) ParamName(index int) string {
	var s ParamString
	p.Dispatch(sdk.EffGetParamName, sdk.Index(index), 0, sdk.Ptr(&s), 0)
	return s.String()
}

// ParamValueName returns the parameter value label: "0.5", "HALL", etc.
func (p *Plugin) ParamValueName(index int) string {
	var s ParamString
	p.Dispatch(sdk.EffGetParamDisplay, sdk.Index(index), 0, sdk.Ptr(&s), 0)
	return s.String()
}

// ParamUnitName returns the parameter unit label: "db", "ms", etc.
func (p *Plugin) ParamUnitName(index int) string {
	var s ParamString
	p.Dispatch(sdk.EffGetParamLabel, sdk.Index(index), 0, sdk.Ptr(&s), 0)
	return s.String()
}

// Program returns current program number.
func (p *Plugin) Program() int {
	return int(p.Dispatch(sdk.EffGetProgram, 0, 0, nil, 0))
}

// SetProgram changes current program index.
func (p *Plugin) SetProgram(index int) {
	p.Dispatch(sdk.EffSetProgram, 0, sdk.Value(index), nil, 0)
}

// CurrentProgramName returns current program name.
func (p *Plugin) CurrentProgramName() string {
	var s ProgramString
	p.Dispatch(sdk.EffGetProgramName, 0, 0, sdk.Ptr(&s), 0)
	return s.String()
}

// ProgramName returns program name for provided program index.
func (p *Plugin) ProgramName(index int) string {
	var s ProgramString
	p.Dispatch(sdk.EffGetProgramNameIndexed, sdk.Index(index), 0, sdk.Ptr(&s), 0)
	return s.String()
}

// SetProgramName sets new name to the current program.
func (p *Plugin) SetProgramName(name string) {
	var s ProgramString
	copy(s[:], []byte(name))
	p.Dispatch(sdk.EffSetProgramName, 0, 0, sdk.Ptr(&s), 0)
}

// GetProgramData returns current preset data. Plugin allocates required
// memory, then this function allocates new byte slice of required length
// where data is copied.
func (p *Plugin) GetProgramData() []byte {
	var ptr unsafe.Pointer
	length := C.int(p.Dispatch(sdk.EffGetChunk, 1, 0, sdk.Ptr(&ptr), 0))
	return C.GoBytes(ptr, length)
}

// SetProgramData sets preset data to the plugin. Data is the full preset
// including chunk header.
func (p *Plugin) SetProgramData(data []byte) {
	ptr := C.CBytes(data)
	p.Dispatch(sdk.EffSetChunk, 1, sdk.Value(len(data)), sdk.Ptr(ptr), 0)
	C.free(ptr)
}

// GetBankData returns current bank data. Plugin allocates required
// memory, then this function allocates new byte slice of required length
// where data is copied.
func (p *Plugin) GetBankData() []byte {
	var ptr unsafe.Pointer
	length := C.int(p.Dispatch(sdk.EffGetChunk, 0, 0, sdk.Ptr(&ptr), 0))
	return C.GoBytes(ptr, length)
}

// SetBankData sets preset data to the plugin. Data is the full preset
// including chunk header.
func (p *Plugin) SetBankData(data []byte) {
	ptr := C.CBytes(data)
	p.Dispatch(sdk.EffSetChunk, 0, sdk.Value(len(data)), sdk.Ptr(ptr), 0)
	C.free(ptr)
}

func trimNull(s string) string {
	return strings.Trim(s, "\x00")
}
