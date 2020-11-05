package vst2

import "unsafe"

// Value cast used in EffSetSpeakerArrangement call.
func (sa *SpeakerArrangement) Value() Value {
	if sa == nil {
		return 0
	}
	return Value(uintptr(unsafe.Pointer(sa)))
}

// Ptr cast used in EffSetSpeakerArrangement call.
func (sa *SpeakerArrangement) Ptr() unsafe.Pointer {
	if sa == nil {
		return nil
	}
	return unsafe.Pointer(unsafe.Pointer(sa))
}
