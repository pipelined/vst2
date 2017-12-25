package vst2

/*
#cgo CFLAGS: -std=gnu99 -I${SRCDIR}/vst2sdk/

#include "vst2.h"
#include "aeffectx.h"
*/
import "C"

import (
	"log"
	"unsafe"
)

//Plugin type provides interface
type Plugin struct {
	entryPoint uintptr
	effect     *C.AEffect
}

//HostCallbackFunc used as callback from plugin
type HostCallbackFunc func(*Plugin, MasterOpcode, int64, int64, unsafe.Pointer, float64) int

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
	EffNumOpcodes      = pluginOpcode(C.effNumOpcodes)

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

const (
	vstMain string = "VSTPluginMain"
)

//MasterOpcode used to wrap C opcodes values
type MasterOpcode uint64

//Constants for audio master callback opcodes
//Plugin -> Host
const (
	//AudioMasterOpcodes opcodes
	AudioMasterAutomate  = MasterOpcode(C.audioMasterAutomate)
	AudioMasterVersion   = MasterOpcode(C.audioMasterVersion)
	AudioMasterCurrentID = MasterOpcode(C.audioMasterCurrentId)
	AudioMasterIdle      = MasterOpcode(C.audioMasterIdle)

	//AudioMasterOpcodesX opcodes
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

var (
	callback HostCallbackFunc = HostCallback
)

//NewPlugin loads the plugin into memory and stores callback func
//TODO: catch panic
func NewPlugin(path string) (*Plugin, error) {
	//Get pointer to plugin's Main function
	mainEntryPoint, err := getEntryPoint(path)
	if err != nil {
		log.Printf("Failed to obtain VST entry point '%s': %v\n", path, err)
		return nil, err
	}

	return &Plugin{entryPoint: mainEntryPoint}, nil
}

//Dispatch wraps-up C method to dispatch calls to plugin
func (plugin *Plugin) Dispatch(opcode pluginOpcode, index int64, value int64, ptr unsafe.Pointer, opt float64) {
	if plugin.effect != nil {
		C.dispatch(plugin.effect, C.int(opcode), C.int(index), C.int(value), ptr, C.float(opt))
	}
}

//Resume the plugin
func (plugin *Plugin) resume() {
	plugin.Dispatch(EffMainsChanged, 0, 1, nil, 0.0)
}

//Suspend the plugin
func (plugin *Plugin) suspend() {
	plugin.Dispatch(EffMainsChanged, 0, 0, nil, 0.0)
}

//Start the plugin
func (plugin *Plugin) Start() {
	//Convert to C++ pointer type
	vstEntryPoint := (C.vstPluginFuncPtr)(unsafe.Pointer(plugin.entryPoint))
	plugin.effect = C.loadEffect(vstEntryPoint)

	plugin.Dispatch(EffOpen, 0, 0, nil, 0.0)

	// Set default sample rate and block size
	sampleRate := 44100.0
	plugin.Dispatch(EffSetSampleRate, 0, 0, nil, sampleRate)

	blocksize := int64(4096)
	plugin.Dispatch(EffSetBlockSize, 0, blocksize, nil, 0.0)
}

//Process audio with VST plugin
func (plugin *Plugin) Process(samples [][]float64) (processed [][]float64) {
	//convert Samples to C type
	inSamples := (**C.double)(unsafe.Pointer(&samples[0][0]))
	blocksize := C.int(len(samples[0]))
	numChannels := C.int(len(samples))
	//call plugin and convert result to slice of slices
	outSamples := (*[1 << 30]*C.double)(unsafe.Pointer(C.processAudio(plugin.effect, numChannels, blocksize, inSamples)))[:numChannels]
	//convert slices to [][]float64
	processed = make([][]float64, numChannels)
	for channel, data := range outSamples {
		processed[channel] = (*[1 << 30]float64)(unsafe.Pointer(data))[:blocksize]
	}
	return processed
}

//SetHostCallback allows to override default host callback with custom implementation
func SetHostCallback(newCallback HostCallbackFunc) {
	if newCallback != nil {
		callback = newCallback
	}
}

//export hostCallback
//calls real callback
func hostCallback(effect *C.AEffect, opcode int64, index int64, value int64, ptr unsafe.Pointer, opt float64) int {
	if callback == nil {
		panic("Host callback is not defined!")
	}

	return callback(&Plugin{effect: effect}, MasterOpcode(opcode), index, value, ptr, opt)
}

//HostCallback is a default callback, can be overriden with SetHostCallback
func HostCallback(plugin *Plugin, opcode MasterOpcode, index int64, value int64, ptr unsafe.Pointer, opt float64) int {
	switch opcode {
	case AudioMasterVersion:
		return 2400
	case AudioMasterIdle:
		plugin.Dispatch(EffEditIdle, 0, 0, nil, 0)
	default:
		log.Printf("Plugin requested value of opcode %v\n", opcode)
		break
	}
	return 0
}
