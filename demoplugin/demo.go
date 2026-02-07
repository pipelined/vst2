//go:build plugin
// +build plugin

package main

import (
	"fmt"
	"math"

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
			Value: 0.5,
			GetValueLabelFunc: func(value float32) string {
				return fmt.Sprintf("%+.2f", value)
			},
			GetValueFunc: func(value float32) float32 {
				return -20 + (40 * value)
			},
		}
		channels := 2
		return vst2.Plugin{
			UniqueID:       uniqueID,
			Version:        version,
			InputChannels:  channels,
			OutputChannels: channels,
			Name:           "Gain",
			Vendor:         "pipelined/vst2",
			Category:       vst2.PluginCategoryEffect,
			Parameters: []*vst2.Parameter{
				&gain,
			},
			ProcessDoubleFunc: func(in, out vst2.DoubleBuffer) {
				g := math.Pow(10, float64(gain.GetValue())/20)
				for c := 0; c < channels; c++ {
					for i := 0; i < in.Frames; i++ {
						out.Channel(c)[i] = in.Channel(c)[i] * g
					}
				}
			},
			ProcessFloatFunc: func(in, out vst2.FloatBuffer) {
				g := math.Pow(10, float64(gain.GetValue())/20)
				for c := 0; c < channels; c++ {
					for i := 0; i < in.Frames; i++ {
						out.Channel(c)[i] = in.Channel(c)[i] * float32(g)
					}
				}
			},
		}, vst2.Dispatcher{}
	}
}

func main() {}
