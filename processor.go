// +build !plugin

package vst2

import (
	"context"

	"pipelined.dev/pipe"
	"pipelined.dev/pipe/mutable"
	"pipelined.dev/signal"
)

type (
	// Processor is pipe component that wraps.
	Processor struct {
		bufferSize int
		channels   int
		sampleRate signal.Frequency
		plugin     *Plugin
		progressFn HostProgressProcessed
		Parameters []Parameter
		Presets    []Preset
	}

	// ProcessorInitFunc applies configuration on plugin before starting it
	// in the processor routine.
	ProcessorInitFunc func(*Plugin)
)

// Processor represents vst2 sound processor. Processor always overrides
// GetBufferSize and GetSampleRate callbacks, because this vaules are
// injected when processor is allocated by pipe.
func (v *VST) Processor(h Host) *Processor {
	processor := Processor{}
	h.GetBufferSize = func() int {
		return processor.bufferSize
	}
	h.GetSampleRate = func() signal.Frequency {
		return processor.sampleRate
	}
	plugin := v.Plugin(h.Callback())
	numParams := plugin.NumParams()
	params := make([]Parameter, numParams)
	for i := 0; i < numParams; i++ {
		params = append(params, Parameter{
			Name:       plugin.ParamName(i),
			Unit:       plugin.ParamUnitName(i),
			ValueLabel: plugin.ParamValueName(i),
			value:      plugin.ParamValue(i),
		})
	}
	numPresets := plugin.NumPrograms()
	presets := make([]Preset, numPresets)
	for i := 0; i < numPresets; i++ {
		presets = append(presets, Preset{
			name: plugin.ProgramName(i),
		})
	}
	return &Processor{
		plugin:     plugin,
		progressFn: h.ProgressProcessed,
		Parameters: params,
		Presets:    presets,
	}
}

// Allocator returns pipe processor allocator that can be plugged into line.
func (p *Processor) Allocator(init ProcessorInitFunc) pipe.ProcessorAllocatorFunc {
	return func(mctx mutable.Context, bufferSize int, props pipe.SignalProperties) (pipe.Processor, error) {
		p.bufferSize = bufferSize
		p.channels = props.Channels
		p.sampleRate = props.SampleRate
		p.plugin.Start()
		p.plugin.SetSampleRate(props.SampleRate)
		p.plugin.SetBufferSize(bufferSize)
		if init != nil {
			init(p.plugin)
		}
		processFn, flushFn := processorFns(p.plugin, p.channels, p.bufferSize, p.progressFn)
		return pipe.Processor{
			SignalProperties: pipe.SignalProperties{
				Channels:   p.channels,
				SampleRate: p.sampleRate,
			},
			StartFunc: func(context.Context) error {
				p.plugin.Resume()
				return nil
			},
			ProcessFunc: processFn,
			FlushFunc:   flushFn,
		}, nil
	}
}

func processorFns(p *Plugin, channels, bufferSize int, progressFn HostProgressProcessed) (pipe.ProcessFunc, pipe.FlushFunc) {
	if p.CanProcessFloat64() {
		return doubleFns(p, channels, bufferSize, progressFn)
	}
	return floatFns(p, channels, bufferSize, progressFn)
}

func doubleFns(p *Plugin, channels, bufferSize int, progressFn HostProgressProcessed) (pipe.ProcessFunc, pipe.FlushFunc) {
	doubleIn := NewDoubleBuffer(channels, bufferSize)
	doubleOut := NewDoubleBuffer(channels, bufferSize)
	processFn := func(in, out signal.Floating) (int, error) {
		doubleIn.CopyFrom(in)
		p.ProcessDouble(doubleIn, doubleOut)
		doubleOut.CopyTo(out)
		return in.Length(), nil
	}
	if progressFn != nil {
		processFn = func(in, out signal.Floating) (int, error) {
			doubleIn.CopyFrom(in)
			p.ProcessDouble(doubleIn, doubleOut)
			doubleOut.CopyTo(out)
			progressFn(in.Length())
			return in.Length(), nil
		}
	}
	return processFn,
		func(context.Context) error {
			doubleIn.Free()
			doubleOut.Free()
			p.Suspend()
			return nil
		}
}

func floatFns(p *Plugin, channels, bufferSize int, progressFn HostProgressProcessed) (pipe.ProcessFunc, pipe.FlushFunc) {
	floatIn := NewFloatBuffer(channels, bufferSize)
	floatOut := NewFloatBuffer(channels, bufferSize)
	processFn := func(in, out signal.Floating) (int, error) {
		floatIn.CopyFrom(in)
		p.ProcessFloat(floatIn, floatOut)
		floatOut.CopyTo(out)
		return in.Length(), nil
	}
	if progressFn != nil {
		processFn = func(in, out signal.Floating) (int, error) {
			floatIn.CopyFrom(in)
			p.ProcessFloat(floatIn, floatOut)
			floatOut.CopyTo(out)
			progressFn(in.Length())
			return in.Length(), nil
		}
	}
	return processFn,
		func(context.Context) error {
			floatIn.Free()
			floatOut.Free()
			p.Suspend()
			return nil
		}
}
