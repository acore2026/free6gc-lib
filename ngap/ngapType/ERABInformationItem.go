package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type ERABInformationItem struct {
	ERABID       ERABID
	DLForwarding *DLForwarding                                        `aper:"optional"`
	IEExtensions *ProtocolExtensionContainerERABInformationItemExtIEs `aper:"optional"`
}
