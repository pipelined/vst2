package vst2

import "pipelined.dev/signal"

type (
	// Host handles all callbacks from plugin.
	Host struct {
		GetSampleRate        HostGetSampleRateFunc
		GetBufferSize        HostGetBufferSizeFunc
		GetProcessLevel      HostGetProcessLevelFunc
		GetTimeInfo          HostGetTimeInfoFunc
		UpdateDisplay        HostUpdateDisplayFunc
		Automate             HostAutomateFunc
		Idle                 HostIdleFunc
		ProcessEvents        HostProcessEventsFunc
		IOChanged            HostIOChangedFunc
		SizeWindow           HostSizeWindowFunc
		GetInputLatency      HostGetInputLatencyFunc
		GetOutputLatency     HostGetOutputLatencyFunc
		GetAutomationState   HostGetAutomationStateFunc
		GetVendorString      HostGetVendorStringFunc
		GetProductString     HostGetProductStringFunc
		GetVendorVersion     HostGetVendorVersionFunc
		CanDo                HostCanDoFunc
		GetLanguage          HostGetLanguageFunc
		GetDirectory         HostGetDirectoryFunc
		BeginEdit            HostBeginEditFunc
		EndEdit              HostEndEditFunc
	}

	// HostGetSampleRateFunc returns host sample rate.
	HostGetSampleRateFunc func() signal.Frequency
	// HostGetBufferSizeFunc returns host buffer size.
	HostGetBufferSizeFunc func() int
	// HostGetProcessLevelFunc returns the context of execution.
	HostGetProcessLevelFunc func() ProcessLevel
	// HostGetTimeInfoFunc returns current time info.
	HostGetTimeInfoFunc func(flags TimeInfoFlag) *TimeInfo
	// HostUpdateDisplayFunc tells there are changes & requests GUI redraw. Returns true on success.
	HostUpdateDisplayFunc func() bool
	// HostAutomateFunc is called when the plugin has automated a parameter value.
	HostAutomateFunc func(index int32, value float32)
	// HostIdleFunc is called when the plugin requests idle from host.
	HostIdleFunc func()
	// HostProcessEventsFunc is called when the plugin sends MIDI events to the host.
	HostProcessEventsFunc func(events *EventsPtr)
	// HostIOChangedFunc is called when plugin I/O configuration has changed. Returns true if supported.
	HostIOChangedFunc func() bool
	// HostSizeWindowFunc is called when plugin editor wants to resize. Returns true if supported.
	HostSizeWindowFunc func(width, height int32) bool
	// HostGetInputLatencyFunc returns input latency in samples.
	HostGetInputLatencyFunc func() int32
	// HostGetOutputLatencyFunc returns output latency in samples.
	HostGetOutputLatencyFunc func() int32
	// HostGetAutomationStateFunc returns the current automation state.
	HostGetAutomationStateFunc func() AutomationState
	// HostGetVendorStringFunc returns the host vendor name.
	HostGetVendorStringFunc func() string
	// HostGetProductStringFunc returns the host product name.
	HostGetProductStringFunc func() string
	// HostGetVendorVersionFunc returns the host version.
	HostGetVendorVersionFunc func() int32
	// HostCanDoFunc queries host capabilities.
	HostCanDoFunc func(HostCanDoString) CanDoResponse
	// HostGetLanguageFunc returns the host UI language.
	HostGetLanguageFunc func() HostLanguage
	// HostGetDirectoryFunc returns the plugin install directory.
	HostGetDirectoryFunc func() string
	// HostBeginEditFunc is called when a parameter edit has started. Returns true on success.
	HostBeginEditFunc func(index int32) bool
	// HostEndEditFunc is called when a parameter edit has ended. Returns true on success.
	HostEndEditFunc func(index int32) bool
)
