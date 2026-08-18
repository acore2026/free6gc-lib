package nasType_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/acore2026/free6gc-lib/nas/nasType"
)

const disasterReturnWaitRangeIEI uint8 = 0x2C

var validDisasterReturnWaitRangeRaw = []byte{
	disasterReturnWaitRangeIEI,
	0x02,
	0x33,
	0x66,
}

func TestDisasterReturnWaitRangeDecode(t *testing.T) {
	ie := nasType.NewDisasterReturnWaitRange(0x00)

	err := ie.Decode(validDisasterReturnWaitRangeRaw)
	require.NoError(t, err)

	assert.Equal(t, disasterReturnWaitRangeIEI, ie.GetIei())
	assert.Equal(t, uint8(0x02), ie.GetLen())
	assert.Equal(t, uint8(0x33), ie.GetMinimumRegistrationWaitTime())
	assert.Equal(t, uint8(0x66), ie.GetMaximumRegistrationWaitTime())
}

func TestDisasterReturnWaitRangeEncode(t *testing.T) {
	ie := nasType.NewDisasterReturnWaitRange(disasterReturnWaitRangeIEI)
	ie.SetMinimumRegistrationWaitTime(0x33)
	ie.SetMaximumRegistrationWaitTime(0x66)

	encoded, err := ie.Encode()
	require.NoError(t, err)

	assert.Equal(t, validDisasterReturnWaitRangeRaw, encoded)
}
