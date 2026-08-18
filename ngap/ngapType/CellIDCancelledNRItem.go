package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type CellIDCancelledNRItem struct {
	NRCGI              NRCGI `aper:"valueExt"`
	NumberOfBroadcasts NumberOfBroadcasts
	IEExtensions       *ProtocolExtensionContainerCellIDCancelledNRItemExtIEs `aper:"optional"`
}
