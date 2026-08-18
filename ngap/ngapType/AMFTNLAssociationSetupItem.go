package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type AMFTNLAssociationSetupItem struct {
	AMFTNLAssociationAddress CPTransportLayerInformation                                 `aper:"valueLB:0,valueUB:1"`
	IEExtensions             *ProtocolExtensionContainerAMFTNLAssociationSetupItemExtIEs `aper:"optional"`
}
