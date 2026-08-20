package acn_test

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/acore2026/free6gc-lib/nas/acn"
)

func TestGroupingInvitationWireFormat(t *testing.T) {
	message := &acn.AgentGroupingInvitation{
		Header:             acn.Header{TransactionID: 0x31},
		GroupConfig:        json.RawMessage(`{ "group_name" : "task-patrol", "scope":"private", "max_members":5 }`),
		GroupAdministrator: json.RawMessage(`{ "agent_id" : "a1", "skills":["AR"], "future":true }`),
	}
	wire, err := acn.Marshal(message)
	if err != nil {
		t.Fatal(err)
	}
	wantHex := "011031" + lveHex(message.GroupConfig) + lveHex(message.GroupAdministrator)
	want, err := hex.DecodeString(wantHex)
	if err != nil {
		t.Fatal(err)
	}
	if string(wire) != string(want) {
		t.Fatalf("wire = %x, want %x", wire, want)
	}
	decoded, err := acn.Unmarshal(acn.Downlink, wire)
	if err != nil {
		t.Fatal(err)
	}
	invitation := decoded.(*acn.AgentGroupingInvitation)
	if string(invitation.GroupConfig) != string(message.GroupConfig) || string(invitation.GroupAdministrator) != string(message.GroupAdministrator) {
		t.Fatalf("unexpected decoded invitation: %#v", invitation)
	}
}

func lveHex(value []byte) string {
	length := make([]byte, 2)
	binary.BigEndian.PutUint16(length, uint16(len(value)))
	return hex.EncodeToString(append(length, value...))
}
