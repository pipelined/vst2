package vst2

import (
	"pipelined.dev/audio/vst2/sdk"
)

// Plugin is VST2 plugin instance.
type Plugin struct {
	*sdk.Effect
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

// ParamName returns the parameter label: "Release", "Gain", etc.
func (p *Plugin) ParamName(index int) string {
	var s sdk.ParamString
	p.Effect.ParamName(index, &s)
	return s.String()
}

// ParamValueName returns the parameter value label: "0.5", "HALL", etc.
func (p *Plugin) ParamValueName(index int) string {
	var s sdk.ParamString
	p.Effect.ParamValueName(index, &s)
	return s.String()
}

// ParamUnitName returns the parameter unit label: "db", "ms", etc.
func (p *Plugin) ParamUnitName(index int) string {
	var s sdk.ParamString
	p.Effect.ParamUnitName(index, &s)
	return s.String()
}

// CurrentProgramName returns current program name.
func (p *Plugin) CurrentProgramName() string {
	var s sdk.ProgramString
	p.Effect.CurrentProgramName(&s)
	return s.String()
}

// ProgramName returns program name for provided program index.
func (p *Plugin) ProgramName(index int) string {
	var s sdk.ProgramString
	p.Effect.ProgramName(index, &s)
	return s.String()
}

// SetProgramName sets new name to the current program.
func (p *Plugin) SetProgramName(name string) {
	var s sdk.ProgramString
	copy(s[:], []byte(name))
	p.Effect.SetCurrentProgramName(&s)
}
