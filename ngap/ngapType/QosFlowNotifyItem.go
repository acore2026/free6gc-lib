package ngapType

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type QosFlowNotifyItem struct {
	QosFlowIdentifier QosFlowIdentifier
	NotificationCause NotificationCause
	IEExtensions      *ProtocolExtensionContainerQosFlowNotifyItemExtIEs `aper:"optional"`
}
