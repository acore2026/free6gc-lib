package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct PLMNSupportList */
/* PLMNSupportItem */
type PLMNSupportList struct {
	List []PLMNSupportItem `aper:"valueExt,sizeLB:1,sizeUB:12"`
}
