package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type PDUSessionResourceSetupResponseTransfer struct {
	DLQosFlowPerTNLInformation           QosFlowPerTNLInformation                                                 `aper:"valueExt"`
	AdditionalDLQosFlowPerTNLInformation *QosFlowPerTNLInformationList                                            `aper:"optional"`
	SecurityResult                       *SecurityResult                                                          `aper:"valueExt,optional"`
	QosFlowFailedToSetupList             *QosFlowListWithCause                                                    `aper:"optional"`
	IEExtensions                         *ProtocolExtensionContainerPDUSessionResourceSetupResponseTransferExtIEs `aper:"optional"`
}
