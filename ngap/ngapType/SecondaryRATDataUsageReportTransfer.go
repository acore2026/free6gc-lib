package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type SecondaryRATDataUsageReportTransfer struct {
	SecondaryRATUsageInformation *SecondaryRATUsageInformation                                        `aper:"valueExt,optional"`
	IEExtensions                 *ProtocolExtensionContainerSecondaryRATDataUsageReportTransferExtIEs `aper:"optional"`
}
