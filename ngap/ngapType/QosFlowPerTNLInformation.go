package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type QosFlowPerTNLInformation struct {
	UPTransportLayerInformation UPTransportLayerInformation `aper:"valueLB:0,valueUB:1"`
	AssociatedQosFlowList       AssociatedQosFlowList
	IEExtensions                *ProtocolExtensionContainerQosFlowPerTNLInformationExtIEs `aper:"optional"`
}
