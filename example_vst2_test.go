//go:build !plugin
// +build !plugin

package vst2_test

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"unsafe"

	"pipelined.dev/audio/vst2"
	"pipelined.dev/signal"
)

// PrinterHostCallback returns closure that prints received opcode with
// provided prefix. This technique allows to provide callback with any
// context needed.
func PrinterHostCallback(prefix string) vst2.HostCallbackFunc {
	return func(code vst2.HostOpcode, index int32, value int64, ptr unsafe.Pointer, opt float32) int64 {
		fmt.Printf("%s: %v\n", prefix, code)
		return 0
	}
}

func Example_plugin() {
	// sample data that we will process: 2 channels, 32 samples
	data := [][]float64{
		{
			-0.0027160645,
			-0.0039978027,
			-0.0071411133,
			-0.0065307617,
			0.0038757324,
			0.021972656,
			0.041229248,
			0.055511475,
			0.064971924,
			0.07342529,
			0.08300781,
			0.092681885,
			0.10070801,
			0.110809326,
			0.12677002,
			0.15231323,
			0.19058228,
			0.24459839,
			0.3140869,
			0.38861084,
			0.44683838,
			0.47177124,
			0.46643066,
			0.45007324,
			0.4449768,
			0.45724487,
			0.47451782,
			0.48321533,
			0.47824097,
			0.46679688,
			0.45999146,
			0.46765137,
		},
		{
			-0.0027160645,
			-0.0039978027,
			-0.0071411133,
			-0.0065307617,
			0.0038757324,
			0.021972656,
			0.041229248,
			0.055511475,
			0.064971924,
			0.07342529,
			0.08300781,
			0.092681885,
			0.10070801,
			0.110809326,
			0.12677002,
			0.15231323,
			0.19058228,
			0.24459839,
			0.3140869,
			0.38861084,
			0.44683838,
			0.47177124,
			0.46643066,
			0.45007324,
			0.4449768,
			0.45724487,
			0.47451782,
			0.48321533,
			0.47824097,
			0.46679688,
			0.45999146,
			0.46765137,
		},
	}

	buffer := signal.Allocator{
		Channels: len(data),
		Length:   len(data[0]),
		Capacity: len(data[0]),
	}.Float64()
	signal.WriteStripedFloat64(data, buffer)

	// Open VST library. Library contains a reference to
	// OS-specific handle, that needs to be freed with Close.
	vst, err := vst2.Open(pluginPath())
	if err != nil {
		log.Panicf("failed to open VST library: %v", err)
	}
	defer vst.Close()

	// Load VST plugin with example callback.
	plugin := vst.Plugin(PrinterHostCallback("Received opcode"))
	defer plugin.Close()

	plugin.Start()
	// Set sample rate in Hertz.
	plugin.SetSampleRate(44100)
	// Set channels information.
	plugin.SetSpeakerArrangement(
		&vst2.SpeakerArrangement{
			Type:        vst2.SpeakerArrMono,
			NumChannels: int32(buffer.Channels()),
		},
		&vst2.SpeakerArrangement{
			Type:        vst2.SpeakerArrMono,
			NumChannels: int32(buffer.Channels()),
		},
	)
	// Set buffer size.
	plugin.SetBufferSize(buffer.Length())
	// Start the plugin.
	plugin.Resume()

	// To process data with plugin, we need to use VST2 buffers.
	// It's needed because VST SDK was written in C and expected
	// memory layout differs from Golang slices.
	// We need two buffers for input and output.
	in := vst2.NewDoubleBuffer(buffer.Channels(), buffer.Length())
	defer in.Free()
	out := vst2.NewDoubleBuffer(buffer.Channels(), buffer.Length())
	defer out.Free()

	// Fill input with data values.
	in.Write(buffer)

	// Process data.
	plugin.ProcessDouble(in, out)
	// Copy processed data.
	out.Read(buffer)

	// Output:
	// Received opcode: HostGetCurrentProcessLevel
	// Received opcode: hostWantMidi
	// Received opcode: HostGetCurrentProcessLevel
}

// pluginPath returns a path to OS-specific plugin. It will panic if OS is
// not supported.
func pluginPath() string {
	os := runtime.GOOS
	var path string
	switch os {
	case "windows":
		path, _ = filepath.Abs("_testdata\\TAL-Noisemaker.dll")
	case "darwin":
		path, _ = filepath.Abs("_testdata/TAL-Noisemaker.vst")
	default:
		panic(fmt.Sprintf("unsupported OS: %v", os))
	}
	return path
}
