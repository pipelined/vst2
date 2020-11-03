package vst2_test

import (
	"reflect"
	"strings"
	"testing"

	"pipelined.dev/audio/vst2"
)

func TestPluginParameters(t *testing.T) {
	v, err := vst2.Open(pluginPath())
	assertEqual(t, "vst error", err, nil)
	defer v.Close()

	host := vst2.HostProperties{
		BufferSize: 512,
		Channels:   2,
		SampleRate: 44100,
	}

	testPlugin := func(fn func(p *vst2.Effect)) func(*testing.T) {
		return func(t *testing.T) {
			p := v.Load(vst2.DefaultHostCallback(&host))
			defer p.Close()
			fn(p)
		}
	}
	t.Run("get params", testPlugin(func(p *vst2.Effect) {
		numParams := p.NumParams()
		assertEqual(t, "num params", numParams > 0, true)
		for i := 0; i < numParams; i++ {
			p.ParamName(i)
			p.ParamUnitName(i)
			p.ParamValueName(i)
		}
	}))
	t.Run("set param", testPlugin(func(p *vst2.Effect) {
		assertEqual(t, "resonance before", p.ParamValue(4), float32(0))
		p.SetParamValue(4, 0.5)
		assertEqual(t, "resonance after", p.ParamValue(4), float32(0.5))
	}))
	t.Run("get programs", testPlugin(func(p *vst2.Effect) {
		numPrograms := p.NumPrograms()
		assertEqual(t, "num programs", numPrograms > 0, true)
		for i := 0; i < numPrograms; i++ {
			name := p.ProgramName(i)
			assertEqual(t, "name not empty", len(name) > 0, true)
		}
	}))
	t.Run("set program name", testPlugin(func(p *vst2.Effect) {
		assertEqual(t, "program name before", p.CurrentProgramName(), "! Startup Juno Osc TAL")
		newProgName := "test"
		p.SetCurrentProgramName(newProgName)
		assertEqual(t, "program name after", p.CurrentProgramName(), newProgName)
	}))
	t.Run("set program", testPlugin(func(p *vst2.Effect) {
		assertEqual(t, "program name before", p.CurrentProgramName(), "! Startup Juno Osc TAL")
		p.SetProgram(2)
		assertEqual(t, "program name after", p.CurrentProgramName(), "ARP 303 Like FN")
	}))
	t.Run("set program data", testPlugin(func(p *vst2.Effect) {
		assertEqual(t, "resonance before", p.ParamValue(4), float32(0))
		progData := string(p.GetProgramData())
		progData = strings.ReplaceAll(progData, "resonance=\"0.0\"", "resonance=\"1.0\"")
		p.SetProgramData([]byte(progData))
		assertEqual(t, "resonance before", p.ParamValue(4), float32(1))
	}))
	t.Run("set bank data", testPlugin(func(p *vst2.Effect) {
		assertEqual(t, "resonance before", p.ParamValue(4), float32(0))
		progData := string(p.GetBankData())
		progData = strings.ReplaceAll(progData, "resonance=\"0.0\"", "resonance=\"1.0\"")
		p.SetBankData([]byte(progData))
		assertEqual(t, "resonance before", p.ParamValue(4), float32(1))
	}))
}

func assertEqual(t *testing.T, name string, result, expected interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("%v\nresult: \t%T\t%+v \nexpected: \t%T\t%+v", name, result, result, expected, expected)
	}
}

func assertNotNil(t *testing.T, name string, result interface{}) {
	t.Helper()
	if reflect.DeepEqual(nil, result) {
		t.Fatalf("%v\nresult: \t%T\t%+v \nexpected: \t%T\t%+v", name, result, result, nil, nil)
	}
}

func assertPanic(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic")
		}
	}()
	fn()
}
