package vst2

/*
#cgo CFLAGS: -std=gnu99 -I${SRCDIR}/sdk
#include <stdlib.h>
#include "sdk.h"
*/
import "C"
import "unsafe"

// Event interface is used to provide safety for events mapping. It's
// implemented by MIDIEvent and SysExMIDIEvent types.
type Event interface {
	sealedEvent()
}

func (*MIDIEvent) sealedEvent() {}

// MIDIEvents are MIDI events to be processed by plugin or host.
type MIDIEvents C.Events

// Events allocates new MIDI events container and place there provided
// events.
func Events(events ...Event) *MIDIEvents {
	CEvents := C.newEvents(C.int32_t(len(events)))

	for i := range events {
		switch e := events[i].(type) {
		case *MIDIEvent:
			e.eventType = MIDI
			C.setEvent(CEvents, unsafe.Pointer(e), C.int32_t(i))
		}
	}
	return (*MIDIEvents)(CEvents)
}

// Free memory allocated for events.
func (e *MIDIEvents) Free() {
	C.free(unsafe.Pointer(e.events))
	C.free(unsafe.Pointer(e))
}

// TestEvents is a helper function to test events.
func TestEvents(e *MIDIEvents) {
	C.testEvents((*C.Events)(e))
}
