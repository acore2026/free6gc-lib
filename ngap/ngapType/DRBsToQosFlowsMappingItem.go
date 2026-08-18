package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type DRBsToQosFlowsMappingItem struct {
	DRBID                 DRBID
	AssociatedQosFlowList AssociatedQosFlowList
	IEExtensions          *ProtocolExtensionContainerDRBsToQosFlowsMappingItemExtIEs `aper:"optional"`
}
