package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type AreaOfInterestItem struct {
	AreaOfInterest               AreaOfInterest `aper:"valueExt"`
	LocationReportingReferenceID LocationReportingReferenceID
	IEExtensions                 *ProtocolExtensionContainerAreaOfInterestItemExtIEs `aper:"optional"`
}
