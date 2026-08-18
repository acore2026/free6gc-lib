package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type CompletedCellsInEAINRItem struct {
	NRCGI        NRCGI                                                      `aper:"valueExt"`
	IEExtensions *ProtocolExtensionContainerCompletedCellsInEAINRItemExtIEs `aper:"optional"`
}
