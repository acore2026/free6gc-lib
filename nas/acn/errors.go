package acn

import (
	"errors"
	"fmt"
)

var (
	ErrTruncated             = errors.New("truncated ACN message")
	ErrInvalidLength         = errors.New("invalid ACN length")
	ErrPayloadTooLarge       = errors.New("ACN payload too large")
	ErrUnsupportedVersion    = errors.New("unsupported ACN version")
	ErrUnknownMessageType    = errors.New("unknown ACN message type")
	ErrWrongDirection        = errors.New("ACN message has wrong direction")
	ErrInvalidValue          = errors.New("invalid ACN value")
	ErrInvalidUTF8           = errors.New("invalid UTF-8")
	ErrInvalidJSON           = errors.New("invalid ACN JSON container")
	ErrInvalidPresenceBitmap = errors.New("invalid ACN presence bitmap")
	ErrCountMismatch         = errors.New("ACN count mismatch")
	ErrTrailingData          = errors.New("trailing ACN data")
	ErrInvalidNASEnvelope    = errors.New("invalid ACN NAS envelope")
	ErrProtectedNAS          = errors.New("protected NAS must be unprotected before ACN decoding")
)

type ErrorCode string

const (
	ErrorCodeTruncated             ErrorCode = "truncated"
	ErrorCodeInvalidLength         ErrorCode = "invalid_length"
	ErrorCodePayloadTooLarge       ErrorCode = "payload_too_large"
	ErrorCodeUnsupportedVersion    ErrorCode = "unsupported_version"
	ErrorCodeUnknownMessageType    ErrorCode = "unknown_message_type"
	ErrorCodeWrongDirection        ErrorCode = "wrong_direction"
	ErrorCodeInvalidValue          ErrorCode = "invalid_value"
	ErrorCodeInvalidUTF8           ErrorCode = "invalid_utf8"
	ErrorCodeInvalidJSON           ErrorCode = "invalid_json"
	ErrorCodeInvalidPresenceBitmap ErrorCode = "invalid_presence_bitmap"
	ErrorCodeCountMismatch         ErrorCode = "count_mismatch"
	ErrorCodeTrailingData          ErrorCode = "trailing_data"
	ErrorCodeInvalidNASEnvelope    ErrorCode = "invalid_nas_envelope"
)

// ProtocolError describes a wire-level ACN failure.
type ProtocolError struct {
	Code   ErrorCode
	Field  string
	Offset int
	Header *Header
	Cause  error
}

func (e *ProtocolError) Error() string {
	if e == nil {
		return "<nil>"
	}
	location := ""
	if e.Field != "" {
		location = " field " + e.Field
	}
	if e.Offset >= 0 {
		location += fmt.Sprintf(" at offset %d", e.Offset)
	}
	if e.Cause == nil {
		return string(e.Code) + location
	}
	return fmt.Sprintf("%s%s: %v", e.Code, location, e.Cause)
}

func (e *ProtocolError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

func newProtocolError(code ErrorCode, field string, offset int, cause error) error {
	return &ProtocolError{
		Code:   code,
		Field:  field,
		Offset: offset,
		Cause:  cause,
	}
}

func attachHeader(err error, header Header) error {
	var protocolErr *ProtocolError
	if errors.As(err, &protocolErr) {
		copyErr := *protocolErr
		if copyErr.Header == nil {
			copyHeader := header
			copyErr.Header = &copyHeader
		}
		return &copyErr
	}
	return err
}
