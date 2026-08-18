package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type GlobalTWIFID struct { /* Sequence Type (Extensible) */
	PLMNIdentity PLMNIdentity
	TWIFID       TWIFID                                        `aper:"valueLB:0,valueUB:1"`
	IEExtensions *ProtocolExtensionContainerGlobalTWIFIDExtIEs `aper:"optional"`
}
