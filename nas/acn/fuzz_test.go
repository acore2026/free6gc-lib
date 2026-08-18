package acn_test

import (
	"testing"

	"github.com/acore2026/free6gc-lib/nas/acn"
)

func FuzzUnmarshalACN(f *testing.F) {
	for _, vector := range goldenVectors() {
		wire := decodeHex(f, vector.wireHex)
		f.Add(uint8(vector.direction), wire[6:])
	}

	f.Fuzz(func(t *testing.T, rawDirection uint8, payload []byte) {
		direction := acn.Direction(rawDirection)
		message, err := acn.Unmarshal(direction, payload)
		if err != nil {
			return
		}

		encoded, err := acn.Marshal(message)
		if err != nil {
			t.Fatalf("successfully decoded message could not be encoded: %v", err)
		}
		if _, err := acn.Unmarshal(message.Direction(), encoded); err != nil {
			t.Fatalf("encoded message could not be decoded: %v", err)
		}
	})
}

func FuzzDecodePlainACNNAS(f *testing.F) {
	for _, vector := range goldenVectors() {
		f.Add(uint8(vector.direction), decodeHex(f, vector.wireHex))
	}

	f.Fuzz(func(t *testing.T, rawDirection uint8, wire []byte) {
		direction := acn.Direction(rawDirection)
		message, err := acn.DecodePlainNAS(direction, wire)
		if err != nil {
			return
		}

		encoded, err := acn.EncodePlainNAS(message)
		if err != nil {
			t.Fatalf("successfully decoded message could not be wrapped: %v", err)
		}
		if _, err := acn.DecodePlainNAS(message.Direction(), encoded); err != nil {
			t.Fatalf("wrapped message could not be decoded: %v", err)
		}
	})
}
