package ngapType

import "github.com/acore2026/free6gc-lib/aper"

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type AMFPointer struct {
	Value aper.BitString `aper:"sizeLB:6,sizeUB:6"`
}
