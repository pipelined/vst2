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
	ProcessorInitFunc func(*Plugin)
)

// Processor represents vst2 sound processor. It loads plugin from the
// provided vst and applies the following configuration: set up sample
// rate, set up buffer size, calls provided init function and then starts
// the plugin.
func Processor(vst VST, callback HostCallbackAllocator, init ProcessorInitFunc) pipe.ProcessorAllocatorFunc {
	return func(mctx mutable.Context, bufferSize int, props pipe.SignalProperties) (pipe.Processor, error) {
		host := HostProperties{
			BufferSize: bufferSize,
			Channels:   props.Channels,
			SampleRate: props.SampleRate,
		}
		plugin := vst.Load(callback(&host))
		plugin.SetSampleRate(int(props.SampleRate))
		plugin.SetBufferSize(bufferSize)
		if init != nil {
			init(plugin)
		}
		plugin.Start()

		p := processor(plugin, &host)
		p.Output = pipe.SignalProperties{
			Channels:   props.Channels,
			SampleRate: props.SampleRate,
		}
		return p, nil
	}
}

func processor(plugin *Plugin, host *HostProperties) pipe.Processor {
	if plugin.CanProcessFloat64() {
		return doubleProcessor(plugin, host)
	}
	return floatProcessor(plugin, host)
}

func doubleProcessor(plugin *Plugin, host *HostProperties) pipe.Processor {
	doubleIn := sdk.NewDoubleBuffer(host.Channels, host.BufferSize)
	doubleOut := sdk.NewDoubleBuffer(host.Channels, host.BufferSize)
	return pipe.Processor{
		ProcessFunc: func(in, out signal.Floating) error {
			doubleIn.CopyFrom(in)
			plugin.ProcessDouble(doubleIn, doubleOut)
			host.CurrentPosition += int64(in.Length())
			doubleOut.CopyTo(out)
			return nil
		},
		FlushFunc: func(context.Context) error {
			doubleIn.Free()
			doubleOut.Free()
			plugin.Stop()
			return nil
		},
	}
}

func floatProcessor(plugin *Plugin, host *HostProperties) pipe.Processor {
	floatIn := sdk.NewFloatBuffer(host.Channels, host.BufferSize)
	floatOut := sdk.NewFloatBuffer(host.Channels, host.BufferSize)
	return pipe.Processor{
		ProcessFunc: func(in, out signal.Floating) error {
			floatIn.CopyFrom(in)
			plugin.ProcessFloat(floatIn, floatOut)
			host.CurrentPosition += int64(in.Length())
			floatOut.CopyTo(out)
			return nil
		},
		FlushFunc: func(context.Context) error {
			floatIn.Free()
			floatOut.Free()
			plugin.Stop()
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
