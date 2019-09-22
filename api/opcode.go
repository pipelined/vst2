//go:generate stringer -type=EffectOpcode,HostOpcode -output=opcode_string.go
package api

// EffectOpcode is sent by host in dispatch call to effect.
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
	// [return]: current program number.
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

	// deprecated in VST v2.4
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

	// deprecated in VST v2.4
	effEditDraw
	// deprecated in VST v2.4
	effEditMouse
	// deprecated in VST v2.4
	effEditKey

	// EffEditIdle passed to notify effect that host goes idle.
	EffEditIdle

	// deprecated in VST v2.4
	effEditTop
	// deprecated in VST v2.4
	effEditSleep
	// deprecated in VST v2.4
	effIdentify

	// EffGetChunk passed to get chunk data.
	// [ptr]: pointer for chunk data address (void**) uint8.
	// [index]: 0 for bank, 1 for program.
	// [return]: length of data.
	EffGetChunk
	// EffSetChunk passed to set chunk data.
	// [ptr]: pointer for chunk data address (void*).
	// [value]: data size in bytes.
	// [index]: 0 for bank, 1 for program.
	EffSetChunk

	// EffProcessEvents passed to communicate events.
	// [ptr]: *Events.
	EffProcessEvents
	// EffCanBeAutomated passed to check if parameter could be automated.
	// [index]: parameter index.
	// [return]: 1 for true, 0 for false.
	EffCanBeAutomated
	// EffString2Parameter passed to convert parameter value to string: "mono" to "1".
	// [index]: parameter index.
	// [ptr]: parameter string.
	// [return]: true for success.
	EffString2Parameter

	// deprecated in VST v2.4
	effGetNumProgramCategories

	// EffGetProgramNameIndexed passed to get program name by index.
	// [index]: program index.
	// [ptr]: *[maxProgNameLen]uint8 buffer for program name.
	// [return]: true for success.
	EffGetProgramNameIndexed

	// deprecated in VST v2.4
	effCopyProgram
	// deprecated in VST v2.4
	effConnectInput
	// deprecated in VST v2.4
	effConnectOutput

	// EffGetInputProperties passed to check if certain input configuration is supported.
	// [index]: input index.
	// [ptr]: *PinProperties.
	// [return]: 1 if supported.
	EffGetInputProperties
	// EffGetOutputProperties passed to check if certain output configuration is supported.
	// [index]: output index.
	// [ptr]: *PinProperties.
	// [return]: 1 if supported.
	EffGetOutputProperties
	// EffGetPlugCategory passed to get plugin's category.
	// [return]: VstPlugCategory value.
	EffGetPlugCategory

	// deprecated in VST v2.4
	effGetCurrentPosition
	// deprecated in VST v2.4
	effGetDestinationBuffer

	// EffOfflineNotify passed to notify about offline file processing.
	// [ptr]: []AudioFile.
	// [value]: count.
	// [index]: start flag.
	EffOfflineNotify
	// EffOfflinePrepare passed to trigger offline processing preparation.
	// [ptr]: []OfflineTask.
	// [value]: count.
	EffOfflinePrepare
	// EffOfflineRun passed to trigger offline processing execution.
	// [ptr]: []OfflineTask.
	// [value]: count.
	EffOfflineRun

	// EffProcessVarIo passed to provide variable I/O processing (offline e.g. timestretching).
	// [ptr]: *VariableIo.
	EffProcessVarIo
	// EffSetSpeakerArrangement passed to set speakers configuration.
	// [value]: input *SpeakerArrangement.
	// [ptr]: output *SpeakerArrangement.
	EffSetSpeakerArrangement

	// deprecated in VST v2.4
	effSetBlockSizeAndSampleRate

	// EffSetBypass passed to make effect bypassed.
	// [value]: 1 is bypass, 0 is no bypass.
	EffSetBypass
	// EffGetEffectName passed to get a name of the effect.
	// [ptr]: *[maxEffectNameLen]uint8 buffer for effect name.
	EffGetEffectName

	// deprecated in VST v2.4
	effGetErrorText

	// EffGetVendorString passed to get vendor string.
	// *[maxVendorStrLen]uint8 buffer for effect vendor string.
	EffGetVendorString
	// EffGetProductString passed to get product string.
	// *[maxProductStrLen]uint8 buffer for effect product string.
	EffGetProductString
	// EffGetVendorVersion passed to get vendor-specific version.
	// [return]: vendor-specific version.
	EffGetVendorVersion
	// EffVendorSpecific passed to get vendor-specific string.
	// No definition, vendor specific handling.
	EffVendorSpecific
	// EffCanDo passed to check capabilities of effect.
	// [ptr]: "can do" string
	// [return]: 0 is don't know, -1 is no, 1 is yes.
	EffCanDo

	// EffGetTailSize passed to check if "tail" data is expected.
	// [return]: tail size (e.g. reverb time). 0 is defualt, 1 means no tail.
	EffGetTailSize

	// deprecated in VST v2.4
	effIdle
	// deprecated in VST v2.4
	effGetIcon
	// deprecated in VST v2.4
	effSetViewPosition

	// EffGetParameterProperties passed to get parameter's properties.
	// [index]: parameter index.
	// [ptr]: *ParameterProperties.
	// [return]: 1 if supported
	EffGetParameterProperties

	// deprecated in VST v2.4
	effKeysRequired

	// EffGetVstVersion passed to get VST version of effect.
	// [return]: VST version, 2400 for VST 2.4.
	EffGetVstVersion

	// EffEditKeyDown passed when key is pressed.
	// [index]: ASCII character.
	// [value]: virtual key.
	// [opt]: modifiers.
	// [return]: 1 if key used.
	EffEditKeyDown
	// EffEditKeyUp passed when key is released.
	// [index]: ASCII character.
	// [value]: virtual key.
	// [opt]: modifiers.
	// [return]: 1 if key used.
	EffEditKeyUp
	// EffSetEditKnobMode passed to set knob Mode.
	// [value]: knob mode 0 is circular, 1 is circular relative, 2 is linear.
	EffSetEditKnobMode

	// EffGetMidiProgramName passed to get a name of used MIDI program.
	// [index]: MIDI channel.
	// [ptr]: *MidiProgramName.
	// [return]: number of used programs, 0 if unsupported.
	EffGetMidiProgramName
	// EffGetCurrentMidiProgram passed to get a name of current MIDI program.
	// [index]: MIDI channel.
	// [ptr]: *MidiProgramName.
	// [return]: index of current program .
	EffGetCurrentMidiProgram
	// EffGetMidiProgramCategory passed to get a category of MIDI program.
	// [index]: MIDI channel.
	// [ptr]: *MidiProgramCategory.
	// [return]: number of used categories, 0 if unsupported.
	EffGetMidiProgramCategory
	// EffHasMidiProgramsChanged passed to check if MIDI program has changed.
	// [index]: MIDI channel.
	// [return]: 1 if the MidiProgramNames or MidiKeyNames have changed.
	EffHasMidiProgramsChanged
	// EffGetMidiKeyName passed to
	// [index]: MIDI channel.
	// [ptr]: *MidiKeyName.
	// [return]: true if supported, false otherwise.
	EffGetMidiKeyName

	// EffBeginSetProgram passed before preset is loaded.
	EffBeginSetProgram
	// EffEndSetProgram passed after preset is loaded.
	EffEndSetProgram

	// EffGetSpeakerArrangement passed to get a speaker configuration of plugin.
	// [value]: input *SpeakerArrangement.
	// [ptr]: output *SpeakerArrangement.
	EffGetSpeakerArrangement
	// EffShellGetNextPlugin passed to get unique id of next plugin.
	// [ptr]: *[maxProductStrLen]uint8 buffer for plug-in name.
	// [return]: next plugin's unique ID.
	EffShellGetNextPlugin

	// EffStartProcess passed to indicate that the process call might be interrupted.
	EffStartProcess
	// EffStopProcess passed to indicate that process call is stopped.
	EffStopProcess
	// EffSetTotalSampleToProcess passed to identify a number of samples to process.
	// [value]: number of samples to process. Called in offline mode before processing.
	EffSetTotalSampleToProcess
	// EffSetPanLaw passed to set pan law type and gain values.
	// [value]: PanLawType value.
	// [opt]: gain value.
	EffSetPanLaw

	// EffBeginLoadBank is passed when VST bank loaded.
	// [ptr]: *PatchChunkInfo.
	// [return]: -1 is bank can't be loaded, 1 is bank can be loaded, 0 is unsupported.
	EffBeginLoadBank
	// EffBeginLoadProgram is passed when VST preset loaded.
	// [ptr]: *PatchChunkInfo.
	// [return]: -1 is bank can't be loaded, 1 is bank can be loaded, 0 is unsupported.
	EffBeginLoadProgram

	// EffSetProcessPrecision passed to set processing precision.
	// [value]: 0 if 32 bit, anything else if 64 bit.
	EffSetProcessPrecision

	// EffGetNumMidiInputChannels passed to get a number of used MIDI inputs.
	// [return]: number of used MIDI input channels (1-15).
	EffGetNumMidiInputChannels
	// EffGetNumMidiOutputChannels passed to get a number of used MIDI outputs.
	// [return]: number of used MIDI output channels (1-15).
	EffGetNumMidiOutputChannels
)

// HostOpcode is sent by plugin in dispatch call to host.
// It relfects AudioMasterOpcodes and AudioMasterOpcodesX opcodes values.
type HostOpcode uint64

const (
	// HostAutomate passed to when parameter value is automated.
	// [index]: parameter index.
	// [opt]: parameter value.
	HostAutomate HostOpcode = iota
	// HostVersion passed to get VST version of host.
	// [return]: host VST version (for example 2400 for VST 2.4).
	HostVersion
	// HostCurrentID passed when unique ID is requested.
	// [return]: current unique identifier on shell plug-in.
	HostCurrentID
	// HostIdle passed to indicate that plugin does some modal action.
	HostIdle

	// deprecated in VST v2.4
	hostPinConnected
	_
	// deprecated in VST v2.4
	hostWantMidi

	// HostGetTime passed when plugin needs time info.
	// [return]: *TimeInfo or null if not supported.
	// [value]: request mask.
	HostGetTime
	// HostProcessEvents passed when plugin has MIDI events to process.
	// [ptr]: *Events the events to be processed.
	// [return]: 1 if supported and processed.
	HostProcessEvents

	// deprecated in VST v2.4
	hostSetTime
	// deprecated in VST v2.4
	hostTempoAt
	// deprecated in VST v2.4
	hostGetNumAutomatableParameters
	// deprecated in VST v2.4
	hostGetParameterQuantization

	// HostIOChanged passed when plugin's IO setup has changed.
	// [return]: 1 if supported.
	HostIOChanged

	// deprecated in VST v2.4
	hostNeedIdle

	// HostSizeWindow passed when host needs to resize plugin window.
	// [index]: new width.
	// [value]: new height.
	HostSizeWindow
	// HostGetSampleRate passed when plugin needs sample rate.
	// [return]: current sample rate.
	HostGetSampleRate
	// HostGetBlockSize passed when plugin needs buffer size.
	// [return]: current buffer size.
	HostGetBlockSize
	// HostGetInputLatency passed when plugin needs input latency.
	// [return]: input latency in samples.
	HostGetInputLatency
	// HostGetOutputLatency passed when plugin needs output latency.
	// [return]: output latency in samples.
	HostGetOutputLatency

	// deprecated in VST v2.4
	hostGetPreviousPlug
	// deprecated in VST v2.4
	hostGetNextPlug
	// deprecated in VST v2.4
	hostWillReplaceOrAccumulate

	// HostGetCurrentProcessLevel passed to get current process level.
	// [return]: ProcessLevel value.
	HostGetCurrentProcessLevel
	// HostGetAutomationState passed to get current automation state.
	// [return]: AutomationState value.
	HostGetAutomationState

	// HostOfflineStart is sent when plugin is ready for offline processing.
	// [index]: number of new audio files.
	// [value]: number of audio files.
	// [ptr]: *AudioFile the host audio files. Flags can be updated from plugin.
	HostOfflineStart
	// HostOfflineRead is sent when plugin reads the data.
	// [index]: boolean - if this value is true then the host will read the original
	//	file's samples, but if it is false it will read the samples which the plugin
	//	has written via HostOfflineWrite.
	// [value]: see OfflineOption
	// [ptr]: *OfflineTask describing the task.
	// [return]: 1 on success.
	HostOfflineRead
	// HostOfflineWrite is sent when plugin writes the data.
	// [value]: see OfflineOption
	// [ptr]: *OfflineTask describing the task.
	HostOfflineWrite
	// HostOfflineGetCurrentPass is unknown.
	HostOfflineGetCurrentPass
	// HostOfflineGetCurrentMetaPass is unknown.
	HostOfflineGetCurrentMetaPass

	// deprecated in VST v2.4
	hostSetOutputSampleRate
	// deprecated in VST v2.4
	hostGetOutputSpeakerArrangement

	// HostGetVendorString is sent to get host vendor string.
	// [ptr]: *[maxVendorStrLen]uint8 buffer for host vendor name.
	HostGetVendorString
	// HostGetProductString is sent to get host product string.
	// [ptr]: *[maxProductStrLen]uint8 buffer for host product name.
	HostGetProductString
	// HostGetVendorVersion is sent to get host version.
	// [return]: vendor-specific version.
	HostGetVendorVersion
	// HostVendorSpecific is sent vendor-specific handling.
	HostVendorSpecific

	// deprecated in VST v2.4
	hostSetIcon

	// HostCanDo passed to check capabilities of host.
	// [ptr]: "can do" string
	// [return]: 0 is don't know, -1 is no, 1 is yes.
	HostCanDo
	// HostGetLanguage passed to get a language of the host.
	// [return]: HostLanguage value.
	HostGetLanguage

	// deprecated in VST v2.4
	hostOpenWindow
	// deprecated in VST v2.4
	hostCloseWindow

	// HostGetDirectory passed to get the current directory.
	// [return]: *[]uint8 with path.
	HostGetDirectory
	// HostUpdateDisplay passed to request host screen refresh.
	HostUpdateDisplay
	// HostBeginEdit passed to notify host that it should capture parameter changes.
	// [index]: index of the control.
	// [return]: true on success.
	HostBeginEdit
	// HostEndEdit passed to notify that control is no longer being changed.
	// [index]: index of the control.
	// [return]: true on success.
	HostEndEdit
	// HostOpenFileSelector passed to open the host file selector.
	// [ptr]: *FileSelect.
	// [return]: true on success.
	HostOpenFileSelector
	// HostCloseFileSelector passed to close the host file selector.
	// [ptr]: *FileSelect.
	HostCloseFileSelector

	// deprecated in VST v2.4
	hostEditFile
	// deprecated in VST v2.4
	hostGetChunkFile
	// deprecated in VST v2.4
	hostGetInputSpeakerArrangement
)
