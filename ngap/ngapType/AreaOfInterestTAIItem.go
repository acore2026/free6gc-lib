package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type AreaOfInterestTAIItem struct {
	TAI          TAI                                                    `aper:"valueExt"`
	IEExtensions *ProtocolExtensionContainerAreaOfInterestTAIItemExtIEs `aper:"optional"`
}
