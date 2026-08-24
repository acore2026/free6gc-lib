package acn_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/acore2026/free6gc-lib/nas/acn"
)

func TestCurrentMessagesRoundTrip(t *testing.T) {
	header := acn.Header{TransactionID: 33}
	proof := json.RawMessage(`{ "type":"JsonWebSignature2020", "verification_method":"did:key:k1", "jws":"sig", "future":true }`)
	groupConfig := json.RawMessage(`{ "group_name" : "task-patrol", "scope":"private", "max_members":5, "future":true }`)
	messages := []acn.Message{
		&acn.AgentRegisterRequest{Header: header, Owner: "u1", AgentName: "Alice", PublicKey: []byte{1, 2, 3}, Description: "Model-X", Timestamp: 1, Signature: []byte("sig"), Metadata: json.RawMessage(`{ "region":"CN", "os":"Linux", "version":"1.0.0", "future":true }`)},
		&acn.AgentRegisterAccept{Header: header, AgentID: "a1", VC0: json.RawMessage(`{"id":"vc0"}`)},
		&acn.AgentRegisterReject{Header: header, Cause: acn.RegisterRejectSignatureInvalid, FailedField: acn.RegisterFieldSignature},
		&acn.AgentDeregisterRequest{Header: header, AgentID: "a1", Reason: acn.DeregistrationReasonRetired, Timestamp: 2, Proof: proof},
		&acn.AgentDeregisterAccept{Header: header},
		&acn.AgentDeregisterReject{Header: header, Cause: acn.DeregisterRejectAgentNotFound, FailedField: acn.DeregisterFieldAgentID},
		&acn.AgentNetworkAbilityRequest{Header: header, AgentID: "a1", Intent: "Issue Network Ability Credential", Timestamp: 3, Proof: proof},
		&acn.AgentNetworkAbilityAccept{Header: header, Timestamp: 4, VC1: json.RawMessage(`{"id":"vc1"}`)},
		&acn.AgentNetworkAbilityReject{Header: header, Cause: acn.NetworkAbilityRejectNotAuthorized, FailedField: acn.NetworkAbilityFieldProof},
		&acn.AgentPublishRequest{Header: header, AgentID: "a1", Priority: acn.PriorityNormal, Timestamp: 5, Proof: proof, VCList: json.RawMessage(`[{"id":"vc2"}]`)},
		&acn.AgentPublishAccept{Header: header},
		&acn.AgentPublishReject{Header: header, Cause: acn.PublishRejectVCInvalid, FailedField: acn.PublishFieldVCList},
		&acn.AgentProfileUpdateRequest{Header: header, AgentID: "a1", UpdateItems: json.RawMessage(`[{"update_type":"add_skill","skill_name":"Driving","reference_vc_id":"vc2"},{"update_type":"remove_skill","skill_name":"Vision","future":true}]`), Credentials: json.RawMessage(`[{"id":"vc2"}]`), Timestamp: 6, Proof: proof},
		&acn.AgentProfileUpdateAccept{Header: header, OperationID: "op-1", Message: "Agent card updated"},
		&acn.AgentProfileUpdateReject{Header: header, Cause: acn.ProfileUpdateRejectCredentialInvalid, FailedField: acn.ProfileUpdateFieldCredentials, RelatedItemIndex: 1},
		&acn.AgentSearchRequest{Header: header, RequesterAgentID: "a1", RequiredSkills: []string{"camera", "radar"}, DiscoveryScope: acn.DiscoveryScopeIntraPLMN, MaxResults: 10, Timestamp: 7, Proof: proof},
		&acn.AgentSearchResponse{Header: header, Results: []acn.SearchResult{{AgentCard: json.RawMessage(`{ "agent_id":"a2", "ipv4":"10.60.0.12", "skills":["patrol"], "future":true }`), Priority: acn.PriorityHigh}, {AgentCard: json.RawMessage(`{"agent_id":"a3","ipv6":"2001:db8::3","skills":["drive"]}`), Priority: acn.PriorityNormal}}, Timestamp: 8},
		&acn.AgentSearchReject{Header: header, Cause: acn.SearchRejectProofInvalid, FailedField: acn.SearchFieldProof},
		&acn.AgentGroupingRequest{Header: header, AgentID: "a1", TargetAgentIDs: []string{"a2", "a3"}, GroupConfig: groupConfig, Timestamp: 9, Proof: proof},
		&acn.AgentGroupingAccept{Header: header, GroupID: "g1"},
		&acn.AgentGroupingReject{Header: header, Cause: acn.GroupingRejectTargetAgentRejected, FailedField: acn.GroupingFieldTargetAgents, RelatedAgentID: "a2"},
		&acn.AgentGroupingInvitation{Header: header, GroupConfig: groupConfig, GroupAdministrator: json.RawMessage(`{ "display_name":"Alice", "future":true }`)},
		&acn.AgentGroupingInvitationResponse{Header: header, Decision: acn.GroupingDecisionAccept},
		&acn.AgentGroupInfoNotification{Header: header, GroupID: "g1", Members: json.RawMessage(`{"agent1":{"agent_id":"a1"}}`)},
		&acn.AgentGroupInfoNotificationResponse{Header: header, GroupID: "g1", AgentID: "a1", Status: acn.GroupInfoStatusAccepted, Detail: "connected", Timestamp: 11, Proof: proof},
	}

	for _, message := range messages {
		t.Run(message.MessageType().String(), func(t *testing.T) {
			wire, err := acn.Marshal(message)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}
			decoded, err := acn.Unmarshal(message.Direction(), wire)
			if err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}
			if !reflect.DeepEqual(decoded, message) {
				t.Fatalf("round trip mismatch\n got: %#v\nwant: %#v", decoded, message)
			}
		})
	}
}

func TestCurrentMessageAssignments(t *testing.T) {
	if acn.Version1 != 0x01 {
		t.Fatalf("Version1 = %#x, want 0x01", acn.Version1)
	}
	want := map[acn.MessageType]uint8{
		acn.MessageTypeAgentGroupingRequest:       0x0d,
		acn.MessageTypeAgentGroupingAccept:        0x0e,
		acn.MessageTypeAgentGroupingReject:        0x0f,
		acn.MessageTypeAgentNetworkAbilityRequest: 0x14,
		acn.MessageTypeAgentNetworkAbilityAccept:  0x15,
		acn.MessageTypeAgentNetworkAbilityReject:  0x16,
		acn.MessageTypeAgentPublishRequest:        0x17,
		acn.MessageTypeAgentPublishAccept:         0x18,
		acn.MessageTypeAgentPublishReject:         0x19,
	}
	for messageType, value := range want {
		if uint8(messageType) != value {
			t.Fatalf("%s = %#x, want %#x", messageType, messageType, value)
		}
	}
	if got := acn.MessageTypeAgentNetworkAbilityAccept.String(); got != "ACN_AGENT_NETWORK_ABILITY_ACCEPT" {
		t.Fatalf("network ability message name = %q", got)
	}
}
