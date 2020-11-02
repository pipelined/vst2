package sdk

/*
#cgo CFLAGS: -std=gnu99 -I${SRCDIR}/sdk
#include "vst.h"
*/
import "C"
import "unsafe"

const (
	// VST main function name.
	main = "VSTPluginMain"
	// VST API version.
	version = 2400
)

type (
	// Effect is an alias on C effect type.
	Effect C.Effect

	// effectMain is a reference to VST main function.
	// wrapper on C entry point.
	effectMain C.EntryPoint

	// EntryPoint is a reference to VST main function.
	// wrapper on C entry point.
	EntryPoint struct {
		main   effectMain
		handle uintptr
		Name   string
	}
)

// Load new instance of VST plugin with provided callback.
// This function also calls dispatch with EffOpen opcode.
func (m *EntryPoint) Load(c HostCallbackFunc) *Effect {
	if m.main == nil || c == nil {
		return nil
	}
	e := (*Effect)(C.loadEffect(m.main))
	mutex.Lock()
	callbacks[e] = c
	mutex.Unlock()

	return e
}

// Close cleans up C refs for plugin
func (e *Effect) Close() error {
	if e == nil {
		return nil
	}
	mutex.Lock()
	delete(callbacks, e)
	mutex.Unlock()
	return nil
}

// Dispatch wraps-up C method to dispatch calls to plugin
func (e *Effect) Dispatch(opcode EffectOpcode, index Index, value Value, ptr Ptr, opt Opt) Return {
	return Return(C.dispatch((*C.Effect)(e), C.int32_t(opcode), C.int32_t(index), C.int64_t(value), unsafe.Pointer(ptr), C.float(opt)))
}

// ScanPaths returns a slice of default vst2 locations.
// Locations are OS-specific.
func ScanPaths() []string {
	paths := make([]string, 0, len(scanPaths))
	paths = append(paths, scanPaths...)
	return paths
}

// NumParams returns the number of parameters.
func (e *Effect) NumParams() int {
	return int(e.numParams)
}

// NumPrograms returns the number of programs.
func (e *Effect) NumPrograms() int {
	return int(e.numPrograms)
}

// Flags returns the effect flags.
func (e *Effect) Flags() EffectFlags {
	return EffectFlags(e.flags)
}

// ProcessDouble audio with VST plugin.
func (e *Effect) ProcessDouble(in, out DoubleBuffer) {
	C.processDouble(
		(*C.Effect)(e),
		C.int32_t(in.numChannels),
		C.int32_t(in.size),
		&in.data[0],
		&out.data[0],
	)
}

// ProcessFloat audio with VST plugin.
func (e *Effect) ProcessFloat(in, out FloatBuffer) {
	C.processFloat(
		(*C.Effect)(e),
		C.int32_t(in.numChannels),
		C.int32_t(in.size),
		&in.data[0],
		&out.data[0],
	)
}

// ParamValue returns the value of parameter.
func (e *Effect) ParamValue(index int) float32 {
	return float32(C.getParameter((*C.Effect)(e), C.int32_t(index)))
}

// SetParamValue sets new value for parameter.
func (e *Effect) SetParamValue(index int, value float32) {
	C.setParameter((*C.Effect)(e), C.int32_t(index), C.float(value))
}
