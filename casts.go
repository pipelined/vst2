package vst2

import "unsafe"

func (sa *SpeakerArrangement) Value() Value {
	if sa == nil {
		return 0
	}
	return Value(uintptr(unsafe.Pointer(sa)))
}

func (sa *SpeakerArrangement) Ptr() Ptr {
	if sa == nil {
		return nil
	}
	return Ptr(unsafe.Pointer(sa))
}

func (ti *TimeInfo) Return() Return {
	if ti == nil {
		return 0
	}
	return Return(uintptr(unsafe.Pointer(ti)))
}
