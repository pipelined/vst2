package vst2

import (
	"sync"
	"testing"
)

func TestParameterRawValue(t *testing.T) {
	var p Parameter
	assertEqual(t, "zero value", p.RawValue(), float32(0))

	p.SetValue(0.25)
	assertEqual(t, "raw value", p.RawValue(), float32(0.25))
	assertEqual(t, "unmapped get value", p.GetValue(), float32(0.25))

	p.GetValueFunc = func(value float32) float32 { return -20 + 40*value }
	assertEqual(t, "raw value", p.RawValue(), float32(0.25))
	assertEqual(t, "mapped get value", p.GetValue(), float32(-10))
	assertEqual(t, "value label", p.GetValueLabel(), "-10.000000")

	p.GetValueLabelFunc = func(value float32) string {
		if value < 0 {
			return "below"
		}
		return "above"
	}
	assertEqual(t, "value label", p.GetValueLabel(), "below")
}

// TestParameterConcurrentAccess mimics the host writing a parameter while
// the audio callback reads it. It only fails under -race.
func TestParameterConcurrentAccess(t *testing.T) {
	var (
		p  Parameter
		wg sync.WaitGroup
	)
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			p.SetValue(float32(i) / 1000)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			_ = p.GetValue()
		}
	}()
	wg.Wait()
	if v := p.RawValue(); v < 0 || v > 1 {
		t.Fatalf("raw value out of range: %v", v)
	}
}
