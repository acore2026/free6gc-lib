package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type NonDynamic5QIDescriptor struct {
	FiveQI                 FiveQI
	PriorityLevelQos       *PriorityLevelQos                                        `aper:"optional"`
	AveragingWindow        *AveragingWindow                                         `aper:"optional"`
	MaximumDataBurstVolume *MaximumDataBurstVolume                                  `aper:"optional"`
	IEExtensions           *ProtocolExtensionContainerNonDynamic5QIDescriptorExtIEs `aper:"optional"`
}
