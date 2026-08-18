package ngapType

import "github.com/acore2026/free6gc-lib/aper"

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type PDUSessionUsageReport struct {
	RATType                   aper.Enumerated
	PDUSessionTimedReportList VolumeTimedReportList
	IEExtensions              *ProtocolExtensionContainerPDUSessionUsageReportExtIEs `aper:"optional"`
}
