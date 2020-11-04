package vst2

import (
	"context"
	"time"
	"unsafe"

	"pipelined.dev/pipe"
	"pipelined.dev/pipe/mutable"
	"pipelined.dev/signal"
)

type (
	// Processor is pipe component that wraps.
	Processor struct {
		host       *HostProperties
		Plugin     *Plugin
		Parameters []Parameter
		Presets    []Preset
	}

	// Parameter refers to plugin parameter that can be mutated in the pipe.
	Parameter struct {
		name       string
		unit       string
		value      float32
		valueLabel string
	}

	// Preset refers to plugin presets.
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
	ProcessorInitFunc func(*Plugin)
)

// Processor represents vst2 sound processor.
func (v *VST) Processor(callback HostCallbackAllocator) Processor {
	host := &HostProperties{}
	p := v.Plugin(callback(host))
	numParams := p.NumParams()
	params := make([]Parameter, numParams)
	for i := 0; i < numParams; i++ {
		params = append(params, Parameter{
			name:       p.ParamName(i),
			unit:       p.ParamUnitName(i),
			value:      p.ParamValue(i),
			valueLabel: p.ParamValueName(i),
		})
	}
	numPresets := p.NumPrograms()
	presets := make([]Preset, numPresets)
	for i := 0; i < numPresets; i++ {
		presets = append(presets, Preset{
			name: p.ProgramName(i),
		})
	}
	return Processor{
		Plugin:     p,
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
		p.Plugin.Start()
		p.Plugin.SetSampleRate(int(props.SampleRate))
		p.Plugin.SetBufferSize(bufferSize)
		if init != nil {
			init(p.Plugin)
		}
		processFn, flushFn := processorFns(p.Plugin, p.host)
		return pipe.Processor{
			Output: pipe.SignalProperties{
				Channels:   props.Channels,
				SampleRate: props.SampleRate,
			},
			StartFunc: func(context.Context) error {
				p.Plugin.Resume()
				return nil
			},
			ProcessFunc: processFn,
			FlushFunc:   flushFn,
		}, nil
	}
}

func processorFns(p *Plugin, host *HostProperties) (pipe.ProcessFunc, pipe.FlushFunc) {
	if p.CanProcessFloat64() {
		return doubleFns(p, host)
	}
	return floatFns(p, host)
}

func doubleFns(p *Plugin, host *HostProperties) (pipe.ProcessFunc, pipe.FlushFunc) {
	doubleIn := NewDoubleBuffer(host.Channels, host.BufferSize)
	doubleOut := NewDoubleBuffer(host.Channels, host.BufferSize)
	return func(in, out signal.Floating) error {
			doubleIn.CopyFrom(in)
			p.ProcessDouble(doubleIn, doubleOut)
			host.CurrentPosition += int64(in.Length())
			doubleOut.CopyTo(out)
			return nil
		},
		func(context.Context) error {
			doubleIn.Free()
			doubleOut.Free()
			p.Suspend()
			return nil
		}
}

func floatFns(p *Plugin, host *HostProperties) (pipe.ProcessFunc, pipe.FlushFunc) {
	floatIn := NewFloatBuffer(host.Channels, host.BufferSize)
	floatOut := NewFloatBuffer(host.Channels, host.BufferSize)
	return func(in, out signal.Floating) error {
			floatIn.CopyFrom(in)
			p.ProcessFloat(floatIn, floatOut)
			host.CurrentPosition += int64(in.Length())
			floatOut.CopyTo(out)
			return nil
		},
		func(context.Context) error {
			floatIn.Free()
			floatOut.Free()
			p.Suspend()
			return nil
		}
}

// DefaultHostCallback returns default vst2 host callback.
func DefaultHostCallback(props *HostProperties) HostCallbackFunc {
	return func(opcode HostOpcode, index Index, value Value, ptr unsafe.Pointer, opt Opt) Return {
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
