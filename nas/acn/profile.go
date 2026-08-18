package acn

func (c *Codec) encodeAgentProfileUpdateRequest(
	output *encoder,
	message *AgentProfileUpdateRequest,
) error {
	if err := validateString(message.AgentID, "agent_id", c.limits.MaxAgentID); err != nil {
		return err
	}
	if !validPriority(message.Priority) {
		return newProtocolError(ErrorCodeInvalidValue, "priority", -1, ErrInvalidValue)
	}
	if err := validateBytes(message.Signature, "signature", c.limits.MaxSignature); err != nil {
		return err
	}
	if err := validateJSONArray(message.VCList, "vc_list", c.limits.MaxJSONContainer); err != nil {
		return err
	}
	output.lve([]byte(message.AgentID))
	output.uint8(uint8(message.Priority))
	output.uint64(message.Timestamp)
	output.lve(message.Signature)
	output.lve(message.VCList)
	return nil
}

func (c *Codec) decodeAgentProfileUpdateRequest(
	input *decoder,
	header Header,
) (Message, error) {
	message := &AgentProfileUpdateRequest{Header: header}
	var err error
	if message.AgentID, err = input.lveString("agent_id", c.limits.MaxAgentID); err != nil {
		return nil, err
	}
	priority, err := input.uint8("priority")
	if err != nil {
		return nil, err
	}
	message.Priority = Priority(priority)
	if !validPriority(message.Priority) {
		return nil, newProtocolError(ErrorCodeInvalidValue, "priority", input.offset-1, ErrInvalidValue)
	}
	if message.Timestamp, err = input.uint64("timestamp"); err != nil {
		return nil, err
	}
	signature, err := input.lveBytes("signature", c.limits.MaxSignature)
	if err != nil {
		return nil, err
	}
	message.Signature = cloneBytes(signature)
	vcList, err := input.lveBytes("vc_list", c.limits.MaxJSONContainer)
	if err != nil {
		return nil, err
	}
	if err := validateJSONArray(vcList, "vc_list", c.limits.MaxJSONContainer); err != nil {
		return nil, err
	}
	message.VCList = cloneBytes(vcList)
	return message, nil
}

func validPriority(value Priority) bool {
	return value <= PriorityLow
}

func encodeAgentProfileUpdateReject(output *encoder, message *AgentProfileUpdateReject) error {
	if !validProfileUpdateRejectCause(message.Cause) {
		return newProtocolError(ErrorCodeInvalidValue, "reject_cause", -1, ErrInvalidValue)
	}
	if !validProfileUpdateFailedField(message.FailedField) {
		return newProtocolError(ErrorCodeInvalidValue, "failed_field", -1, ErrInvalidValue)
	}
	output.uint8(uint8(message.Cause))
	output.uint8(uint8(message.FailedField))
	return nil
}

func decodeAgentProfileUpdateReject(input *decoder, header Header) (Message, error) {
	cause, err := input.uint8("reject_cause")
	if err != nil {
		return nil, err
	}
	if !validProfileUpdateRejectCause(ProfileUpdateRejectCause(cause)) {
		return nil, newProtocolError(ErrorCodeInvalidValue, "reject_cause", input.offset-1, ErrInvalidValue)
	}
	failedField, err := input.uint8("failed_field")
	if err != nil {
		return nil, err
	}
	if !validProfileUpdateFailedField(ProfileUpdateFailedField(failedField)) {
		return nil, newProtocolError(ErrorCodeInvalidValue, "failed_field", input.offset-1, ErrInvalidValue)
	}
	return &AgentProfileUpdateReject{
		Header:      header,
		Cause:       ProfileUpdateRejectCause(cause),
		FailedField: ProfileUpdateFailedField(failedField),
	}, nil
}

func validProfileUpdateRejectCause(value ProfileUpdateRejectCause) bool {
	return value >= ProfileUpdateRejectInvalidAgentID && value <= ProfileUpdateRejectInternalError
}

func validProfileUpdateFailedField(value ProfileUpdateFailedField) bool {
	return value >= ProfileUpdateFieldUnspecified && value <= ProfileUpdateFieldVCList
}
