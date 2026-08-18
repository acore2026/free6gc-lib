package ngapType

import "github.com/acore2026/free6gc-lib/aper"

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

const (
	CauseTransportPresentTransportResourceUnavailable aper.Enumerated = 0
	CauseTransportPresentUnspecified                  aper.Enumerated = 1
)

type CauseTransport struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:1"`
}
