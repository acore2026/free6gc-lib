package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

/* Sequence of = 35, FULL Name = struct AssociatedQosFlowList */
/* AssociatedQosFlowItem */
type AssociatedQosFlowList struct {
	List []AssociatedQosFlowItem `aper:"valueExt,sizeLB:1,sizeUB:64"`
}
