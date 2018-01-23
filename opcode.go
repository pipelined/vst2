package vst2

/*
#include "aeffectx.h"
*/
import "C"

import (
	"fmt"
)

//pluginOpcode used to wrap C opcodes values
type pluginOpcode uint64

//Constants for audio master callback opcodes
//Host -> Plugin
const (
	//AEffectOpcodes opcodes
	EffOpen            = pluginOpcode(C.effOpen)
	EffClose           = pluginOpcode(C.effClose)
	EffSetProgram      = pluginOpcode(C.effSetProgram)
	EffGetProgram      = pluginOpcode(C.effGetProgram)
	EffSetProgramName  = pluginOpcode(C.effSetProgramName)
	EffGetProgramName  = pluginOpcode(C.effGetProgramName)
	EffGetParamLabel   = pluginOpcode(C.effGetParamLabel)
	EffGetParamDisplay = pluginOpcode(C.effGetParamDisplay)
	EffGetParamName    = pluginOpcode(C.effGetParamName)
	EffSetSampleRate   = pluginOpcode(C.effSetSampleRate)
	EffSetBlockSize    = pluginOpcode(C.effSetBlockSize)
	EffMainsChanged    = pluginOpcode(C.effMainsChanged)
	EffEditGetRect     = pluginOpcode(C.effEditGetRect)
	EffEditOpen        = pluginOpcode(C.effEditOpen)
	EffEditClose       = pluginOpcode(C.effEditClose)
	EffEditIdle        = pluginOpcode(C.effEditIdle)
	EffGetChunk        = pluginOpcode(C.effGetChunk)
	EffSetChunk        = pluginOpcode(C.effSetChunk)

	//AEffectXOpcodes opcodes
	EffProcessEvents            = pluginOpcode(C.effProcessEvents)
	EffCanBeAutomated           = pluginOpcode(C.effCanBeAutomated)
	EffString2Parameter         = pluginOpcode(C.effString2Parameter)
	EffGetProgramNameIndexed    = pluginOpcode(C.effGetProgramNameIndexed)
	EffGetInputProperties       = pluginOpcode(C.effGetInputProperties)
	EffGetOutputProperties      = pluginOpcode(C.effGetOutputProperties)
	EffGetPlugCategory          = pluginOpcode(C.effGetPlugCategory)
	EffOfflineNotify            = pluginOpcode(C.effOfflineNotify)
	EffOfflinePrepare           = pluginOpcode(C.effOfflinePrepare)
	EffOfflineRun               = pluginOpcode(C.effOfflineRun)
	EffProcessVarIo             = pluginOpcode(C.effProcessVarIo)
	EffSetSpeakerArrangement    = pluginOpcode(C.effSetSpeakerArrangement)
	EffSetBypass                = pluginOpcode(C.effSetBypass)
	EffGetEffectName            = pluginOpcode(C.effGetEffectName)
	EffGetVendorString          = pluginOpcode(C.effGetVendorString)
	EffGetProductString         = pluginOpcode(C.effGetProductString)
	EffGetVendorVersion         = pluginOpcode(C.effGetVendorVersion)
	EffVendorSpecific           = pluginOpcode(C.effVendorSpecific)
	EffCanDo                    = pluginOpcode(C.effCanDo)
	EffGetTailSize              = pluginOpcode(C.effGetTailSize)
	EffGetParameterProperties   = pluginOpcode(C.effGetParameterProperties)
	EffGetVstVersion            = pluginOpcode(C.effGetVstVersion)
	EffEditKeyDown              = pluginOpcode(C.effEditKeyDown)
	EffEditKeyUp                = pluginOpcode(C.effEditKeyUp)
	EffSetEditKnobMode          = pluginOpcode(C.effSetEditKnobMode)
	EffGetMidiProgramName       = pluginOpcode(C.effGetMidiProgramName)
	EffGetCurrentMidiProgram    = pluginOpcode(C.effGetCurrentMidiProgram)
	EffGetMidiProgramCategory   = pluginOpcode(C.effGetMidiProgramCategory)
	EffHasMidiProgramsChanged   = pluginOpcode(C.effHasMidiProgramsChanged)
	EffGetMidiKeyName           = pluginOpcode(C.effGetMidiKeyName)
	EffBeginSetProgram          = pluginOpcode(C.effBeginSetProgram)
	EffEndSetProgram            = pluginOpcode(C.effEndSetProgram)
	EffGetSpeakerArrangement    = pluginOpcode(C.effGetSpeakerArrangement)
	EffShellGetNextPlugin       = pluginOpcode(C.effShellGetNextPlugin)
	EffStartProcess             = pluginOpcode(C.effStartProcess)
	EffStopProcess              = pluginOpcode(C.effStopProcess)
	EffSetTotalSampleToProcess  = pluginOpcode(C.effSetTotalSampleToProcess)
	EffSetPanLaw                = pluginOpcode(C.effSetPanLaw)
	EffBeginLoadBank            = pluginOpcode(C.effBeginLoadBank)
	EffBeginLoadProgram         = pluginOpcode(C.effBeginLoadProgram)
	EffSetProcessPrecision      = pluginOpcode(C.effSetProcessPrecision)
	EffGetNumMidiInputChannels  = pluginOpcode(C.effGetNumMidiInputChannels)
	EffGetNumMidiOutputChannels = pluginOpcode(C.effGetNumMidiOutputChannels)
)

//masterOpcode used to wrap C opcodes values
type masterOpcode uint64

//Constants for audio master callback opcodes
//Plugin -> Host
const (
	//AudioMasterOpcodes opcodes
	AudioMasterAutomate  = masterOpcode(C.audioMasterAutomate)
	AudioMasterVersion   = masterOpcode(C.audioMasterVersion)
	AudioMasterCurrentID = masterOpcode(C.audioMasterCurrentId)
	AudioMasterIdle      = masterOpcode(C.audioMasterIdle)

	//AudioMasterOpcodesX opcodes
	AudioMasterGetTime                   = masterOpcode(C.audioMasterGetTime)
	AudioMasterProcessEvents             = masterOpcode(C.audioMasterProcessEvents)
	AudioMasterIOChanged                 = masterOpcode(C.audioMasterIOChanged)
	AudioMasterSizeWindow                = masterOpcode(C.audioMasterSizeWindow)
	AudioMasterGetSampleRate             = masterOpcode(C.audioMasterGetSampleRate)
	AudioMasterGetBlockSize              = masterOpcode(C.audioMasterGetBlockSize)
	AudioMasterGetInputLatency           = masterOpcode(C.audioMasterGetInputLatency)
	AudioMasterGetOutputLatency          = masterOpcode(C.audioMasterGetOutputLatency)
	AudioMasterGetCurrentProcessLevel    = masterOpcode(C.audioMasterGetCurrentProcessLevel)
	AudioMasterGetAutomationState        = masterOpcode(C.audioMasterGetAutomationState)
	AudioMasterOfflineStart              = masterOpcode(C.audioMasterOfflineStart)
	AudioMasterOfflineRead               = masterOpcode(C.audioMasterOfflineRead)
	AudioMasterOfflineWrite              = masterOpcode(C.audioMasterOfflineWrite)
	AudioMasterOfflineGetCurrentPass     = masterOpcode(C.audioMasterOfflineGetCurrentPass)
	AudioMasterOfflineGetCurrentMetaPass = masterOpcode(C.audioMasterOfflineGetCurrentMetaPass)
	AudioMasterGetVendorString           = masterOpcode(C.audioMasterGetVendorString)
	AudioMasterGetProductString          = masterOpcode(C.audioMasterGetProductString)
	AudioMasterGetVendorVersion          = masterOpcode(C.audioMasterGetVendorVersion)
	AudioMasterVendorSpecific            = masterOpcode(C.audioMasterVendorSpecific)
	AudioMasterCanDo                     = masterOpcode(C.audioMasterCanDo)
	AudioMasterGetLanguage               = masterOpcode(C.audioMasterGetLanguage)
	AudioMasterGetDirectory              = masterOpcode(C.audioMasterGetDirectory)
	AudioMasterUpdateDisplay             = masterOpcode(C.audioMasterUpdateDisplay)
	AudioMasterBeginEdit                 = masterOpcode(C.audioMasterBeginEdit)
	AudioMasterEndEdit                   = masterOpcode(C.audioMasterEndEdit)
	AudioMasterOpenFileSelector          = masterOpcode(C.audioMasterOpenFileSelector)
	AudioMasterCloseFileSelector         = masterOpcode(C.audioMasterCloseFileSelector)
)

func (p masterOpcode) String() string {
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

func (p pluginOpcode) String() string {
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
