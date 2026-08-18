package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type EmergencyAreaIDBroadcastNRItem struct {
	EmergencyAreaID       EmergencyAreaID
	CompletedCellsInEAINR CompletedCellsInEAINR
	IEExtensions          *ProtocolExtensionContainerEmergencyAreaIDBroadcastNRItemExtIEs `aper:"optional"`
}
