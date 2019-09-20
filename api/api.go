package api

type TimeInfo struct {
	SamplePos          float64
	SampleRate         float64
	NanoSeconds        float64
	PpqPos             float64
	Tempo              float64
	BarStartPos        float64
	CycleStartPos      float64
	CycleEndPos        float64
	TimeSigNumerator   int32
	TimeSigDenominator int32
	SmpteOffset        int32
	SmpteFrameRate     int32
	SamplesToNextClock int32
	Flags              TimeInfoFlags
}

type SpeakerProperties struct {
	Azimuth   float32
	Elevation float32
	Radius    float32
	Reserved  float32
	Name      [64]byte
	SpeakerType
	Future [28]byte
}

type SpeakerArrangement struct {
	SpeakerArrangementType
	NumChannels int32
	Speakers    [8]SpeakerProperties
}

type EffectFlags int32

const (
	EffFlagsHasEditor EffectFlags = 1 << iota
	EffFlagsCanReplacing
	EffFlagsProgramChunks
	EffFlagsIsSynth
	EffFlagsNoSoundInStop
	EffFlagsCanDoubleReplacing
	effFlagsHasClipDeprecated
	effFlagsHasVuDeprecated
	effFlagsCanMonoDeprecated
	effFlagsExtIsAsyncDeprecated
	effFlagsExtHasBufferDeprecated
)

type SpeakerType int32

const (
	// SpeakerUndefined is undefined
	SpeakerUndefined SpeakerType = 2147483647
	// SpeakerM is Mono (M)
	SpeakerM = iota
	// SpeakerL is Left (L)
	SpeakerL
	// SpeakerR is Right (R)
	SpeakerR
	// Center (C)
	SpeakerC
	// Subbass (Lfe)
	SpeakerLfe
	// Left Surround (Ls)
	SpeakerLs
	// Right Surround (Rs)
	SpeakerRs
	// Left of Center (Lc)
	SpeakerLc
	// Right of Center (Rc)
	SpeakerRc
	// Surround (S)
	SpeakerS
	// Center of Surround (Cs) = Surround (S)
	SpeakerCs = SpeakerS
	// Side Left (Sl)
	SpeakerSl
	// Side Right (Sr)
	SpeakerSr
	// Top Middle (Tm)
	SpeakerTm
	// Top Front Left (Tfl)
	SpeakerTfl
	// Top Front Center (Tfc)
	SpeakerTfc
	// Top Front Right (Tfr)
	SpeakerTfr
	// Top Rear Left (Trl)
	SpeakerTrl
	// Top Rear Center (Trc)
	SpeakerTrc
	// Top Rear Right (Trr)
	SpeakerTrr
	// Subbass 2 (Lfe2)
	SpeakerLfe2
)

type SpeakerArrangementType int32

const (
	SpeakerArrUserDefined SpeakerArrangementType = iota - 2
	SpeakerArrEmpty
	SpeakerArrMono
	SpeakerArrStereo
	SpeakerArrStereoSurround
	SpeakerArrStereoCenter
	SpeakerArrStereoSide
	SpeakerArrStereoCLfe
	SpeakerArr30Cine
	SpeakerArr30Music
	SpeakerArr31Cine
	SpeakerArr31Music
	SpeakerArr40Cine
	SpeakerArr40Music
	SpeakerArr41Cine
	SpeakerArr41Music
	SpeakerArr50
	SpeakerArr51
	SpeakerArr60Cine
	SpeakerArr60Music
	SpeakerArr61Cine
	SpeakerArr61Music
	SpeakerArr70Cine
	SpeakerArr70Music
	SpeakerArr71Cine
	SpeakerArr71Music
	SpeakerArr80Cine
	SpeakerArr80Music
	SpeakerArr81Cine
	SpeakerArr81Music
	SpeakerArr102
	NumSpeakerArr
)

type TimeInfoFlags int32

const (
	TransportChanged TimeInfoFlags = iota + 1
	TransportPlaying
	TransportCycleActive
	TransportRecording
	AutomationWriting
	AutomationReading
	NanosValid
	PpqPosValid
	TempoValid
	BarsValid
	CyclePosValid
	TimeSigValid
	SmpteValid
	ClockValid
)

type ProcessLevels int32

const (
	ProcessLevelUnknown ProcessLevels = iota
	ProcessLevelUser
	ProcessLevelRealtime
	ProcessLevelPrefetch
	ProcessLevelOffline
)
