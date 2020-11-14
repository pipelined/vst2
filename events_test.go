package vst2_test

import (
	"testing"

	"pipelined.dev/audio/vst2"
)

func TestEvents(t *testing.T) {
	events := vst2.Events(
		&vst2.MIDIEvent{},
		&vst2.MIDIEvent{},
		&vst2.MIDIEvent{},
	)
	assertNotNil(t, "events", events)
	vst2.TestEvents(events)
	events.Free()
}
