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
func (sa *SpeakerArrangement) Ptr() Ptr {
	if sa == nil {
		return nil
	}
	return Ptr(unsafe.Pointer(sa))
}

// Return cast used in HostGetTime call.
func (ti *TimeInfo) Return() Return {
	if ti == nil {
		return 0
	}
	return Return(uintptr(unsafe.Pointer(ti)))
}
