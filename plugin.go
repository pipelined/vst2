package vst2

import (
	"pipelined.dev/audio/vst2/sdk"
)

// Plugin is VST2 plugin instance.
type Plugin struct {
	*sdk.Effect
	Parameters []Parameter
	Programs   []Program
}

type Parameter struct {
	name       string
	unit       string
	value      float32
	valueLabel string
}

type Program struct {
	name string
}
