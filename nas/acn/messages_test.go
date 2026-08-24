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
	groupConfig := json.RawMessage(`{ "group_id":"g1", "members":["a1","a2"], "future":true }`)
	messages := []acn.Message{
		&acn.AgentRegisterRequest{Header: header, Owner: "u1", AgentName: "Alice", PublicKey: []byte{1, 2, 3}, Description: "Model-X", Timestamp: 1, Signature: []byte("sig"), Metadata: json.RawMessage(`{ "region":"CN", "os":"Linux", "version":"1.0.0", "future":true }`)},
		&acn.AgentRegisterAccept{Header: header, AgentID: "a1", VC0: json.RawMessage(`{"id":"vc0"}`)},
		&acn.AgentRegisterReject{Header: header, Cause: acn.RegisterRejectSignatureInvalid, FailedField: acn.RegisterFieldSignature},
		&acn.AgentDeregisterRequest{Header: header, AgentID: "a1", Reason: acn.DeregistrationReasonRetired, Timestamp: 2, Proof: proof},
		&acn.AgentDeregisterAccept{Header: header},
		&acn.AgentDeregisterReject{Header: header, Cause: acn.DeregisterRejectAgentNotFound, FailedField: acn.DeregisterFieldAgentID},
		&acn.AgentNetworkAbilityRequest{Header: header, NetworkAbilityRequest: json.RawMessage(`{"agent_id":"a1","intent":"Issue Network Ability Credential","timestamp":3,"proof":{"jws":"sig"}}`)},
		&acn.AgentNetworkAbilityAccept{Header: header, NetworkAbilityResult: json.RawMessage(`{"timestamp":4,"vc1":{"id":"vc1"}}`)},
		&acn.AgentNetworkAbilityReject{Header: header, Cause: acn.NetworkAbilityRejectNotAuthorized, FailedField: acn.NetworkAbilityFieldProof},
		&acn.AgentPublishRequest{Header: header, ProfilePublish: json.RawMessage(`{"agent_id":"a1","priority":2,"timestamp":5,"proof":{"jws":"sig"},"vc_list":[{"id":"vc2"}]}`)},
		&acn.AgentPublishAccept{Header: header},
		&acn.AgentPublishReject{Header: header, Cause: acn.PublishRejectVCInvalid, FailedField: acn.PublishFieldVCList},
		&acn.AgentProfileUpdateRequest{Header: header, ProfileUpdate: json.RawMessage(`{"agent_id":"a1","updates":[{"skill":"Driving"}],"timestamp":6}`)},
		&acn.AgentProfileUpdateAccept{Header: header, ProfileUpdateResult: json.RawMessage(`{"operation_id":"op-1","message":"Agent card updated"}`)},
		&acn.AgentProfileUpdateReject{Header: header, Cause: acn.ProfileUpdateRejectCredentialInvalid, FailedField: acn.ProfileUpdateFieldCredentials, RelatedItemIndex: 1},
		&acn.AgentSearchRequest{Header: header, DiscoveryRequest: json.RawMessage(`{"requester_agent_id":"a1","required_skills":["camera","radar"],"max_results":10,"timestamp":7}`)},
		&acn.AgentSearchResponse{Header: header, DiscoveryResult: json.RawMessage(`{"results":[{"agent_id":"a2","skills":["patrol"]},{"agent_id":"a3","skills":["drive"]}],"timestamp":8}`)},
		&acn.AgentSearchReject{Header: header, Cause: acn.SearchRejectProofInvalid, FailedField: acn.SearchFieldProof},
		&acn.AgentGroupingRequest{Header: header, GroupCreation: json.RawMessage(`{"source_agent_id":"a1","target_agent_ids":["a2","a3"],"timestamp":9}`)},
		&acn.AgentGroupingAccept{Header: header, GroupID: "g1"},
		&acn.AgentGroupingReject{Header: header, Cause: acn.GroupingRejectTargetAgentRejected, FailedField: acn.GroupingFieldTargetAgents, RelatedAgentID: "a2"},
		&acn.AgentGroupingInvitation{Header: header, GroupInvitation: json.RawMessage(`{"group_id":"g1","group_administrator":{"display_name":"Alice"}}`)},
		&acn.AgentGroupingInvitationResponse{Header: header, Decision: acn.GroupingDecisionAccept},
		&acn.AgentGroupInfoNotification{Header: header, GroupConfig: groupConfig},
		&acn.AgentGroupInfoNotificationResponse{Header: header, ApplyResult: acn.GroupInfoApplyResultACK},
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

func TestCodecCarriesCompleteMessageProofUnchanged(t *testing.T) {
	request := json.RawMessage(`{"agent_id":"agent-1","intent":"Issue Network Ability Credential","timestamp":"2026-08-24T12:00:00.000Z","proof":{"type":"JsonWebSignature2020","verification_method":"did:key:test#test","proof_purpose":"authentication","created":"2026-08-24T12:00:00.000Z","jws":"protected..signature"}}`)
	message := &acn.AgentNetworkAbilityRequest{
		Header: acn.Header{TransactionID: 7}, NetworkAbilityRequest: request,
	}
	wire, err := acn.Marshal(message)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := acn.Unmarshal(acn.Uplink, wire)
	if err != nil {
		t.Fatal(err)
	}
	got := decoded.(*acn.AgentNetworkAbilityRequest)
	if string(got.NetworkAbilityRequest) != string(request) {
		t.Fatalf("network ability request = %s, want exact %s", got.NetworkAbilityRequest, request)
	}
}
