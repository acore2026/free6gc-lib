package acn

func (c *Codec) encodeAgentProfileUpdateRequest(output *encoder, message *AgentProfileUpdateRequest) error {
	if err := validateString(message.RequestID, "request_id", c.limits.MaxTaskID); err != nil {
		return err
	}
	if err := validateString(message.AgentID, "agent_id", c.limits.MaxAgentID); err != nil {
		return err
	}
	if err := validateJSONArray(message.UpdateItems, "update_items", c.limits.MaxJSONContainer); err != nil {
		return err
	}
	if err := validateOptionalJSONArray(message.Credentials, "credentials", c.limits.MaxJSONContainer); err != nil {
		return err
	}
	output.lve([]byte(message.RequestID))
	output.lve([]byte(message.AgentID))
	output.lve(message.UpdateItems)
	output.lve(message.Credentials)
	output.uint64(message.Timestamp)
	return c.encodeProof(output, message.Proof, "proof")
}

func (c *Codec) decodeAgentProfileUpdateRequest(input *decoder, header Header) (Message, error) {
	message := &AgentProfileUpdateRequest{Header: header}
	var err error
	if message.RequestID, err = input.lveString("request_id", c.limits.MaxTaskID); err != nil {
		return nil, err
	}
	if message.AgentID, err = input.lveString("agent_id", c.limits.MaxAgentID); err != nil {
		return nil, err
	}
	updateItems, err := input.lveBytes("update_items", c.limits.MaxJSONContainer)
	if err != nil {
		return nil, err
	}
	if err := validateJSONArray(updateItems, "update_items", c.limits.MaxJSONContainer); err != nil {
		return nil, err
	}
	message.UpdateItems = cloneBytes(updateItems)
	credentials, err := input.lveOptionalBytes("credentials", c.limits.MaxJSONContainer)
	if err != nil {
		return nil, err
	}
	if err := validateOptionalJSONArray(credentials, "credentials", c.limits.MaxJSONContainer); err != nil {
		return nil, err
	}
	message.Credentials = cloneBytes(credentials)
	if message.Timestamp, err = input.uint64("timestamp"); err != nil {
		return nil, err
	}
	if message.Proof, err = c.decodeProof(input, "proof"); err != nil {
		return nil, err
	}
	return message, nil
}

func encodeAgentProfileUpdateReject(output *encoder, message *AgentProfileUpdateReject) error {
	if message.Cause < ProfileUpdateRejectAgentInvalid || message.Cause > ProfileUpdateRejectInternalError {
		return newProtocolError(ErrorCodeInvalidValue, "reject_cause", -1, ErrInvalidValue)
	}
	if message.FailedField < ProfileUpdateFieldUnspecified || message.FailedField > ProfileUpdateFieldProof {
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
	field, err := input.uint8("failed_field")
	if err != nil {
		return nil, err
	}
	message := &AgentProfileUpdateReject{
		Header: header, Cause: ProfileUpdateRejectCause(cause), FailedField: ProfileUpdateFailedField(field),
	}
	if err := encodeAgentProfileUpdateReject(&encoder{}, message); err != nil {
		return nil, err
	}
	return message, nil
}
