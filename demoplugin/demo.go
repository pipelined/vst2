// +build plugin

package main

import (
	"unsafe"

	"pipelined.dev/audio/vst2"
)

type Gain struct {
	Gain float64
}

func init() {
	vst2.PluginAllocator = func(vst2.HostCallback) vst2.Plugin {
		gain := vst2.Parameter{
			Name:  "Gain",
			Unit:  "db",
			Value: 1,
		}
		return vst2.Plugin{
			InputChannels:  2,
			OutputChannels: 2,
			Parameters: []*vst2.Parameter{
				&gain,
			},
			ProcessDoubleFunc: func(in, out vst2.DoubleBuffer) {
				in1 := in.Channel(0)
				in2 := in.Channel(1)
				out1 := out.Channel(0)
				out2 := out.Channel(1)
				for i := 0; i < in.Frames; i++ {
					out1[i] = in1[i] * float64(gain.Value)
					out2[i] = in2[i] * float64(gain.Value)
				}
			},
			ProcessFloatFunc: func(in, out vst2.FloatBuffer) {
				in1 := in.Channel(0)
				in2 := in.Channel(1)
				out1 := out.Channel(0)
				out2 := out.Channel(1)
				for i := 0; i < in.Frames; i++ {
					out1[i] = in1[i] * float32(gain.Value)
					out2[i] = in2[i] * float32(gain.Value)
				}
			},
			DispatchFunc: func(op vst2.PluginOpcode, index int32, value int64, ptr unsafe.Pointer, opt float32) int64 {
				return 0
			},
		}
	}
}

func main() {}
