package vst2_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/cwbudde/vst2"
)

func TestEvents(t *testing.T) {
	dump := vst2.SysExData([]byte("this is a test"))
	defer dump.Free()
	events := vst2.Events(
		&vst2.MIDIEvent{},
		&vst2.SysExMIDIEvent{
			SysExDump: dump,
		},
	)
	defer events.Free()

	assertNotNil(t, "events", events)

	for i := 0; i < events.NumEvents(); i++ {
		fmt.Printf("%v", events.Event(i))
		switch ev := events.Event(i).(type) {
		case *vst2.MIDIEvent:
		case *vst2.SysExMIDIEvent:
			fmt.Printf("dump %v\n", string(ev.SysExDump.Bytes()))
		}
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
