package vst2

import "os"

const (
	// Extension of Vst2 files
	Extension = ".dll"
)

var (
	// ScanPaths of Vst2 files
	scanPaths = []string{
		"C:\\Program Files (x86)\\Steinberg\\VSTPlugins",
		"C:\\Program Files\\Steinberg\\VSTPlugins ",
	}
)

func init() {
	envVstPath := os.Getenv("VST_PATH")
	if len(envVstPath) > 0 {
		scanPaths = append(paths, envVstPath)
	}
}
