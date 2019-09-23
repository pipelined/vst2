package vst2

import (
	"fmt"
	"math"
	"time"

	"github.com/pipelined/signal"
	"github.com/pipelined/vst2/api"
)

// Processor represents vst2 sound processor
type Processor struct {
	*Vst
	Plugin *Plugin

	bufferSize    int
	numChannels   int
	sampleRate    int
	tempo         float64
	timeSignature TimeSignature

	currentPosition int64
}

// Process returns processor function with default settings initialized.
func (p *Processor) Process(pipeID string, sampleRate, numChannels int) (func([][]float64) ([][]float64, error), error) {
	p.sampleRate = sampleRate
	p.numChannels = numChannels
	plugin, err := p.Vst.Open(p.callback())
	if err != nil {
		return nil, err
	}
	p.Plugin = plugin

	// p.Plugin.SetCallback(p.callback())
	p.Plugin.SetSampleRate(p.sampleRate)
	p.Plugin.SetSpeakerArrangement(p.numChannels)
	p.Plugin.Resume()
	var currentSize int
	return func(b [][]float64) ([][]float64, error) {
		if bufferSize := signal.Float64(b).Size(); currentSize != bufferSize {
			p.Plugin.SetBufferSize(p.bufferSize)
			currentSize = bufferSize
		}
		b = p.Plugin.Process(b)
		p.currentPosition += int64(signal.Float64(b).Size())
		return b, nil
	}, nil
}

// Flush suspends plugin.
func (p *Processor) Flush(string) error {
	p.Plugin.Suspend()
	return nil
}

// wraped callback with session.
func (p *Processor) callback() api.HostCallbackFunc {
	return func(e *api.Effect, opcode api.HostOpcode, index api.Index, value api.Value, ptr api.Ptr, opt api.Opt) int {
		fmt.Printf("Callback: %v\n", opcode)
		switch opcode {
		case api.HostIdle:
			p.Plugin.e.Dispatch(api.EffEditIdle, 0, 0, nil, 0)
		case api.HostGetCurrentProcessLevel:
			return int(api.ProcessLevelRealtime)
		case api.HostGetSampleRate:
			return int(p.sampleRate)
		case api.HostGetBlockSize:
			return int(p.bufferSize)
		case api.HostGetTime:
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

			return int(p.Plugin.SetTimeInfo(int(p.sampleRate), samplePos, tempo, p.timeSignature, float64(nanoseconds), ppqPos, barPos))
		default:
			// log.Printf("Plugin requested value of opcode %v\n", opcode)
			break
		}
		return 0
	}
}
