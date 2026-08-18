package acn

import (
	"encoding/binary"
	"fmt"

	naslib "github.com/acore2026/free6gc-lib/nas"
	"github.com/acore2026/free6gc-lib/nas/nasMessage"
)

const nasEnvelopeLength = 6

// EncodePlainNAS wraps an ACN PDU in a plain UL or DL NAS TRANSPORT message
// using default limits.
func EncodePlainNAS(message Message) ([]byte, error) {
	return defaultCodec.EncodePlainNAS(message)
}

// EncodePlainNAS wraps an ACN PDU in a plain UL or DL NAS TRANSPORT message.
func (c *Codec) EncodePlainNAS(message Message) ([]byte, error) {
	payload, err := c.Marshal(message)
	if err != nil {
		return nil, err
	}

	outer := naslib.NewMessage()
	outer.GmmMessage = naslib.NewGmmMessage()

	switch message.Direction() {
	case Uplink:
		outer.GmmHeader.SetMessageType(naslib.MsgTypeULNASTransport)
		transport := nasMessage.NewULNASTransport(0)
		transport.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(
			nasMessage.Epd5GSMobilityManagementMessage,
		)
		transport.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(
			naslib.SecurityHeaderTypePlainNas,
		)
		transport.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0)
		transport.ULNASTRANSPORTMessageIdentity.SetMessageType(naslib.MsgTypeULNASTransport)
		transport.SpareHalfOctetAndPayloadContainerType.SetPayloadContainerType(
			PayloadContainerTypeACN,
		)
		transport.PayloadContainer.SetLen(uint16(len(payload)))
		transport.PayloadContainer.SetPayloadContainerContents(payload)
		outer.GmmMessage.ULNASTransport = transport
	case Downlink:
		outer.GmmHeader.SetMessageType(naslib.MsgTypeDLNASTransport)
		transport := nasMessage.NewDLNASTransport(0)
		transport.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(
			nasMessage.Epd5GSMobilityManagementMessage,
		)
		transport.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(
			naslib.SecurityHeaderTypePlainNas,
		)
		transport.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0)
		transport.DLNASTRANSPORTMessageIdentity.SetMessageType(naslib.MsgTypeDLNASTransport)
		transport.SpareHalfOctetAndPayloadContainerType.SetPayloadContainerType(
			PayloadContainerTypeACN,
		)
		transport.PayloadContainer.SetLen(uint16(len(payload)))
		transport.PayloadContainer.SetPayloadContainerContents(payload)
		outer.GmmMessage.DLNASTransport = transport
	default:
		return nil, newProtocolError(ErrorCodeInvalidValue, "direction", -1, ErrInvalidValue)
	}

	wire, err := outer.PlainNasEncode()
	if err != nil {
		return nil, newProtocolError(ErrorCodeInvalidNASEnvelope, "nas", -1, err)
	}
	if len(wire) != nasEnvelopeLength+len(payload) {
		return nil, newProtocolError(
			ErrorCodeInvalidNASEnvelope,
			"nas.length",
			-1,
			fmt.Errorf("%w: encoded %d bytes, want %d",
				ErrInvalidNASEnvelope, len(wire), nasEnvelopeLength+len(payload)),
		)
	}
	return wire, nil
}

// DecodePlainNAS extracts and decodes an ACN PDU from a plain UL or DL NAS
// TRANSPORT message using default limits.
func DecodePlainNAS(direction Direction, wire []byte) (Message, error) {
	return defaultCodec.DecodePlainNAS(direction, wire)
}

// DecodePlainNAS extracts and decodes an ACN PDU from a plain UL or DL NAS
// TRANSPORT message.
func (c *Codec) DecodePlainNAS(direction Direction, wire []byte) (Message, error) {
	if c == nil {
		return nil, newProtocolError(ErrorCodeInvalidValue, "codec", -1, ErrInvalidValue)
	}
	payload, err := c.extractPlainNASPayload(direction, wire)
	if err != nil {
		return nil, err
	}
	return c.Unmarshal(direction, payload)
}

func (c *Codec) extractPlainNASPayload(direction Direction, wire []byte) ([]byte, error) {
	if direction != Uplink && direction != Downlink {
		return nil, newProtocolError(ErrorCodeInvalidValue, "direction", -1, ErrInvalidValue)
	}
	if len(wire) < nasEnvelopeLength {
		return nil, newProtocolError(
			ErrorCodeInvalidNASEnvelope,
			"nas",
			0,
			fmt.Errorf("%w: %w", ErrInvalidNASEnvelope, ErrTruncated),
		)
	}
	if wire[0] != nasMessage.Epd5GSMobilityManagementMessage {
		return nil, newProtocolError(
			ErrorCodeInvalidNASEnvelope,
			"nas.epd",
			0,
			ErrInvalidNASEnvelope,
		)
	}
	if wire[1] != naslib.SecurityHeaderTypePlainNas {
		cause := ErrInvalidNASEnvelope
		if wire[1]&0x0f != naslib.SecurityHeaderTypePlainNas {
			cause = ErrProtectedNAS
		}
		return nil, newProtocolError(ErrorCodeInvalidNASEnvelope, "nas.security_header", 1, cause)
	}
	expectedMessageType := uint8(naslib.MsgTypeULNASTransport)
	if direction == Downlink {
		expectedMessageType = naslib.MsgTypeDLNASTransport
	}
	if wire[2] != expectedMessageType {
		return nil, newProtocolError(
			ErrorCodeWrongDirection,
			"nas.message_type",
			2,
			ErrWrongDirection,
		)
	}
	if wire[3] != PayloadContainerTypeACN {
		return nil, newProtocolError(
			ErrorCodeInvalidNASEnvelope,
			"nas.payload_container_type",
			3,
			ErrInvalidNASEnvelope,
		)
	}
	declaredLength := int(binary.BigEndian.Uint16(wire[4:6]))
	if declaredLength < 3 {
		return nil, newProtocolError(
			ErrorCodeInvalidLength,
			"nas.payload_length",
			4,
			ErrInvalidLength,
		)
	}
	if declaredLength > c.limits.MaxACNPayload {
		return nil, newProtocolError(
			ErrorCodePayloadTooLarge,
			"nas.payload_length",
			4,
			ErrPayloadTooLarge,
		)
	}
	if len(wire) != nasEnvelopeLength+declaredLength {
		return nil, newProtocolError(
			ErrorCodeInvalidLength,
			"nas.payload_length",
			4,
			fmt.Errorf("%w: declared %d, wire carries %d",
				ErrInvalidLength, declaredLength, len(wire)-nasEnvelopeLength),
		)
	}

	// Decode through the existing NAS implementation after the strict raw-wire
	// checks so future changes to the transport structures remain exercised.
	decodedWire := cloneBytes(wire)
	outer := naslib.NewMessage()
	if err := outer.PlainNasDecode(&decodedWire); err != nil {
		return nil, newProtocolError(ErrorCodeInvalidNASEnvelope, "nas", 0, err)
	}

	var payload []byte
	switch direction {
	case Uplink:
		if outer.GmmMessage == nil || outer.GmmMessage.ULNASTransport == nil {
			return nil, newProtocolError(
				ErrorCodeInvalidNASEnvelope,
				"nas.uplink_transport",
				2,
				ErrInvalidNASEnvelope,
			)
		}
		transport := outer.GmmMessage.ULNASTransport
		if transport.SpareHalfOctetAndPayloadContainerType.GetPayloadContainerType() !=
			PayloadContainerTypeACN {
			return nil, newProtocolError(
				ErrorCodeInvalidNASEnvelope,
				"nas.payload_container_type",
				3,
				ErrInvalidNASEnvelope,
			)
		}
		payload = transport.PayloadContainer.Buffer
	case Downlink:
		if outer.GmmMessage == nil || outer.GmmMessage.DLNASTransport == nil {
			return nil, newProtocolError(
				ErrorCodeInvalidNASEnvelope,
				"nas.downlink_transport",
				2,
				ErrInvalidNASEnvelope,
			)
		}
		transport := outer.GmmMessage.DLNASTransport
		if transport.SpareHalfOctetAndPayloadContainerType.GetPayloadContainerType() !=
			PayloadContainerTypeACN {
			return nil, newProtocolError(
				ErrorCodeInvalidNASEnvelope,
				"nas.payload_container_type",
				3,
				ErrInvalidNASEnvelope,
			)
		}
		payload = transport.PayloadContainer.Buffer
	}
	if len(payload) != declaredLength {
		return nil, newProtocolError(
			ErrorCodeInvalidNASEnvelope,
			"nas.payload_length",
			4,
			ErrInvalidNASEnvelope,
		)
	}
	return cloneBytes(payload), nil
}
