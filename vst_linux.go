package vst2

import (
	"os"
	"path/filepath"
)

const (
	// FileExtension of VST2 files in Windows.
	FileExtension = ".so"
)

var (
	// ScanPaths of Vst2 files
	scanPaths = []string{
		"/usr/local/lib/vst",
		"/usr/lib/vst",
	}
)

func init() {
	home := os.Getenv("HOME")
	scanPaths = append(scanPaths,
		filepath.Join(home, ".vst"),
		filepath.Join(home, "vst"),
	)
}

func Open(path string) (*VST, error) {
	return nil, nil
}

// Close frees plugin handle.
func (m *VST) Close() error {
	return nil
}
