package vst2

/*
#include "aeffectx.h"
*/
import "C"
import "unsafe"

//export hostCallback
//calls real callback
func hostCallback(effect *C.AEffect, opcode int64, index int64, value int64, ptr unsafe.Pointer, opt float64) int {
	if callback == nil {
		panic("Host callback is not defined!")
	}

	return callback(&Plugin{effect: effect}, MasterOpcode(opcode), index, value, ptr, opt)
}
