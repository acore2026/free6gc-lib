package acn_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/acore2026/free6gc-lib/nas/acn"
)

func TestGroupMessageRoundTrips(t *testing.T) {
	tests := []acn.Message{
		&acn.AgentGroupingInvitation{
			Header:        acn.Header{TransactionID: 49},
			GroupID:       "g1",
			SourceAgentID: "a1",
			TargetAgentID: "a2",
			TaskID:        "task-1",
			ExpiresAt:     1786608030000,
			Proof:         []byte("sig"),
		},
		&acn.AgentGroupingInvitationResponse{
			Header:       acn.Header{TransactionID: 49},
			GroupID:      "g1",
			AgentID:      "a2",
			Decision:     acn.GroupingDecisionAccept,
			RejectReason: acn.GroupingRejectReasonNone,
			Timestamp:    1786608010000,
			Proof:        []byte("sig"),
		},
		&acn.AgentGroupInfoNotification{
			Header:        acn.Header{TransactionID: 50},
			GroupID:       "g1",
			TargetAgentID: "a2",
			GroupConfig:   json.RawMessage(`{"relay":"r1"}`),
		},
		&acn.AgentGroupInfoNotificationResponse{
			Header:       acn.Header{TransactionID: 50},
			GroupID:      "g1",
			AgentID:      "a2",
			Result:       acn.GroupInfoApplyResultSuccess,
			FailureCause: acn.GroupInfoFailureCauseNone,
		},
	}

	for _, message := range tests {
		t.Run(message.MessageType().String(), func(t *testing.T) {
			payload, err := acn.Marshal(message)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}
			decoded, err := acn.Unmarshal(message.Direction(), payload)
			if err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}
			if !reflect.DeepEqual(decoded, message) {
				t.Fatalf("round trip mismatch:\n got: %#v\nwant: %#v", decoded, message)
			}
		})
	}
}
