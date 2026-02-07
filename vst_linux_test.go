//go:build !plugin
// +build !plugin

package vst2_test

import (
	"fmt"
	"testing"

	"pipelined.dev/audio/vst2"
)

func TestLinux(t *testing.T) {
	fmt.Printf("linux paths: %v\n", vst2.ScanPaths())
}
