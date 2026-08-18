package ngapType

import "github.com/acore2026/free6gc-lib/aper"

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

const (
	WAGFIDPresentNothing int = iota /* No components present */
	WAGFIDPresentWAGFID
	WAGFIDPresentChoiceExtensions
)

type WAGFID struct {
	Present          int             /* Choice Type */
	WAGFID           *aper.BitString `aper:"sizeExt,sizeLB:16,sizeUB:16"`
	ChoiceExtensions *ProtocolIESingleContainerWAGFIDExtIEs
}
