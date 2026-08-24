package acn

func (c *Codec) encodeAgentPublishRequest(output *encoder, message *AgentPublishRequest) error {
	if err := validateString(message.AgentID, "agent_id", c.limits.MaxAgentID); err != nil {
		return err
	}
	if !validPriority(message.Priority) {
		return newProtocolError(ErrorCodeInvalidValue, "priority", -1, ErrInvalidValue)
	}
	if err := validateJSONArray(message.VCList, "vc_list", c.limits.MaxJSONContainer); err != nil {
		return err
	}
	output.lve([]byte(message.AgentID))
	output.uint8(uint8(message.Priority))
	output.uint64(message.Timestamp)
	if err := c.encodeProof(output, message.Proof, "proof"); err != nil {
		return err
	}
	output.lve(message.VCList)
	return nil
}

func (c *Codec) decodeAgentPublishRequest(input *decoder, header Header) (Message, error) {
	message := &AgentPublishRequest{Header: header}
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
	if message.Proof, err = c.decodeProof(input, "proof"); err != nil {
		return nil, err
	}
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

func encodeAgentPublishReject(output *encoder, message *AgentPublishReject) error {
	if message.Cause < PublishRejectAgentInvalid || message.Cause > PublishRejectInternalError {
		return newProtocolError(ErrorCodeInvalidValue, "reject_cause", -1, ErrInvalidValue)
	}
	if message.FailedField < PublishFieldUnspecified || message.FailedField > PublishFieldVCList {
		return newProtocolError(ErrorCodeInvalidValue, "failed_field", -1, ErrInvalidValue)
	}
	output.uint8(uint8(message.Cause))
	output.uint8(uint8(message.FailedField))
	return nil
}

func decodeAgentPublishReject(input *decoder, header Header) (Message, error) {
	cause, err := input.uint8("reject_cause")
	if err != nil {
		return nil, err
	}
	field, err := input.uint8("failed_field")
	if err != nil {
		return nil, err
	}
	message := &AgentPublishReject{Header: header, Cause: PublishRejectCause(cause), FailedField: PublishFailedField(field)}
	if err := encodeAgentPublishReject(&encoder{}, message); err != nil {
		return nil, err
	}
	return message, nil
}

func validPriority(value Priority) bool {
	return value >= PriorityUnspecified && value <= PriorityLow
}
