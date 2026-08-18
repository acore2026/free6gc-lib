package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type SecurityIndication struct {
	IntegrityProtectionIndication       IntegrityProtectionIndication
	ConfidentialityProtectionIndication ConfidentialityProtectionIndication
	MaximumIntegrityProtectedDataRateUL *MaximumIntegrityProtectedDataRate                  `aper:"optional"`
	IEExtensions                        *ProtocolExtensionContainerSecurityIndicationExtIEs `aper:"optional"`
}
