package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type CoreNetworkAssistanceInformation struct {
	UEIdentityIndexValue            UEIdentityIndexValue `aper:"valueLB:0,valueUB:1"`
	UESpecificDRX                   *PagingDRX           `aper:"optional"`
	PeriodicRegistrationUpdateTimer PeriodicRegistrationUpdateTimer
	MICOModeIndication              *MICOModeIndication `aper:"optional"`
	TAIListForInactive              TAIListForInactive
	ExpectedUEBehaviour             *ExpectedUEBehaviour                                              `aper:"valueExt,optional"`
	IEExtensions                    *ProtocolExtensionContainerCoreNetworkAssistanceInformationExtIEs `aper:"optional"`
}
