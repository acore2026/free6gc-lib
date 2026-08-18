package acn

func (c *Codec) encodeAgentDeregisterRequest(
	output *encoder,
	message *AgentDeregisterRequest,
) error {
	if err := validateString(message.AgentID, "agent_id", c.limits.MaxAgentID); err != nil {
		return err
	}
	if !validDeregistrationReason(message.Reason) {
		return newProtocolError(ErrorCodeInvalidValue, "reason", -1, ErrInvalidValue)
	}
	if err := validateBytes(message.Signature, "signature", c.limits.MaxSignature); err != nil {
		return err
	}
	output.lve([]byte(message.AgentID))
	output.uint8(uint8(message.Reason))
	output.uint64(message.Timestamp)
	output.lve(message.Signature)
	return nil
}

func (c *Codec) decodeAgentDeregisterRequest(
	input *decoder,
	header Header,
) (Message, error) {
	message := &AgentDeregisterRequest{Header: header}
	var err error
	if message.AgentID, err = input.lveString("agent_id", c.limits.MaxAgentID); err != nil {
		return nil, err
	}
	reason, err := input.uint8("reason")
	if err != nil {
		return nil, err
	}
	message.Reason = DeregistrationReason(reason)
	if !validDeregistrationReason(message.Reason) {
		return nil, newProtocolError(ErrorCodeInvalidValue, "reason", input.offset-1, ErrInvalidValue)
	}
	if message.Timestamp, err = input.uint64("timestamp"); err != nil {
		return nil, err
	}
	signature, err := input.lveBytes("signature", c.limits.MaxSignature)
	if err != nil {
		return nil, err
	}
	message.Signature = cloneBytes(signature)
	return message, nil
}

func validDeregistrationReason(value DeregistrationReason) bool {
	return value <= DeregistrationReasonRetired || value == DeregistrationReasonOther
}

func encodeAgentDeregisterReject(output *encoder, message *AgentDeregisterReject) error {
	if !validDeregisterRejectCause(message.Cause) {
		return newProtocolError(ErrorCodeInvalidValue, "reject_cause", -1, ErrInvalidValue)
	}
	if !validDeregisterFailedField(message.FailedField) {
		return newProtocolError(ErrorCodeInvalidValue, "failed_field", -1, ErrInvalidValue)
	}
	output.uint8(uint8(message.Cause))
	output.uint8(uint8(message.FailedField))
	return nil
}

func decodeAgentDeregisterReject(input *decoder, header Header) (Message, error) {
	cause, err := input.uint8("reject_cause")
	if err != nil {
		return nil, err
	}
	if !validDeregisterRejectCause(DeregisterRejectCause(cause)) {
		return nil, newProtocolError(ErrorCodeInvalidValue, "reject_cause", input.offset-1, ErrInvalidValue)
	}
	failedField, err := input.uint8("failed_field")
	if err != nil {
		return nil, err
	}
	if !validDeregisterFailedField(DeregisterFailedField(failedField)) {
		return nil, newProtocolError(ErrorCodeInvalidValue, "failed_field", input.offset-1, ErrInvalidValue)
	}
	return &AgentDeregisterReject{
		Header:      header,
		Cause:       DeregisterRejectCause(cause),
		FailedField: DeregisterFailedField(failedField),
	}, nil
}

func validDeregisterRejectCause(value DeregisterRejectCause) bool {
	return value >= DeregisterRejectUnknownAgentID && value <= DeregisterRejectBackendFailure
}

func validDeregisterFailedField(value DeregisterFailedField) bool {
	return value >= DeregisterFieldUnspecified && value <= DeregisterFieldUEAgentBinding
}
