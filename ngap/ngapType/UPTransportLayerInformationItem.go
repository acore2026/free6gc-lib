package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type UPTransportLayerInformationItem struct {
	NGUUPTNLInformation UPTransportLayerInformation                                      `aper:"valueLB:0,valueUB:1"`
	IEExtensions        *ProtocolExtensionContainerUPTransportLayerInformationItemExtIEs `aper:"optional"`
}
