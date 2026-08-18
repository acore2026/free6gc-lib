package nasType_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/acore2026/free6gc-lib/nas"
	"github.com/acore2026/free6gc-lib/nas/nasType"
)

type nasTypeDeregistrationAcceptMessageIdentityData struct {
	in  uint8
	out uint8
}

var nasTypeDeregistrationAcceptMessageIdentityTable = []nasTypeDeregistrationAcceptMessageIdentityData{
	{nas.MsgTypeDeregistrationAcceptUETerminatedDeregistration, nas.MsgTypeDeregistrationAcceptUETerminatedDeregistration},
}

func TestNasTypeNewDeregistrationAcceptMessageIdentity(t *testing.T) {
	a := nasType.NewDeregistrationAcceptMessageIdentity()
	assert.NotNil(t, a)
}

func TestNasTypeGetSetDeregistrationAcceptMessageIdentity(t *testing.T) {
	a := nasType.NewDeregistrationAcceptMessageIdentity()
	for _, table := range nasTypeDeregistrationAcceptMessageIdentityTable {
		a.SetMessageType(table.in)
		assert.Equal(t, table.out, a.GetMessageType())
	}
}
