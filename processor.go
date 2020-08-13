package vst2

import (
	"context"
	"time"

	"pipelined.dev/pipe"
	"pipelined.dev/signal"
)

// Processor represents vst2 sound processor.
func Processor(vst VST, callback HostCallbackAllocator) pipe.ProcessorAllocatorFunc {
	return func(bufferSize int, props pipe.SignalProperties) (pipe.Processor, pipe.SignalProperties, error) {
		host := HostProperties{
			BufferSize: bufferSize,
			Channels:   props.Channels,
			SampleRate: props.SampleRate,
		}
		if callback == nil {
			callback = DefaultHostCallback
		}
		plugin := vst.Load(callback(&host))
		plugin.SetSampleRate(int(props.SampleRate))
		plugin.SetSpeakerArrangement(newSpeakerArrangement(props.Channels), newSpeakerArrangement(props.Channels))
		plugin.Start()

		return processor(plugin, &host),
			pipe.SignalProperties{
				Channels:   props.Channels,
				SampleRate: props.SampleRate,
			},
			nil
	}
}

func processor(plugin *Plugin, host *HostProperties) pipe.Processor {
	if plugin.CanProcessFloat64() {
		return doubleProcessor(plugin, host)
	}
	return floatProcessor(plugin, host)
}

func doubleProcessor(plugin *Plugin, host *HostProperties) pipe.Processor {
	doubleIn := NewDoubleBuffer(host.Channels, host.BufferSize)
	doubleOut := NewDoubleBuffer(host.Channels, host.BufferSize)
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
			doubleIn.Free()
			plugin.Stop()
			return nil
		},
	}
}

func floatProcessor(plugin *Plugin, host *HostProperties) pipe.Processor {
	floatIn := NewFloatBuffer(host.Channels, host.BufferSize)
	floatOut := NewFloatBuffer(host.Channels, host.BufferSize)
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
			floatIn.Free()
			plugin.Stop()
			return nil
		},
	}
}

// HostProperties contains values required to handle plugin-to-host
// callbacks. It must be modified in the processing goroutine, otherwise
// race condition might happen.
type HostProperties struct {
	BufferSize int
	Channels   int
	signal.SampleRate
	CurrentPosition int64
}

// HostCallbackAllocator returns new host callback function that uses host
// properties to interact with the plugin.
type HostCallbackAllocator func(*HostProperties) HostCallbackFunc

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
