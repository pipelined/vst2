package vst2

import (
	"log"
	"math"
	"time"
	"unsafe"

	"github.com/dudk/phono/pipe/runner"

	"github.com/dudk/phono"
	"github.com/dudk/phono/pipe"
)

// Processor represents vst2 sound processor
type Processor struct {
	phono.UID
	plugin *Plugin

	bufferSize    phono.BufferSize
	numChannels   phono.NumChannels
	sampleRate    phono.SampleRate
	tempo         phono.Tempo
	timeSignature phono.TimeSignature

	currentPosition int64
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

// RunProcess returns configured processor runner
func (p *Processor) RunProcess() pipe.ProcessRunner {
	return &runner.Process{
		Processor: p,
		Before: func() error {
			p.plugin.SetCallback(p.callback())
			p.plugin.SetBufferSize(p.bufferSize)
			p.plugin.SetSampleRate(p.sampleRate)
			p.plugin.SetSpeakerArrangement(p.numChannels)
			p.plugin.Resume()
			return nil
		},
		After: func() error {
			p.plugin.Suspend()
			return nil
		},
	}
}

// Process buffer
func (p *Processor) Process(buf phono.Buffer) (phono.Buffer, error) {
	buf = p.plugin.Process(buf)
	p.currentPosition += int64(p.bufferSize)
	return buf, nil
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
