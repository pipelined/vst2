package vst2

import (
	"context"
	"time"

	"pipelined.dev/audio/vst2/sdk"
	"pipelined.dev/pipe"
	"pipelined.dev/pipe/mutable"
	"pipelined.dev/signal"
)

type (
	// Processor is pipe component that wraps.
	Processor struct {
		Effect     *sdk.Effect
		Parameters []Parameter
		Programs   []Program
	}

	Parameter struct {
		name       string
		unit       string
		value      float32
		valueLabel string
	}

	Program struct {
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
	HostCallbackAllocator func(*HostProperties) sdk.HostCallbackFunc

	// ProcessorInitFunc applies configuration on plugin before starting it
	// in the processor routine.
	ProcessorInitFunc func(*sdk.Effect)
)

// Processor represents vst2 sound processor. It loads plugin from the
// provided vst and applies the following configuration: set up sample
// rate, set up buffer size, calls provided init function and then starts
// the plugin.
func (v *VST) Processor(callback HostCallbackAllocator, init ProcessorInitFunc) pipe.ProcessorAllocatorFunc {
	return func(mctx mutable.Context, bufferSize int, props pipe.SignalProperties) (pipe.Processor, error) {
		host := HostProperties{
			BufferSize: bufferSize,
			Channels:   props.Channels,
			SampleRate: props.SampleRate,
		}
		e := (*sdk.EntryPoint)(v).Load(callback(&host))
		e.SetSampleRate(int(props.SampleRate))
		e.SetBufferSize(bufferSize)
		if init != nil {
			init(e)
		}
		e.Start()

		p := processor(e, &host)
		p.Output = pipe.SignalProperties{
			Channels:   props.Channels,
			SampleRate: props.SampleRate,
		}
		return p, nil
	}
}

func processor(e *sdk.Effect, host *HostProperties) pipe.Processor {
	if e.CanProcessFloat64() {
		return doubleProcessor(e, host)
	}
	return floatProcessor(e, host)
}

func doubleProcessor(e *sdk.Effect, host *HostProperties) pipe.Processor {
	doubleIn := sdk.NewDoubleBuffer(host.Channels, host.BufferSize)
	doubleOut := sdk.NewDoubleBuffer(host.Channels, host.BufferSize)
	return pipe.Processor{
		ProcessFunc: func(in, out signal.Floating) error {
			doubleIn.CopyFrom(in)
			e.ProcessDouble(doubleIn, doubleOut)
			host.CurrentPosition += int64(in.Length())
			doubleOut.CopyTo(out)
			return nil
		},
		FlushFunc: func(context.Context) error {
			doubleIn.Free()
			doubleOut.Free()
			e.Stop()
			return nil
		},
	}
}

func floatProcessor(e *sdk.Effect, host *HostProperties) pipe.Processor {
	floatIn := sdk.NewFloatBuffer(host.Channels, host.BufferSize)
	floatOut := sdk.NewFloatBuffer(host.Channels, host.BufferSize)
	return pipe.Processor{
		ProcessFunc: func(in, out signal.Floating) error {
			floatIn.CopyFrom(in)
			e.ProcessFloat(floatIn, floatOut)
			host.CurrentPosition += int64(in.Length())
			floatOut.CopyTo(out)
			return nil
		},
		FlushFunc: func(context.Context) error {
			floatIn.Free()
			floatOut.Free()
			e.Stop()
			return nil
		},
	}
}

// DefaultHostCallback returns default vst2 host callback.
func DefaultHostCallback(props *HostProperties) sdk.HostCallbackFunc {
	return func(opcode sdk.HostOpcode, index sdk.Index, value sdk.Value, ptr sdk.Ptr, opt sdk.Opt) sdk.Return {
		switch opcode {
		case sdk.HostGetCurrentProcessLevel:
			return sdk.Return(sdk.ProcessLevelRealtime)
		case sdk.HostGetSampleRate:
			return sdk.Return(props.SampleRate)
		case sdk.HostGetBlockSize:
			return sdk.Return(props.SampleRate)
		case sdk.HostGetTime:
			ti := &sdk.TimeInfo{
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
