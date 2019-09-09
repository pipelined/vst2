package vst2

import (
	"fmt"
)

// PluginOpcode used for audio master callback opcodes.
// They are sent when host communicates with plugin.
// It relfects AEffectOpcodes and AEffectXOpcodes opcodes values.
type PluginOpcode uint64

const (
	EffOpen PluginOpcode = iota
	EffClose
	EffSetProgram
	EffGetProgram
	EffSetProgramName
	EffGetProgramName
	EffGetParamLabel
	EffGetParamDisplay
	EffGetParamName
	// EffGetVu is deprecated in VST v2.4
	EffGetVu
	EffSetSampleRate
	EffSetBlockSize
	EffMainsChanged
	EffEditGetRect
	EffEditOpen
	EffEditClose
	// EffEditDraw is deprecated in VST v2.4
	EffEditDraw
	// EffEditMouse is deprecated in VST v2.4
	EffEditMouse
	// EffEditKey is deprecated in VST v2.4
	EffEditKey
	EffEditIdle
	// EffEditTop is deprecated in VST v2.4
	EffEditTop
	// EffEditSleep is deprecated in VST v2.4
	EffEditSleep
	// EffIdentify is deprecated in VST v2.4
	EffIdentify
	EffGetChunk
	EffSetChunk

	//AEffectXOpcodes opcode
	EffProcessEvents
	EffCanBeAutomated
	EffString2Parameter
	// EffGetNumProgramCategories is deprecated in VST v2.4
	EffGetNumProgramCategories
	EffGetProgramNameIndexed
	// EffCopyProgram is deprecated in VST v2.4
	EffCopyProgram
	// EffConnectInput is deprecated in VST v2.4
	EffConnectInput
	// EffConnectOutput is deprecated in VST v2.4
	EffConnectOutput
	EffGetInputProperties
	EffGetOutputProperties
	EffGetPlugCategory
	// EffGetCurrentPosition is deprecated in VST v2.4
	EffGetCurrentPosition
	// EffGetDestinationBuffer is deprecated in VST v2.4
	EffGetDestinationBuffer
	EffOfflineNotify
	EffOfflinePrepare
	EffOfflineRun
	EffProcessVarIo
	EffSetSpeakerArrangement
	// EffSetBlockSizeAndSampleRate is deprecated in VST v2.4
	EffSetBlockSizeAndSampleRate
	EffSetBypass
	EffGetEffectName
	// EffGetErrorText is deprecated in VST v2.4
	EffGetErrorText
	EffGetVendorString
	EffGetProductString
	EffGetVendorVersion
	EffVendorSpecific
	EffCanDo
	EffGetTailSize
	// EffIdle is deprecated in VST v2.4
	EffIdle
	// EffGetIcon is deprecated in VST v2.4
	EffGetIcon
	// EffSetViewPosition is deprecated in VST v2.4
	EffSetViewPosition
	EffGetParameterProperties
	// EffKeysRequired is deprecated in VST v2.4
	EffKeysRequired
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

// MasterOpcode used to wrap C opcodes values
type MasterOpcode uint64

// Constants for audio master callback opcodes
// Plugin -> Host
const (
	// AudioMasterOpcodes opcodes
	AudioMasterAutomate MasterOpcode = iota
	AudioMasterVersion
	AudioMasterCurrentID
	AudioMasterIdle
	// AudioMasterPinConnected is deprecated in VST v2.4
	AudioMasterPinConnected

	_
	// AudioMasterOpcodesX opcodes
	AudioMasterWantMidi
	AudioMasterGetTime
	AudioMasterProcessEvents
	// AudioMasterSetTime is deprecated in VST v2.4
	AudioMasterSetTime
	// AudioMasterTempoAt is deprecated in VST v2.4
	AudioMasterTempoAt
	// AudioMasterGetNumAutomatableParameters is deprecated in VST v2.4
	AudioMasterGetNumAutomatableParameters
	// AudioMasterGetParameterQuantization is deprecated in VST v2.4
	AudioMasterGetParameterQuantization
	AudioMasterIOChanged
	// AudioMasterNeedIdle is deprecated in VST v2.4
	AudioMasterNeedIdle
	AudioMasterSizeWindow
	AudioMasterGetSampleRate
	AudioMasterGetBlockSize
	AudioMasterGetInputLatency
	AudioMasterGetOutputLatency
	// AudioMasterGetPreviousPlug is deprecated in VST v2.4
	AudioMasterGetPreviousPlug
	// AudioMasterGetNextPlug is deprecated in VST v2.4
	AudioMasterGetNextPlug
	// AudioMasterWillReplaceOrAccumulate is deprecated in VST v2.4
	AudioMasterWillReplaceOrAccumulate
	AudioMasterGetCurrentProcessLevel
	AudioMasterGetAutomationState
	AudioMasterOfflineStart
	AudioMasterOfflineRead
	AudioMasterOfflineWrite
	AudioMasterOfflineGetCurrentPass
	AudioMasterOfflineGetCurrentMetaPass
	// AudioMasterSetOutputSampleRate is deprecated in VST v2.4
	AudioMasterSetOutputSampleRate
	// AudioMasterGetOutputSpeakerArrangement is deprecated in VST v2.4
	AudioMasterGetOutputSpeakerArrangement
	AudioMasterGetVendorString
	AudioMasterGetProductString
	AudioMasterGetVendorVersion
	AudioMasterVendorSpecific
	// AudioMasterSetIcon is deprecated in VST v2.4
	AudioMasterSetIcon
	AudioMasterCanDo
	AudioMasterGetLanguage
	// AudioMasterOpenWindow is deprecated in VST v2.4
	AudioMasterOpenWindow
	// AudioMasterCloseWindow is deprecated in VST v2.4
	AudioMasterCloseWindow
	AudioMasterGetDirectory
	AudioMasterUpdateDisplay
	AudioMasterBeginEdit
	AudioMasterEndEdit
	AudioMasterOpenFileSelector
	AudioMasterCloseFileSelector
	// AudioMasterEditFile is deprecated in VST v2.4
	AudioMasterEditFile
	// AudioMasterGetChunkFile is deprecated in VST v2.4
	AudioMasterGetChunkFile
	// AudioMasterGetInputSpeakerArrangement is deprecated in VST v2.4
	AudioMasterGetInputSpeakerArrangement
)

func (p MasterOpcode) String() string {
	switch p {
	case AudioMasterAutomate:
		return "AudioMasterAutomate"
	case AudioMasterVersion:
		return "AudioMasterVersion"

	case AudioMasterCurrentID:
		return "AudioMasterCurrentID"
	case AudioMasterIdle:
		return "AudioMasterIdle"
	case AudioMasterWantMidi:
		return "AudioMasterWantMidi"
	case AudioMasterGetTime:
		return "AudioMasterGetTime"
	case AudioMasterProcessEvents:
		return "AudioMasterProcessEvents"
	case AudioMasterIOChanged:
		return "AudioMasterIOChanged"
	case AudioMasterSizeWindow:
		return "AudioMasterSizeWindow"
	case AudioMasterGetSampleRate:
		return "AudioMasterGetSampleRate"
	case AudioMasterGetBlockSize:
		return "AudioMasterGetBlockSize"
	case AudioMasterGetInputLatency:
		return "AudioMasterGetInputLatency"
	case AudioMasterGetOutputLatency:
		return "AudioMasterGetOutputLatency"
	case AudioMasterGetCurrentProcessLevel:
		return "AudioMasterGetCurrentProcessLevel"
	case AudioMasterGetAutomationState:
		return "AudioMasterGetAutomationState"
	case AudioMasterOfflineStart:
		return "AudioMasterOfflineStart"
	case AudioMasterOfflineRead:
		return "AudioMasterOfflineRead"
	case AudioMasterOfflineWrite:
		return "AudioMasterOfflineWrite"
	case AudioMasterOfflineGetCurrentPass:
		return "AudioMasterOfflineGetCurrentPass"
	case AudioMasterOfflineGetCurrentMetaPass:
		return "AudioMasterOfflineGetCurrentMetaPass"
	case AudioMasterGetVendorString:
		return "AudioMasterGetVendorString"
	case AudioMasterGetProductString:
		return "AudioMasterGetProductString"
	case AudioMasterGetVendorVersion:
		return "AudioMasterGetVendorVersion"
	case AudioMasterVendorSpecific:
		return "AudioMasterVendorSpecific"
	case AudioMasterCanDo:
		return "AudioMasterCanDo"
	case AudioMasterGetLanguage:
		return "AudioMasterGetLanguage"
	case AudioMasterGetDirectory:
		return "AudioMasterGetDirectory"
	case AudioMasterUpdateDisplay:
		return "AudioMasterUpdateDisplay"
	case AudioMasterBeginEdit:
		return "AudioMasterBeginEdit"
	case AudioMasterEndEdit:
		return "AudioMasterEndEdit"
	case AudioMasterOpenFileSelector:
		return "AudioMasterOpenFileSelector"
	case AudioMasterCloseFileSelector:
		return "AudioMasterCloseFileSelector"
	default:
		return fmt.Sprintf("Unknown master opode %d\n", p)
	}
}

func (p PluginOpcode) String() string {
	switch p {
	case EffOpen:
		return "EffOpen"
	case EffClose:
		return "EffClose"
	case EffSetProgram:
		return "EffSetProgram"
	case EffGetProgram:
		return "EffGetProgram"
	case EffSetProgramName:
		return "EffSetProgramName"
	case EffGetProgramName:
		return "EffGetProgramName"
	case EffGetParamLabel:
		return "EffGetParamLabel"
	case EffGetParamDisplay:
		return "EffGetParamDisplay"
	case EffGetParamName:
		return "EffGetParamName"
	case EffSetSampleRate:
		return "EffSetSampleRate"
	case EffSetBlockSize:
		return "EffSetBlockSize"
	case EffMainsChanged:
		return "EffMainsChanged"
	case EffEditGetRect:
		return "EffEditGetRect"
	case EffEditOpen:
		return "EffEditOpen"
	case EffEditClose:
		return "EffEditClose"
	case EffEditIdle:
		return "EffEditIdle"
	case EffGetChunk:
		return "EffGetChunk"
	case EffSetChunk:
		return "EffSetChunk"
	case EffProcessEvents:
		return "EffProcessEvents"
	case EffCanBeAutomated:
		return "EffCanBeAutomated"
	case EffString2Parameter:
		return "EffString2Parameter"
	case EffGetProgramNameIndexed:
		return "EffGetProgramNameIndexed"
	case EffGetInputProperties:
		return "EffGetInputProperties"
	case EffGetOutputProperties:
		return "EffGetOutputProperties"
	case EffGetPlugCategory:
		return "EffGetPlugCategory"
	case EffOfflineNotify:
		return "EffOfflineNotify"
	case EffOfflinePrepare:
		return "EffOfflinePrepare"
	case EffOfflineRun:
		return "EffOfflineRun"
	case EffProcessVarIo:
		return "EffProcessVarIo"
	case EffSetSpeakerArrangement:
		return "EffSetSpeakerArrangement"
	case EffSetBypass:
		return "EffSetBypass"
	case EffGetEffectName:
		return "EffGetEffectName"
	case EffGetVendorString:
		return "EffGetVendorString"
	case EffGetProductString:
		return "EffGetProductString"
	case EffGetVendorVersion:
		return "EffGetVendorVersion"
	case EffVendorSpecific:
		return "EffVendorSpecific"
	case EffCanDo:
		return "EffCanDo"
	case EffGetTailSize:
		return "EffGetTailSize"
	case EffGetParameterProperties:
		return "EffGetParameterProperties"
	case EffGetVstVersion:
		return "EffGetVstVersion"
	case EffEditKeyDown:
		return "EffEditKeyDown"
	case EffEditKeyUp:
		return "EffEditKeyUp"
	case EffSetEditKnobMode:
		return "EffSetEditKnobMode"
	case EffGetMidiProgramName:
		return "EffGetMidiProgramName"
	case EffGetCurrentMidiProgram:
		return "EffGetCurrentMidiProgram"
	case EffGetMidiProgramCategory:
		return "EffGetMidiProgramCategory"
	case EffHasMidiProgramsChanged:
		return "EffHasMidiProgramsChanged"
	case EffGetMidiKeyName:
		return "EffGetMidiKeyName"
	case EffBeginSetProgram:
		return "EffBeginSetProgram"
	case EffEndSetProgram:
		return "EffEndSetProgram"
	case EffGetSpeakerArrangement:
		return "EffGetSpeakerArrangement"
	case EffShellGetNextPlugin:
		return "EffShellGetNextPlugin"
	case EffStartProcess:
		return "EffStartProcess"
	case EffStopProcess:
		return "EffStopProcess"
	case EffSetTotalSampleToProcess:
		return "EffSetTotalSampleToProcess"
	case EffSetPanLaw:
		return "EffSetPanLaw"
	case EffBeginLoadBank:
		return "EffBeginLoadBank"
	case EffBeginLoadProgram:
		return "EffBeginLoadProgram"
	case EffSetProcessPrecision:
		return "EffSetProcessPrecision"
	case EffGetNumMidiInputChannels:
		return "EffGetNumMidiInputChannels"
	case EffGetNumMidiOutputChannels:
		return "EffGetNumMidiOutputChannels"
	default:
		return fmt.Sprintf("Unknown plugin opcode %d", p)
	}
}
