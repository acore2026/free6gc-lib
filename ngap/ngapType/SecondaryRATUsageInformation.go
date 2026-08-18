package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type SecondaryRATUsageInformation struct {
	PDUSessionUsageReport   *PDUSessionUsageReport                                        `aper:"valueExt,optional"`
	QosFlowsUsageReportList *QoSFlowsUsageReportList                                      `aper:"optional"`
	IEExtension             *ProtocolExtensionContainerSecondaryRATUsageInformationExtIEs `aper:"optional"`
}
