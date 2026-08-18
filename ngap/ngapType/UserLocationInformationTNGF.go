package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type UserLocationInformationTNGF struct {
	TNAPID       TNAPID
	IPAddress    TransportLayerAddress
	PortNumber   *PortNumber                                                  `aper:"optional"`
	IEExtensions *ProtocolExtensionContainerUserLocationInformationTNGFExtIEs `aper:"optional"`
}
