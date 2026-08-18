package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type XnTNLConfigurationInfo struct {
	XnTransportLayerAddresses         XnTLAs
	XnExtendedTransportLayerAddresses *XnExtTLAs                                              `aper:"optional"`
	IEExtensions                      *ProtocolExtensionContainerXnTNLConfigurationInfoExtIEs `aper:"optional"`
}
