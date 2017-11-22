package loader

import (
	"testing"
)

var (
	path = "../vst2/_testdata/"
)

func TestLoadAll(t *testing.T) {
	LoadAll(path)
}
