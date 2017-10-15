package vst2

import (
	"fmt"
	"testing"
)

func TestLoadPlugin(t *testing.T) {
	fmt.Print("Hello vst host!\n")
	plugin := LoadPlugin("D:\\vst_home\\ValhallaFreqEcho_x64.dll")
	plugin.start()
	fmt.Printf("%v\n", plugin)
}
