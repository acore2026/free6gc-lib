package acn

func (c *Codec) encodeAgentNetworkAbilityRequest(output *encoder, message *AgentNetworkAbilityRequest) error {
	if err := validateJSONObject(message.NetworkAbilityRequest, "network_ability_request", c.limits.MaxJSONContainer); err != nil {
		return err
	}
	output.lve(message.NetworkAbilityRequest)
	return nil
}

func (c *Codec) decodeAgentNetworkAbilityRequest(input *decoder, header Header) (Message, error) {
	value, err := input.lveBytes("network_ability_request", c.limits.MaxJSONContainer)
	if err != nil {
		return nil, err
	}
	if err := validateJSONObject(value, "network_ability_request", c.limits.MaxJSONContainer); err != nil {
		return nil, err
	}
	return &AgentNetworkAbilityRequest{Header: header, NetworkAbilityRequest: cloneBytes(value)}, nil
}

func (c *Codec) encodeAgentNetworkAbilityAccept(output *encoder, message *AgentNetworkAbilityAccept) error {
	if err := validateJSONObject(message.NetworkAbilityResult, "network_ability_result", c.limits.MaxJSONContainer); err != nil {
		return err
	}
	output.lve(message.NetworkAbilityResult)
	return nil
}

func (c *Codec) decodeAgentNetworkAbilityAccept(input *decoder, header Header) (Message, error) {
	value, err := input.lveBytes("network_ability_result", c.limits.MaxJSONContainer)
	if err != nil {
		return nil, err
	}
	if err := validateJSONObject(value, "network_ability_result", c.limits.MaxJSONContainer); err != nil {
		return nil, err
	}
	return &AgentNetworkAbilityAccept{Header: header, NetworkAbilityResult: cloneBytes(value)}, nil
}

func encodeAgentNetworkAbilityReject(output *encoder, message *AgentNetworkAbilityReject) error {
	if message.Cause < NetworkAbilityRejectAgentInvalid || message.Cause > NetworkAbilityRejectInternalError {
		return newProtocolError(ErrorCodeInvalidValue, "reject_cause", -1, ErrInvalidValue)
	}
	if message.FailedField < NetworkAbilityFieldUnspecified || message.FailedField > NetworkAbilityFieldProof {
		return newProtocolError(ErrorCodeInvalidValue, "failed_field", -1, ErrInvalidValue)
	}
	output.uint8(uint8(message.Cause))
	output.uint8(uint8(message.FailedField))
	return nil
}

func decodeAgentNetworkAbilityReject(input *decoder, header Header) (Message, error) {
	cause, err := input.uint8("reject_cause")
	if err != nil {
		return nil, err
	}
	field, err := input.uint8("failed_field")
	if err != nil {
		return nil, err
	}
	message := &AgentNetworkAbilityReject{
		Header: header, Cause: NetworkAbilityRejectCause(cause), FailedField: NetworkAbilityFailedField(field),
	}
	if err := encodeAgentNetworkAbilityReject(&encoder{}, message); err != nil {
		return nil, err
	}
	return message, nil
}
