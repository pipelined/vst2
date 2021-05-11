// +build plugin

package main

import (
	"pipelined.dev/audio/vst2"
)

func init() {
	var (
		uniqueID = [4]byte{'d', 'u', 'd', 'k'}
		version  = int32(1000)
	)
	vst2.PluginAllocator = func(h vst2.Host) (vst2.Plugin, vst2.Dispatcher) {
		gain := vst2.Parameter{
			Name:  "Gain",
			Unit:  "db",
			Value: 1,
		}
		channels := 2
		return vst2.Plugin{
			UniqueID:       uniqueID,
			Version:        version,
			InputChannels:  channels,
			OutputChannels: channels,
			Parameters: []*vst2.Parameter{
				&gain,
			},
			ProcessDoubleFunc: func(in, out vst2.DoubleBuffer) {
				for c := 0; c < channels; c++ {
					for i := 0; i < in.Frames; i++ {
						out.Channel(c)[i] = in.Channel(c)[i] * float64(gain.Value)
					}
				}
			},
			ProcessFloatFunc: func(in, out vst2.FloatBuffer) {
				for c := 0; c < channels; c++ {
					for i := 0; i < in.Frames; i++ {
						out.Channel(c)[i] = in.Channel(c)[i] * float32(gain.Value)
					}
				}
			},
		}, vst2.Dispatcher{}
	}
}

func main() {}
