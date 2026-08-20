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
			name: "network ability request", direction: acn.Uplink,
			hex:      "7e00670e00a2011425000261310000019fc6f1c04000917b2274797065223a224a736f6e5765625369676e617475726532303230222c22766572696669636174696f6e5f6d6574686f64223a226469643a6b65793a6b31222c2270726f6f665f707572706f7365223a2261757468656e7469636174696f6e222c2263726561746564223a22323032362d30382d30335430393a32363a30305a222c226a7773223a226a777331227d",
			wantType: acn.MessageTypeAgentNetworkAbilityRequest,
		},
		{
			name: "search response", direction: acn.Downlink,
			hex:      "7e00680e0107010a2600027431000d506174726f6c20417265612041020072006f7b226167656e745f6964223a226132222c226167656e745f6970223a22382e382e382e38222c227463705f706f7274223a2234303031222c227564705f706f7274223a223238343433222c22736b696c6c73223a5b22706174726f6c222c226472697665222c22646f672d41225d7d010072006f7b226167656e745f6964223a226133222c226167656e745f6970223a22382e382e382e37222c227463705f706f7274223a2234303032222c227564705f706f7274223a223238343435222c22736b696c6c73223a5b22706174726f6c222c226472697665222c22646f672d42225d7d020000018e6b2e9a00",
			wantType: acn.MessageTypeAgentSearchResponse,
		},
		{
			name: "grouping invitation", direction: acn.Downlink,
			hex:      "7e00680e0066011031003e7b2267726f75705f6e616d65223a227461736b2d706174726f6c222c2273636f7065223a2270726976617465222c226d61785f6d656d62657273223a357d00217b226167656e745f6964223a226131222c22736b696c6c73223a5b224152225d7d",
			wantType: acn.MessageTypeAgentGroupingInvitation,
		},
		{
			name: "group info response", direction: acn.Uplink,
			hex:      "7e00670e00b20113320002673100026131000009636f6e6e65637465640000019fd5198f9000917b2274797065223a224a736f6e5765625369676e617475726532303230222c22766572696669636174696f6e5f6d6574686f64223a226469643a6b65793a6b31222c2270726f6f665f707572706f7365223a2261757468656e7469636174696f6e222c2263726561746564223a22323032362d30382d30365430333a32343a31305a222c226a7773223a226a777334227d",
			wantType: acn.MessageTypeAgentGroupInfoNotificationResponse,
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
