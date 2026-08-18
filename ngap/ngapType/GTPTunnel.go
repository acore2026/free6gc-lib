package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type GTPTunnel struct {
	TransportLayerAddress TransportLayerAddress
	GTPTEID               GTPTEID
	IEExtensions          *ProtocolExtensionContainerGTPTunnelExtIEs `aper:"optional"`
}
