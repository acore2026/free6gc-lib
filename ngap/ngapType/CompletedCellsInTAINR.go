package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct CompletedCellsInTAI_NR */
/* CompletedCellsInTAINRItem */
type CompletedCellsInTAINR struct {
	List []CompletedCellsInTAINRItem `aper:"valueExt,sizeLB:1,sizeUB:65535"`
}
