package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type CancelledCellsInTAIEUTRAItem struct {
	EUTRACGI           EUTRACGI `aper:"valueExt"`
	NumberOfBroadcasts NumberOfBroadcasts
	IEExtensions       *ProtocolExtensionContainerCancelledCellsInTAIEUTRAItemExtIEs `aper:"optional"`
}
