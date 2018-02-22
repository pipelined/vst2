package vst2

/*
#include "aeffectx.h"
*/
import "C"

import (
	"fmt"
)

// PluginOpcode used to wrap C opcodes values
type PluginOpcode uint64

//Constants for audio master callback opcodes
//Host -> Plugin
const (
	//AEffectOpcodes opcodes
	EffOpen            = PluginOpcode(C.effOpen)
	EffClose           = PluginOpcode(C.effClose)
	EffSetProgram      = PluginOpcode(C.effSetProgram)
	EffGetProgram      = PluginOpcode(C.effGetProgram)
	EffSetProgramName  = PluginOpcode(C.effSetProgramName)
	EffGetProgramName  = PluginOpcode(C.effGetProgramName)
	EffGetParamLabel   = PluginOpcode(C.effGetParamLabel)
	EffGetParamDisplay = PluginOpcode(C.effGetParamDisplay)
	EffGetParamName    = PluginOpcode(C.effGetParamName)
	EffSetSampleRate   = PluginOpcode(C.effSetSampleRate)
	EffSetBlockSize    = PluginOpcode(C.effSetBlockSize)
	EffMainsChanged    = PluginOpcode(C.effMainsChanged)
	EffEditGetRect     = PluginOpcode(C.effEditGetRect)
	EffEditOpen        = PluginOpcode(C.effEditOpen)
	EffEditClose       = PluginOpcode(C.effEditClose)
	EffEditIdle        = PluginOpcode(C.effEditIdle)
	EffGetChunk        = PluginOpcode(C.effGetChunk)
	EffSetChunk        = PluginOpcode(C.effSetChunk)

	//AEffectXOpcodes opcodes
	EffProcessEvents            = PluginOpcode(C.effProcessEvents)
	EffCanBeAutomated           = PluginOpcode(C.effCanBeAutomated)
	EffString2Parameter         = PluginOpcode(C.effString2Parameter)
	EffGetProgramNameIndexed    = PluginOpcode(C.effGetProgramNameIndexed)
	EffGetInputProperties       = PluginOpcode(C.effGetInputProperties)
	EffGetOutputProperties      = PluginOpcode(C.effGetOutputProperties)
	EffGetPlugCategory          = PluginOpcode(C.effGetPlugCategory)
	EffOfflineNotify            = PluginOpcode(C.effOfflineNotify)
	EffOfflinePrepare           = PluginOpcode(C.effOfflinePrepare)
	EffOfflineRun               = PluginOpcode(C.effOfflineRun)
	EffProcessVarIo             = PluginOpcode(C.effProcessVarIo)
	EffSetSpeakerArrangement    = PluginOpcode(C.effSetSpeakerArrangement)
	EffSetBypass                = PluginOpcode(C.effSetBypass)
	EffGetEffectName            = PluginOpcode(C.effGetEffectName)
	EffGetVendorString          = PluginOpcode(C.effGetVendorString)
	EffGetProductString         = PluginOpcode(C.effGetProductString)
	EffGetVendorVersion         = PluginOpcode(C.effGetVendorVersion)
	EffVendorSpecific           = PluginOpcode(C.effVendorSpecific)
	EffCanDo                    = PluginOpcode(C.effCanDo)
	EffGetTailSize              = PluginOpcode(C.effGetTailSize)
	EffGetParameterProperties   = PluginOpcode(C.effGetParameterProperties)
	EffGetVstVersion            = PluginOpcode(C.effGetVstVersion)
	EffEditKeyDown              = PluginOpcode(C.effEditKeyDown)
	EffEditKeyUp                = PluginOpcode(C.effEditKeyUp)
	EffSetEditKnobMode          = PluginOpcode(C.effSetEditKnobMode)
	EffGetMidiProgramName       = PluginOpcode(C.effGetMidiProgramName)
	EffGetCurrentMidiProgram    = PluginOpcode(C.effGetCurrentMidiProgram)
	EffGetMidiProgramCategory   = PluginOpcode(C.effGetMidiProgramCategory)
	EffHasMidiProgramsChanged   = PluginOpcode(C.effHasMidiProgramsChanged)
	EffGetMidiKeyName           = PluginOpcode(C.effGetMidiKeyName)
	EffBeginSetProgram          = PluginOpcode(C.effBeginSetProgram)
	EffEndSetProgram            = PluginOpcode(C.effEndSetProgram)
	EffGetSpeakerArrangement    = PluginOpcode(C.effGetSpeakerArrangement)
	EffShellGetNextPlugin       = PluginOpcode(C.effShellGetNextPlugin)
	EffStartProcess             = PluginOpcode(C.effStartProcess)
	EffStopProcess              = PluginOpcode(C.effStopProcess)
	EffSetTotalSampleToProcess  = PluginOpcode(C.effSetTotalSampleToProcess)
	EffSetPanLaw                = PluginOpcode(C.effSetPanLaw)
	EffBeginLoadBank            = PluginOpcode(C.effBeginLoadBank)
	EffBeginLoadProgram         = PluginOpcode(C.effBeginLoadProgram)
	EffSetProcessPrecision      = PluginOpcode(C.effSetProcessPrecision)
	EffGetNumMidiInputChannels  = PluginOpcode(C.effGetNumMidiInputChannels)
	EffGetNumMidiOutputChannels = PluginOpcode(C.effGetNumMidiOutputChannels)
)

// MasterOpcode used to wrap C opcodes values
type MasterOpcode uint64

// Constants for audio master callback opcodes
// Plugin -> Host
const (
	// AudioMasterOpcodes opcodes
	AudioMasterAutomate  = MasterOpcode(C.audioMasterAutomate)
	AudioMasterVersion   = MasterOpcode(C.audioMasterVersion)
	AudioMasterCurrentID = MasterOpcode(C.audioMasterCurrentId)
	AudioMasterIdle      = MasterOpcode(C.audioMasterIdle)

	// AudioMasterOpcodesX opcodes
	AudioMasterGetTime                   = MasterOpcode(C.audioMasterGetTime)
	AudioMasterProcessEvents             = MasterOpcode(C.audioMasterProcessEvents)
	AudioMasterIOChanged                 = MasterOpcode(C.audioMasterIOChanged)
	AudioMasterSizeWindow                = MasterOpcode(C.audioMasterSizeWindow)
	AudioMasterGetSampleRate             = MasterOpcode(C.audioMasterGetSampleRate)
	AudioMasterGetBlockSize              = MasterOpcode(C.audioMasterGetBlockSize)
	AudioMasterGetInputLatency           = MasterOpcode(C.audioMasterGetInputLatency)
	AudioMasterGetOutputLatency          = MasterOpcode(C.audioMasterGetOutputLatency)
	AudioMasterGetCurrentProcessLevel    = MasterOpcode(C.audioMasterGetCurrentProcessLevel)
	AudioMasterGetAutomationState        = MasterOpcode(C.audioMasterGetAutomationState)
	AudioMasterOfflineStart              = MasterOpcode(C.audioMasterOfflineStart)
	AudioMasterOfflineRead               = MasterOpcode(C.audioMasterOfflineRead)
	AudioMasterOfflineWrite              = MasterOpcode(C.audioMasterOfflineWrite)
	AudioMasterOfflineGetCurrentPass     = MasterOpcode(C.audioMasterOfflineGetCurrentPass)
	AudioMasterOfflineGetCurrentMetaPass = MasterOpcode(C.audioMasterOfflineGetCurrentMetaPass)
	AudioMasterGetVendorString           = MasterOpcode(C.audioMasterGetVendorString)
	AudioMasterGetProductString          = MasterOpcode(C.audioMasterGetProductString)
	AudioMasterGetVendorVersion          = MasterOpcode(C.audioMasterGetVendorVersion)
	AudioMasterVendorSpecific            = MasterOpcode(C.audioMasterVendorSpecific)
	AudioMasterCanDo                     = MasterOpcode(C.audioMasterCanDo)
	AudioMasterGetLanguage               = MasterOpcode(C.audioMasterGetLanguage)
	AudioMasterGetDirectory              = MasterOpcode(C.audioMasterGetDirectory)
	AudioMasterUpdateDisplay             = MasterOpcode(C.audioMasterUpdateDisplay)
	AudioMasterBeginEdit                 = MasterOpcode(C.audioMasterBeginEdit)
	AudioMasterEndEdit                   = MasterOpcode(C.audioMasterEndEdit)
	AudioMasterOpenFileSelector          = MasterOpcode(C.audioMasterOpenFileSelector)
	AudioMasterCloseFileSelector         = MasterOpcode(C.audioMasterCloseFileSelector)
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
