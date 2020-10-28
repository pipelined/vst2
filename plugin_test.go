package vst2_test

import (
	"fmt"
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
	p := v.Load(vst2.DefaultHostCallback(&host))
	defer p.Close()

	p.SetParamValue(0, 0)
	for i := 0; i < p.NumParams(); i++ {
		fmt.Printf("param %d \tname: %v \tdisplay: %v \tlabel: %v \tvalue: %v\n", i, p.ParamName(i), p.ParamValueName(i), p.ParamUnitName(i), p.ParamValue(i))
	}

	fmt.Printf("programs: %v\n", p.NumPrograms())
	fmt.Printf("current program: %v\n", p.Program())
	fmt.Printf("current program name: %v\n", p.CurrentProgramName())
	p.SetProgramName("test")
	fmt.Printf("current program name: %v\n", p.CurrentProgramName())
	for i := 0; i < p.NumPrograms(); i++ {
		fmt.Printf("prog name: \t%v\n", p.ProgramName(i))
	}
	assertEqual(t, "resonance", p.ParamValue(4), float32(0))
	prog := p.GetProgramData()
	// fmt.Printf("program data before: %v\n", string(prog))
	newProg := strings.ReplaceAll(string(prog), "resonance=\"0.0\"", "resonance=\"1.0\"")

	// preset := getPreset(t)
	p.SetProgramData(([]byte)(newProg))

	prog = p.GetProgramData()
	// fmt.Printf("program data after: %v\n", string(prog))
	// prog = p.GetProgramData()
	// fmt.Printf("program data: %v\n", string(prog))
	assertEqual(t, "resonance", p.ParamValue(4), float32(1))
}

// func getPreset(t *testing.T) []byte {
// 	preset, err := os.Open("_testdata/preset")
// 	assertEqual(t, "open preset error", err, nil)
// 	b, err := ioutil.ReadAll(preset)
// 	assertEqual(t, "load preset error", err, nil)
// 	return b
// }

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