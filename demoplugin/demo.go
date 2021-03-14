package demoplugin

import (
	"fmt"

	"pipelined.dev/signal"
)

type GainPlugin struct {
	Gain float64
}

func (p *GainPlugin) Process(in, out signal.Floating) {
	// p.p.ProcessDouble()
	fmt.Println("process")
}
