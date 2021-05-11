package vst2

import (
	"reflect"
	"testing"

	"pipelined.dev/signal"
)

func TestBuffer(t *testing.T) {
	type writeFn func(f signal.Floating, b DoubleBuffer)
	testBuffer := func(floats [][]float64, fn writeFn) func(*testing.T) {
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

			fn(f, b)

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
	write := func(f signal.Floating, b DoubleBuffer) {
		b.Write(f)
	}
	t.Run("mono write", testBuffer([][]float64{{1, 2, 3}}, write))
	t.Run("stereo write", testBuffer([][]float64{{11, 12, 13}, {21, 22, 23}}, write))
	iterate := func(f signal.Floating, b DoubleBuffer) {
		for c := 0; c < f.Channels(); c++ {
			for i := 0; i < b.Frames; i++ {
				b.Channel(c)[i] = f.Sample(f.BufferIndex(c, i))
			}
		}
	}
	t.Run("mono iterate", testBuffer([][]float64{{1, 2, 3}}, iterate))
	t.Run("stereo iterate", testBuffer([][]float64{{11, 12, 13}, {21, 22, 23}}, iterate))
}

func assertEqual(t *testing.T, name string, result, expected interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("%v\nresult: \t%T\t%+v \nexpected: \t%T\t%+v", name, result, result, expected, expected)
	}
}
