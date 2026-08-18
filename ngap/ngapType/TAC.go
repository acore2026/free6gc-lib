package ngapType

import "github.com/acore2026/free6gc-lib/aper"

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type TAC struct {
	Value aper.OctetString `aper:"sizeLB:3,sizeUB:3"`
}
