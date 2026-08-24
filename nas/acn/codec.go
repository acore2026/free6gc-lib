package acn

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"reflect"
	"unicode/utf8"
)

type encoder struct {
	data []byte
}

func (e *encoder) uint8(value uint8) {
	e.data = append(e.data, value)
}

func (e *encoder) uint16(value uint16) {
	e.data = binary.BigEndian.AppendUint16(e.data, value)
}

func (e *encoder) uint64(value uint64) {
	e.data = binary.BigEndian.AppendUint64(e.data, value)
}

func (e *encoder) bytes(value []byte) {
	e.data = append(e.data, value...)
}

func (e *encoder) lv(value []byte) {
	e.uint8(uint8(len(value)))
	e.bytes(value)
}

func (e *encoder) lve(value []byte) {
	e.uint16(uint16(len(value)))
	e.bytes(value)
}

type decoder struct {
	data   []byte
	offset int
}

func (d *decoder) remaining() int {
	return len(d.data) - d.offset
}

func (d *decoder) uint8(field string) (uint8, error) {
	if d.remaining() < 1 {
		return 0, newProtocolError(ErrorCodeTruncated, field, d.offset, ErrTruncated)
	}
	value := d.data[d.offset]
	d.offset++
	return value, nil
}

func (d *decoder) uint16(field string) (uint16, error) {
	if d.remaining() < 2 {
		return 0, newProtocolError(ErrorCodeTruncated, field, d.offset, ErrTruncated)
	}
	value := binary.BigEndian.Uint16(d.data[d.offset : d.offset+2])
	d.offset += 2
	return value, nil
}

func (d *decoder) uint64(field string) (uint64, error) {
	if d.remaining() < 8 {
		return 0, newProtocolError(ErrorCodeTruncated, field, d.offset, ErrTruncated)
	}
	value := binary.BigEndian.Uint64(d.data[d.offset : d.offset+8])
	d.offset += 8
	return value, nil
}

func (d *decoder) bytes(field string, length int) ([]byte, error) {
	if length < 0 || length > d.remaining() {
		return nil, newProtocolError(ErrorCodeInvalidLength, field, d.offset, ErrInvalidLength)
	}
	value := d.data[d.offset : d.offset+length]
	d.offset += length
	return value, nil
}

func (d *decoder) lvBytes(field string, maximum int) ([]byte, error) {
	length, err := d.uint8(field + ".length")
	if err != nil {
		return nil, err
	}
	return d.variableBytes(field, int(length), maximum)
}

func (d *decoder) lveBytes(field string, maximum int) ([]byte, error) {
	length, err := d.uint16(field + ".length")
	if err != nil {
		return nil, err
	}
	return d.variableBytes(field, int(length), maximum)
}

func (d *decoder) lveOptionalBytes(field string, maximum int) ([]byte, error) {
	length, err := d.uint16(field + ".length")
	if err != nil {
		return nil, err
	}
	if length == 0 {
		return nil, nil
	}
	return d.variableBytes(field, int(length), maximum)
}

func (d *decoder) variableBytes(field string, length, maximum int) ([]byte, error) {
	if length < 1 || length > maximum {
		return nil, newProtocolError(
			ErrorCodeInvalidLength,
			field,
			d.offset,
			fmt.Errorf("%w: got %d, want 1..%d", ErrInvalidLength, length, maximum),
		)
	}
	return d.bytes(field, length)
}

func (d *decoder) lvString(field string, maximum int) (string, error) {
	value, err := d.lvBytes(field, maximum)
	if err != nil {
		return "", err
	}
	return decodeUTF8(value, field, d.offset-len(value))
}

func (d *decoder) lveString(field string, maximum int) (string, error) {
	value, err := d.lveBytes(field, maximum)
	if err != nil {
		return "", err
	}
	return decodeUTF8(value, field, d.offset-len(value))
}

func (d *decoder) lveOptionalString(field string, maximum int) (string, error) {
	value, err := d.lveOptionalBytes(field, maximum)
	if err != nil || value == nil {
		return "", err
	}
	return decodeUTF8(value, field, d.offset-len(value))
}

func (d *decoder) finish(field string) error {
	if d.remaining() != 0 {
		return newProtocolError(
			ErrorCodeTrailingData,
			field,
			d.offset,
			fmt.Errorf("%w: %d bytes", ErrTrailingData, d.remaining()),
		)
	}
	return nil
}

func decodeUTF8(value []byte, field string, offset int) (string, error) {
	if !utf8.Valid(value) {
		return "", newProtocolError(ErrorCodeInvalidUTF8, field, offset, ErrInvalidUTF8)
	}
	return string(value), nil
}

func validateString(value, field string, maximum int) error {
	length := len(value)
	if length < 1 || length > maximum {
		return newProtocolError(
			ErrorCodeInvalidLength,
			field,
			-1,
			fmt.Errorf("%w: got %d, want 1..%d", ErrInvalidLength, length, maximum),
		)
	}
	if !utf8.ValidString(value) {
		return newProtocolError(ErrorCodeInvalidUTF8, field, -1, ErrInvalidUTF8)
	}
	return nil
}

func validateBytes(value []byte, field string, maximum int) error {
	if len(value) < 1 || len(value) > maximum {
		return newProtocolError(
			ErrorCodeInvalidLength,
			field,
			-1,
			fmt.Errorf("%w: got %d, want 1..%d", ErrInvalidLength, len(value), maximum),
		)
	}
	return nil
}

func validateJSONObject(value json.RawMessage, field string, maximum int) error {
	if err := validateBytes(value, field, maximum); err != nil {
		return err
	}
	if !utf8.Valid(value) {
		return newProtocolError(ErrorCodeInvalidUTF8, field, -1, ErrInvalidUTF8)
	}
	if !json.Valid(value) {
		return newProtocolError(ErrorCodeInvalidJSON, field, -1, ErrInvalidJSON)
	}
	var object map[string]json.RawMessage
	if err := json.Unmarshal(value, &object); err != nil || object == nil || len(object) == 0 {
		return newProtocolError(ErrorCodeInvalidJSON, field, -1, ErrInvalidJSON)
	}
	return nil
}

func validateJSONArray(value json.RawMessage, field string, maximum int) error {
	return validateJSONArrayLength(value, field, maximum, false)
}

func validateJSONArrayAllowEmpty(value json.RawMessage, field string, maximum int) error {
	return validateJSONArrayLength(value, field, maximum, true)
}

func validateJSONArrayLength(value json.RawMessage, field string, maximum int, allowEmpty bool) error {
	if err := validateBytes(value, field, maximum); err != nil {
		return err
	}
	if !utf8.Valid(value) {
		return newProtocolError(ErrorCodeInvalidUTF8, field, -1, ErrInvalidUTF8)
	}
	if !json.Valid(value) {
		return newProtocolError(ErrorCodeInvalidJSON, field, -1, ErrInvalidJSON)
	}
	var array []json.RawMessage
	if err := json.Unmarshal(value, &array); err != nil || array == nil || (!allowEmpty && len(array) == 0) {
		return newProtocolError(ErrorCodeInvalidJSON, field, -1, ErrInvalidJSON)
	}
	return nil
}

func validateOptionalJSONArray(value json.RawMessage, field string, maximum int) error {
	if len(value) == 0 {
		return nil
	}
	return validateJSONArray(value, field, maximum)
}

func cloneBytes(value []byte) []byte {
	return append([]byte(nil), value...)
}

func typedNil(value any) bool {
	if value == nil {
		return true
	}
	reflected := reflect.ValueOf(value)
	switch reflected.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return reflected.IsNil()
	default:
		return false
	}
}

// Marshal encodes an ACN control PDU with default limits.
func Marshal(message Message) ([]byte, error) {
	return defaultCodec.Marshal(message)
}

// Marshal encodes an ACN control PDU.
func (c *Codec) Marshal(message Message) ([]byte, error) {
	if c == nil {
		return nil, newProtocolError(ErrorCodeInvalidValue, "codec", -1, ErrInvalidValue)
	}
	if typedNil(message) {
		return nil, newProtocolError(ErrorCodeInvalidValue, "message", -1, ErrInvalidValue)
	}
	direction, known := message.MessageType().direction()
	if !known {
		return nil, newProtocolError(
			ErrorCodeUnknownMessageType,
			"message_type",
			-1,
			ErrUnknownMessageType,
		)
	}
	if message.Direction() != direction {
		return nil, newProtocolError(ErrorCodeWrongDirection, "direction", -1, ErrWrongDirection)
	}
	header := message.GetHeader()
	if header.TransactionID == 0 {
		return nil, newProtocolError(ErrorCodeInvalidValue, "transaction_id", -1, ErrInvalidValue)
	}

	output := &encoder{data: make([]byte, 0, 64)}
	output.uint8(Version1)
	output.uint8(uint8(message.MessageType()))
	output.uint8(header.TransactionID)

	var err error
	switch typed := message.(type) {
	case *AgentRegisterRequest:
		err = c.encodeAgentRegisterRequest(output, typed)
	case *AgentRegisterAccept:
		err = c.encodeAgentRegisterAccept(output, typed)
	case *AgentRegisterReject:
		err = encodeAgentRegisterReject(output, typed)
	case *AgentDeregisterRequest:
		err = c.encodeAgentDeregisterRequest(output, typed)
	case *AgentDeregisterAccept:
	case *AgentDeregisterReject:
		err = encodeAgentDeregisterReject(output, typed)
	case *AgentNetworkAbilityRequest:
		err = c.encodeAgentNetworkAbilityRequest(output, typed)
	case *AgentNetworkAbilityAccept:
		err = c.encodeAgentNetworkAbilityAccept(output, typed)
	case *AgentNetworkAbilityReject:
		err = encodeAgentNetworkAbilityReject(output, typed)
	case *AgentPublishRequest:
		err = c.encodeAgentPublishRequest(output, typed)
	case *AgentPublishAccept:
	case *AgentPublishReject:
		err = encodeAgentPublishReject(output, typed)
	case *AgentProfileUpdateRequest:
		err = c.encodeAgentProfileUpdateRequest(output, typed)
	case *AgentProfileUpdateAccept:
		err = c.encodeAgentProfileUpdateAccept(output, typed)
	case *AgentProfileUpdateReject:
		err = encodeAgentProfileUpdateReject(output, typed)
	case *AgentSearchRequest:
		err = c.encodeAgentSearchRequest(output, typed)
	case *AgentSearchResponse:
		err = c.encodeAgentSearchResponse(output, typed)
	case *AgentSearchReject:
		err = encodeAgentSearchReject(output, typed)
	case *AgentGroupingRequest:
		err = c.encodeAgentGroupingRequest(output, typed)
	case *AgentGroupingAccept:
		err = c.encodeAgentGroupingAccept(output, typed)
	case *AgentGroupingReject:
		err = c.encodeAgentGroupingReject(output, typed)
	case *AgentGroupingInvitation:
		err = c.encodeAgentGroupingInvitation(output, typed)
	case *AgentGroupingInvitationResponse:
		err = encodeAgentGroupingInvitationResponse(output, typed)
	case *AgentGroupInfoNotification:
		err = c.encodeAgentGroupInfoNotification(output, typed)
	case *AgentGroupInfoNotificationResponse:
		err = c.encodeAgentGroupInfoNotificationResponse(output, typed)
	default:
		err = newProtocolError(ErrorCodeUnknownMessageType, "message", -1, ErrUnknownMessageType)
	}
	if err != nil {
		return nil, attachHeader(err, header)
	}
	if len(output.data) > c.limits.MaxACNPayload {
		return nil, attachHeader(newProtocolError(
			ErrorCodePayloadTooLarge,
			"message",
			-1,
			fmt.Errorf("%w: got %d, maximum %d",
				ErrPayloadTooLarge, len(output.data), c.limits.MaxACNPayload),
		), header)
	}
	return cloneBytes(output.data), nil
}

// Unmarshal decodes an ACN control PDU with default limits.
func Unmarshal(direction Direction, payload []byte) (Message, error) {
	return defaultCodec.Unmarshal(direction, payload)
}

// Unmarshal decodes an ACN control PDU.
func (c *Codec) Unmarshal(direction Direction, payload []byte) (Message, error) {
	if c == nil {
		return nil, newProtocolError(ErrorCodeInvalidValue, "codec", -1, ErrInvalidValue)
	}
	if direction != Uplink && direction != Downlink {
		return nil, newProtocolError(ErrorCodeInvalidValue, "direction", -1, ErrInvalidValue)
	}
	if len(payload) > c.limits.MaxACNPayload {
		return nil, newProtocolError(
			ErrorCodePayloadTooLarge,
			"message",
			0,
			fmt.Errorf("%w: got %d, maximum %d",
				ErrPayloadTooLarge, len(payload), c.limits.MaxACNPayload),
		)
	}
	input := &decoder{data: payload}
	version, err := input.uint8("version")
	if err != nil {
		return nil, err
	}
	if version != Version1 {
		return nil, newProtocolError(
			ErrorCodeUnsupportedVersion,
			"version",
			0,
			fmt.Errorf("%w: 0x%02x", ErrUnsupportedVersion, version),
		)
	}
	rawType, err := input.uint8("message_type")
	if err != nil {
		return nil, err
	}
	messageType := MessageType(rawType)
	expectedDirection, known := messageType.direction()
	if !known {
		return nil, newProtocolError(
			ErrorCodeUnknownMessageType,
			"message_type",
			1,
			fmt.Errorf("%w: 0x%02x", ErrUnknownMessageType, rawType),
		)
	}
	transactionID, err := input.uint8("transaction_id")
	if err != nil {
		return nil, err
	}
	header := Header{TransactionID: transactionID}
	if transactionID == 0 {
		return nil, attachHeader(
			newProtocolError(ErrorCodeInvalidValue, "transaction_id", 2, ErrInvalidValue),
			header,
		)
	}
	if direction != expectedDirection {
		return nil, attachHeader(newProtocolError(
			ErrorCodeWrongDirection,
			"direction",
			1,
			fmt.Errorf("%w: message is %s, decoder is %s",
				ErrWrongDirection, expectedDirection, direction),
		), header)
	}

	var message Message
	switch messageType {
	case MessageTypeAgentRegisterRequest:
		message, err = c.decodeAgentRegisterRequest(input, header)
	case MessageTypeAgentRegisterAccept:
		message, err = c.decodeAgentRegisterAccept(input, header)
	case MessageTypeAgentRegisterReject:
		message, err = decodeAgentRegisterReject(input, header)
	case MessageTypeAgentDeregisterRequest:
		message, err = c.decodeAgentDeregisterRequest(input, header)
	case MessageTypeAgentDeregisterAccept:
		message = &AgentDeregisterAccept{Header: header}
	case MessageTypeAgentDeregisterReject:
		message, err = decodeAgentDeregisterReject(input, header)
	case MessageTypeAgentNetworkAbilityRequest:
		message, err = c.decodeAgentNetworkAbilityRequest(input, header)
	case MessageTypeAgentNetworkAbilityAccept:
		message, err = c.decodeAgentNetworkAbilityAccept(input, header)
	case MessageTypeAgentNetworkAbilityReject:
		message, err = decodeAgentNetworkAbilityReject(input, header)
	case MessageTypeAgentPublishRequest:
		message, err = c.decodeAgentPublishRequest(input, header)
	case MessageTypeAgentPublishAccept:
		message = &AgentPublishAccept{Header: header}
	case MessageTypeAgentPublishReject:
		message, err = decodeAgentPublishReject(input, header)
	case MessageTypeAgentProfileUpdateRequest:
		message, err = c.decodeAgentProfileUpdateRequest(input, header)
	case MessageTypeAgentProfileUpdateAccept:
		message, err = c.decodeAgentProfileUpdateAccept(input, header)
	case MessageTypeAgentProfileUpdateReject:
		message, err = decodeAgentProfileUpdateReject(input, header)
	case MessageTypeAgentSearchRequest:
		message, err = c.decodeAgentSearchRequest(input, header)
	case MessageTypeAgentSearchResponse:
		message, err = c.decodeAgentSearchResponse(input, header)
	case MessageTypeAgentSearchReject:
		message, err = decodeAgentSearchReject(input, header)
	case MessageTypeAgentGroupingRequest:
		message, err = c.decodeAgentGroupingRequest(input, header)
	case MessageTypeAgentGroupingAccept:
		message, err = c.decodeAgentGroupingAccept(input, header)
	case MessageTypeAgentGroupingReject:
		message, err = c.decodeAgentGroupingReject(input, header)
	case MessageTypeAgentGroupingInvitation:
		message, err = c.decodeAgentGroupingInvitation(input, header)
	case MessageTypeAgentGroupingInvitationResponse:
		message, err = decodeAgentGroupingInvitationResponse(input, header)
	case MessageTypeAgentGroupInfoNotification:
		message, err = c.decodeAgentGroupInfoNotification(input, header)
	case MessageTypeAgentGroupInfoNotificationResponse:
		message, err = c.decodeAgentGroupInfoNotificationResponse(input, header)
	default:
		err = newProtocolError(ErrorCodeUnknownMessageType, "message_type", 1, ErrUnknownMessageType)
	}
	if err != nil {
		return nil, attachHeader(err, header)
	}
	if err := input.finish("message"); err != nil {
		return nil, attachHeader(err, header)
	}
	return message, nil
}
