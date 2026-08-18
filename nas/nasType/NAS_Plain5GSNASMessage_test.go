package nasType_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/acore2026/free6gc-lib/nas/nasType"
)

func TestNasTypeNewPlain5GSNASMessage(t *testing.T) {
	a := nasType.NewPlain5GSNASMessage()
	assert.NotNil(t, a)
}
