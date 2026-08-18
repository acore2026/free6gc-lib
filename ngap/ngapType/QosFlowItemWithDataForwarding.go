package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type QosFlowItemWithDataForwarding struct {
	QosFlowIdentifier      QosFlowIdentifier
	DataForwardingAccepted *DataForwardingAccepted                                        `aper:"optional"`
	IEExtensions           *ProtocolExtensionContainerQosFlowItemWithDataForwardingExtIEs `aper:"optional"`
}
