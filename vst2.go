// Package vst2 provides interface to VST2 plugins.
package vst2

import (
	"fmt"

	"pipelined.dev/audio/vst2/sdk"
)

type (
	// VST used to create new instances of plugin.
	// It also keeps reference to VST handle to clean up on Close.
	VST sdk.EntryPoint
)

// Open loads the VST into memory and stores entry point func.
func Open(path string) (*VST, error) {
	ep, err := sdk.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load VST '%s': %w", path, err)
	}

	return (*VST)(ep), nil
}

// Close cleans up VST resources.
func (v *VST) Close() error {
	return (*sdk.EntryPoint)(v).Close()
}

// Load new instance of VST plugin with provided callback.
// This function also calls dispatch with EffOpen opcode.
// func (v VST) Load(c sdk.HostCallbackFunc) *Plugin {
// 	e := v.EntryPoint.Load(c)
// 	e.Dispatch(sdk.EffOpen, 0, 0, nil, 0.0)

// 	numParams := e.NumParams()
// 	params := make([]Parameter, numParams)
// 	for i := 0; i < numParams; i++ {
// 		params = append(params, Parameter{
// 			name:       e.ParamName(i),
// 			unit:       e.ParamUnitName(i),
// 			value:      e.ParamValue(i),
// 			valueLabel: e.ParamValueName(i),
// 		})
// 	}
// 	numPresets := e.NumPrograms()
// 	presets := make([]Program, numPresets)
// 	for i := 0; i < numPresets; i++ {
// 		presets = append(presets, Program{
// 			name: e.ProgramName(i),
// 		})
// 	}
// 	return &Plugin{
// 		Effect:     e,
// 		Parameters: params,
// 		Programs:   presets,
// 	}
// }
