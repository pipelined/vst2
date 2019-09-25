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

	doubleIn  api.DoubleBuffer
	doubleOut api.DoubleBuffer
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
	var out signal.Float64
	return func(b [][]float64) ([][]float64, error) {
		// new buffer size.
		if currentSize != signal.Float64(b).Size() {
			currentSize = signal.Float64(b).Size()
			p.plugin.e.SetBufferSize(currentSize)

			// reset buffers.
			p.doubleIn.Free()
			p.doubleOut.Free()
			p.doubleIn = api.NewDoubleBuffer(numChannels, currentSize)
			p.doubleOut = api.NewDoubleBuffer(numChannels, currentSize)
			out = signal.Float64Buffer(numChannels, currentSize, 0)
		}
		p.plugin.e.ProcessDouble(p.doubleIn, p.doubleOut)
		p.currentPosition += int64(signal.Float64(b).Size())

		api.CopyDouble(p.doubleOut, out)
		return out, nil
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
