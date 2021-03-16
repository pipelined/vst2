package demoplugin

import (
	"pipelined.dev/signal"
)

func init() {}

type Gain struct {
	Gain float64
}

// type Param struct {
// 	value float64
// }

// type Params struct {
// 	Params []*Param
// }

// func (p Params) SetParamValue(i int, value float64) {
// 	p.Params[i].value = value
// }

// func (p Params) GetParamValue(i int) float64 {
// 	return p.Params[i].value
// }

func (p *Gain) Process(in, out signal.Floating) {
	for i := 0; i < in.Len(); i++ {
		out.SetSample(i, in.Sample(i)*p.Gain)
	}
}
