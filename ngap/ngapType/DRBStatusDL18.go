package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type DRBStatusDL18 struct {
	DLCOUNTValue COUNTValueForPDCPSN18                          `aper:"valueExt"`
	IEExtension  *ProtocolExtensionContainerDRBStatusDL18ExtIEs `aper:"optional"`
}
