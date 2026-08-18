package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type UserPlaneSecurityInformation struct {
	SecurityResult     SecurityResult                                                `aper:"valueExt"`
	SecurityIndication SecurityIndication                                            `aper:"valueExt"`
	IEExtensions       *ProtocolExtensionContainerUserPlaneSecurityInformationExtIEs `aper:"optional"`
}
