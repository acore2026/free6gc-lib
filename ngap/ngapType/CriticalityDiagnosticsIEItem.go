package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type CriticalityDiagnosticsIEItem struct {
	IECriticality Criticality
	IEID          ProtocolIEID
	TypeOfError   TypeOfError
	IEExtensions  *ProtocolExtensionContainerCriticalityDiagnosticsIEItemExtIEs `aper:"optional"`
}
