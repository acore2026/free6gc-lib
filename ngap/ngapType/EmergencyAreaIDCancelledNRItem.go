package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type EmergencyAreaIDCancelledNRItem struct {
	EmergencyAreaID       EmergencyAreaID
	CancelledCellsInEAINR CancelledCellsInEAINR
	IEExtensions          *ProtocolExtensionContainerEmergencyAreaIDCancelledNRItemExtIEs `aper:"optional"`
}
