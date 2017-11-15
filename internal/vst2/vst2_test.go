package vst2

import (
	// "fmt"
	"github.com/youpy/go-wav"
	"io"
	"os"
	"testing"
)

const (
	pluginPath = "_testdata/ValhallaFreqEcho_x64.dll"
	wavPath    = "_testdata/test.wav"
)

//Test load plugin
func TestLoadPlugin(t *testing.T) {
	plugin, err := LoadPlugin(pluginPath)
	if err != nil {
		t.Fatalf("Failed LoadPlugin: %v\n", err)
	}
	plugin.start()
}

//Test processAudio function
func TestProcessAudio(t *testing.T) {
	samples := readWav(wavPath)

	plugin, _ := LoadPlugin(pluginPath)
	plugin.start()

	plugin.processAudio(samples)
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
