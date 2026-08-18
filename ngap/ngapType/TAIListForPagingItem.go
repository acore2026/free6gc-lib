package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type TAIListForPagingItem struct {
	TAI          TAI                                                   `aper:"valueExt"`
	IEExtensions *ProtocolExtensionContainerTAIListForPagingItemExtIEs `aper:"optional"`
}
