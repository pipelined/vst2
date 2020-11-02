//go:generate stringer -type=EffectOpcode,HostOpcode -output=opcode_string.go

package sdk

/*
#include <stdlib.h>
*/
import "C"
import "unsafe"

// EffectOpcode is sent by host in dispatch call to effect.
// It reflects AEffectOpcodes and AEffectXOpcodes opcodes values.
type EffectOpcode uint64

const (
	// EffOpen passed to open the plugin.
	EffOpen EffectOpcode = iota
	// EffClose passed to close the plugin.
	EffClose

	// EffSetProgram passed to set program.
	// Value: new program number.
	EffSetProgram
	// EffGetProgram passed to get program.
	// Return: current program number.
	EffGetProgram
	// EffSetProgramName passed to set new program name.
	// Ptr: *[maxProgNameLen]byte buffer with new program name.
	EffSetProgramName
	// EffGetProgramName passed to get current program name.
	// Ptr: *[maxProgNameLen]byte buffer for current program name.
	EffGetProgramName

	// EffGetParamLabel passed to get parameter unit label: "db", "ms", etc.
	// Index: parameter index.
	// Ptr: *[maxParamStrLen]byte buffer for parameter unit label.
	EffGetParamLabel
	// EffGetParamDisplay passed to get parameter value label: "0.5", "HALL", etc.
	// Index: parameter index.
	// Ptr: *[maxParamStrLen]byte buffer for parameter value label.
	EffGetParamDisplay
	// EffGetParamName passed to get parameter label: "Release", "Gain", etc.
	// Index: parameter index.
	// Ptr: *[maxParamStrLen]byte buffer for parameter label.
	EffGetParamName

	// deprecated in VST v2.4
	effGetVu

	// EffSetSampleRate passed to set new sample rate.
	// Opt: new sample rate value.
	EffSetSampleRate
	// EffSetBufferSize passed to set new buffer size.
	// Value: new buffer size value.
	EffSetBufferSize
	// EffStateChanged passed when plugin's state changed.
	// Value: 0 means disabled, 1 means enabled.
	EffStateChanged

	// EffEditGetRect passed to get editor size.
	// Ptr: ERect** receiving pointer to editor size.
	EffEditGetRect
	// EffEditOpen passed to get system dependent window pointer, eg HWND on Windows.
	// Ptr: window pointer.
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
	// Ptr: pointer for chunk data address (void**) uint8.
	// Index: 0 for bank, 1 for program.
	// Return: length of data.
	EffGetChunk
	// EffSetChunk passed to set chunk data.
	// Ptr: pointer for chunk data address (void*).
	// Value: data size in bytes.
	// Index: 0 for bank, 1 for program.
	EffSetChunk

	// EffProcessEvents passed to communicate events.
	// Ptr: *Events.
	EffProcessEvents
	// EffCanBeAutomated passed to check if parameter could be automated.
	// Index: parameter index.
	// Return: 1 for true, 0 for false.
	EffCanBeAutomated
	// EffString2Parameter passed to convert parameter value to string: "mono" to "1".
	// Index: parameter index.
	// Ptr: parameter string.
	// Return: true for success.
	EffString2Parameter

	// deprecated in VST v2.4
	effGetNumProgramCategories

	// EffGetProgramNameIndexed passed to get program name by index.
	// Index: program index.
	// Ptr: *[maxProgNameLen]byte buffer for program name.
	// Return: true for success.
	EffGetProgramNameIndexed

	// deprecated in VST v2.4
	effCopyProgram
	// deprecated in VST v2.4
	effConnectInput
	// deprecated in VST v2.4
	effConnectOutput

	// EffGetInputProperties passed to check if certain input configuration is supported.
	// Index: input index.
	// Ptr: *PinProperties.
	// Return: 1 if supported.
	EffGetInputProperties
	// EffGetOutputProperties passed to check if certain output configuration is supported.
	// Index: output index.
	// Ptr: *PinProperties.
	// Return: 1 if supported.
	EffGetOutputProperties
	// EffGetPlugCategory passed to get plugin's category.
	// Return: VstPlugCategory value.
	EffGetPlugCategory

	// deprecated in VST v2.4
	effGetCurrentPosition
	// deprecated in VST v2.4
	effGetDestinationBuffer

	// EffOfflineNotify passed to notify about offline file processing.
	// Ptr: []AudioFile.
	// Value: count.
	// Index: start flag.
	EffOfflineNotify
	// EffOfflinePrepare passed to trigger offline processing preparation.
	// Ptr: []OfflineTask.
	// Value: count.
	EffOfflinePrepare
	// EffOfflineRun passed to trigger offline processing execution.
	// Ptr: []OfflineTask.
	// Value: count.
	EffOfflineRun

	// EffProcessVarIo passed to provide variable I/O processing (offline e.g. timestretching).
	// Ptr: *VariableIo.
	EffProcessVarIo
	// EffSetSpeakerArrangement passed to set speakers configuration.
	// Value: input *SpeakerArrangement.
	// Ptr: output *SpeakerArrangement.
	EffSetSpeakerArrangement

	// deprecated in VST v2.4
	effSetBlockSizeAndSampleRate

	// EffSetBypass passed to make effect bypassed.
	// Value: 1 is bypass, 0 is no bypass.
	EffSetBypass
	// EffGetEffectName passed to get a name of the effect.
	// Ptr: *[maxEffectNameLen]byte buffer for effect name.
	EffGetEffectName

	// deprecated in VST v2.4
	effGetErrorText

	// EffGetVendorString passed to get vendor string.
	// *[maxVendorStrLen]byte buffer for effect vendor string.
	EffGetVendorString
	// EffGetProductString passed to get product string.
	// *[maxProductStrLen]byte buffer for effect product string.
	EffGetProductString
	// EffGetVendorVersion passed to get vendor-specific version.
	// Return: vendor-specific version.
	EffGetVendorVersion
	// EffVendorSpecific passed to get vendor-specific string.
	// No definition, vendor specific handling.
	EffVendorSpecific
	// EffCanDo passed to check capabilities of effect.
	// Ptr: "can do" string
	// Return: 0 is don't know, -1 is no, 1 is yes.
	EffCanDo

	// EffGetTailSize passed to check if "tail" data is expected.
	// Return: tail size (e.g. reverb time). 0 is default, 1 means no tail.
	EffGetTailSize

	// deprecated in VST v2.4
	effIdle
	// deprecated in VST v2.4
	effGetIcon
	// deprecated in VST v2.4
	effSetViewPosition

	// EffGetParameterProperties passed to get parameter's properties.
	// Index: parameter index.
	// Ptr: *ParameterProperties.
	// Return: 1 if supported
	EffGetParameterProperties

	// deprecated in VST v2.4
	effKeysRequired

	// EffGetVstVersion passed to get VST version of effect.
	// Return: VST version, 2400 for VST 2.4.
	EffGetVstVersion

	// EffEditKeyDown passed when key is pressed.
	// Index: ASCII character.
	// Value: virtual key.
	// Opt: ModifierKey flags.
	// Return: 1 if key used.
	EffEditKeyDown
	// EffEditKeyUp passed when key is released.
	// Index: ASCII character.
	// Value: virtual key.
	// Opt: ModifierKey flags.
	// Return: 1 if key used.
	EffEditKeyUp
	// EffSetEditKnobMode passed to set knob Mode.
	// Value: knob mode 0 is circular, 1 is circular relative, 2 is linear.
	EffSetEditKnobMode

	// EffGetMidiProgramName passed to get a name of used MIDI program.
	// Index: MIDI channel.
	// Ptr: *MidiProgramName.
	// Return: number of used programs, 0 if unsupported.
	EffGetMidiProgramName
	// EffGetCurrentMidiProgram passed to get a name of current MIDI program.
	// Index: MIDI channel.
	// Ptr: *MidiProgramName.
	// Return: index of current program .
	EffGetCurrentMidiProgram
	// EffGetMidiProgramCategory passed to get a category of MIDI program.
	// Index: MIDI channel.
	// Ptr: *MidiProgramCategory.
	// Return: number of used categories, 0 if unsupported.
	EffGetMidiProgramCategory
	// EffHasMidiProgramsChanged passed to check if MIDI program has changed.
	// Index: MIDI channel.
	// Return: 1 if the MidiProgramNames or MidiKeyNames have changed.
	EffHasMidiProgramsChanged
	// EffGetMidiKeyName passed to
	// Index: MIDI channel.
	// Ptr: *MidiKeyName.
	// Return: true if supported, false otherwise.
	EffGetMidiKeyName

	// EffBeginSetProgram passed before preset is loaded.
	EffBeginSetProgram
	// EffEndSetProgram passed after preset is loaded.
	EffEndSetProgram

	// EffGetSpeakerArrangement passed to get a speaker configuration of plugin.
	// Value: input *SpeakerArrangement.
	// Ptr: output *SpeakerArrangement.
	EffGetSpeakerArrangement
	// EffShellGetNextPlugin passed to get unique id of next plugin.
	// Ptr: *[maxProductStrLen]byte buffer for plug-in name.
	// Return: next plugin's unique ID.
	EffShellGetNextPlugin

	// EffStartProcess passed to indicate that the process call might be interrupted.
	EffStartProcess
	// EffStopProcess passed to indicate that process call is stopped.
	EffStopProcess
	// EffSetTotalSampleToProcess passed to identify a number of samples to process.
	// Value: number of samples to process. Called in offline mode before processing.
	EffSetTotalSampleToProcess
	// EffSetPanLaw passed to set pan law type and gain values.
	// Value: PanLawType value.
	// Opt: gain value.
	EffSetPanLaw

	// EffBeginLoadBank is passed when VST bank loaded.
	// Ptr: *PatchChunkInfo.
	// Return: -1 is bank can't be loaded, 1 is bank can be loaded, 0 is unsupported.
	EffBeginLoadBank
	// EffBeginLoadProgram is passed when VST preset loaded.
	// Ptr: *PatchChunkInfo.
	// Return: -1 is bank can't be loaded, 1 is bank can be loaded, 0 is unsupported.
	EffBeginLoadProgram

	// EffSetProcessPrecision passed to set processing precision.
	// Value: 0 if 32 bit, anything else if 64 bit.
	EffSetProcessPrecision

	// EffGetNumMidiInputChannels passed to get a number of used MIDI inputs.
	// Return: number of used MIDI input channels (1-15).
	EffGetNumMidiInputChannels
	// EffGetNumMidiOutputChannels passed to get a number of used MIDI outputs.
	// Return: number of used MIDI output channels (1-15).
	EffGetNumMidiOutputChannels
)

// HostOpcode is sent by plugin in dispatch call to host.
// It reflects AudioMasterOpcodes and AudioMasterOpcodesX opcodes values.
type HostOpcode uint64

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
	// HostGetBlockSize passed when plugin needs buffer size.
	// Return: current buffer size.
	HostGetBlockSize
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

// Open executes the EffOpen opcode.
func (e *Effect) Open() {
	e.Dispatch(EffOpen, 0, 0, nil, 0.0)
}

// Close cleans up C refs for plugin
func (e *Effect) Close() {
	e.Dispatch(EffClose, 0, 0, nil, 0.0)
	mutex.Lock()
	delete(callbacks, e)
	mutex.Unlock()
}

// Start the plugin.
func (e *Effect) Start() {
	e.Dispatch(EffStateChanged, 0, 1, nil, 0)
}

// Stop the plugin.
func (e *Effect) Stop() {
	e.Dispatch(EffStateChanged, 0, 0, nil, 0)
}

// SetBufferSize sets a buffer size per channel.
func (e *Effect) SetBufferSize(bufferSize int) {
	e.Dispatch(EffSetBufferSize, 0, Value(bufferSize), nil, 0)
}

// SetSampleRate sets a sample rate for plugin.
func (e *Effect) SetSampleRate(sampleRate int) {
	e.Dispatch(EffSetSampleRate, 0, 0, nil, Opt(sampleRate))
}

// SetSpeakerArrangement creates and passes SpeakerArrangement structures to plugin
func (e *Effect) SetSpeakerArrangement(in, out *SpeakerArrangement) {
	e.Dispatch(EffSetSpeakerArrangement, 0, in.Value(), out.Ptr(), 0)
}

// ParamName returns the parameter label: "Release", "Gain", etc.
func (e *Effect) ParamName(index int) string {
	var s ParamString
	e.Dispatch(EffGetParamName, Index(index), 0, Ptr(&s), 0)
	return s.String()
}

// ParamValueName returns the parameter value label: "0.5", "HALL", etc.
func (e *Effect) ParamValueName(index int) string {
	var s ParamString
	e.Dispatch(EffGetParamDisplay, Index(index), 0, Ptr(&s), 0)
	return s.String()
}

// ParamUnitName returns the parameter unit label: "db", "ms", etc.
func (e *Effect) ParamUnitName(index int) string {
	var s ParamString
	e.Dispatch(EffGetParamLabel, Index(index), 0, Ptr(&s), 0)
	return s.String()
}

// CurrentProgramName returns current program name.
func (e *Effect) CurrentProgramName() string {
	var s ProgramString
	e.Dispatch(EffGetProgramName, 0, 0, Ptr(&s), 0)
	return s.String()
}

// ProgramName returns program name for provided program index.
func (e *Effect) ProgramName(index int) string {
	var s ProgramString
	e.Dispatch(EffGetProgramNameIndexed, Index(index), 0, Ptr(&s), 0)
	return s.String()
}

// SetCurrentProgramName sets new name to the current program.
func (e *Effect) SetCurrentProgramName(name string) {
	var s ProgramString
	copy(s[:], []byte(name))
	e.Dispatch(EffSetProgramName, 0, 0, Ptr(&s), 0)
}

// Program returns current program number.
func (e *Effect) Program() int {
	return int(e.Dispatch(EffGetProgram, 0, 0, nil, 0))
}

// SetProgram changes current program index.
func (e *Effect) SetProgram(index int) {
	e.Dispatch(EffSetProgram, 0, Value(index), nil, 0)
}

// ParamProperties returns parameter properties for provided parameter
// index. If opcode is not supported, boolean result is false.
func (e *Effect) ParamProperties(index int) (*ParameterProperties, bool) {
	var props ParameterProperties
	r := e.Dispatch(EffGetParameterProperties, Index(index), 0, Ptr(&props), 0)
	if r > 0 {
		return &props, true
	}
	return nil, false
}

// GetProgramData returns current preset data. Plugin allocates required
// memory, then this function allocates new byte slice of required length
// where data is copied.
func (e *Effect) GetProgramData() []byte {
	var ptr unsafe.Pointer
	length := C.int(e.Dispatch(EffGetChunk, 1, 0, Ptr(&ptr), 0))
	return C.GoBytes(ptr, length)
}

// SetProgramData sets preset data to the plugin. Data is the full preset
// including chunk header.
func (e *Effect) SetProgramData(data []byte) {
	e.Dispatch(EffSetChunk, 1, Value(len(data)), Ptr(&data[0]), 0)
}

// GetBankData returns current bank data. Plugin allocates required
// memory, then this function allocates new byte slice of required length
// where data is copied.
func (e *Effect) GetBankData() []byte {
	var ptr unsafe.Pointer
	length := C.int(e.Dispatch(EffGetChunk, 0, 0, Ptr(&ptr), 0))
	return C.GoBytes(ptr, length)
}

// SetBankData sets preset data to the plugin. Data is the full preset
// including chunk header.
func (e *Effect) SetBankData(data []byte) {
	ptr := C.CBytes(data)
	e.Dispatch(EffSetChunk, 0, Value(len(data)), Ptr(ptr), 0)
	C.free(ptr)
}
