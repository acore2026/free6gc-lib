package ngapType

import "github.com/acore2026/free6gc-lib/aper"

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type PortNumber struct {
	Value aper.OctetString `aper:"sizeLB:2,sizeUB:2"`
}
