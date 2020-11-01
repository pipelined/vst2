// Package vst2 provides interface to VST2 plugins.
package vst2

import (
	"fmt"
	"path/filepath"

	"pipelined.dev/audio/vst2/sdk"
)

type (
	// VST used to create new instances of plugin.
	// It also keeps reference to VST handle to clean up on Close.
	VST struct {
		entryPoint *sdk.EntryPoint
		Name       string
		Path       string
	}
)

// Open loads the VST into memory and stores entry point func.
func Open(path string) (VST, error) {
	p, err := filepath.Abs(path)
	if err != nil {
		return VST{}, fmt.Errorf("failed to get absolute path: %w", err)
	}

	ep, err := sdk.Open(path)
	if err != nil {
		return VST{}, fmt.Errorf("failed to load VST '%s': %w", path, err)
	}

	return VST{
		entryPoint: ep,
		Path:       p,
	}, nil
}

// Close cleans up VST resources.
func (v VST) Close() error {
	return v.entryPoint.Close()
}

// Load new instance of VST plugin with provided callback.
// This function also calls dispatch with EffOpen opcode.
func (v VST) Load(c sdk.HostCallbackFunc) *Plugin {
	e := v.entryPoint.Load(c)
	e.Dispatch(sdk.EffOpen, 0, 0, nil, 0.0)
	return &Plugin{
		Effect: e,
	}
}
