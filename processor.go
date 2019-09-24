package vst2

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/pipelined/signal"
	"github.com/pipelined/vst2/api"
)

// Processor represents vst2 sound processor
type Processor struct {
	*VST
	plugin *Plugin

	bufferSize  int
	numChannels int
	sampleRate  int

	currentPosition int64
}

// Process returns processor function with default settings initialized.
func (p *Processor) Process(pipeID string, sampleRate, numChannels int) (func([][]float64) ([][]float64, error), error) {
	p.sampleRate = sampleRate
	p.numChannels = numChannels
	p.plugin = p.VST.Load(p.callback())

	p.plugin.e.SetSampleRate(p.sampleRate)
	p.plugin.e.SetSpeakerArrangement(newSpeakerArrangement(p.numChannels), newSpeakerArrangement(p.numChannels))
	p.plugin.e.Start()
	var currentSize int
	return func(b [][]float64) ([][]float64, error) {
		if bufferSize := signal.Float64(b).Size(); currentSize != bufferSize {
			p.plugin.e.SetBufferSize(p.bufferSize)
			currentSize = bufferSize
		}
		b = p.plugin.Process(b)
		p.currentPosition += int64(signal.Float64(b).Size())
		return b, nil
	}, nil
}

// Flush suspends plugin.
func (p *Processor) Flush(string) error {
	p.plugin.e.Stop()
	return nil
}

// wraped callback with session.
func (p *Processor) callback() api.HostCallbackFunc {
	return func(e *api.Effect, opcode api.HostOpcode, index api.Index, value api.Value, ptr api.Ptr, opt api.Opt) api.Return {
		fmt.Printf("Callback: %v\n", opcode)
		switch opcode {
		case api.HostIdle:
			p.plugin.e.Dispatch(api.EffEditIdle, 0, 0, nil, 0)
		case api.HostGetCurrentProcessLevel:
			return api.Return(api.ProcessLevelRealtime)
		case api.HostGetSampleRate:
			return api.Return(p.sampleRate)
		case api.HostGetBlockSize:
			return api.Return(p.bufferSize)
		case api.HostGetTime:
			nanoseconds := time.Now().UnixNano()

			return api.Return(uintptr(unsafe.Pointer(&api.TimeInfo{
				SampleRate:         float64(p.sampleRate),
				SamplePos:          float64(p.currentPosition),
				NanoSeconds:        float64(nanoseconds),
				TimeSigNumerator:   4,
				TimeSigDenominator: 4,
			})))
		default:
			// log.Printf("Plugin requested value of opcode %v\n", opcode)
			break
		}
		return 0
	}
}
