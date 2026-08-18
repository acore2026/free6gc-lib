package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type TAICancelledNRItem struct {
	TAI                   TAI `aper:"valueExt"`
	CancelledCellsInTAINR CancelledCellsInTAINR
	IEExtensions          *ProtocolExtensionContainerTAICancelledNRItemExtIEs `aper:"optional"`
}
