//go:generate stringer -type=PluginOpcode,HostOpcode -output=opcode_string.go

package vst2

const (
	maxVendorStrLen  = 64 // used for #plugGetVendorString, #audioMasterGetVendorString
	maxProductStrLen = 64 // used for #plugGetProductString, #audioMasterGetProductString
	maxPluginNameLen = 32 // used for #plugGetPluginName

	maxFileNameLen = 100 // used for #VstAudioFile->name
)

type (
	// 8 bytes ascii string.
	ascii8 [8]byte

	// 24 bytes ascii string.
	ascii24 [24]byte

	// 32 bytes ascii string.
	ascii32 [32]byte

	// 64 bytes ascii string.
	ascii64 [64]byte
)

func (s ascii8) String() string {
	return trimNull(string(s[:]))
}

func (s ascii24) String() string {
	return trimNull(string(s[:]))
}

func (s ascii32) String() string {
	return trimNull(string(s[:]))
}

func (s ascii64) String() string {
	return trimNull(string(s[:]))
}

// PluginOpcode is sent by host in dispatch call to plugin.
// It reflects APluginOpcodes and APluginXOpcodes opcodes values.
type PluginOpcode uint32

const (
	// PlugOpen passed to open the plugin.
	plugOpen PluginOpcode = iota
	// PlugClose passed to close the plugin.
	plugClose

	// PlugSetProgram passed to set program.
	// Value: new program number.
	plugSetProgram
	// PlugGetProgram passed to get program.
	// Return: current program number.
	plugGetProgram
	// PlugSetProgramName passed to set new program name.
	// Ptr: *[maxProgNameLen]byte buffer with new program name.
	plugSetProgramName
	// PlugGetProgramName passed to get current program name.
	// Ptr: *[maxProgNameLen]byte buffer for current program name.
	plugGetProgramName

	// PlugGetParamLabel passed to get parameter unit label: "db", "ms", etc.
	// Index: parameter index.
	// Ptr: *[maxParamStrLen]byte buffer for parameter unit label.
	plugGetParamLabel
	// PlugGetParamDisplay passed to get parameter value label: "0.5", "HALL", etc.
	// Index: parameter index.
	// Ptr: *[maxParamStrLen]byte buffer for parameter value label.
	plugGetParamDisplay
	// PlugGetParamName passed to get parameter label: "Release", "Gain", etc.
	// Index: parameter index.
	// Ptr: *[maxParamStrLen]byte buffer for parameter label.
	plugGetParamName

	// deprecated in VST v2.4
	plugGetVu

	// PlugSetSampleRate passed to set new sample rate.
	// Opt: new sample rate value.
	plugSetSampleRate
	// PlugSetBufferSize passed to set new buffer size.
	// Value: new buffer size value.
	plugSetBufferSize
	// PlugStateChanged passed when plugin's state changed.
	// Value: 0 means disabled, 1 means enabled.
	plugStateChanged

	// PlugEditGetRect passed to get editor size.
	// Ptr: ERect** receiving pointer to editor size.
	PlugEditGetRect
	// PlugEditOpen passed to get system dependent window pointer, eg HWND on Windows.
	// Ptr: window pointer.
	PlugEditOpen
	// PlugEditClose passed to close editor window.
	PlugEditClose

	// deprecated in VST v2.4
	plugEditDraw
	// deprecated in VST v2.4
	plugEditMouse
	// deprecated in VST v2.4
	plugEditKey

	// PlugEditIdle passed to notify plugin that host goes idle.
	PlugEditIdle

	// deprecated in VST v2.4
	plugEditTop
	// deprecated in VST v2.4
	plugEditSleep
	// deprecated in VST v2.4
	plugIdentify

	// PlugGetChunk passed to get chunk data.
	// Ptr: pointer for chunk data address (void**) uint8.
	// Index: 0 for bank, 1 for program.
	// Return: length of data.
	plugGetChunk
	// PlugSetChunk passed to set chunk data.
	// Ptr: pointer for chunk data address (void*).
	// Value: data size in bytes.
	// Index: 0 for bank, 1 for program.
	plugSetChunk

	// PlugProcessEvents passed to communicate events.
	// Ptr: *Events.
	PlugProcessEvents
	// PlugCanBeAutomated passed to check if parameter could be automated.
	// Index: parameter index.
	// Return: 1 for true, 0 for false.
	PlugCanBeAutomated
	// PlugString2Parameter passed to convert parameter value to string: "mono" to "1".
	// Index: parameter index.
	// Ptr: parameter string.
	// Return: true for success.
	PlugString2Parameter

	// deprecated in VST v2.4
	plugGetNumProgramCategories

	// PlugGetProgramNameIndexed passed to get program name by index.
	// Index: program index.
	// Ptr: *[maxProgNameLen]byte buffer for program name.
	// Return: true for success.
	plugGetProgramNameIndexed

	// deprecated in VST v2.4
	plugCopyProgram
	// deprecated in VST v2.4
	plugConnectInput
	// deprecated in VST v2.4
	plugConnectOutput

	// PlugGetInputProperties passed to check if certain input configuration is supported.
	// Index: input index.
	// Ptr: *PinProperties.
	// Return: 1 if supported.
	PlugGetInputProperties
	// PlugGetOutputProperties passed to check if certain output configuration is supported.
	// Index: output index.
	// Ptr: *PinProperties.
	// Return: 1 if supported.
	PlugGetOutputProperties
	// PlugGetPlugCategory passed to get plugin's category.
	// Return: VstPlugCategory value.
	PlugGetPlugCategory

	// deprecated in VST v2.4
	plugGetCurrentPosition
	// deprecated in VST v2.4
	plugGetDestinationBuffer

	// PlugOfflineNotify passed to notify about offline file processing.
	// Ptr: []AudioFile.
	// Value: count.
	// Index: start flag.
	PlugOfflineNotify
	// PlugOfflinePrepare passed to trigger offline processing preparation.
	// Ptr: []OfflineTask.
	// Value: count.
	PlugOfflinePrepare
	// PlugOfflineRun passed to trigger offline processing execution.
	// Ptr: []OfflineTask.
	// Value: count.
	PlugOfflineRun

	// PlugProcessVarIo passed to provide variable I/O processing (offline p.g. timestretching).
	// Ptr: *VariableIo.
	PlugProcessVarIo
	// PlugSetSpeakerArrangement passed to set speakers configuration.
	// Value: input *SpeakerArrangement.
	// Ptr: output *SpeakerArrangement.
	plugSetSpeakerArrangement

	// deprecated in VST v2.4
	plugSetBlockSizeAndSampleRate

	// PlugSetBypass passed to make plugin bypassed.
	// Value: 1 is bypass, 0 is no bypass.
	PlugSetBypass
	// PlugGetPluginName passed to get a name of the plugin.
	// Ptr: *[maxPluginNameLen]byte buffer for plugin name.
	PlugGetPluginName

	// deprecated in VST v2.4
	plugGetErrorText

	// PlugGetVendorString passed to get vendor string.
	// *[maxVendorStrLen]byte buffer for plugin vendor string.
	PlugGetVendorString
	// PlugGetProductString passed to get product string.
	// *[maxProductStrLen]byte buffer for plugin product string.
	PlugGetProductString
	// PlugGetVendorVersion passed to get vendor-specific version.
	// Return: vendor-specific version.
	PlugGetVendorVersion
	// PlugVendorSpecific passed to get vendor-specific string.
	// No definition, vendor specific handling.
	PlugVendorSpecific
	// PlugCanDo passed to check capabilities of plugin.
	// Ptr: "can do" string
	// Return: 0 is don't know, -1 is no, 1 is yes.
	PlugCanDo

	// PlugGetTailSize passed to check if "tail" data is expected.
	// Return: tail size (p.g. reverb time). 0 is default, 1 means no tail.
	PlugGetTailSize

	// deprecated in VST v2.4
	plugIdle
	// deprecated in VST v2.4
	plugGetIcon
	// deprecated in VST v2.4
	plugSetViewPosition

	// PlugGetParameterProperties passed to get parameter's properties.
	// Index: parameter index.
	// Ptr: *ParameterProperties.
	// Return: 1 if supported
	plugGetParameterProperties

	// deprecated in VST v2.4
	plugKeysRequired

	// PlugGetVstVersion passed to get VST version of plugin.
	// Return: VST version, 2400 for VST 2.4.
	PlugGetVstVersion

	// PlugEditKeyDown passed when key is pressed.
	// Index: ASCII character.
	// Value: virtual key.
	// Opt: ModifierKey flags.
	// Return: 1 if key used.
	PlugEditKeyDown
	// PlugEditKeyUp passed when key is released.
	// Index: ASCII character.
	// Value: virtual key.
	// Opt: ModifierKey flags.
	// Return: 1 if key used.
	PlugEditKeyUp
	// PlugSetEditKnobMode passed to set knob Mode.
	// Value: knob mode 0 is circular, 1 is circular relative, 2 is linear.
	PlugSetEditKnobMode

	// PlugGetMidiProgramName passed to get a name of used MIDI program.
	// Index: MIDI channel.
	// Ptr: *MidiProgramName.
	// Return: number of used programs, 0 if unsupported.
	PlugGetMidiProgramName
	// PlugGetCurrentMidiProgram passed to get a name of current MIDI program.
	// Index: MIDI channel.
	// Ptr: *MidiProgramName.
	// Return: index of current program .
	PlugGetCurrentMidiProgram
	// PlugGetMidiProgramCategory passed to get a category of MIDI program.
	// Index: MIDI channel.
	// Ptr: *MidiProgramCategory.
	// Return: number of used categories, 0 if unsupported.
	PlugGetMidiProgramCategory
	// PlugHasMidiProgramsChanged passed to check if MIDI program has changed.
	// Index: MIDI channel.
	// Return: 1 if the MidiProgramNames or MidiKeyNames have changed.
	PlugHasMidiProgramsChanged
	// PlugGetMidiKeyName passed to
	// Index: MIDI channel.
	// Ptr: *MidiKeyName.
	// Return: true if supported, false otherwise.
	PlugGetMidiKeyName

	// PlugBeginSetProgram passed before preset is loaded.
	PlugBeginSetProgram
	// PlugEndSetProgram passed after preset is loaded.
	PlugEndSetProgram

	// PlugGetSpeakerArrangement passed to get a speaker configuration of plugin.
	// Value: input *SpeakerArrangement.
	// Ptr: output *SpeakerArrangement.
	PlugGetSpeakerArrangement
	// PlugShellGetNextPlugin passed to get unique id of next plugin.
	// Ptr: *[maxProductStrLen]byte buffer for plug-in name.
	// Return: next plugin's unique ID.
	PlugShellGetNextPlugin

	// PlugStartProcess passed to indicate that the process call might be interrupted.
	PlugStartProcess
	// PlugStopProcess passed to indicate that process call is stopped.
	PlugStopProcess
	// PlugSetTotalSampleToProcess passed to identify a number of samples to process.
	// Value: number of samples to process. Called in offline mode before processing.
	PlugSetTotalSampleToProcess
	// PlugSetPanLaw passed to set pan law type and gain values.
	// Value: PanLawType value.
	// Opt: gain value.
	PlugSetPanLaw

	// PlugBeginLoadBank is passed when VST bank loaded.
	// Ptr: *PatchChunkInfo.
	// Return: -1 is bank can't be loaded, 1 is bank can be loaded, 0 is unsupported.
	PlugBeginLoadBank
	// PlugBeginLoadProgram is passed when VST preset loaded.
	// Ptr: *PatchChunkInfo.
	// Return: -1 is bank can't be loaded, 1 is bank can be loaded, 0 is unsupported.
	PlugBeginLoadProgram

	// PlugSetProcessPrecision passed to set processing precision.
	// Value: 0 if 32 bit, anything else if 64 bit.
	PlugSetProcessPrecision

	// PlugGetNumMidiInputChannels passed to get a number of used MIDI inputs.
	// Return: number of used MIDI input channels (1-15).
	PlugGetNumMidiInputChannels
	// PlugGetNumMidiOutputChannels passed to get a number of used MIDI outputs.
	// Return: number of used MIDI output channels (1-15).
	PlugGetNumMidiOutputChannels
)

// HostOpcode is sent by plugin in dispatch call to host.
// It reflects AudioMasterOpcodes and AudioMasterOpcodesX opcodes values.
type HostOpcode uint32

const (
	// HostAutomate passed to when parameter value is automated.
	// Index: parameter index.
	// Opt: parameter value.
	HostAutomate HostOpcode = iota
	// HostVersion passed to get VST version of host.
	// Return: host VST version (for example 2400 for VST 2.4).
	HostVersion
	// HostCurrentID passed when unique ID is requested.
	// Return: current unique identifier on shell plug-in.
	HostCurrentID
	// HostIdle passed to indicate that plugin does some modal action.
	HostIdle

	// deprecated in VST v2.4
	hostPinConnected
	_
	// deprecated in VST v2.4
	hostWantMidi

	// HostGetTime passed when plugin needs time info.
	// Return: *TimeInfo or null if not supported.
	// Value: request mask.
	HostGetTime
	// HostProcessEvents passed when plugin has MIDI events to process.
	// Ptr: *Events the events to be processed.
	// Return: 1 if supported and processed.
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
	// Return: 1 if supported.
	HostIOChanged

	// deprecated in VST v2.4
	hostNeedIdle

	// HostSizeWindow passed when host needs to resize plugin window.
	// Index: new width.
	// Value: new height.
	HostSizeWindow
	// HostGetSampleRate passed when plugin needs sample rate.
	// Return: current sample rate.
	HostGetSampleRate
	// HostGetBufferSize passed when plugin needs buffer size.
	// Return: current buffer size.
	HostGetBufferSize
	// HostGetInputLatency passed when plugin needs input latency.
	// Return: input latency in samples.
	HostGetInputLatency
	// HostGetOutputLatency passed when plugin needs output latency.
	// Return: output latency in samples.
	HostGetOutputLatency

	// deprecated in VST v2.4
	hostGetPreviousPlug
	// deprecated in VST v2.4
	hostGetNextPlug
	// deprecated in VST v2.4
	hostWillReplaceOrAccumulate

	// HostGetCurrentProcessLevel passed to get current process level.
	// Return: ProcessLevel value.
	HostGetCurrentProcessLevel
	// HostGetAutomationState passed to get current automation state.
	// Return: AutomationState value.
	HostGetAutomationState

	// HostOfflineStart is sent when plugin is ready for offline processing.
	// Index: number of new audio files.
	// Value: number of audio files.
	// Ptr: *AudioFile the host audio files. Flags can be updated from plugin.
	HostOfflineStart
	// HostOfflineRead is sent when plugin reads the data.
	// Index: boolean - if this value is true then the host will read the original
	//	file's samples, but if it is false it will read the samples which the plugin
	//	has written via HostOfflineWrite.
	// Value: see OfflineOption
	// Ptr: *OfflineTask describing the task.
	// Return: 1 on success.
	HostOfflineRead
	// HostOfflineWrite is sent when plugin writes the data.
	// Value: see OfflineOption
	// Ptr: *OfflineTask describing the task.
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
	// Ptr: *[maxVendorStrLen]byte buffer for host vendor name.
	HostGetVendorString
	// HostGetProductString is sent to get host product string.
	// Ptr: *[maxProductStrLen]byte buffer for host product name.
	HostGetProductString
	// HostGetVendorVersion is sent to get host version.
	// Return: vendor-specific version.
	HostGetVendorVersion
	// HostVendorSpecific is sent vendor-specific handling.
	HostVendorSpecific

	// deprecated in VST v2.4
	hostSetIcon

	// HostCanDo passed to check capabilities of host.
	// Ptr: "can do" string
	// Return: 0 is don't know, -1 is no, 1 is yes.
	HostCanDo
	// HostGetLanguage passed to get a language of the host.
	// Return: HostLanguage value.
	HostGetLanguage

	// deprecated in VST v2.4
	hostOpenWindow
	// deprecated in VST v2.4
	hostCloseWindow

	// HostGetDirectory passed to get the current directory.
	// Return: *[]byte with path.
	HostGetDirectory
	// HostUpdateDisplay passed to request host screen refresh.
	HostUpdateDisplay
	// HostBeginEdit passed to notify host that it should capture parameter changes.
	// Index: index of the control.
	// Return: true on success.
	HostBeginEdit
	// HostEndEdit passed to notify that control is no longer being changed.
	// Index: index of the control.
	// Return: true on success.
	HostEndEdit
	// HostOpenFileSelector passed to open the host file selector.
	// Ptr: *FileSelect.
	// Return: true on success.
	HostOpenFileSelector
	// HostCloseFileSelector passed to close the host file selector.
	// Ptr: *FileSelect.
	HostCloseFileSelector

	// deprecated in VST v2.4
	hostEditFile
	// deprecated in VST v2.4
	hostGetChunkFile
	// deprecated in VST v2.4
	hostGetInputSpeakerArrangement
)

func copyASCII(dst []byte, src string) {
	var read int
	for i := 0; i < len(src); i++ {
		if read == len(dst)-1 {
			break
		}
		if r := src[i]; r <= 127 {
			dst[read] = r
			read++
		}
	}
	dst[read] = '\x00'
}
