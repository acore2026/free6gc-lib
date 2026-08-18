package ngapType

import "github.com/acore2026/free6gc-lib/aper"

// Need to import "github.com/acore2026/free6gc-lib/aper" if it uses "aper"

const (
	NotificationControlPresentNotificationRequested aper.Enumerated = 0
)

type NotificationControl struct {
	Value aper.Enumerated `aper:"valueExt,valueLB:0,valueUB:0"`
}
