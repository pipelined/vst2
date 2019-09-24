package vst2

import (
	"fmt"
	"path/filepath"

	"github.com/pipelined/vst2/api"
)

// ScanPaths returns a slice of default vst2 locations.
// Locations are OS-specific.
func ScanPaths() (paths []string) {
	return append([]string{}, scanPaths...)
}

// VST used to create new instances of plugin.
// TODO: make a list of references to opened plugins and close them when VST is closed.
type (
	VST struct {
		entryPoint api.EntryPoint
		Name       string
		Path       string
	}

	// Plugin type provides interface
	Plugin struct {
		e    *api.Effect
		Name string
		Path string
	}
)

// Open loads the VST into memory and stores entry point func.
func Open(path string) (*VST, error) {
	p, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	ep, err := api.Open(p)
	if err != nil {
		return nil, fmt.Errorf("failed to load VST '%s': %w", path, err)
	}

	return &VST{
		Path:       p,
		entryPoint: ep,
	}, nil
}

// Close cleans up VST resoures.
func (v *VST) Close() error {
	if err := v.entryPoint.Close(); err != nil {
		return fmt.Errorf("failed close VST %s: %w", v.Name, err)
	}
	return nil
}

// Load creates new instance of plugin.
func (v *VST) Load(c api.HostCallbackFunc) *Plugin {
	return &Plugin{
		e:    v.entryPoint.Load(c),
		Path: v.Path,
		Name: v.Name,
	}
}

// Close cleans up C refs for plugin
func (p *Plugin) Close() error {
	p.e.Dispatch(api.EffClose, 0, 0, nil, 0.0)
	p.e = nil
	return nil
}

// Process is a wrapper over ProcessFloat64 and ProcessFloat32
// in case if plugin supports only ProcessFloat32, conversion is done
func (p *Plugin) Process(buffer [][]float64) (result [][]float64) {
	if buffer == nil || len(buffer) == 0 || buffer[0] == nil {
		return
	}

	if p.e.CanProcessFloat32() {

		in32 := make([][]float32, len(buffer))
		for i := range buffer {
			in32[i] = make([]float32, len(buffer[0]))
			for j, v := range buffer[i] {
				in32[i][j] = float32(v)
			}
		}

		out32 := p.e.ProcessFloat32(in32)

		result = make([][]float64, len(out32))
		for i := range out32 {
			result[i] = make([]float64, len(out32[i]))
			for j, v := range out32[i] {
				result[i][j] = float64(v)
			}
		}
	} else {
		result = p.e.ProcessFloat64([][]float64(buffer))
	}
	return
}

func newSpeakerArrangement(numChannels int) *api.SpeakerArrangement {
	sa := api.SpeakerArrangement{}
	sa.NumChannels = int32(numChannels)
	switch numChannels {
	case 0:
		sa.Type = api.SpeakerArrEmpty
	case 1:
		sa.Type = api.SpeakerArrMono
	case 2:
		sa.Type = api.SpeakerArrStereo
	case 3:
		sa.Type = api.SpeakerArr30Music
	case 4:
		sa.Type = api.SpeakerArr40Music
	case 5:
		sa.Type = api.SpeakerArr50
	case 6:
		sa.Type = api.SpeakerArr60Music
	case 7:
		sa.Type = api.SpeakerArr70Music
	case 8:
		sa.Type = api.SpeakerArr80Music
	default:
		sa.Type = api.SpeakerArrUserDefined
	}

	for i := 0; i < int(numChannels); i++ {
		sa.Speakers[i].Type = api.SpeakerUndefined
	}
	return &sa
}

// DefaultHostCallback is a default callback, just prints incoming opcodes should be overridden with SetHostCallback
func DefaultHostCallback(print bool) api.HostCallbackFunc {
	return func(e *api.Effect, opcode api.HostOpcode, index api.Index, value api.Value, ptr api.Ptr, opt api.Opt) api.Return {
		fmt.Printf("Callback called with opcode: %v\n", opcode)
		switch opcode {
		case api.HostVersion:
			return 2400
		default:
			break
		}
		return 0
	}
}
