package vst2

import (
	"fmt"
	"io"
	"os"
	"testing"

	wav "github.com/youpy/go-wav"
)

const (
	wavPath = "_testdata/test.wav"
)

var (
	samples64 [][]float64 //to test processDoubleReplacing
	samples32 [][]float32 //to test processReplacing
	// reader    wav.Reader
	wavFormat *wav.WavFormat
)

//Read wav for test
func init() {
	var wavSamples []wav.Sample
	inFile, _ := os.Open(wavPath)
	defer inFile.Close()
	reader := wav.NewReader(inFile)
	wavFormat, _ = reader.Format()

	for {
		read, err := reader.ReadSamples()
		if err == io.EOF {
			break
		}
		wavSamples = append(wavSamples, read...)
	}

	samples64 = convertWavSamplesToFloat64(wavSamples)
	samples32 = convertWavSamplesToFloat32(wavSamples)
}

//Test load plugin
func TestNewPlugin(t *testing.T) {
	_, err := NewPlugin(pluginPath)
	if err != nil {
		t.Fatalf("Failed NewPlugin: %v\n", err)
	}
}

//Test start plugin
func TestStartPlugin(t *testing.T) {
	plugin, err := NewPlugin(pluginPath)
	if err != nil {
		t.Fatalf("Failed NewPlugin: %v\n", err)
	}
	plugin.Start()
}

//Test Process function
func TestProcess(t *testing.T) {

	plugin, _ := NewPlugin(pluginPath)
	plugin.Start()
	plugin.Resume()
	processedSamples := plugin.ProcessFloat(samples32)

	if processedSamples == nil {
		return
	}
	for i := 0; i < 20; i++ {
		fmt.Printf("Sample %d: [%.6f][%.6f]\n", i, processedSamples[0][i], processedSamples[1][i])
	}

	fmt.Println(pluginOpcode(25))

	// outFile, err := ioutil.TempFile("/tmp", "outfile.wav")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// defer outFile.Close()

	// numChannels := uint16(len(processedSamples))
	// numSamples := uint32(len(processedSamples[0]))
	// writer := wav.NewWriter(outFile, numSamples, numChannels, wavFormat.SampleRate, wavFormat.BitsPerSample)
	// writer.WriteSamples(convertFloat32ToWavSamples(processedSamples))
}

//convert WAV samples to float64 slice
func convertWavSamplesToFloat64(wavSamples []wav.Sample) (samples [][]float64) {

	samples = make([][]float64, 2)

	samples[0] = make([]float64, 0, len(wavSamples))
	samples[1] = make([]float64, 0, len(wavSamples))

	for _, sample := range wavSamples {
		samples[0] = append(samples[0], float64(sample.Values[0])/0x8000)
		samples[1] = append(samples[1], float64(sample.Values[1])/0x8000)
	}
	return samples
}

//convert WAV samples to float32 slice
func convertWavSamplesToFloat32(wavSamples []wav.Sample) (samples [][]float32) {

	samples = make([][]float32, 2)

	samples[0] = make([]float32, 0, len(wavSamples))
	samples[1] = make([]float32, 0, len(wavSamples))

	for _, sample := range wavSamples {
		//if i < 512 {
		samples[0] = append(samples[0], float32(sample.Values[0])/0x8000)
		samples[1] = append(samples[1], float32(sample.Values[1])/0x8000)
		//}
	}
	return samples
}

func convertFloat64ToWavSamples(samples [][]float64) (wavSamples []wav.Sample) {
	wavSamples = make([]wav.Sample, len(samples[0]))
	for i := 0; i < len(samples[0]); i++ {
		wavSamples[i].Values[0] = int(samples[0][i] * 0x7FFF)
		wavSamples[i].Values[1] = int(samples[1][i] * 0x7FFF)
	}
	return
}

func convertFloat32ToWavSamples(samples [][]float32) (wavSamples []wav.Sample) {
	wavSamples = make([]wav.Sample, len(samples[0]))
	for i := 0; i < len(samples[0]); i++ {
		wavSamples[i].Values[0] = int(samples[0][i] * 0x7FFF)
		wavSamples[i].Values[1] = int(samples[1][i] * 0x7FFF)
	}
	return
}
