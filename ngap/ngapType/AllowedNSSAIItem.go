package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type AllowedNSSAIItem struct {
	SNSSAI       SNSSAI                                            `aper:"valueExt"`
	IEExtensions *ProtocolExtensionContainerAllowedNSSAIItemExtIEs `aper:"optional"`
}
