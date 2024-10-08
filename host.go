package vst2

import "pipelined.dev/signal"

type (
	// Host handles all callbacks from plugin.
	Host struct {
		GetSampleRate   HostGetSampleRateFunc
		GetBufferSize   HostGetBufferSizeFunc
		GetProcessLevel HostGetProcessLevelFunc
		GetTimeInfo     HostGetTimeInfoFunc
		UpdateDisplay   HostUpdateDisplayFunc
	}

	// HostGetSampleRateFunc returns host sample rate.
	HostGetSampleRateFunc func() signal.Frequency
	// HostGetBufferSizeFunc returns host buffer size.
	HostGetBufferSizeFunc func() int
	// HostGetProcessLevel returns the context of execution.
	HostGetProcessLevelFunc func() ProcessLevel
	// HostGetTimeInfo returns current time info.
	HostGetTimeInfoFunc func(flags TimeInfoFlag) *TimeInfo
	// HostUpdateDisplay tells there are changes & requests GUI redraw. Returns true on success
	HostUpdateDisplayFunc func() bool
)
