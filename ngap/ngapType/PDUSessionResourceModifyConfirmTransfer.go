package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type PDUSessionResourceModifyConfirmTransfer struct {
	QosFlowModifyConfirmList      QosFlowModifyConfirmList
	ULNGUUPTNLInformation         UPTransportLayerInformation                                              `aper:"valueLB:0,valueUB:1"`
	AdditionalNGUUPTNLInformation *UPTransportLayerInformationPairList                                     `aper:"optional"`
	QosFlowFailedToModifyList     *QosFlowListWithCause                                                    `aper:"optional"`
	IEExtensions                  *ProtocolExtensionContainerPDUSessionResourceModifyConfirmTransferExtIEs `aper:"optional"`
}
