package acn_test

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/acore2026/free6gc-lib/nas/acn"
)

func TestMarshalRejectsInvalidMessages(t *testing.T) {
	var typedNil *acn.AgentDeregisterAccept
	invalidUTF8 := string([]byte{0xff})
	emptyCapabilities := []string{}

	tests := []struct {
		name    string
		message acn.Message
		target  error
	}{
		{
			name:   "nil message",
			target: acn.ErrInvalidValue,
		},
		{
			name:    "typed nil message",
			message: typedNil,
			target:  acn.ErrInvalidValue,
		},
		{
			name: "zero transaction ID",
			message: &acn.AgentDeregisterAccept{
				Header: acn.Header{},
			},
			target: acn.ErrInvalidValue,
		},
		{
			name: "empty mandatory string",
			message: &acn.AgentDeregisterRequest{
				Header:    acn.Header{TransactionID: 1},
				Reason:    acn.DeregistrationReasonNormal,
				Signature: []byte("sig"),
			},
			target: acn.ErrInvalidLength,
		},
		{
			name: "invalid UTF-8",
			message: &acn.AgentSearchRequest{
				Header: acn.Header{TransactionID: 1},
				Query:  &acn.AgentInfoQuery{TargetAgentID: invalidUTF8},
			},
			target: acn.ErrInvalidUTF8,
		},
		{
			name: "reserved deregistration reason",
			message: &acn.AgentDeregisterRequest{
				Header:    acn.Header{TransactionID: 1},
				AgentID:   "a1",
				Reason:    acn.DeregistrationReason(6),
				Signature: []byte("sig"),
			},
			target: acn.ErrInvalidValue,
		},
		{
			name: "reserved reject cause",
			message: &acn.AgentRegisterReject{
				Header:      acn.Header{TransactionID: 1},
				Cause:       acn.RegisterRejectCause(0xff),
				FailedField: acn.RegisterFieldUnspecified,
			},
			target: acn.ErrInvalidValue,
		},
		{
			name: "nil search query",
			message: &acn.AgentSearchRequest{
				Header: acn.Header{TransactionID: 1},
			},
			target: acn.ErrInvalidValue,
		},
		{
			name: "empty required capability list",
			message: &acn.AgentSearchRequest{
				Header: acn.Header{TransactionID: 1},
				Query: &acn.CapabilityDiscoveryQuery{
					SourceAgentID:        "a1",
					TaskID:               "task",
					RequiredCapabilities: emptyCapabilities,
				},
			},
			target: acn.ErrInvalidLength,
		},
		{
			name: "agent info has multiple records",
			message: &acn.AgentSearchResponse{
				Header:     acn.Header{TransactionID: 1},
				SearchType: acn.SearchTypeAgentInfo,
				Records: []acn.AgentRecord{
					{AgentID: "a1"},
					{AgentID: "a2"},
				},
			},
			target: acn.ErrCountMismatch,
		},
		{
			name: "present capabilities are empty",
			message: &acn.AgentSearchResponse{
				Header:     acn.Header{TransactionID: 1},
				SearchType: acn.SearchTypeOwnerAgents,
				Records: []acn.AgentRecord{{
					AgentID:      "a1",
					Capabilities: emptyCapabilities,
				}},
			},
			target: acn.ErrInvalidLength,
		},
		{
			name: "VC0 is not JSON",
			message: &acn.AgentRegisterAccept{
				Header:  acn.Header{TransactionID: 1},
				AgentID: "a1",
				VC0:     json.RawMessage(`{`),
			},
			target: acn.ErrInvalidJSON,
		},
		{
			name: "VC0 contains invalid UTF-8",
			message: &acn.AgentRegisterAccept{
				Header:  acn.Header{TransactionID: 1},
				AgentID: "a1",
				VC0: json.RawMessage{
					'{', '"', 'x', '"', ':', '"', 0xff, '"', '}',
				},
			},
			target: acn.ErrInvalidUTF8,
		},
		{
			name: "VC0 is not an object",
			message: &acn.AgentRegisterAccept{
				Header:  acn.Header{TransactionID: 1},
				AgentID: "a1",
				VC0:     json.RawMessage(`[1]`),
			},
			target: acn.ErrInvalidJSON,
		},
		{
			name: "VC0 object is empty",
			message: &acn.AgentRegisterAccept{
				Header:  acn.Header{TransactionID: 1},
				AgentID: "a1",
				VC0:     json.RawMessage(`{}`),
			},
			target: acn.ErrInvalidJSON,
		},
		{
			name: "VC list is not an array",
			message: &acn.AgentProfileUpdateRequest{
				Header:    acn.Header{TransactionID: 1},
				AgentID:   "a1",
				Priority:  acn.PriorityNormal,
				Signature: []byte("sig"),
				VCList:    json.RawMessage(`{"id":"vc1"}`),
			},
			target: acn.ErrInvalidJSON,
		},
		{
			name: "VC list is empty",
			message: &acn.AgentProfileUpdateRequest{
				Header:    acn.Header{TransactionID: 1},
				AgentID:   "a1",
				Priority:  acn.PriorityNormal,
				Signature: []byte("sig"),
				VCList:    json.RawMessage(`[]`),
			},
			target: acn.ErrInvalidJSON,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			_, err := acn.Marshal(test.message)
			require.ErrorIs(t, err, test.target)

			var protocolErr *acn.ProtocolError
			require.ErrorAs(t, err, &protocolErr)
			require.NotEmpty(t, protocolErr.Code)
			require.NotEmpty(t, protocolErr.Field)
			if test.name != "nil message" &&
				test.name != "typed nil message" &&
				test.message.GetHeader().TransactionID != 0 {
				require.NotNil(t, protocolErr.Header)
				require.Equal(t, test.message.GetHeader(), *protocolErr.Header)
			}
		})
	}
}

func TestUnmarshalRejectsInvalidPayloads(t *testing.T) {
	tests := []struct {
		name      string
		direction acn.Direction
		payload   []byte
		target    error
	}{
		{
			name:      "empty",
			direction: acn.Uplink,
			target:    acn.ErrTruncated,
		},
		{
			name:      "unsupported version",
			direction: acn.Uplink,
			payload:   []byte{2, byte(acn.MessageTypeAgentDeregisterRequest), 1},
			target:    acn.ErrUnsupportedVersion,
		},
		{
			name:      "unknown message type",
			direction: acn.Uplink,
			payload:   []byte{acn.Version1, 0xff, 1},
			target:    acn.ErrUnknownMessageType,
		},
		{
			name:      "zero transaction ID",
			direction: acn.Downlink,
			payload:   []byte{acn.Version1, byte(acn.MessageTypeAgentDeregisterAccept), 0},
			target:    acn.ErrInvalidValue,
		},
		{
			name:      "wrong direction",
			direction: acn.Uplink,
			payload:   []byte{acn.Version1, byte(acn.MessageTypeAgentDeregisterAccept), 1},
			target:    acn.ErrWrongDirection,
		},
		{
			name:      "trailing data",
			direction: acn.Downlink,
			payload:   []byte{acn.Version1, byte(acn.MessageTypeAgentDeregisterAccept), 1, 0},
			target:    acn.ErrTrailingData,
		},
		{
			name:      "reserved enum",
			direction: acn.Downlink,
			payload: []byte{
				acn.Version1, byte(acn.MessageTypeAgentRegisterReject), 1, 0xff, 0,
			},
			target: acn.ErrInvalidValue,
		},
		{
			name:      "unknown request search type",
			direction: acn.Uplink,
			payload: []byte{
				acn.Version1, byte(acn.MessageTypeAgentSearchRequest), 1, 0xff,
			},
			target: acn.ErrInvalidValue,
		},
		{
			name:      "invalid UTF-8",
			direction: acn.Uplink,
			payload: []byte{
				acn.Version1, byte(acn.MessageTypeAgentSearchRequest), 1,
				byte(acn.SearchTypeAgentInfo), 0, 1, 0xff,
			},
			target: acn.ErrInvalidUTF8,
		},
		{
			name:      "zero capability count",
			direction: acn.Uplink,
			payload: []byte{
				acn.Version1, byte(acn.MessageTypeAgentSearchRequest), 1,
				byte(acn.SearchTypeCapabilityDiscovery),
				0, 1, 'a',
				0, 1, 't',
				0, 0, 0, 0, 0, 0, 0, 0,
				0,
			},
			target: acn.ErrInvalidLength,
		},
		{
			name:      "result count exceeds available records",
			direction: acn.Downlink,
			payload: []byte{
				acn.Version1, byte(acn.MessageTypeAgentSearchResponse), 1,
				byte(acn.SearchTypeOwnerAgents), 0, 1,
			},
			target: acn.ErrCountMismatch,
		},
		{
			name:      "record is shorter than minimum",
			direction: acn.Downlink,
			payload: []byte{
				acn.Version1, byte(acn.MessageTypeAgentSearchResponse), 1,
				byte(acn.SearchTypeOwnerAgents), 0, 1,
				0, 3, 0, 1, 'a',
				0, 0, 0,
			},
			target: acn.ErrInvalidLength,
		},
		{
			name:      "record length overruns message",
			direction: acn.Downlink,
			payload: []byte{
				acn.Version1, byte(acn.MessageTypeAgentSearchResponse), 1,
				byte(acn.SearchTypeOwnerAgents), 0, 1,
				0, 7, 0, 1, 'a', 0,
				0, 0,
			},
			target: acn.ErrInvalidLength,
		},
		{
			name:      "record has trailing data",
			direction: acn.Downlink,
			payload: []byte{
				acn.Version1, byte(acn.MessageTypeAgentSearchResponse), 1,
				byte(acn.SearchTypeOwnerAgents), 0, 1,
				0, 5, 0, 1, 'a', 0, 0,
			},
			target: acn.ErrTrailingData,
		},
		{
			name:      "agent info multiple records",
			direction: acn.Downlink,
			payload: []byte{
				acn.Version1, byte(acn.MessageTypeAgentSearchResponse), 1,
				byte(acn.SearchTypeAgentInfo), 0, 2,
			},
			target: acn.ErrCountMismatch,
		},
		{
			name:      "reserved record presence bits",
			direction: acn.Downlink,
			payload: []byte{
				acn.Version1, byte(acn.MessageTypeAgentSearchResponse), 1,
				byte(acn.SearchTypeOwnerAgents), 0, 1,
				0, 4,
				0, 1, 'a', 0x80,
			},
			target: acn.ErrInvalidPresenceBitmap,
		},
		{
			name:      "present record capabilities have zero count",
			direction: acn.Downlink,
			payload: []byte{
				acn.Version1, byte(acn.MessageTypeAgentSearchResponse), 1,
				byte(acn.SearchTypeOwnerAgents), 0, 1,
				0, 5,
				0, 1, 'a', 0x10, 0,
			},
			target: acn.ErrInvalidLength,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			_, err := acn.Unmarshal(test.direction, test.payload)
			require.ErrorIs(t, err, test.target)
			var protocolErr *acn.ProtocolError
			require.ErrorAs(t, err, &protocolErr)
		})
	}
}

func TestDecodePlainNASRejectsInvalidEnvelope(t *testing.T) {
	base := decodeHex(t, `7E 00 68 0E 00 03 01 05 22`)

	tests := []struct {
		name      string
		direction acn.Direction
		mutate    func([]byte) []byte
		target    error
	}{
		{
			name:      "truncated envelope",
			direction: acn.Downlink,
			mutate:    func(wire []byte) []byte { return wire[:5] },
			target:    acn.ErrInvalidNASEnvelope,
		},
		{
			name:      "wrong EPD",
			direction: acn.Downlink,
			mutate:    func(wire []byte) []byte { wire[0] = 0x2e; return wire },
			target:    acn.ErrInvalidNASEnvelope,
		},
		{
			name:      "protected NAS",
			direction: acn.Downlink,
			mutate:    func(wire []byte) []byte { wire[1] = 1; return wire },
			target:    acn.ErrProtectedNAS,
		},
		{
			name:      "nonzero spare half octet",
			direction: acn.Downlink,
			mutate:    func(wire []byte) []byte { wire[1] = 0xf0; return wire },
			target:    acn.ErrInvalidNASEnvelope,
		},
		{
			name:      "wrong transport direction",
			direction: acn.Uplink,
			mutate:    func(wire []byte) []byte { return wire },
			target:    acn.ErrWrongDirection,
		},
		{
			name:      "wrong payload container type",
			direction: acn.Downlink,
			mutate:    func(wire []byte) []byte { wire[3] = 1; return wire },
			target:    acn.ErrInvalidNASEnvelope,
		},
		{
			name:      "declared payload too short",
			direction: acn.Downlink,
			mutate:    func(wire []byte) []byte { wire[5] = 2; return wire },
			target:    acn.ErrInvalidLength,
		},
		{
			name:      "under-declared payload",
			direction: acn.Downlink,
			mutate:    func(wire []byte) []byte { wire[5] = 3; return append(wire, 0) },
			target:    acn.ErrInvalidLength,
		},
		{
			name:      "over-declared payload",
			direction: acn.Downlink,
			mutate:    func(wire []byte) []byte { wire[5] = 4; return wire },
			target:    acn.ErrInvalidLength,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			wire := append([]byte(nil), base...)
			_, err := acn.DecodePlainNAS(test.direction, test.mutate(wire))
			require.ErrorIs(t, err, test.target)
			var protocolErr *acn.ProtocolError
			require.ErrorAs(t, err, &protocolErr)
		})
	}
}

func TestEveryGoldenVectorRejectsEveryTruncatedPrefix(t *testing.T) {
	for _, vector := range goldenVectors() {
		vector := vector
		t.Run(vector.name, func(t *testing.T) {
			wire := decodeHex(t, vector.wireHex)
			payload := wire[6:]

			for length := 0; length < len(payload); length++ {
				_, err := acn.Unmarshal(vector.direction, payload[:length])
				require.Error(t, err, "payload prefix length %d was accepted", length)
			}
			for length := 0; length < len(wire); length++ {
				_, err := acn.DecodePlainNAS(vector.direction, wire[:length])
				require.Error(t, err, "NAS prefix length %d was accepted", length)
			}
		})
	}
}

func TestCodecLimits(t *testing.T) {
	defaults := acn.DefaultLimits()
	require.Equal(t, math.MaxUint16, defaults.MaxACNPayload)
	require.Equal(t, math.MaxUint8, defaults.MaxCapability)

	invalid := defaults
	invalid.MaxAgentID = 0
	_, err := acn.NewCodec(invalid)
	require.Error(t, err)

	invalid = defaults
	invalid.MaxJSONContainer = invalid.MaxACNPayload + 1
	_, err = acn.NewCodec(invalid)
	require.Error(t, err)

	limits := defaults
	limits.MaxAgentID = 2
	codec, err := acn.NewCodec(limits)
	require.NoError(t, err)
	require.Equal(t, limits, codec.Limits())

	message := &acn.AgentSearchRequest{
		Header: acn.Header{TransactionID: 1},
		Query:  &acn.AgentInfoQuery{TargetAgentID: "abc"},
	}
	_, err = codec.Marshal(message)
	require.ErrorIs(t, err, acn.ErrInvalidLength)

	payload := []byte{
		acn.Version1, byte(acn.MessageTypeAgentSearchRequest), 1,
		byte(acn.SearchTypeAgentInfo), 0, 3, 'a', 'b', 'c',
	}
	_, err = codec.Unmarshal(acn.Uplink, payload)
	require.ErrorIs(t, err, acn.ErrInvalidLength)

	var nilCodec *acn.Codec
	_, err = nilCodec.Marshal(&acn.AgentDeregisterAccept{
		Header: acn.Header{TransactionID: 1},
	})
	require.ErrorIs(t, err, acn.ErrInvalidValue)
	_, err = nilCodec.Unmarshal(acn.Downlink, []byte{1, 5, 1})
	require.ErrorIs(t, err, acn.ErrInvalidValue)
	_, err = nilCodec.DecodePlainNAS(acn.Downlink, nil)
	require.ErrorIs(t, err, acn.ErrInvalidValue)
}

func TestCodecPayloadLimitBoundary(t *testing.T) {
	limits := acn.DefaultLimits()
	limits.MaxACNPayload = 20
	limits.MaxJSONContainer = 20
	codec, err := acn.NewCodec(limits)
	require.NoError(t, err)

	exact := &acn.AgentSearchRequest{
		Header: acn.Header{TransactionID: 1},
		Query:  &acn.AgentInfoQuery{TargetAgentID: "12345678901234"},
	}
	payload, err := codec.Marshal(exact)
	require.NoError(t, err)
	require.Len(t, payload, limits.MaxACNPayload)
	_, err = codec.Unmarshal(acn.Uplink, payload)
	require.NoError(t, err)

	over := &acn.AgentSearchRequest{
		Header: acn.Header{TransactionID: 1},
		Query:  &acn.AgentInfoQuery{TargetAgentID: "123456789012345"},
	}
	_, err = codec.Marshal(over)
	require.ErrorIs(t, err, acn.ErrPayloadTooLarge)

	overPayload := append(append([]byte(nil), payload...), 0)
	_, err = codec.Unmarshal(acn.Uplink, overPayload)
	require.ErrorIs(t, err, acn.ErrPayloadTooLarge)

	wire := decodeHex(t, goldenVectors()[0].wireHex)
	_, err = codec.DecodePlainNAS(acn.Uplink, wire)
	require.ErrorIs(t, err, acn.ErrPayloadTooLarge)
}

func TestSearchResponseStopsAtPayloadLimit(t *testing.T) {
	limits := acn.DefaultLimits()
	limits.MaxACNPayload = 12
	limits.MaxJSONContainer = 12
	codec, err := acn.NewCodec(limits)
	require.NoError(t, err)

	message := &acn.AgentSearchResponse{
		Header:     acn.Header{TransactionID: 1},
		SearchType: acn.SearchTypeOwnerAgents,
		Records: []acn.AgentRecord{
			{AgentID: "a"},
			{AgentID: ""},
		},
	}
	_, err = codec.Marshal(message)
	require.ErrorIs(t, err, acn.ErrPayloadTooLarge)
}

func TestDecodedAndEncodedBuffersDoNotAliasInputs(t *testing.T) {
	registerWire := decodeHex(t, goldenVectors()[0].wireHex)
	registerMessage, err := acn.DecodePlainNAS(acn.Uplink, registerWire)
	require.NoError(t, err)
	register := registerMessage.(*acn.AgentRegisterRequest)
	publicKey := append([]byte(nil), register.PublicKey...)
	signature := append([]byte(nil), register.Signature...)

	for index := range registerWire {
		registerWire[index] = 0
	}
	require.Equal(t, publicKey, register.PublicKey)
	require.Equal(t, signature, register.Signature)

	acceptWire := decodeHex(t, goldenVectors()[1].wireHex)
	acceptMessage, err := acn.DecodePlainNAS(acn.Downlink, acceptWire)
	require.NoError(t, err)
	accept := acceptMessage.(*acn.AgentRegisterAccept)
	vc0 := append([]byte(nil), accept.VC0...)
	for index := range acceptWire {
		acceptWire[index] = 0
	}
	require.Equal(t, vc0, []byte(accept.VC0))

	sourceVC0 := json.RawMessage(`{"id":"vc1"}`)
	source := &acn.AgentRegisterAccept{
		Header:  acn.Header{TransactionID: 1},
		AgentID: "a1",
		VC0:     sourceVC0,
	}
	encoded, err := acn.Marshal(source)
	require.NoError(t, err)
	sourceVC0[2] = 'X'
	require.Contains(t, string(encoded), `"id":"vc1"`)

	encodedCopy := append([]byte(nil), encoded...)
	for index := range encoded {
		encoded[index] = 0
	}
	require.NotEqual(t, encodedCopy, encoded)
	require.Equal(t, byte('{'), source.VC0[0])
	require.Equal(t, byte('X'), source.VC0[2])
}
