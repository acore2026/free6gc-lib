package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type TAICancelledEUTRAItem struct {
	TAI                      TAI `aper:"valueExt"`
	CancelledCellsInTAIEUTRA CancelledCellsInTAIEUTRA
	IEExtensions             *ProtocolExtensionContainerTAICancelledEUTRAItemExtIEs `aper:"optional"`
}
