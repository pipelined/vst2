package audio

type Processor interface {
	Process(samples [][]float64) [][]float64
}
