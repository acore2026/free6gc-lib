package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type PDUSessionResourceInformationItem struct {
	PDUSessionID              PDUSessionID
	QosFlowInformationList    QosFlowInformationList
	DRBsToQosFlowsMappingList *DRBsToQosFlowsMappingList                                         `aper:"optional"`
	IEExtensions              *ProtocolExtensionContainerPDUSessionResourceInformationItemExtIEs `aper:"optional"`
}
