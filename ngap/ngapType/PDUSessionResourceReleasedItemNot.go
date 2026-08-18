package ngapType

import "github.com/acore2026/free6gc-lib/aper"

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

type PDUSessionResourceReleasedItemNot struct {
	PDUSessionID                             PDUSessionID
	PDUSessionResourceNotifyReleasedTransfer aper.OctetString
	IEExtensions                             *ProtocolExtensionContainerPDUSessionResourceReleasedItemNotExtIEs `aper:"optional"`
}
