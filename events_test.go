package vst2_test

import (
	"testing"

	"pipelined.dev/audio/vst2"
)

func TestEvents(t *testing.T) {
	dump := vst2.SysExDump([]byte("this is a test"))
	defer dump.Free()
	events := vst2.Events(
		&vst2.SysExMIDIEvent{
			SysExDump: dump,
		},
	)

	assertNotNil(t, "events", events)
	vst2.TestEvents(events)
	events.Free()
}
