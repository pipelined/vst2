package vst2

import (
	"github.com/youpy/go-wav"
	"io"
	"os"
	"testing"
)

const (
	pluginPath = "_testdata/ValhallaFreqEcho_x64.dll"
	wavPath    = "_testdata/test.wav"
)

func TestLoadPlugin(t *testing.T) {
	plugin := LoadPlugin(pluginPath)
	plugin.start()
}

func TestProcessWav(t *testing.T) {

	file, _ := os.Open(wavPath)
	reader := wav.NewReader(file)

	defer file.Close()

	var samples []wav.Sample

	for {
		read, err := reader.ReadSamples()
		if err == io.EOF {
			break
		}

		samples = append(samples, read...)
		/*for _, sample := range samples {
			fmt.Printf("L/R: %d/%d\n", reader.IntValue(sample, 0), reader.IntValue(sample, 1))
		}*/
	}

	plugin := LoadPlugin(pluginPath)
	plugin.start()

	plugin.processAudio(samples)
}
