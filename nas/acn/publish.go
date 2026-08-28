package acn

func (c *Codec) encodeAgentPublishRequest(output *encoder, message *AgentPublishRequest) error {
	if err := validateJSONObject(message.ProfilePublish, "profile_publish", c.limits.MaxJSONContainer); err != nil {
		return err
	}
	output.lve(message.ProfilePublish)
	return nil
}

func (c *Codec) decodeAgentPublishRequest(input *decoder, header Header) (Message, error) {
	value, err := input.lveBytes("profile_publish", c.limits.MaxJSONContainer)
	if err != nil {
		return nil, err
	}
	if err := validateJSONObject(value, "profile_publish", c.limits.MaxJSONContainer); err != nil {
		return nil, err
	}
	return &AgentPublishRequest{Header: header, ProfilePublish: cloneBytes(value)}, nil
}

func encodeAgentPublishReject(output *encoder, message *AgentPublishReject) error {
	if message.Cause < PublishRejectAgentInvalid || message.Cause > PublishRejectInternalError {
		return newProtocolError(ErrorCodeInvalidValue, "reject_cause", -1, ErrInvalidValue)
	}
	if message.FailedField < PublishFieldUnspecified || message.FailedField > PublishFieldServiceEndpoints {
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
