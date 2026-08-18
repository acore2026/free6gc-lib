package ngapType

import "github.com/acore2026/free6gc-lib/aper"

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type VolumeTimedReportItem struct {
	StartTimeStamp aper.OctetString                                       `aper:"sizeLB:4,sizeUB:4"`
	EndTimeStamp   aper.OctetString                                       `aper:"sizeLB:4,sizeUB:4"`
	UsageCountUL   int64                                                  `aper:"valueLB:0,valueUB:18446744073709551615"`
	UsageCountDL   int64                                                  `aper:"valueLB:0,valueUB:18446744073709551615"`
	IEExtensions   *ProtocolExtensionContainerVolumeTimedReportItemExtIEs `aper:"optional"`
}
