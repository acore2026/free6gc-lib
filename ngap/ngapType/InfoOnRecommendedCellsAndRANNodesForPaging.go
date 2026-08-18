package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type InfoOnRecommendedCellsAndRANNodesForPaging struct {
	RecommendedCellsForPaging  RecommendedCellsForPaging                                                   `aper:"valueExt"`
	RecommendRANNodesForPaging RecommendedRANNodesForPaging                                                `aper:"valueExt"`
	IEExtensions               *ProtocolExtensionContainerInfoOnRecommendedCellsAndRANNodesForPagingExtIEs `aper:"optional"`
}
