package acn_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/acore2026/free6gc-lib/nas/acn"
)

func TestCodecRejectsInvalidEnvelopeFields(t *testing.T) {
	tests := []struct {
		name      string
		direction acn.Direction
		payload   []byte
		want      error
	}{
		{name: "version", direction: acn.Downlink, payload: []byte{2, 0x05, 1}, want: acn.ErrUnsupportedVersion},
		{name: "message type", direction: acn.Downlink, payload: []byte{1, 0xff, 1}, want: acn.ErrUnknownMessageType},
		{name: "transaction id", direction: acn.Downlink, payload: []byte{1, 0x05, 0}, want: acn.ErrInvalidValue},
		{name: "direction", direction: acn.Uplink, payload: []byte{1, 0x05, 1}, want: acn.ErrWrongDirection},
		{name: "trailing data", direction: acn.Downlink, payload: []byte{1, 0x05, 1, 0}, want: acn.ErrTrailingData},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := acn.Unmarshal(test.direction, test.payload)
			if !errors.Is(err, test.want) {
				t.Fatalf("error = %v, want %v", err, test.want)
			}
		})
	}
}

func TestCodecRejectsInvalidCurrentMessages(t *testing.T) {
	tests := []acn.Message{
		&acn.AgentNetworkAbilityRequest{Header: acn.Header{TransactionID: 1}, AgentID: "a1"},
		&acn.AgentPublishRequest{Header: acn.Header{TransactionID: 1}, AgentID: "a1", Priority: acn.Priority(9), Timestamp: 1, Signature: []byte("sig"), VCList: []byte(`[{"id":"vc"}]`)},
		&acn.AgentSearchRequest{Header: acn.Header{TransactionID: 1}, TaskID: "t1", AgentID: "a1", TaskDescription: "task", DiscoveryScope: acn.DiscoveryScopeIntraPLMN, MaxResults: 1},
		&acn.AgentRegisterRequest{Header: acn.Header{TransactionID: 1}, Owner: "u1", AgentName: "a1", PublicKey: []byte("key"), Description: "agent", Timestamp: 1, Signature: []byte("sig"), Metadata: json.RawMessage(`[]`)},
		&acn.AgentProfileUpdateRequest{Header: acn.Header{TransactionID: 1}, RequestID: "r1", AgentID: "a1", UpdateItems: json.RawMessage(`{}`), Timestamp: 1, Proof: json.RawMessage(`{"jws":"x"}`)},
		&acn.AgentGroupingInvitation{Header: acn.Header{TransactionID: 1}, GroupConfig: json.RawMessage(`{}`), GroupAdministrator: json.RawMessage(`{"agent_id":"a1"}`)},
		&acn.AgentGroupingInvitationResponse{Header: acn.Header{TransactionID: 1}, Decision: acn.GroupingDecision(9)},
	}
	for _, message := range tests {
		t.Run(message.MessageType().String(), func(t *testing.T) {
			if _, err := acn.Marshal(message); err == nil {
				t.Fatal("Marshal() accepted invalid message")
			}
		})
	}
}
