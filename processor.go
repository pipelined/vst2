package vst2

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/pipelined/signal"
)

// Processor represents vst2 sound processor
type Processor struct {
	*VST
	plugin *Plugin

	bufferSize  int
	numChannels int
	sampleRate  int

	currentPosition int64

	doubleIn  DoubleBuffer
	doubleOut DoubleBuffer
}

// Process returns processor function with default settings initialized.
func (p *Processor) Process(pipeID string, sampleRate, numChannels int) (func([][]float64) ([][]float64, error), error) {
	p.sampleRate = sampleRate
	p.numChannels = numChannels
	p.plugin = p.VST.Load(p.callback())

	p.plugin.SetSampleRate(p.sampleRate)
	p.plugin.SetSpeakerArrangement(newSpeakerArrangement(p.numChannels), newSpeakerArrangement(p.numChannels))
	p.plugin.Start()
	var currentSize int
	var out signal.Float64
	return func(b [][]float64) ([][]float64, error) {
		// new buffer size.
		if currentSize != signal.Float64(b).Size() {
			currentSize = signal.Float64(b).Size()
			p.plugin.SetBufferSize(currentSize)

			// reset buffers.
			p.doubleIn.Free()
			p.doubleOut.Free()
			p.doubleIn = NewDoubleBuffer(numChannels, currentSize)
			p.doubleOut = NewDoubleBuffer(numChannels, currentSize)
			out = signal.Float64Buffer(numChannels, currentSize, 0)
		}
		p.plugin.ProcessDouble(p.doubleIn, p.doubleOut)
		p.currentPosition += int64(signal.Float64(b).Size())

		CopyDouble(p.doubleOut, out)
		return out, nil
	}, nil
}

// Flush suspends plugin.
func (p *Processor) Flush(string) error {
	p.plugin.Stop()
	return nil
}

// wraped callback with session.
func (p *Processor) callback() HostCallbackFunc {
	return func(e *Plugin, opcode HostOpcode, index Index, value Value, ptr Ptr, opt Opt) Return {
		fmt.Printf("Callback: %v\n", opcode)
		switch opcode {
		case HostIdle:
			p.plugin.Dispatch(EffEditIdle, 0, 0, nil, 0)
		case HostGetCurrentProcessLevel:
			return Return(ProcessLevelRealtime)
		case HostGetSampleRate:
			return Return(p.sampleRate)
		case HostGetBlockSize:
			return Return(p.bufferSize)
		case HostGetTime:
			nanoseconds := time.Now().UnixNano()

			return Return(uintptr(unsafe.Pointer(&TimeInfo{
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
