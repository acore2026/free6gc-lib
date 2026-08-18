package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type EmergencyFallbackIndicator struct {
	EmergencyFallbackRequestIndicator EmergencyFallbackRequestIndicator
	EmergencyServiceTargetCN          *EmergencyServiceTargetCN                                   `aper:"optional"`
	IEExtensions                      *ProtocolExtensionContainerEmergencyFallbackIndicatorExtIEs `aper:"optional"`
}
