package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type CellIDBroadcastEUTRAItem struct {
	EUTRACGI     EUTRACGI                                                  `aper:"valueExt"`
	IEExtensions *ProtocolExtensionContainerCellIDBroadcastEUTRAItemExtIEs `aper:"optional"`
}
