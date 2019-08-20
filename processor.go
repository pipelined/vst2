package vst2

import (
	"log"
	"math"
	"time"
	"unsafe"

	"github.com/pipelined/signal"
)

// Processor represents vst2 sound processor
type Processor struct {
	Plugin *Plugin

	bufferSize    int
	numChannels   int
	sampleRate    int
	tempo         float32
	timeSignature TimeSignature

	currentPosition int64
}

// Process returns processor function with default settings initialized.
func (p *Processor) Process(pipeID string, sampleRate, numChannels int) (func([][]float64) ([][]float64, error), error) {
	p.sampleRate = sampleRate
	p.numChannels = numChannels
	p.Plugin.SetCallback(p.callback())
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

			return int(plugin.SetTimeInfo(int(p.sampleRate), samplePos, float32(tempo), p.timeSignature, nanoseconds, ppqPos, barPos))
		default:
			// log.Printf("Plugin requested value of opcode %v\n", opcode)
			break
		}
		return 0
	}
}
