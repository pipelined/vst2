package vst2

import (
	"context"
	"time"

	"pipelined.dev/pipe"
	"pipelined.dev/pipe/mutable"
	"pipelined.dev/signal"
)

type (
	// Processor is pipe component that wraps.
	Processor struct {
		host       *HostProperties
		Effect     *Effect
		Parameters []Parameter
		Presets    []Preset
	}

	Parameter struct {
		name       string
		unit       string
		value      float32
		valueLabel string
	}

	Preset struct {
		name string
	}

	// HostProperties contains values required to handle plugin-to-host
	// callbacks. It must be modified in the processing goroutine, otherwise
	// race condition might happen.
	HostProperties struct {
		BufferSize      int
		Channels        int
		SampleRate      signal.Frequency
		CurrentPosition int64
	}

	// HostCallbackAllocator returns new host callback function that uses host
	// properties to interact with the plugin.
	HostCallbackAllocator func(*HostProperties) HostCallbackFunc

	// ProcessorInitFunc applies configuration on plugin before starting it
	// in the processor routine.
	ProcessorInitFunc func(*Effect)
)

// Processor represents vst2 sound processor.
func (v *VST) Processor(callback HostCallbackAllocator) Processor {
	host := &HostProperties{}
	e := v.Plugin(callback(host))
	numParams := e.NumParams()
	params := make([]Parameter, numParams)
	for i := 0; i < numParams; i++ {
		params = append(params, Parameter{
			name:       e.ParamName(i),
			unit:       e.ParamUnitName(i),
			value:      e.ParamValue(i),
			valueLabel: e.ParamValueName(i),
		})
	}
	numPresets := e.NumPrograms()
	presets := make([]Preset, numPresets)
	for i := 0; i < numPresets; i++ {
		presets = append(presets, Preset{
			name: e.ProgramName(i),
		})
	}
	return Processor{
		Effect:     e,
		host:       host,
		Parameters: params,
		Presets:    presets,
	}
}

// Allocator returns pipe processor allocator that can be plugged into line.
func (p *Processor) Allocator(init ProcessorInitFunc) pipe.ProcessorAllocatorFunc {
	return func(mctx mutable.Context, bufferSize int, props pipe.SignalProperties) (pipe.Processor, error) {
		p.host.BufferSize = bufferSize
		p.host.Channels = props.Channels
		p.host.SampleRate = props.SampleRate
		p.Effect.Start()
		p.Effect.SetSampleRate(int(props.SampleRate))
		p.Effect.SetBufferSize(bufferSize)
		if init != nil {
			init(p.Effect)
		}
		processFn, flushFn := processorFns(p.Effect, p.host)
		return pipe.Processor{
			Output: pipe.SignalProperties{
				Channels:   props.Channels,
				SampleRate: props.SampleRate,
			},
			StartFunc: func(context.Context) error {
				p.Effect.Resume()
				return nil
			},
			ProcessFunc: processFn,
			FlushFunc:   flushFn,
		}, nil
	}
}

func processorFns(e *Effect, host *HostProperties) (pipe.ProcessFunc, pipe.FlushFunc) {
	if e.CanProcessFloat64() {
		return doubleFns(e, host)
	}
	return floatFns(e, host)
}

func doubleFns(e *Effect, host *HostProperties) (pipe.ProcessFunc, pipe.FlushFunc) {
	doubleIn := NewDoubleBuffer(host.Channels, host.BufferSize)
	doubleOut := NewDoubleBuffer(host.Channels, host.BufferSize)
	return func(in, out signal.Floating) error {
			doubleIn.CopyFrom(in)
			e.ProcessDouble(doubleIn, doubleOut)
			host.CurrentPosition += int64(in.Length())
			doubleOut.CopyTo(out)
			return nil
		},
		func(context.Context) error {
			doubleIn.Free()
			doubleOut.Free()
			e.Suspend()
			return nil
		}
}

func floatFns(e *Effect, host *HostProperties) (pipe.ProcessFunc, pipe.FlushFunc) {
	floatIn := NewFloatBuffer(host.Channels, host.BufferSize)
	floatOut := NewFloatBuffer(host.Channels, host.BufferSize)
	return func(in, out signal.Floating) error {
			floatIn.CopyFrom(in)
			e.ProcessFloat(floatIn, floatOut)
			host.CurrentPosition += int64(in.Length())
			floatOut.CopyTo(out)
			return nil
		},
		func(context.Context) error {
			floatIn.Free()
			floatOut.Free()
			e.Suspend()
			return nil
		}
}

// DefaultHostCallback returns default vst2 host callback.
func DefaultHostCallback(props *HostProperties) HostCallbackFunc {
	return func(opcode HostOpcode, index Index, value Value, ptr Ptr, opt Opt) Return {
		switch opcode {
		case HostGetCurrentProcessLevel:
			return Return(ProcessLevelRealtime)
		case HostGetSampleRate:
			return Return(props.SampleRate)
		case HostGetBlockSize:
			return Return(props.SampleRate)
		case HostGetTime:
			ti := &TimeInfo{
				SampleRate:         float64(props.SampleRate),
				SamplePos:          float64(props.CurrentPosition),
				NanoSeconds:        float64(time.Now().UnixNano()),
				TimeSigNumerator:   4,
				TimeSigDenominator: 4,
			}
			return ti.Return()
		default:
			break
		}
		return 0
	}
}
