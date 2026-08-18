package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type NetworkInstance struct {
	Value int64 `aper:"valueExt,valueLB:1,valueUB:256"`
}
