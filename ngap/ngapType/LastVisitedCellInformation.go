package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

const (
	LastVisitedCellInformationPresentNothing int = iota /* No components present */
	LastVisitedCellInformationPresentNGRANCell
	LastVisitedCellInformationPresentEUTRANCell
	LastVisitedCellInformationPresentUTRANCell
	LastVisitedCellInformationPresentGERANCell
	LastVisitedCellInformationPresentChoiceExtensions
)

type LastVisitedCellInformation struct {
	Present          int
	NGRANCell        *LastVisitedNGRANCellInformation `aper:"valueExt"`
	EUTRANCell       *LastVisitedEUTRANCellInformation
	UTRANCell        *LastVisitedUTRANCellInformation
	GERANCell        *LastVisitedGERANCellInformation
	ChoiceExtensions *ProtocolIESingleContainerLastVisitedCellInformationExtIEs
}
