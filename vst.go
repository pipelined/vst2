package vst2

import (
	"fmt"
	"strings"
)

// EffectMagic is constant in every plugin.
const EffectMagic int32 = 'V'<<24 | 's'<<16 | 't'<<8 | 'P'<<0

type (
	// ParameterProperties contains the information about parameter.
	ParameterProperties struct {
		// valid if ParameterUsesIntegerMinMax is set
		StepFloat      float32
		SmallStepFloat float32
		LargeStepFloat float32

		Label ascii64
		Flags ParameterFlag

		// valid if ParameterUsesIntegerMinMax is set
		MinInteger int32
		MaxInteger int32

		// valid if ParameterUsesIntStep is set
		StepInteger      int32
		LargeStepInteger int32
		ShortLabel       ascii8

		// valid if ParameterSupportsDisplayIndex is set
		DisplayIndex int16 // Index where parameter should be displayed, starts with 0

		// valid if ParameterSupportsDisplayCategory is set
		Category             int16
		ParametersInCategory int16
		reserved             int16 // always zero
		CategoryLabel        ascii24

		future [16]byte // Reserved for future.
	}

	// ParameterFlag is used to describe ParameterProperties struct.
	ParameterFlag int32
)

const (
	// ParameterIsSwitch is set if parameter is a switch (on/off).
	ParameterIsSwitch ParameterFlag = 1 << iota
	// ParameterUsesIntegerMinMax is set if parameter has min/max int
	// values.
	ParameterUsesIntegerMinMax
	// ParameterUsesFloatStep is set if parameter uses float steps.
	ParameterUsesFloatStep
	// ParameterUsesIntStep is set if parameter uses int steps.
	ParameterUsesIntStep
	// ParameterSupportsDisplayIndex is set if parameter should be
	// displayed at certain position.
	ParameterSupportsDisplayIndex
	// ParameterSupportsDisplayCategory is set if parameter should be
	// displayed under specific category.
	ParameterSupportsDisplayCategory
	// ParameterCanRamp is set if parameter can ramp up/down.
	ParameterCanRamp
)

type (
	// TimeInfo describes the time at the start of the block currently
	// being processed.
	TimeInfo struct {
		// Current Position in audio samples.
		SamplePos float64
		// Current Sample Rate in Herz.
		SampleRate float64
		// System Time in nanoseconds.
		NanoSeconds float64
		// Musical Position, in Quarter Note (1.0 equals 1 Quarter Note).
		PpqPos float64
		// Current Tempo in BPM (Beats Per Minute).
		Tempo float64
		// Last Bar Start Position, in Quarter Note.
		BarStartPos float64
		// Cycle Start (left locator), in Quarter Note.
		CycleStartPos float64
		// Cycle End (right locator), in Quarter Note.
		CycleEndPos float64
		// Time Signature Numerator (e.g. 3 for 3/4).
		TimeSigNumerator int32
		// Time Signature Denominator (e.g. 4 for 3/4).
		TimeSigDenominator int32
		// SMPTE offset in SMPTE subframes (bits; 1/80 of a frame). The
		// current SMPTE position can be calculated using SamplePos,
		// SampleRate, and SMPTEFrameRate.
		SMPTEOffset int32
		// SMPTEFrameRate value.
		SMPTEFrameRate
		// MIDI Clock Resolution (24 Per Quarter Note), can be negative
		// (nearest clock).
		SamplesToNextClock int32
		// TimeInfoFlags values.
		Flags TimeInfoFlag
	}

	// TimeInfoFlag used in TimeInfo.
	TimeInfoFlag int32

	// SMPTEFrameRate values, used in TimeInfo.
	SMPTEFrameRate int32
)

const (
	// TransportChanged is set if play, cycle or record state has changed.
	TransportChanged TimeInfoFlag = 1 << iota
	// TransportPlaying is set if Host sequencer is currently playing
	TransportPlaying
	// TransportCycleActive is set if Host sequencer is in cycle mode.
	TransportCycleActive
	// TransportRecording is set if Host sequencer is in record mode.
	TransportRecording
	_
	_
	// AutomationWriting is set if automation write mode active.
	AutomationWriting
	// AutomationReading is set if automation read mode active.
	AutomationReading
	// NanosValid is set if TimeInfo.NanoSeconds are valid.
	NanosValid
	// PpqPosValid is set if TimeInfo.PpqPos is valid.
	PpqPosValid
	// TempoValid is set if TimeInfo.Tempo is valid.
	TempoValid
	// BarsValid is set if TimeInfo.BarStartPos is valid.
	BarsValid
	// CyclePosValid is set if both TimeInfo.CycleStartPos and
	// TimeInfo.CycleEndPos are valid.
	CyclePosValid
	// TimeSigValid is set if both TimeInfo.TimeSigNumerator and
	// TimeInfo.TimeSigDenominator are valid.
	TimeSigValid
	// SMPTEValid is set if both TimeInfo.SMPTEOffset and
	// TimeInfo.SMPTEFrameRate are valid.
	SMPTEValid
	// ClockValid is set if TimeInfo.SamplesToNextClock are valid.
	ClockValid
)

const (
	// SMPTE24fps is 24 fps.
	SMPTE24fps SMPTEFrameRate = iota
	// SMPTE25fps is 25 fps.
	SMPTE25fps
	// SMPTE2997fps is 29.97 fps.
	SMPTE2997fps
	// SMPTE30fps is 30 fps.
	SMPTE30fps
	// SMPTE2997dfps is 29.97 drop.
	SMPTE2997dfps
	// SMPTE30dfps is 30 drop.
	SMPTE30dfps
	// SMPTEFilm16mm is Film 16mm.
	SMPTEFilm16mm
	// SMPTEFilm35mm is Film 35mm.
	SMPTEFilm35mm
	_
	_
	// SMPTE239fps is HDTV 23.976 fps.
	SMPTE239fps
	// SMPTE249fps is HDTV 24.976 fps.
	SMPTE249fps
	// SMPTE599fps is HDTV 59.94 fps.
	SMPTE599fps
	// SMPTE60fps is HDTV 60 fps.
	SMPTE60fps
)

type (
	// SpeakerArrangement contains information about a channel.
	SpeakerArrangement struct {
		Type        SpeakerArrangementType
		NumChannels int32
		Speakers    [8]Speaker
	}

	// SpeakerArrangementType indicates how the channels are intended to be
	// used in the plugin. Only useful for some hosts.
	SpeakerArrangementType int32

	// Speaker configuration.
	Speaker struct {
		Azimuth   float32
		Elevation float32
		Radius    float32
		Reserved  float32
		Name      ascii64
		Type      SpeakerType
		Future    [28]byte
	}

	// SpeakerType of particular speaker.
	SpeakerType int32
)

const (
	// SpeakerArrUserDefined is user defined.
	SpeakerArrUserDefined SpeakerArrangementType = iota - 2
	// SpeakerArrEmpty is empty arrangement.
	SpeakerArrEmpty
	// SpeakerArrMono is M.
	SpeakerArrMono
	// SpeakerArrStereo is L R.
	SpeakerArrStereo
	// SpeakerArrStereoSurround is Ls Rs.
	SpeakerArrStereoSurround
	// SpeakerArrStereoCenter is Lc Rc.
	SpeakerArrStereoCenter
	// SpeakerArrStereoSide is Sl Sr.
	SpeakerArrStereoSide
	// SpeakerArrStereoCLfe is C Lfe.
	SpeakerArrStereoCLfe
	// SpeakerArr30Cine is L R C.
	SpeakerArr30Cine
	// SpeakerArr30Music is L R S.
	SpeakerArr30Music
	// SpeakerArr31Cine is L R C Lfe.
	SpeakerArr31Cine
	// SpeakerArr31Music is L R Lfe S.
	SpeakerArr31Music
	// SpeakerArr40Cine is L R C S (LCRS).
	SpeakerArr40Cine
	// SpeakerArr40Music is L R Ls Rs (Quadro).
	SpeakerArr40Music
	// SpeakerArr41Cine is L R C Lfe S (LCRS+Lfe).
	SpeakerArr41Cine
	// SpeakerArr41Music is L R Lfe Ls Rs (Quadro+Lfe).
	SpeakerArr41Music
	// SpeakerArr50 is L R C Ls Rs.
	SpeakerArr50
	// SpeakerArr51 is L R C Lfe Ls Rs.
	SpeakerArr51
	// SpeakerArr60Cine is L R C Ls Rs Cs.
	SpeakerArr60Cine
	// SpeakerArr60Music is L R Ls Rs Sl Sr.
	SpeakerArr60Music
	// SpeakerArr61Cine is L R C Lfe Ls Rs Cs.
	SpeakerArr61Cine
	// SpeakerArr61Music is L R Lfe Ls Rs Sl Sr.
	SpeakerArr61Music
	// SpeakerArr70Cine is L R C Ls Rs Lc Rc.
	SpeakerArr70Cine
	// SpeakerArr70Music is L R C Ls Rs Sl Sr.
	SpeakerArr70Music
	// SpeakerArr71Cine is L R C Lfe Ls Rs Lc Rc.
	SpeakerArr71Cine
	// SpeakerArr71Music is L R C Lfe Ls Rs Sl Sr.
	SpeakerArr71Music
	// SpeakerArr80Cine is L R C Ls Rs Lc Rc Cs.
	SpeakerArr80Cine
	// SpeakerArr80Music is L R C Ls Rs Cs Sl Sr.
	SpeakerArr80Music
	// SpeakerArr81Cine is L R C Lfe Ls Rs Lc Rc Cs.
	SpeakerArr81Cine
	// SpeakerArr81Music is L R C Lfe Ls Rs Cs Sl Sr.
	SpeakerArr81Music
	// SpeakerArr102 is L R C Lfe Ls Rs Tfl Tfc Tfr Trl Trr Lfe2.
	SpeakerArr102
	// numSpeakerArr not defined.
	_
)

const (
	// SpeakerUndefined is undefined.
	SpeakerUndefined SpeakerType = 0x7fffffff
	// SpeakerM is Mono (M).
	SpeakerM = iota
	// SpeakerL is Left (L).
	SpeakerL
	// SpeakerR is Right (R).
	SpeakerR
	// SpeakerC is Center (C).
	SpeakerC
	// SpeakerLfe is Subbass (Lfe).
	SpeakerLfe
	// SpeakerLs is Left Surround (Ls).
	SpeakerLs
	// SpeakerRs is Right Surround (Rs).
	SpeakerRs
	// SpeakerLc is Left of Center (Lc).
	SpeakerLc
	// SpeakerRc is Right of Center (Rc).
	SpeakerRc
	// SpeakerS is Surround (S).
	SpeakerS
	// SpeakerCs is Center of Surround (Cs) = Surround (S).
	SpeakerCs = SpeakerS
	// SpeakerSl is Side Left (Sl).
	SpeakerSl
	// SpeakerSr is Side Right (Sr).
	SpeakerSr
	// SpeakerTm is Top Middle (Tm).
	SpeakerTm
	// SpeakerTfl is Top Front Left (Tfl).
	SpeakerTfl
	// SpeakerTfc is Top Front Center (Tfc).
	SpeakerTfc
	// SpeakerTfr is Top Front Right (Tfr).
	SpeakerTfr
	// SpeakerTrl is Top Rear Left (Trl).
	SpeakerTrl
	// SpeakerTrc is Top Rear Center (Trc).
	SpeakerTrc
	// SpeakerTrr is Top Rear Right (Trr).
	SpeakerTrr
	// SpeakerLfe2 is Subbass 2 (Lfe2).
	SpeakerLfe2
)

// PluginFlag values.
type PluginFlag int32

const (
	// PluginHasEditor is set if the plugin provides a custom editor.
	PluginHasEditor PluginFlag = 1 << iota
	_
	_
	_
	// PluginFloatProcessing is set if plugin supports replacing process
	// mode.
	PluginFloatProcessing
	// PluginProgramChunks is set if preset data is handled in formatless
	// chunks.
	PluginProgramChunks
	_
	_
	// PluginIsSynth is set if plugin is a synth.
	PluginIsSynth
	// PluginNoSoundInStop is set if plugin does not produce sound when
	// input is silence.
	PluginNoSoundInStop
	_
	_
	// PluginDoubleProcessing is set if plugin supports double precision
	// processing.
	PluginDoubleProcessing

	// pluginHasClip deprecated in VST v2.4
	_
	// pluginHasVu deprecated in VST v2.4
	_
	// pluginCanMono deprecated in VST v2.4
	_
	// pluginExtIsAsync deprecated in VST v2.4
	_
	// pluginExtHasBuffer deprecated in VST v2.4
	_
)

// ProcessLevel is used as result for in HostGetCurrentProcessLevel call.
// It tells the plugin in which thread host is right now.
type ProcessLevel uintptr

const (
	// ProcessLevelUnknown is returned when not supported by host.
	ProcessLevelUnknown ProcessLevel = iota
	// ProcessLevelUser is returned when in user thread (GUI).
	ProcessLevelUser
	// ProcessLevelRealtime is returned when in audio thread (where process
	// is called).
	ProcessLevelRealtime
	// ProcessLevelPrefetch is returned when in sequencer thread (MIDI,
	// timer etc).
	ProcessLevelPrefetch
	// ProcessLevelOffline is returned when in offline processing and thus
	// in user thread.
	ProcessLevelOffline
)

// HostLanguage is the language of the host.
type HostLanguage uintptr

const (
	// HostLanguageEnglish English.
	HostLanguageEnglish HostLanguage = iota + 1
	// HostLanguageGerman German.
	HostLanguageGerman
	// HostLanguageFrench French.
	HostLanguageFrench
	// HostLanguageItalian Italian.
	HostLanguageItalian
	// HostLanguageSpanish Spanish.
	HostLanguageSpanish
	// HostLanguageJapanese Japanese.
	HostLanguageJapanese
)

type (
	// PinProperties provides info about about plugin connectivity.
	PinProperties struct {
		Label ascii64
		Flags PinPropertiesFlag
		SpeakerArrangementType
		ShortLabel ascii8   // Short name, recommended 6 chars + delimiter.
		_          [48]byte // reserved not used.
	}

	// PinPropertiesFlag values.
	PinPropertiesFlag int32
)

const (
	// PinIsActive is ignored by Host.
	PinIsActive PinPropertiesFlag = 1 << iota
	// PinIsStereo means that pin is first of a stereo pair.
	PinIsStereo
	// PinUseSpeaker means that arrangement type is valid and pin can be
	// used for arrangement setup.
	PinUseSpeaker
)

// PluginCategory denotes the category of plugin.
type PluginCategory int64

const (
	// PluginCategoryUnknown means category not implemented.
	PluginCategoryUnknown PluginCategory = iota
	// PluginCategoryEffect simple Effect.
	PluginCategoryEffect
	// PluginCategorySynth VST Instrument: synth, sampler, etc.
	PluginCategorySynth
	// PluginCategoryAnalysis scope, tuner.
	PluginCategoryAnalysis
	// PluginCategoryMastering dynamics control.
	PluginCategoryMastering
	// PluginCategorySpacializer panner.
	PluginCategorySpacializer
	// PluginCategoryRoomFx delay and reverb.
	PluginCategoryRoomFx
	// PluginCategorySurroundFx dedicated surround.
	PluginCategorySurroundFx
	// PluginCategoryRestoration denoiser.
	PluginCategoryRestoration
	// PluginCategoryOfflineProcess offline processor.
	PluginCategoryOfflineProcess
	// PluginCategoryShell plugin is a shell for other plugins.
	PluginCategoryShell
	// PluginCategoryGenerator tone generator.
	PluginCategoryGenerator
	pluginCategoryMaxCount
)

// ProcessPrecision allows to set processing precision of plugin.
type ProcessPrecision int64

const (
	// ProcessFloat is 32 bits processing.
	ProcessFloat ProcessPrecision = iota
	// ProcessDouble is 64 bits processing.
	ProcessDouble
)

// MIDIProgram describes the MIDI program.
type MIDIProgram struct {
	Index       int32
	Name        ascii64
	MIDIProgram int8  // -1:off [-1;-127]
	MIDIBankMsb int8  // -1:off [-1;-127]
	MIDIBankLsb int8  // -1:off [-1;-127]
	_           int8  // reserved not used.
	ParentIndex int32 // -1 means there is no parent category.
	Flags       MIDIProgramFlag
}

// MIDIProgramFlag values.
type MIDIProgramFlag int32

// MIDIProgramIsOmni program is in omni mode, channel 0 is used for
// inquiries and program changes.
const MIDIProgramIsOmni MIDIProgramFlag = 1

// MIDIProgramCategory describes the MIDI program category.
type MIDIProgramCategory struct {
	Index       int32
	Name        ascii64 ///< name
	ParentIndex int32   // -1 means there is no parent category.
	_           int32   // flags not used.
}

// MIDIKey describes the MIDI key.
type MIDIKey struct {
	Index     int32
	KeyNumber int32 // [0; 127]
	Name      ascii64
	_         int32 // reserved not used.
	_         int32 // flags not used.
}

// PatchChunk is used to communicate preset or bank properties with plugin
// before uploading it.
type PatchChunk struct {
	version        int32 // Always equals 1.
	PluginUniqueID int32
	PluginVersion  int32
	NumElements    int32    // Number of presets for bank or number of parameters for preset.
	_              [48]byte // reserved not used.
}

// PanningLaw determines the algorithm panning happens.
type PanningLaw float32

const (
	// PanningLawLinear uses the following formula: L = pan * M; R = (1 -
	// pan) * M.
	PanningLawLinear PanningLaw = 0
	// PanningLawEqualPower uses the following formula: L = pow (pan, 0.5)
	// * M; R = pow ((1 - pan), 0.5) * M.
	PanningLawEqualPower
)

// AutomationState communicates the host state of automation.
type AutomationState uintptr

const (
	// AutomationUnsupported returned when host doesn't support automation.
	AutomationUnsupported AutomationState = 0
	// AutomationOff returned when automation is switched off.
	AutomationOff
	// AutomationRead returned when host is reading the automation.
	AutomationRead
	// AutomationWrite returned when host is writing the automation.
	AutomationWrite
	// AutomationReadWrite returned when host is reading and writing the
	// automation.
	AutomationReadWrite
)

// KeyCode used to pass information about key presses.
type KeyCode struct {
	Character int32 // ASCII character.
	VirtualKey
	ModifierKeyFlag // Bit flags.
}

// VirtualKey is platform-independent definition of Virtual Keys used in
// KeyCode messages.
type VirtualKey uint8

// ModifierKeyFlag are flags used in KeyCode messages.
type ModifierKeyFlag uint8

const (
	// VirtualKeyBack is Backspace key.
	VirtualKeyBack VirtualKey = iota + 1
	// VirtualKeyTab is Tab key.
	VirtualKeyTab
	// VirtualKeyClear is Clear key.
	VirtualKeyClear
	// VirtualKeyReturn is Return key.
	VirtualKeyReturn
	// VirtualKeyPause is Pause key.
	VirtualKeyPause
	// VirtualKeyEscape is Escape key.
	VirtualKeyEscape
	// VirtualKeySpace is Space key.
	VirtualKeySpace
	// VirtualKeyNext is Next key.
	VirtualKeyNext
	// VirtualKeyEnd is End key.
	VirtualKeyEnd
	// VirtualKeyHome is Home key.
	VirtualKeyHome
	// VirtualKeyLeft is Left key.
	VirtualKeyLeft
	// VirtualKeyUp is Up key.
	VirtualKeyUp
	// VirtualKeyRight is Right key.
	VirtualKeyRight
	// VirtualKeyDown is Down key.
	VirtualKeyDown
	// VirtualKeyPageUp is PageUp key.
	VirtualKeyPageUp
	// VirtualKeyPageDown is PageDown key.
	VirtualKeyPageDown
	// VirtualKeySelect is Select key.
	VirtualKeySelect
	// VirtualKeyPrint is Print key.
	VirtualKeyPrint
	// VirtualKeyEnter is Enter key.
	VirtualKeyEnter
	// VirtualKeySnapshot is Snapshot key.
	VirtualKeySnapshot
	// VirtualKeyInsert is Insert key.
	VirtualKeyInsert
	// VirtualKeyDelete is Delete key.
	VirtualKeyDelete
	// VirtualKeyHelp is Help key.
	VirtualKeyHelp
	// VirtualKeyNumpad0 is Numpad0 key.
	VirtualKeyNumpad0
	// VirtualKeyNumpad1 is Numpad1 key.
	VirtualKeyNumpad1
	// VirtualKeyNumpad2 is Numpad2 key.
	VirtualKeyNumpad2
	// VirtualKeyNumpad3 is Numpad3 key.
	VirtualKeyNumpad3
	// VirtualKeyNumpad4 is Numpad4 key.
	VirtualKeyNumpad4
	// VirtualKeyNumpad5 is Numpad5 key.
	VirtualKeyNumpad5
	// VirtualKeyNumpad6 is Numpad6 key.
	VirtualKeyNumpad6
	// VirtualKeyNumpad7 is Numpad7 key.
	VirtualKeyNumpad7
	// VirtualKeyNumpad8 is Numpad8 key.
	VirtualKeyNumpad8
	// VirtualKeyNumpad9 is Numpad9 key.
	VirtualKeyNumpad9
	// VirtualKeyMultiply is Multiply key.
	VirtualKeyMultiply
	// VirtualKeyAdd is Add key.
	VirtualKeyAdd
	// VirtualKeySeparator is Separator key.
	VirtualKeySeparator
	// VirtualKeySubtract is Subtract key.
	VirtualKeySubtract
	// VirtualKeyDecimal is Decimal key.
	VirtualKeyDecimal
	// VirtualKeyDivide is Divide key.
	VirtualKeyDivide
	// VirtualKeyF1 is F1 key.
	VirtualKeyF1
	// VirtualKeyF2 is F2 key.
	VirtualKeyF2
	// VirtualKeyF3 is F3 key.
	VirtualKeyF3
	// VirtualKeyF4 is F4 key.
	VirtualKeyF4
	// VirtualKeyF5 is F5 key.
	VirtualKeyF5
	// VirtualKeyF6 is F6 key.
	VirtualKeyF6
	// VirtualKeyF7 is F7 key.
	VirtualKeyF7
	// VirtualKeyF8 is F8 key.
	VirtualKeyF8
	// VirtualKeyF9 is F9 key.
	VirtualKeyF9
	// VirtualKeyF10 is F10 key.
	VirtualKeyF10
	// VirtualKeyF11 is F11 key.
	VirtualKeyF11
	// VirtualKeyF12 is F12 key.
	VirtualKeyF12
	// VirtualKeyNumlock is Numlock key.
	VirtualKeyNumlock
	// VirtualKeyScroll is Scroll key.
	VirtualKeyScroll
	// VirtualKeyShift is Shift key.
	VirtualKeyShift
	// VirtualKeyControl is Control key.
	VirtualKeyControl
	// VirtualKeyAlt is Alt key.
	VirtualKeyAlt
	// VirtualKeyEquals is Equals key.
	VirtualKeyEquals
)

const (
	// ModifierKeyShift is Shift key.
	ModifierKeyShift ModifierKeyFlag = 1 << iota
	// ModifierKeyAlternate is Alt key.
	ModifierKeyAlternate
	// ModifierKeyCommand is Command key on Mac.
	ModifierKeyCommand
	// ModifierKeyControl is Control key.
	ModifierKeyControl
)

// EditorRectangle holds the information about plugin editor window.
type EditorRectangle struct {
	Top    int16
	Left   int16
	Bottom int16
	Right  int16
}

type (
	// HostCanDoString are constants that can be used to check host
	// capabilities.
	HostCanDoString string
	// PluginCanDoString are constants that can be used to check plugin
	// capabilities.
	PluginCanDoString string
	// CanDoResponse lists the possible return values for CanDo queries
	CanDoResponse int64
)

const (
	// HostCanSendEvents host can send events.
	HostCanSendEvents HostCanDoString = "sendVstEvents"
	// HostCanSendMIDIEvent host can send MIDI events.
	HostCanSendMIDIEvent HostCanDoString = "sendVstMidiEvent"
	// HostCanSendTimeInfo host can send TimeInfo.
	HostCanSendTimeInfo HostCanDoString = "sendVstTimeInfo"
	// HostCanReceiveEvents host can receive events from plugin.
	HostCanReceiveEvents HostCanDoString = "receiveVstEvents"
	// HostCanReceiveMIDIEvent host can receive MIDI events from plugin.
	HostCanReceiveMIDIEvent HostCanDoString = "receiveVstMidiEvent"
	// HostCanReportConnectionChanges host can notify the plugin when
	// something change in plugin´s routing/connections with
	// Suspend/Resume/SetSpeakerArrangement.
	HostCanReportConnectionChanges HostCanDoString = "reportConnectionChanges"
	// HostCanAcceptIOChanges host can receive HostIOChanged.
	HostCanAcceptIOChanges HostCanDoString = "acceptIOChanges"
	// HostCanSizeWindow used by VSTGUI.
	HostCanSizeWindow HostCanDoString = "sizeWindow"
	// HostCanOffline host supports offline processing feature.
	HostCanOffline HostCanDoString = "offline"
	// HostCanOpenFileSelector host supports opcode HostOpenFileSelector.
	HostCanOpenFileSelector HostCanDoString = "openFileSelector"
	// HostCanCloseFileSelector host supports opcode HostCloseFileSelector.
	HostCanCloseFileSelector HostCanDoString = "closeFileSelector"
	// HostCanStartStopProcess host supports PlugStartProcess and
	// PlugStopProcess functions.
	HostCanStartStopProcess HostCanDoString = "startStopProcess"
	// HostCanShellCategory host supports plugins with PluginCategoryShell.
	HostCanShellCategory HostCanDoString = "shellCategory"
	// HostCanSendRealtimeMIDIEvent host can send realtime MIDI events.
	HostCanSendRealtimeMIDIEvent HostCanDoString = "sendVstMidiEventFlagIsRealtime"
)

const (
	// PluginCanSendEvents plugin can send events.
	PluginCanSendEvents PluginCanDoString = "sendVstEvents"
	// PluginCanSendMIDIEvent plugin can send MIDI events.
	PluginCanSendMIDIEvent PluginCanDoString = "sendVstMidiEvent"
	// PluginCanReceiveEvents plugin can receive events.
	PluginCanReceiveEvents PluginCanDoString = "receiveVstEvents"
	// PluginCanReceiveMIDIEvent plugin can receive MIDI events.
	PluginCanReceiveMIDIEvent PluginCanDoString = "receiveVstMidiEvent"
	// PluginCanReceiveTimeInfo plugin can receive TimeInfo.
	PluginCanReceiveTimeInfo PluginCanDoString = "receiveVstTimeInfo"
	// PluginCanOffline plugin supports offline functions.
	PluginCanOffline PluginCanDoString = "offline"
	// PluginCanMIDIProgramNames plugin supports function
	// GetMIDIProgramName.
	PluginCanMIDIProgramNames PluginCanDoString = "midiProgramNames"
	// PluginCanBypass plugin supports function SetBypass.
	PluginCanBypass PluginCanDoString = "bypass"
)

const (
	// YesCanDo is returned by CanDo queries when plugin/host is certain that it supports a feature.
	YesCanDo CanDoResponse = 1
	// NoCanDo is returned by CanDo queries when plugin/host is certain that it does not support a feature.
	NoCanDo CanDoResponse = -1
	// MaybeCanDo is returned by CanDo queries when plugin/host doesn't know if it supports a particular feature.
	MaybeCanDo CanDoResponse = 0
)

type (
	// Parameter refers to plugin parameter that can be mutated in the pipe.
	Parameter struct {
		Name              string
		Unit              string
		Value             float32
		NotAutomated      bool
		GetValueLabelFunc func(value float32) string
		GetValueFunc      func(value float32) float32
	}

	// Preset refers to plugin presets.
	Preset struct {
		name string
	}
)

// GetValue should be called in ProcessDoubleFunc or ProcessFloatFunc and will be called in GetDisplayVal.
// It returns the plain Value or the return value of GetValueFunc, when it was set for the Parameter
func (e Parameter) GetValue() float32 {
	if e.GetValueFunc == nil {
		return e.Value
	}
	return e.GetValueFunc(e.Value)
}

// GetValueLabel will be called in HostOpcode plugGetParamDisplay. Return a string formatted float value or the return
// value of GetValueLabelFunc, when it was set for the Parameter
func (e Parameter) GetValueLabel() string {
	if e.GetValueLabelFunc == nil {
		return fmt.Sprintf("%f", e.GetValue())
	}
	return e.GetValueLabelFunc(e.GetValue())
}

func trimNull(s string) string {
	return strings.Trim(s, "\x00")
}
