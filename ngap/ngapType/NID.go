package ngapType

import "github.com/acore2026/free6gc-lib/aper"

type NID struct {
	Value aper.BitString `aper:"sizeLB:44,sizeUB:44"`
}
