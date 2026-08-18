package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type TAI struct {
	PLMNIdentity PLMNIdentity
	TAC          TAC
	IEExtensions *ProtocolExtensionContainerTAIExtIEs `aper:"optional"`
}
