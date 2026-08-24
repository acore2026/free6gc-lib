package acn_test

import (
	"encoding/hex"
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

func TestUpdatedDocumentedPlainNASVectors(t *testing.T) {
	const proofHex = "7b2263726561746f72223a226469643a6131236b65792d31222c227369676e61747572655f76616c7565223a22736967227d"
	tests := []struct {
		name      string
		direction acn.Direction
		hex       string
		wantType  acn.MessageType
	}{
		{
			name: "register request", direction: acn.Uplink,
			hex:      "7e00670e005a010121000275310005416c69636500040102030400074d6f64656c2d580000019d2f78d0200003736967002e7b22726567696f6e223a22434e222c226f73223a224c696e7578222c2276657273696f6e223a22312e302e30227d",
			wantType: acn.MessageTypeAgentRegisterRequest,
		},
		{
			name: "deregister request", direction: acn.Uplink,
			hex:      "7e00670e00440104220002613105000001a01d01eaa00032" + proofHex,
			wantType: acn.MessageTypeAgentDeregisterRequest,
		},
		{
			name: "profile update request", direction: acn.Uplink,
			hex: "7e00670e007f010640000261310036" +
				"5b7b227570646174655f74797065223a2272656d6f76655f736b696c6c222c22736b696c6c5f6e616d65223a22766973696f6e227d5d" +
				"00025b5d0000019fc6f1c0400032" + proofHex,
			wantType: acn.MessageTypeAgentProfileUpdateRequest,
		},
		{
			name: "profile update accept", direction: acn.Downlink,
			hex:      "7e00680e0024010740000b6f702d7570646174652d3100124167656e7420636172642075706461746564",
			wantType: acn.MessageTypeAgentProfileUpdateAccept,
		},
		{
			name: "profile update reject", direction: acn.Downlink,
			hex:      "7e00680e000701084006020001",
			wantType: acn.MessageTypeAgentProfileUpdateReject,
		},
		{
			name: "search request", direction: acn.Uplink,
			hex: "7e00670e005d01092400026131030663616d65726105726164617209666f75725f6c656773010a" +
				"0000019fc6f1c0400032" + proofHex,
			wantType: acn.MessageTypeAgentSearchRequest,
		},
		{
			name: "network ability request", direction: acn.Uplink,
			hex:      "7e00670e00650114270002613100204973737565204e6574776f726b204162696c6974792043726564656e7469616c0000019fc6f1c04000327b2263726561746f72223a226469643a6131236b65792d31222c227369676e61747572655f76616c7565223a22736967227d",
			wantType: acn.MessageTypeAgentNetworkAbilityRequest,
		},
		{
			name: "search response", direction: acn.Downlink,
			hex:      "7e00680e0070010a24010062005f7b226167656e745f6964223a226132222c226167656e745f6970223a22382e382e382e38222c227463705f706f7274223a2234303031222c227564705f706f7274223a223238343433222c22736b696c6c73223a5b2263616d657261225d7d010000018e6b2e9a00",
			wantType: acn.MessageTypeAgentSearchResponse,
		},
		{
			name: "grouping invitation", direction: acn.Downlink,
			hex:      "7e00680e002901103100117b2267726f75705f6964223a226731227d00117b226167656e745f6964223a226131227d",
			wantType: acn.MessageTypeAgentGroupingInvitation,
		},
		{
			name: "group info notification", direction: acn.Downlink,
			hex:      "7e00680e002501123200026731001c7b226167656e7431223a7b226167656e745f6964223a226131227d7d",
			wantType: acn.MessageTypeAgentGroupInfoNotification,
		},
		{
			name: "group info response", direction: acn.Uplink,
			hex:      "7e00670e00530113320002673100026131000009636f6e6e65637465640000019fd5198f9000327b2263726561746f72223a226469643a6131236b65792d31222c227369676e61747572655f76616c7565223a22736967227d",
			wantType: acn.MessageTypeAgentGroupInfoNotificationResponse,
		},
		{
			name: "publish request", direction: acn.Uplink,
			hex: "7e00670e005b0117230002613102000001a01d02d5000032" + proofHex +
				"00155b7b226964223a2263616d6572612d303031227d5d",
			wantType: acn.MessageTypeAgentPublishRequest,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wire, err := hex.DecodeString(test.hex)
			if err != nil {
				t.Fatal(err)
			}
			message, err := acn.DecodePlainNAS(test.direction, wire)
			if err != nil {
				t.Fatal(err)
			}
			if message.MessageType() != test.wantType {
				t.Fatalf("message type = %s, want %s", message.MessageType(), test.wantType)
			}
			encoded, err := acn.EncodePlainNAS(message)
			if err != nil {
				t.Fatal(err)
			}
			if string(encoded) != string(wire) {
				t.Fatalf("wire = %x, want %x", encoded, wire)
			}
		})
	}
}
