package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type XnExtTLAItem struct {
	IPsecTLA     *TransportLayerAddress                        `aper:"optional"`
	GTPTLAs      *XnGTPTLAs                                    `aper:"optional"`
	IEExtensions *ProtocolExtensionContainerXnExtTLAItemExtIEs `aper:"optional"`
}
