package ngapType

import "github.com/acore2026/free6gc-lib/aper"

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type SecurityKey struct {
	Value aper.BitString `aper:"sizeLB:256,sizeUB:256"`
}
