package vst2

/*
#cgo CFLAGS: -std=gnu99 -I${SRCDIR}
#include <stdlib.h>
#include "sdk.h"
*/
import "C"
import "unsafe"

const (
	midiEventSize  = 32
	sysExEventSize = 44
)

type (
	// EventType denotes the type of event.
	eventType int32
)

const (
	// MIDI event.
	MIDI eventType = iota + 1
	audio
	video
	parameter
	trigger
	// SysExMIDI system exclusive midi event.
	SysExMIDI
)

// EventsPtr is a container for events to be processed by plugin or host.
type EventsPtr C.Events

// Events allocates new events container and place there provided events.
// It must be freed after use.
func Events(events ...Event) *EventsPtr {
	CEvents := C.newEvents(C.int32_t(len(events)))
	for i := range events {
		switch e := events[i].(type) {
		case *MIDIEvent:
			e.eventType = MIDI
			e.byteSize = midiEventSize
			C.setEvent(CEvents, unsafe.Pointer(e), C.int32_t(i))
		case *SysExMIDIEvent:
			e.eventType = SysExMIDI
			e.byteSize = sysExEventSize
			C.setEvent(CEvents, unsafe.Pointer(e), C.int32_t(i))
		}
	}
	return (*EventsPtr)(CEvents)
}

// NumEvents returns number of events within container.
func (e *EventsPtr) NumEvents() int {
	return int(e.numEvents)
}

// Event returns an event from the container.
func (e *EventsPtr) Event(i int) Event {
	type event struct {
		eventType
	}
	ev := (*event)(C.getEvent((*C.Events)(e), C.int32_t(i)))
	switch ev.eventType {
	case MIDI:
		return (*MIDIEvent)(unsafe.Pointer(ev))
	case SysExMIDI:
		return (*SysExMIDIEvent)(unsafe.Pointer(ev))
	}
	return nil
}

// Free memory allocated for container.
func (e *EventsPtr) Free() {
	C.free(unsafe.Pointer(e.events))
	C.free(unsafe.Pointer(e))
}

// Event interface is used to provide safety for events mapping. It's
// implemented by MIDIEvent and SysExMIDIEvent types.
type Event interface {
	sealedEvent()
}

func (*MIDIEvent) sealedEvent()      {}
func (*SysExMIDIEvent) sealedEvent() {}

type (
	// MIDIEvent contains midi information.
	MIDIEvent struct {
		eventType             // Always MIDI.
		byteSize        int32 // Always 32.
		DeltaFrames     int32 // Number of sample frames into the current processing block that this event occurs on.
		Flags           MIDIEventFlag
		NoteLength      int32   // In sample frames, 0 if not available.
		NoteOffset      int32   // In sample frames from note start, 0 if not available.
		Data            [3]byte // 1 to 3 MIDI bytes.
		dataReserved    byte
		Detune          uint8 // Between -64 to +63 cents, for scales other than 'well-tempered' e.g. 'microtuning'.
		NoteOffVelocity uint8 // Between 0 and 127.
		reserved1       uint8 // Not used.
		reserved2       uint8 // Not used.
	}

	// MIDIEventFlag is set in midi event.
	MIDIEventFlag int32
)

const (
	// MIDIEventRealtime means that this event is played and not coming
	// from sequencer.
	MIDIEventRealtime MIDIEventFlag = 1
)

type (
	// SysExMIDIEvent is system exclusive MIDI event.
	SysExMIDIEvent struct {
		eventType         // Always SysExMIDI.
		byteSize    int32 // Always 44
		DeltaFrames int32 // Number of sample frames into the current processing block that this event occurs on.
		flags       int32 // Not used.
		SysExDump   SysExDataPtr
		resvd2      int64 // Not used.
	}

	// SysExDataPtr holds sysex dump data.
	SysExDataPtr struct {
		length   int32
		reserved int64
		data     *C.char
	}
)

// SysExData allocates new SysExDumpPtr with provided bytes. It must be
// freed after use.
func SysExData(b []byte) SysExDataPtr {
	return SysExDataPtr{
		length: int32(len(b)),
		data:   (*C.char)(C.CBytes(b)),
	}
}

// Bytes returns bytes representation of sysex data.
func (s SysExDataPtr) Bytes() []byte {
	return C.GoBytes(unsafe.Pointer(s.data), C.int(s.length))
}

// Free releases allocated memory.
func (s SysExDataPtr) Free() {
	C.free(unsafe.Pointer(unsafe.Pointer(s.data)))
}

// TestEvents is a helper function to test events.
func TestEvents(e *EventsPtr) {
	C.testEvents((*C.Events)(e))
}
