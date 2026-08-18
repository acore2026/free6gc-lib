package acn_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/acore2026/free6gc-lib/nas/acn"
)

type goldenVector struct {
	name        string
	direction   acn.Direction
	messageType acn.MessageType
	transaction uint8
	wireHex     string
}

func goldenVectors() []goldenVector {
	return []goldenVector{
		{
			name:        "register request",
			direction:   acn.Uplink,
			messageType: acn.MessageTypeAgentRegisterRequest,
			transaction: 0x21,
			wireHex: `
				7E 00 67 0E 00 3C 01 01 21 00 02 75 31 00 05 41
				6C 69 63 65 00 04 01 02 03 04 00 07 4D 6F 64 65
				6C 2D 58 00 00 01 9D 2F 78 D0 20 00 03 73 69 67
				00 02 43 4E 00 05 4C 69 6E 75 78 00 05 31 2E 30
				2E 30`,
		},
		{
			name:        "register accept",
			direction:   acn.Downlink,
			messageType: acn.MessageTypeAgentRegisterAccept,
			transaction: 0x21,
			wireHex: `
				7E 00 68 0E 01 13 01 02 21 00 02 61 31 01 0A 7B
				22 63 6F 6E 74 65 78 74 22 3A 5B 22 63 22 5D 2C
				22 69 64 22 3A 22 63 72 65 64 31 22 2C 22 74 79
				70 65 22 3A 5B 22 56 43 22 2C 22 53 49 4D 22 5D
				2C 22 69 73 73 75 65 72 22 3A 22 69 64 6D 22 2C
				22 76 61 6C 69 64 5F 66 72 6F 6D 22 3A 22 32 30
				32 36 22 2C 22 76 61 6C 69 64 5F 75 6E 74 69 6C
				22 3A 22 32 30 32 37 22 2C 22 63 6C 61 69 6D 73
				22 3A 7B 22 61 67 65 6E 74 5F 6E 61 6D 65 22 3A
				22 41 6C 69 63 65 22 2C 22 61 67 65 6E 74 5F 69
				64 22 3A 22 61 31 22 2C 22 61 67 65 6E 74 5F 61
				74 74 72 69 62 75 74 65 22 3A 22 62 69 6E 64 69
				6E 67 22 2C 22 6D 61 73 74 65 72 5F 69 64 22 3A
				22 6D 31 22 2C 22 73 65 6C 66 5F 69 64 22 3A 22
				73 31 22 7D 2C 22 70 72 6F 6F 66 22 3A 7B 22 63
				72 65 61 74 6F 72 22 3A 22 69 64 6D 23 31 22 2C
				22 73 69 67 6E 61 74 75 72 65 5F 76 61 6C 75 65
				22 3A 22 73 69 67 22 7D 7D`,
		},
		{
			name:        "register reject",
			direction:   acn.Downlink,
			messageType: acn.MessageTypeAgentRegisterReject,
			transaction: 0x21,
			wireHex:     `7E 00 68 0E 00 05 01 03 21 04 06`,
		},
		{
			name:        "deregister request",
			direction:   acn.Uplink,
			messageType: acn.MessageTypeAgentDeregisterRequest,
			transaction: 0x22,
			wireHex: `
				7E 00 67 0E 00 15 01 04 22 00 02 61 31 05 00 00
				01 8E 6B 2E 9A 00 00 03 73 69 67`,
		},
		{
			name:        "deregister accept",
			direction:   acn.Downlink,
			messageType: acn.MessageTypeAgentDeregisterAccept,
			transaction: 0x22,
			wireHex:     `7E 00 68 0E 00 03 01 05 22`,
		},
		{
			name:        "deregister reject",
			direction:   acn.Downlink,
			messageType: acn.MessageTypeAgentDeregisterReject,
			transaction: 0x22,
			wireHex:     `7E 00 68 0E 00 05 01 0C 22 04 04`,
		},
		{
			name:        "profile update request",
			direction:   acn.Uplink,
			messageType: acn.MessageTypeAgentProfileUpdateRequest,
			transaction: 0x23,
			wireHex: `
				7E 00 67 0E 00 9D 01 06 23 00 02 61 31 02 00 00
				01 8E 6B 2E 9A 00 00 03 73 69 67 00 86 5B 7B 22
				69 64 22 3A 22 63 72 65 64 31 22 2C 22 74 79 70
				65 22 3A 5B 22 56 43 22 5D 2C 22 63 6C 61 69 6D
				73 22 3A 7B 22 61 67 65 6E 74 5F 61 74 74 72 69
				62 75 74 65 22 3A 22 73 65 72 76 69 63 65 22 7D
				7D 2C 7B 22 69 64 22 3A 22 63 72 65 64 32 22 2C
				22 74 79 70 65 22 3A 5B 22 56 43 22 5D 2C 22 63
				6C 61 69 6D 73 22 3A 7B 22 61 67 65 6E 74 5F 61
				74 74 72 69 62 75 74 65 22 3A 22 66 61 6C 6C 22
				7D 7D 5D`,
		},
		{
			name:        "profile update accept",
			direction:   acn.Downlink,
			messageType: acn.MessageTypeAgentProfileUpdateAccept,
			transaction: 0x23,
			wireHex:     `7E 00 68 0E 00 03 01 07 23`,
		},
		{
			name:        "profile update reject",
			direction:   acn.Downlink,
			messageType: acn.MessageTypeAgentProfileUpdateReject,
			transaction: 0x23,
			wireHex:     `7E 00 68 0E 00 05 01 08 23 04 05`,
		},
		{
			name:        "capability discovery request",
			direction:   acn.Uplink,
			messageType: acn.MessageTypeAgentSearchRequest,
			transaction: 0x24,
			wireHex: `
				7E 00 67 0E 00 2C 01 09 24 01 00 02 61 31 00 02
				74 31 00 00 01 8E 6B 2E 9A 00 03 06 63 61 6D 65
				72 61 05 72 61 64 61 72 09 66 6F 75 72 5F 6C 65
				67 73`,
		},
		{
			name:        "agent info request",
			direction:   acn.Uplink,
			messageType: acn.MessageTypeAgentSearchRequest,
			transaction: 0x25,
			wireHex: `
				7E 00 67 0E 00 0F 01 09 25 02 00 09 61 67 65 6E
				74 2D 31 31 31`,
		},
		{
			name:        "owner agents request",
			direction:   acn.Uplink,
			messageType: acn.MessageTypeAgentSearchRequest,
			transaction: 0x26,
			wireHex: `
				7E 00 67 0E 00 0F 01 09 26 03 00 09 6F 77 6E 65
				72 2D 31 32 33`,
		},
		{
			name:        "capability discovery response",
			direction:   acn.Downlink,
			messageType: acn.MessageTypeAgentSearchResponse,
			transaction: 0x24,
			wireHex: `
				7E 00 68 0E 00 4B 01 0A 24 01 00 01 00 43 00 09
				61 67 65 6E 74 2D 32 32 32 1F 00 12 49 6E 73 70
				65 63 74 69 6F 6E 20 52 6F 62 6F 74 20 42 00 07
				4D 6F 64 65 6C 2D 58 02 02 03 06 63 61 6D 65 72
				61 05 72 61 64 61 72 09 66 6F 75 72 5F 6C 65 67
				73`,
		},
		{
			name:        "agent info response",
			direction:   acn.Downlink,
			messageType: acn.MessageTypeAgentSearchResponse,
			transaction: 0x25,
			wireHex: `
				7E 00 68 0E 00 4A 01 0A 25 02 00 01 00 42 00 09
				61 67 65 6E 74 2D 31 31 31 1D 00 12 49 6E 73 70
				65 63 74 69 6F 6E 20 52 6F 62 6F 74 20 41 02 01
				03 06 63 61 6D 65 72 61 0A 6D 6F 6E 69 74 6F 72
				69 6E 67 0C 6E 69 67 68 74 5F 76 69 73 69 6F 6E`,
		},
		{
			name:        "owner agents response",
			direction:   acn.Downlink,
			messageType: acn.MessageTypeAgentSearchResponse,
			transaction: 0x26,
			wireHex: `
				7E 00 68 0E 00 5C 01 0A 26 03 00 02 00 29 00 09
				61 67 65 6E 74 2D 31 31 31 03 00 12 49 6E 73 70
				65 63 74 69 6F 6E 20 52 6F 62 6F 74 20 41 00 07
				4D 6F 64 65 6C 2D 58 00 29 00 09 61 67 65 6E 74
				2D 32 32 32 03 00 12 49 6E 73 70 65 63 74 69 6F
				6E 20 52 6F 62 6F 74 20 42 00 07 4D 6F 64 65 6C
				2D 58`,
		},
		{
			name:        "no match response",
			direction:   acn.Downlink,
			messageType: acn.MessageTypeAgentSearchResponse,
			transaction: 0x25,
			wireHex:     `7E 00 68 0E 00 06 01 0A 25 02 00 00`,
		},
		{
			name:        "search reject",
			direction:   acn.Downlink,
			messageType: acn.MessageTypeAgentSearchReject,
			transaction: 0x24,
			wireHex:     `7E 00 68 0E 00 05 01 0B 24 02 05`,
		},
	}
}

func decodeHex(t testing.TB, value string) []byte {
	t.Helper()
	decoded, err := hex.DecodeString(strings.Join(strings.Fields(value), ""))
	require.NoError(t, err)
	return decoded
}

func TestGoldenVectors(t *testing.T) {
	for _, vector := range goldenVectors() {
		vector := vector
		t.Run(vector.name, func(t *testing.T) {
			wire := decodeHex(t, vector.wireHex)

			message, err := acn.DecodePlainNAS(vector.direction, wire)
			require.NoError(t, err)
			require.Equal(t, vector.messageType, message.MessageType())
			require.Equal(t, vector.direction, message.Direction())
			require.Equal(t, vector.transaction, message.GetHeader().TransactionID)

			payload := wire[6:]
			payloadMessage, err := acn.Unmarshal(vector.direction, payload)
			require.NoError(t, err)
			require.Equal(t, reflect.TypeOf(message), reflect.TypeOf(payloadMessage))

			encodedPayload, err := acn.Marshal(message)
			require.NoError(t, err)
			require.Equal(t, payload, encodedPayload)

			encodedWire, err := acn.EncodePlainNAS(message)
			require.NoError(t, err)
			require.Equal(t, wire, encodedWire)

			assertDecodedFields(t, vector.name, message)
		})
	}
}

func assertDecodedFields(t *testing.T, name string, message acn.Message) {
	t.Helper()
	switch name {
	case "register request":
		value := message.(*acn.AgentRegisterRequest)
		require.Equal(t, "u1", value.Owner)
		require.Equal(t, "Alice", value.AgentName)
		require.Equal(t, []byte{1, 2, 3, 4}, value.PublicKey)
		require.Equal(t, "Model-X", value.Description)
		require.Equal(t, uint64(1774617940000), value.Timestamp)
		require.Equal(t, []byte("sig"), value.Signature)
		require.Equal(t, "CN", value.Region)
		require.Equal(t, "Linux", value.OS)
		require.Equal(t, "1.0.0", value.SoftwareVersion)
	case "register accept":
		value := message.(*acn.AgentRegisterAccept)
		require.Equal(t, "a1", value.AgentID)
		require.True(t, json.Valid(value.VC0))
		require.Contains(t, string(value.VC0), `"agent_id":"a1"`)
	case "register reject":
		value := message.(*acn.AgentRegisterReject)
		require.Equal(t, acn.RegisterRejectInvalidSignature, value.Cause)
		require.Equal(t, acn.RegisterFieldSignature, value.FailedField)
	case "deregister request":
		value := message.(*acn.AgentDeregisterRequest)
		require.Equal(t, "a1", value.AgentID)
		require.Equal(t, acn.DeregistrationReasonRetired, value.Reason)
		require.Equal(t, uint64(1711195200000), value.Timestamp)
		require.Equal(t, []byte("sig"), value.Signature)
	case "deregister reject":
		value := message.(*acn.AgentDeregisterReject)
		require.Equal(t, acn.DeregisterRejectInvalidSignature, value.Cause)
		require.Equal(t, acn.DeregisterFieldSignature, value.FailedField)
	case "profile update request":
		value := message.(*acn.AgentProfileUpdateRequest)
		require.Equal(t, "a1", value.AgentID)
		require.Equal(t, acn.PriorityNormal, value.Priority)
		require.Equal(t, uint64(1711195200000), value.Timestamp)
		require.True(t, json.Valid(value.VCList))
		require.True(t, bytes.HasPrefix(value.VCList, []byte("[")))
	case "capability discovery request":
		value := message.(*acn.AgentSearchRequest)
		query := value.Query.(*acn.CapabilityDiscoveryQuery)
		require.Equal(t, "a1", query.SourceAgentID)
		require.Equal(t, "t1", query.TaskID)
		require.Equal(t, []string{"camera", "radar", "four_legs"}, query.RequiredCapabilities)
	case "agent info request":
		value := message.(*acn.AgentSearchRequest)
		require.Equal(t, "agent-111", value.Query.(*acn.AgentInfoQuery).TargetAgentID)
	case "owner agents request":
		value := message.(*acn.AgentSearchRequest)
		require.Equal(t, "owner-123", value.Query.(*acn.OwnerAgentsQuery).OwnerID)
	case "capability discovery response":
		value := message.(*acn.AgentSearchResponse)
		require.Equal(t, acn.SearchTypeCapabilityDiscovery, value.SearchType)
		require.Len(t, value.Records, 1)
		require.Equal(t, []string{"camera", "radar", "four_legs"}, value.Records[0].Capabilities)
	case "agent info response":
		value := message.(*acn.AgentSearchResponse)
		require.Equal(t, acn.SearchTypeAgentInfo, value.SearchType)
		require.Len(t, value.Records, 1)
		require.Equal(t, acn.AgentStatusOnline, *value.Records[0].Status)
		require.Equal(t, acn.PriorityHigh, *value.Records[0].Priority)
	case "owner agents response":
		value := message.(*acn.AgentSearchResponse)
		require.Equal(t, acn.SearchTypeOwnerAgents, value.SearchType)
		require.Len(t, value.Records, 2)
		require.Equal(t, "agent-222", value.Records[1].AgentID)
	case "no match response":
		value := message.(*acn.AgentSearchResponse)
		require.Equal(t, acn.SearchTypeAgentInfo, value.SearchType)
		require.Empty(t, value.Records)
	case "search reject":
		value := message.(*acn.AgentSearchReject)
		require.Equal(t, acn.SearchRejectInvalidMandatory, value.Cause)
		require.Equal(t, acn.SearchFieldRequiredCapability, value.FailedField)
	}
}
