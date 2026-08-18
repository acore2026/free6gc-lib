package ngapType

import "github.com/acore2026/free6gc-lib/aper"

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

const (
	DelayCriticalPresentDelayCritical    aper.Enumerated = 0
	DelayCriticalPresentNonDelayCritical aper.Enumerated = 1
)

type DelayCritical struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:1"`
}
