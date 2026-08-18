package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type NumberOfBroadcastsRequested struct {
	Value int64 `aper:"valueLB:0,valueUB:65535"`
}
