package acn_test

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/acore2026/free6gc-lib/nas/acn"
)

func TestDocumentedControlPDUVectors(t *testing.T) {
	tests := []struct {
		name      string
		direction acn.Direction
		hex       string
		wantType  acn.MessageType
	}{
		{name: "grouping accept", direction: acn.Downlink, hex: "010e30000967726f75702d303031", wantType: acn.MessageTypeAgentGroupingAccept},
		{name: "invitation response", direction: acn.Uplink, hex: "01113100", wantType: acn.MessageTypeAgentGroupingInvitationResponse},
		{name: "group info apply result", direction: acn.Uplink, hex: "01133200", wantType: acn.MessageTypeAgentGroupInfoNotificationResponse},
		{name: "network ability reject", direction: acn.Downlink, hex: "0116250203", wantType: acn.MessageTypeAgentNetworkAbilityReject},
		{name: "publish accept", direction: acn.Downlink, hex: "011823", wantType: acn.MessageTypeAgentPublishAccept},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wire, err := hex.DecodeString(test.hex)
			if err != nil {
				t.Fatal(err)
			}
			message, err := acn.Unmarshal(test.direction, wire)
			if err != nil {
				t.Fatal(err)
			}
			if message.MessageType() != test.wantType {
				t.Fatalf("message type = %s, want %s", message.MessageType(), test.wantType)
			}
			encoded, err := acn.Marshal(message)
			if err != nil {
				t.Fatal(err)
			}
			if string(encoded) != string(wire) {
				t.Fatalf("wire = %x, want %x", encoded, wire)
			}
		})
	}
}

func TestJSONContainerPlainNASVector(t *testing.T) {
	container := json.RawMessage(`{"agent_id":"a1","intent":"Issue Network Ability Credential"}`)
	message := &acn.AgentNetworkAbilityRequest{
		Header:                acn.Header{TransactionID: 0x27},
		NetworkAbilityRequest: container,
	}

	wire, err := acn.EncodePlainNAS(message)
	if err != nil {
		t.Fatal(err)
	}

	acnPayload := []byte{acn.Version1, byte(acn.MessageTypeAgentNetworkAbilityRequest), 0x27}
	length := make([]byte, 2)
	binary.BigEndian.PutUint16(length, uint16(len(container)))
	acnPayload = append(acnPayload, length...)
	acnPayload = append(acnPayload, container...)
	want := []byte{0x7e, 0x00, 0x67, acn.PayloadContainerTypeACN}
	binary.BigEndian.PutUint16(length, uint16(len(acnPayload)))
	want = append(want, length...)
	want = append(want, acnPayload...)
	if string(wire) != string(want) {
		t.Fatalf("wire = %x, want %x", wire, want)
	}

	decoded, err := acn.DecodePlainNAS(acn.Uplink, wire)
	if err != nil {
		t.Fatal(err)
	}
	request := decoded.(*acn.AgentNetworkAbilityRequest)
	if string(request.NetworkAbilityRequest) != string(container) {
		t.Fatalf("container = %s, want %s", request.NetworkAbilityRequest, container)
	}
}
