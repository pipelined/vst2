package vst2

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"time"
	"unsafe"

	"github.com/dudk/phono"
)

// Processor represents vst2 sound processor
type Processor struct {
	plugin *Plugin

	bufferSize    phono.BufferSize
	numChannels   phono.NumChannels
	sampleRate    phono.SampleRate
	tempo         phono.Tempo
	timeSignature phono.TimeSignature

	// possibly mutex is required
	// pulse           phono.Pulse
	currentPosition phono.SamplePosition
}

// NewProcessor creates new vst2 processor
func NewProcessor(plugin *Plugin, bufferSize phono.BufferSize, sampleRate phono.SampleRate, numChannels phono.NumChannels) *Processor {
	return &Processor{
		plugin:          plugin,
		currentPosition: 0,
		bufferSize:      bufferSize,
		sampleRate:      sampleRate,
		numChannels:     numChannels,
	}
}

// Process implements processor.Processor
func (p *Processor) Process() phono.ProcessFunc {
	// p.pulse = pulse
	p.plugin.SetCallback(p.callback())
	return func(ctx context.Context, in <-chan phono.Message) (<-chan phono.Message, <-chan error, error) {
		errc := make(chan error, 1)
		out := make(chan phono.Message)
		go func() {
			defer close(out)
			defer close(errc)
			// pulse := p.pulse
			p.plugin.SetBufferSize(p.bufferSize)
			p.plugin.SetSampleRate(p.sampleRate)
			p.plugin.SetSpeakerArrangement(p.numChannels)
			p.plugin.Resume()
			defer p.plugin.Suspend()
			for in != nil {
				select {
				case m, ok := <-in:
					if !ok {
						in = nil
					} else {
						if m.Options != nil {
							m.Options.ApplyTo(p)
						}
						m.Samples = p.plugin.Process(m.Samples)
						// calculate new position and advance it after processing is done
						p.currentPosition += phono.SamplePosition(p.bufferSize)
						out <- m
					}
				case <-ctx.Done():
					return
				}
			}
		}()
		return out, errc, nil
	}
}

// Validate the processor
func (p *Processor) Validate() error {
	if p.bufferSize == 0 {
		return errors.New("Buffer size is not defined")
	}

	if p.sampleRate == 0 {
		return errors.New("Sample rate is not defined")
	}

	if p.numChannels == 0 {
		return errors.New("Number of channels is not defined")
	}
	return nil
}

// // readCurrentPosition reads a position with read lock
// func (p *Processor) readCurrentPosition() phono.SamplePosition {
// 	return p.currentPosition
// }

// // setPosition sets the position through mutex
// func (p *Processor) setCurrentPosition(pos phono.SamplePosition) {
// 	p.currentPosition = pos
// }

// DefaultScanPaths returns a slice of default vst2 locations
func DefaultScanPaths() (paths []string) {
	switch goos := runtime.GOOS; goos {
	case "darwin":
		paths = []string{
			"~/Library/Audio/Plug-Ins/VST",
			"/Library/Audio/Plug-Ins/VST",
		}
	case "windows":
		paths = []string{
			"C:\\Program Files (x86)\\Steinberg\\VSTPlugins",
			"C:\\Program Files\\Steinberg\\VSTPlugins ",
		}
		envVstPath := os.Getenv("VST_PATH")
		if len(envVstPath) > 0 {
			paths = append(paths, envVstPath)
		}
	}
	return
}

// FileExtension returns default vst2 extension
func FileExtension() string {
	switch os := runtime.GOOS; os {
	case "darwin":
		return ".vst"
	case "windows":
		return ".dll"
	default:
		return ".so"
	}
}

// Resume starts the plugin
func (p *Plugin) Resume() {
	p.Dispatch(EffMainsChanged, 0, 1, nil, 0.0)
}

// Suspend stops the plugin
func (p *Plugin) Suspend() {
	p.Dispatch(EffMainsChanged, 0, 0, nil, 0.0)
}

// SetBufferSize sets a buffer size
func (p *Plugin) SetBufferSize(bufferSize phono.BufferSize) {
	p.Dispatch(EffSetBlockSize, 0, int64(bufferSize), nil, 0.0)
}

// SetSampleRate sets a sample rate for plugin
func (p *Plugin) SetSampleRate(sampleRate phono.SampleRate) {
	p.Dispatch(EffSetSampleRate, 0, 0, nil, float64(sampleRate))
}

func (p *Plugin) defaultCallback() HostCallbackFunc {
	return func(plugin *Plugin, opcode MasterOpcode, index int64, value int64, ptr unsafe.Pointer, opt float64) int {
		fmt.Printf("Call from default callback! Plugin name: %v\n", p.Name)
		return 0
	}
}

// Process is a wrapper over ProcessFloat64 and ProcessFloat32
// in case if plugin supports only ProcessFloat32, conversion is done
func (p *Plugin) Process(in [][]float64) [][]float64 {
	if p.CanProcessFloat32() {

		in32 := make([][]float32, len(in))
		for i := range in {
			in32[i] = make([]float32, len(in[i]))
			for j, v := range in[i] {
				in32[i][j] = float32(v)
			}
		}

		out32 := p.ProcessFloat32(in32)

		out := make([][]float64, len(out32))
		for i := range out32 {
			out[i] = make([]float64, len(out32[i]))
			for j, v := range out32[i] {
				out[i][j] = float64(v)
			}
		}
		return out
	}
	return p.ProcessFloat64(in)
}

// wraped callback with session
func (p *Processor) callback() HostCallbackFunc {
	return func(plugin *Plugin, opcode MasterOpcode, index int64, value int64, ptr unsafe.Pointer, opt float64) int {
		switch opcode {
		case AudioMasterIdle:
			log.Printf("AudioMasterIdle")
			plugin.Dispatch(EffEditIdle, 0, 0, nil, 0)

		case AudioMasterGetCurrentProcessLevel:
			//TODO: return C.kVstProcessLevel
		case AudioMasterGetSampleRate:
			return int(p.sampleRate)
		case AudioMasterGetBlockSize:
			return int(p.bufferSize)
		case AudioMasterGetTime:
			nanoseconds := time.Now().UnixNano()
			notesPerMeasure := p.timeSignature.NotesPerBar
			//TODO: make this dynamic (handle time signature changes)
			// samples position
			samplePos := p.currentPosition
			// todo tempo
			tempo := p.tempo

			samplesPerBeat := (60.0 / float64(tempo)) * float64(p.sampleRate)
			// todo: ppqPos
			ppqPos := float64(samplePos)/samplesPerBeat + 1.0
			// todo: barPos
			barPos := math.Floor(ppqPos / float64(notesPerMeasure))

			return int(plugin.SetTimeInfo(p.sampleRate, samplePos, tempo, p.timeSignature, nanoseconds, ppqPos, barPos))
		default:
			// log.Printf("Plugin requested value of opcode %v\n", opcode)
			break
		}
		return 0
	}
}
