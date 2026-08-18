package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type AMFName struct {
	Value string `aper:"sizeExt,sizeLB:1,sizeUB:150"`
}
