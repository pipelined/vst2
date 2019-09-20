//go:generate stringer -type=EffectOpcode,HostOpcode -output=opcode_string.go
package api

const (
	maxProgNameLen   = 24 ///< used for #effGetProgramName, #effSetProgramName, #effGetProgramNameIndexed
	maxParamStrLen   = 8  ///< used for #effGetParamLabel, #effGetParamDisplay, #effGetParamName
	maxVendorStrLen  = 64 ///< used for #effGetVendorString, #audioMasterGetVendorString
	maxProductStrLen = 64 ///< used for #effGetProductString, #audioMasterGetProductString
	maxEffectNameLen = 32 ///< used for #effGetEffectName
)

// EffectOpcode is sent by host in dispatch call to plugin.
// It relfects AEffectOpcodes and AEffectXOpcodes opcodes values.
type EffectOpcode uint64

const (
	// EffOpen passed to open the plugin.
	EffOpen EffectOpcode = iota
	// EffClose passed to close the plugin.
	EffClose

	// EffSetProgram passed to set program.
	// [value]: new program number.
	EffSetProgram
	// EffGetProgram passed to get program.
	// [return value]: current program number.
	EffGetProgram
	// EffSetProgramName passed to set new program name.
	// [ptr]: *[maxProgNameLen]uint8 buffer with new program name.
	EffSetProgramName
	// EffGetProgramName passed to get current program name.
	// [ptr]: *[maxProgNameLen]uint8 buffer for current program name.
	EffGetProgramName

	// EffGetParamLabel passed to get parameter unit label: "db", "ms", etc.
	// [index]: parameter index.
	// [ptr]: *[maxParamStrLen]uint8 buffer for parameter unit label.
	EffGetParamLabel
	// EffGetParamDisplay passed to get parameter value label: "0.5", "HALL", etc.
	// [index]: parameter index.
	// [ptr]: *[maxParamStrLen]uint8 buffer for parameter value label.
	EffGetParamDisplay
	// EffGetParamName passed to get parameter label: "Release", "Gain", etc.
	// [index]: parameter index.
	// [ptr]: *[maxParamStrLen]uint8 buffer for parameter label.
	EffGetParamName

	// effGetVu is deprecated in VST v2.4
	effGetVu

	// EffSetSampleRate passed to set new sample rate.
	// [opt]: new sample rate value.
	EffSetSampleRate
	// EffSetBufferSize passed to set new buffer size.
	// [value]: new buffer size value.
	EffSetBufferSize
	// EffStateChanged passed when plugin's state changed.
	// [value]: 0 means disabled, 1 means enabled.
	EffStateChanged

	// EffEditGetRect passed to get editor size.
	// [ptr]: ERect** receiving pointer to editor size.
	EffEditGetRect
	// EffEditOpen passed to get system dependent window pointer, eg HWND on Windows.
	// [ptr]: window pointer.
	EffEditOpen
	// EffEditClose passed to close editor window.
	EffEditClose

	// effEditDraw is deprecated in VST v2.4
	effEditDraw
	// effEditMouse is deprecated in VST v2.4
	effEditMouse
	// effEditKey is deprecated in VST v2.4
	effEditKey

	// EffEditIdle passed to notify effect that host goes idle.
	EffEditIdle

	// effEditTop is deprecated in VST v2.4
	effEditTop
	// effEditSleep is deprecated in VST v2.4
	effEditSleep
	// effIdentify is deprecated in VST v2.4
	effIdentify

	// EffGetChunk passed to get chunk data.
	// [ptr]: pointer for chunk data address (void**) uint8.
	// [index]: 0 for bank, 1 for program.
	// [return value]: length of data.
	EffGetChunk
	// EffSetChunk passed to set chunk data.
	// [ptr]: pointer for chunk data address (void*).
	// [value]: data size in bytes.
	// [index]: 0 for bank, 1 for program.
	EffSetChunk

	// EffProcessEvents passed to communicate events.
	// [ptr]: #VstEvents*
	EffProcessEvents
	// EffCanBeAutomated passed to check if parameter could be automated.
	// [index]: parameter index.
	// [return value]: 1 for true, 0 for false.
	EffCanBeAutomated
	// EffString2Parameter passed to convert parameter value to string: "mono" to "1".
	// [index]: parameter index.
	// [ptr]: parameter string.
	// [return value]: true for success.
	EffString2Parameter

	// effGetNumProgramCategories is deprecated in VST v2.4
	effGetNumProgramCategories

	// EffGetProgramNameIndexed passed to get program name by index.
	// [index]: program index.
	// [ptr]: *[maxProgNameLen]uint8 buffer for program name.
	// [return value]: true for success.
	EffGetProgramNameIndexed

	// effCopyProgram is deprecated in VST v2.4
	effCopyProgram
	// effConnectInput is deprecated in VST v2.4
	effConnectInput
	// effConnectOutput is deprecated in VST v2.4
	effConnectOutput

	// EffGetInputProperties passed to check if certain input configuration is supported.
	// [index]: input index.
	// [ptr]: #VstPinProperties*.
	// [return value]: 1 if supported.
	EffGetInputProperties
	// EffGetOutputProperties passed to check if certain output configuration is supported.
	// [index]: output index.
	// [ptr]: #VstPinProperties*.
	// [return value]: 1 if supported.
	EffGetOutputProperties
	// EffGetPlugCategory passed to get plugin's category.
	// [return value]: VstPlugCategory value.
	EffGetPlugCategory

	// effGetCurrentPosition is deprecated in VST v2.4
	effGetCurrentPosition
	// effGetDestinationBuffer is deprecated in VST v2.4
	effGetDestinationBuffer

	// EffOfflineNotify passed to notify about offline file processing.
	// [ptr]: #VstAudioFile array.
	// [value]: count.
	// [index]: start flag.
	EffOfflineNotify
	// EffOfflinePrepare passed to trigger offline processing preparation.
	// [ptr]: #VstOfflineTask array.
	// [value]: count.
	EffOfflinePrepare
	// EffOfflineRun passed to trigger offline processing execution.
	// [ptr]: #VstOfflineTask array.
	// [value]: count.
	EffOfflineRun

	// EffProcessVarIo passed to provide variable I/O processing (offline e.g. timestretching).
	// [ptr]: #VstVariableIo*.
	EffProcessVarIo
	// EffSetSpeakerArrangement passed to set speakers configuration.
	// [value]: input #VstSpeakerArrangement*.
	// [ptr]: output #VstSpeakerArrangement*.
	EffSetSpeakerArrangement

	// effSetBlockSizeAndSampleRate is deprecated in VST v2.4
	effSetBlockSizeAndSampleRate

	EffSetBypass
	EffGetEffectName
	// effGetErrorText is deprecated in VST v2.4
	effGetErrorText
	EffGetVendorString
	EffGetProductString
	EffGetVendorVersion
	EffVendorSpecific
	EffCanDo
	EffGetTailSize
	// effIdle is deprecated in VST v2.4
	effIdle
	// effGetIcon is deprecated in VST v2.4
	effGetIcon
	// effSetViewPosition is deprecated in VST v2.4
	effSetViewPosition
	EffGetParameterProperties
	// effKeysRequired is deprecated in VST v2.4
	effKeysRequired
	EffGetVstVersion
	EffEditKeyDown
	EffEditKeyUp
	EffSetEditKnobMode
	EffGetMidiProgramName
	EffGetCurrentMidiProgram
	EffGetMidiProgramCategory
	EffHasMidiProgramsChanged
	EffGetMidiKeyName
	EffBeginSetProgram
	EffEndSetProgram
	EffGetSpeakerArrangement
	EffShellGetNextPlugin
	EffStartProcess
	EffStopProcess
	EffSetTotalSampleToProcess
	EffSetPanLaw
	EffBeginLoadBank
	EffBeginLoadProgram
	EffSetProcessPrecision
	EffGetNumMidiInputChannels
	EffGetNumMidiOutputChannels
)

// HostOpcode used to wrap C opcodes values
type HostOpcode uint64

// Constants for audio master callback opcodes
// Plugin -> Host
const (
	// hostOpcodes opcodes
	hostAutomate HostOpcode = iota
	HostVersion
	HostCurrentID
	HostIdle
	// hostPinConnected is deprecated in VST v2.4
	hostPinConnected

	_
	// hostOpcodesX opcodes
	hostWantMidi
	HostGetTime
	HostProcessEvents
	// hostSetTime is deprecated in VST v2.4
	hostSetTime
	// hostTempoAt is deprecated in VST v2.4
	hostTempoAt
	// hostGetNumAutomatableParameters is deprecated in VST v2.4
	hostGetNumAutomatableParameters
	// hostGetParameterQuantization is deprecated in VST v2.4
	hostGetParameterQuantization
	HostIOChanged
	// hostNeedIdle is deprecated in VST v2.4
	hostNeedIdle
	HostSizeWindow
	HostGetSampleRate
	HostGetBlockSize
	HostGetInputLatency
	HostGetOutputLatency
	// hostGetPreviousPlug is deprecated in VST v2.4
	hostGetPreviousPlug
	// hostGetNextPlug is deprecated in VST v2.4
	hostGetNextPlug
	// hostWillReplaceOrAccumulate is deprecated in VST v2.4
	hostWillReplaceOrAccumulate
	HostGetCurrentProcessLevel
	HostGetAutomationState
	HostOfflineStart
	HostOfflineRead
	HostOfflineWrite
	HostOfflineGetCurrentPass
	HostOfflineGetCurrentMetaPass
	// hostSetOutputSampleRate is deprecated in VST v2.4
	hostSetOutputSampleRate
	// hostGetOutputSpeakerArrangement is deprecated in VST v2.4
	hostGetOutputSpeakerArrangement
	HostGetVendorString
	HostGetProductString
	HostGetVendorVersion
	HostVendorSpecific
	// hostSetIcon is deprecated in VST v2.4
	hostSetIcon
	HostCanDo
	HostGetLanguage
	// hostOpenWindow is deprecated in VST v2.4
	hostOpenWindow
	// hostCloseWindow is deprecated in VST v2.4
	hostCloseWindow
	HostGetDirectory
	HostUpdateDisplay
	HostBeginEdit
	HostEndEdit
	HostOpenFileSelector
	HostCloseFileSelector
	// hostEditFile is deprecated in VST v2.4
	hostEditFile
	// hostGetChunkFile is deprecated in VST v2.4
	hostGetChunkFile
	// hostGetInputSpeakerArrangement is deprecated in VST v2.4
	hostGetInputSpeakerArrangement
)
