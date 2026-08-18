package ngapType

import "github.com/acore2026/free6gc-lib/aper"

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

const (
	TNGFIDPresentNothing int = iota /* No components present */
	TNGFIDPresentTNGFID
	TNGFIDPresentChoiceExtensions
)

type TNGFID struct {
	Present          int             /* Choice Type */
	TNGFID           *aper.BitString `aper:"sizeLB:32,sizeUB:32"`
	ChoiceExtensions *ProtocolIESingleContainerTNGFIDExtIEs
}
