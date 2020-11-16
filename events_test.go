package vst2_test

import (
	"fmt"
	"testing"

	"pipelined.dev/audio/vst2"
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
	vst2.TestEvents(events)

	for i := 0; i < events.NumEvents(); i++ {
		fmt.Printf("%v", events.Event(i))
		switch ev := events.Event(i).(type) {
		case *vst2.MIDIEvent:
		case *vst2.SysExMIDIEvent:
			fmt.Printf("dump %v\n", ev.SysExDump)
		}
	}
}
