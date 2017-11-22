package vst2

import (
	"github.com/youpy/go-wav"
	// "go/build"
	"io"
	"os"
	"testing"
)

const (
	pluginPath = "_testdata/2-band distortion.dll"
	wavPath    = "_testdata/test.wav"
)

//Test load plugin
func TestLoadPlugin(t *testing.T) {
	plugin, err := LoadPlugin(pluginPath)
	if err != nil {
		t.Fatalf("Failed LoadPlugin: %v\n", err)
	}
	t.Logf("Passed LoadPlugin: %v\n", plugin)
}

//Test load plugin
func TestFailLoadPlugin(t *testing.T) {
	_, err := LoadPlugin(pluginPath)
	if err != nil {
		t.Logf("Passed FailLoadPlugin: %v\n", err)
	}
}

//Test processAudio function
func TestProcess(t *testing.T) {
	samples := ConvertWavSamplesToFloat64(readWav(wavPath))

	plugin, _ := LoadPlugin(pluginPath)
	plugin.start()

	plugin.resume()

	plugin.Process(samples)
}

//Read wav for test
func readWav(wavPath string) (wavSamples []wav.Sample) {

	file, _ := os.Open(wavPath)
	reader := wav.NewReader(file)

	defer file.Close()

	for {
		read, err := reader.ReadSamples()
		if err == io.EOF {
			break
		}

		wavSamples = append(wavSamples, read...)
	}

	return
}

//convert WAV samples to float slice
func ConvertWavSamplesToFloat64(wavSamples []wav.Sample) (samples [][]float64) {
	samples = make([][]float64, 2)

	samples[0] = make([]float64, len(wavSamples))
	samples[1] = make([]float64, len(wavSamples))

	for i, sample := range wavSamples {
		samples[0][i] = float64(sample.Values[0]) / 32768
		samples[1][i] = float64(sample.Values[1]) / 32768
	}
	return samples
}
