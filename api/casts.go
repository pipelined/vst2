package api

import "unsafe"

func (sa *SpeakerArrangement) AsValue() Value {
	return Value(uintptr(unsafe.Pointer(sa)))
}

func (sa *SpeakerArrangement) AsPtr() Ptr {
	return Ptr(unsafe.Pointer(sa))
}
