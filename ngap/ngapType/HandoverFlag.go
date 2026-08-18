package ngapType

import "github.com/acore2026/free6gc-lib/aper"

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

const (
	HandoverFlagPresentHandoverPreparation aper.Enumerated = 0
)

type HandoverFlag struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:0"`
}
