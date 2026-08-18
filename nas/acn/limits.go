package acn

import (
	"fmt"
	"math"
)

// Limits bounds allocations and field sizes accepted by a Codec.
type Limits struct {
	MaxACNPayload      int
	MaxAgentID         int
	MaxGroupID         int
	MaxOwnerID         int
	MaxTaskID          int
	MaxAgentName       int
	MaxDescription     int
	MaxCapability      int
	MaxPublicKey       int
	MaxSignature       int
	MaxRegion          int
	MaxOS              int
	MaxSoftwareVersion int
	MaxJSONContainer   int
}

// DefaultLimits returns the ACN version 1 default limits.
func DefaultLimits() Limits {
	return Limits{
		MaxACNPayload:      math.MaxUint16,
		MaxAgentID:         1024,
		MaxGroupID:         1024,
		MaxOwnerID:         1024,
		MaxTaskID:          1024,
		MaxAgentName:       255,
		MaxDescription:     2048,
		MaxCapability:      math.MaxUint8,
		MaxPublicKey:       8192,
		MaxSignature:       8192,
		MaxRegion:          64,
		MaxOS:              128,
		MaxSoftwareVersion: 64,
		MaxJSONContainer:   math.MaxUint16,
	}
}

func (l Limits) validate() error {
	if err := validateLimit("MaxACNPayload", l.MaxACNPayload, 3, math.MaxUint16); err != nil {
		return err
	}
	for name, value := range map[string]int{
		"MaxAgentID":         l.MaxAgentID,
		"MaxGroupID":         l.MaxGroupID,
		"MaxOwnerID":         l.MaxOwnerID,
		"MaxTaskID":          l.MaxTaskID,
		"MaxAgentName":       l.MaxAgentName,
		"MaxDescription":     l.MaxDescription,
		"MaxPublicKey":       l.MaxPublicKey,
		"MaxSignature":       l.MaxSignature,
		"MaxRegion":          l.MaxRegion,
		"MaxOS":              l.MaxOS,
		"MaxSoftwareVersion": l.MaxSoftwareVersion,
		"MaxJSONContainer":   l.MaxJSONContainer,
	} {
		if err := validateLimit(name, value, 1, math.MaxUint16); err != nil {
			return err
		}
	}
	if err := validateLimit("MaxCapability", l.MaxCapability, 1, math.MaxUint8); err != nil {
		return err
	}
	if l.MaxJSONContainer > l.MaxACNPayload {
		return fmt.Errorf("MaxJSONContainer %d exceeds MaxACNPayload %d",
			l.MaxJSONContainer, l.MaxACNPayload)
	}
	return nil
}

func validateLimit(name string, value, minimum, maximum int) error {
	if value < minimum || value > maximum {
		return fmt.Errorf("%s is %d, want %d..%d", name, value, minimum, maximum)
	}
	return nil
}

// Codec encodes and decodes ACN messages with immutable size limits.
type Codec struct {
	limits Limits
}

// NewCodec constructs a codec with explicit limits.
func NewCodec(limits Limits) (*Codec, error) {
	if err := limits.validate(); err != nil {
		return nil, err
	}
	return &Codec{limits: limits}, nil
}

// Limits returns a copy of the codec configuration.
func (c *Codec) Limits() Limits {
	if c == nil {
		return Limits{}
	}
	return c.limits
}

var defaultCodec = func() *Codec {
	codec, err := NewCodec(DefaultLimits())
	if err != nil {
		panic(err)
	}
	return codec
}()
