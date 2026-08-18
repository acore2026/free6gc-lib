package ngapType

import "github.com/acore2026/free6gc-lib/aper"

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type DRBStatusUL18 struct {
	ULCOUNTValue              COUNTValueForPDCPSN18                          `aper:"valueExt"`
	ReceiveStatusOfULPDCPSDUs *aper.BitString                                `aper:"sizeLB:1,sizeUB:131072,optional"`
	IEExtension               *ProtocolExtensionContainerDRBStatusUL18ExtIEs `aper:"optional"`
}
