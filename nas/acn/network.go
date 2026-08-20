package acn

func (c *Codec) encodeAgentNetworkAbilityRequest(output *encoder, message *AgentNetworkAbilityRequest) error {
	if err := validateString(message.AgentID, "agent_id", c.limits.MaxAgentID); err != nil {
		return err
	}
	output.lve([]byte(message.AgentID))
	output.uint64(message.Timestamp)
	return c.encodeProof(output, message.Proof, "proof")
}

func (c *Codec) decodeAgentNetworkAbilityRequest(input *decoder, header Header) (Message, error) {
	agentID, err := input.lveString("agent_id", c.limits.MaxAgentID)
	if err != nil {
		return nil, err
	}
	timestamp, err := input.uint64("timestamp")
	if err != nil {
		return nil, err
	}
	proof, err := c.decodeProof(input, "proof")
	if err != nil {
		return nil, err
	}
	return &AgentNetworkAbilityRequest{
		Header: header, AgentID: agentID, Timestamp: timestamp, Proof: proof,
	}, nil
}

func (c *Codec) encodeAgentNetworkAbilityResponse(output *encoder, message *AgentNetworkAbilityResponse) error {
	if err := validateJSONObject(message.VC1, "vc1", c.limits.MaxJSONContainer); err != nil {
		return err
	}
	output.uint64(message.Timestamp)
	output.lve(message.VC1)
	return nil
}

func (c *Codec) decodeAgentNetworkAbilityResponse(input *decoder, header Header) (Message, error) {
	timestamp, err := input.uint64("timestamp")
	if err != nil {
		return nil, err
	}
	vc1, err := input.lveBytes("vc1", c.limits.MaxJSONContainer)
	if err != nil {
		return nil, err
	}
	if err := validateJSONObject(vc1, "vc1", c.limits.MaxJSONContainer); err != nil {
		return nil, err
	}
	return &AgentNetworkAbilityResponse{Header: header, Timestamp: timestamp, VC1: cloneBytes(vc1)}, nil
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
