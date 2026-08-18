package ngapType

import "github.com/acore2026/free6gc-lib/aper"

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

const (
	CancelAllWarningMessagesPresentTrue aper.Enumerated = 0
)

type CancelAllWarningMessages struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:0"`
}
