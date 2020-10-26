package vst2_test

import (
	"fmt"
	"reflect"
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
