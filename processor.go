//go:build !plugin
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
		progressFn ProgressProcessedFunc
	}

	// ProcessorInitFunc applies configuration on plugin before starting it
	// in the processor routine.
	ProcessorInitFunc func(*Plugin)

	// HostProgressProcessed is executed by processor after every process
	// call.
	ProgressProcessedFunc func(int)
)

// Processor represents vst2 sound processor. Processor always overrides
// GetBufferSize and GetSampleRate callbacks, because this vaules are
// injected when processor is allocated by pipe.
func (v *VST) Processor(h Host, progressFn ProgressProcessedFunc) *Processor {
	processor := Processor{}
	h.GetBufferSize = func() int {
		return processor.bufferSize
	}
	h.GetSampleRate = func() signal.Frequency {
		return processor.sampleRate
	}
	plugin := v.Plugin(h.Callback())
	return &Processor{
		plugin:     plugin,
		progressFn: progressFn,
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

func processorFns(p *Plugin, channels, bufferSize int, progressFn ProgressProcessedFunc) (pipe.ProcessFunc, pipe.FlushFunc) {
	if p.CanProcessFloat64() {
		return doubleFns(p, channels, bufferSize, progressFn)
	}
	return floatFns(p, channels, bufferSize, progressFn)
}

func doubleFns(p *Plugin, channels, bufferSize int, progressFn ProgressProcessedFunc) (pipe.ProcessFunc, pipe.FlushFunc) {
	doubleIn := NewDoubleBuffer(channels, bufferSize)
	doubleOut := NewDoubleBuffer(channels, bufferSize)
	processFn := func(in, out signal.Floating) (int, error) {
		doubleIn.Write(in)
		p.ProcessDouble(doubleIn, doubleOut)
		doubleOut.Read(out)
		return in.Length(), nil
	}
	if progressFn != nil {
		processFn = func(in, out signal.Floating) (int, error) {
			doubleIn.Write(in)
			p.ProcessDouble(doubleIn, doubleOut)
			doubleOut.Read(out)
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

func floatFns(p *Plugin, channels, bufferSize int, progressFn ProgressProcessedFunc) (pipe.ProcessFunc, pipe.FlushFunc) {
	floatIn := NewFloatBuffer(channels, bufferSize)
	floatOut := NewFloatBuffer(channels, bufferSize)
	processFn := func(in, out signal.Floating) (int, error) {
		floatIn.Write(in)
		p.ProcessFloat(floatIn, floatOut)
		floatOut.Read(out)
		return in.Length(), nil
	}
	if progressFn != nil {
		processFn = func(in, out signal.Floating) (int, error) {
			floatIn.Write(in)
			p.ProcessFloat(floatIn, floatOut)
			floatOut.Read(out)
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
