package vst2

import (
	"reflect"
	"testing"

	"pipelined.dev/signal"
)

func TestBuffer(t *testing.T) {
	testBuffer := func(floats [][]float64) func(*testing.T) {
		return func(t *testing.T) {
			channels := len(floats)
			size := len(floats[0])
			b := NewDoubleBuffer(channels, size)

			f := signal.Allocator{
				Channels: channels,
				Length:   size,
				Capacity: size,
			}.Float64()
			signal.WriteStripedFloat64(floats, f)

			b.CopyFrom(f)

			for i := range floats {
				for j := range floats[i] {
					assertEqual(t, "double value", b.Channel(i)[j], floats[i][j])
				}
			}

			for i := range floats {
				var testBuf DoubleBuffer
				testBuf.data = append(testBuf.data, getDoubleChannel(b.cArray(), i))
				testBuf.Frames = len(floats[i])
				cChannel := testBuf.Channel(0)
				for j := range floats[i] {
					assertEqual(t, "double value", cChannel[j], floats[i][j])
				}
			}
		}
	}
	t.Run("mono", testBuffer([][]float64{{1, 2, 3}}))
	t.Run("stereo", testBuffer([][]float64{{11, 12, 13}, {21, 22, 23}}))
}

func assertEqual(t *testing.T, name string, result, expected interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("%v\nresult: \t%T\t%+v \nexpected: \t%T\t%+v", name, result, result, expected, expected)
	}
}
