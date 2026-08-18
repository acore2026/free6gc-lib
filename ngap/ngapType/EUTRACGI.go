package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type EUTRACGI struct {
	PLMNIdentity      PLMNIdentity
	EUTRACellIdentity EUTRACellIdentity
	IEExtensions      *ProtocolExtensionContainerEUTRACGIExtIEs `aper:"optional"`
}
