package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type PDUSessionResourceModifyIndicationTransfer struct {
	DLQosFlowPerTNLInformation           QosFlowPerTNLInformation                                                    `aper:"valueExt"`
	AdditionalDLQosFlowPerTNLInformation *QosFlowPerTNLInformationList                                               `aper:"optional"`
	IEExtensions                         *ProtocolExtensionContainerPDUSessionResourceModifyIndicationTransferExtIEs `aper:"optional"`
}
