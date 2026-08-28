package acn

func (c *Codec) encodeAgentSearchRequest(output *encoder, message *AgentSearchRequest) error {
	if err := validateJSONObject(message.DiscoveryRequest, "discovery_request", c.limits.MaxJSONContainer); err != nil {
		return err
	}
	output.lve(message.DiscoveryRequest)
	return nil
}

func (c *Codec) decodeAgentSearchRequest(input *decoder, header Header) (Message, error) {
	discoveryRequest, err := input.lveBytes("discovery_request", c.limits.MaxJSONContainer)
	if err != nil {
		return nil, err
	}
	if err := validateJSONObject(discoveryRequest, "discovery_request", c.limits.MaxJSONContainer); err != nil {
		return nil, err
	}
	return &AgentSearchRequest{Header: header, DiscoveryRequest: cloneBytes(discoveryRequest)}, nil
}

func (c *Codec) encodeAgentSearchResponse(output *encoder, message *AgentSearchResponse) error {
	if err := validateJSONObject(message.DiscoveryResult, "discovery_result", c.limits.MaxJSONContainer); err != nil {
		return err
	}
	output.lve(message.DiscoveryResult)
	return nil
}

func (c *Codec) decodeAgentSearchResponse(input *decoder, header Header) (Message, error) {
	discoveryResult, err := input.lveBytes("discovery_result", c.limits.MaxJSONContainer)
	if err != nil {
		return nil, err
	}
	if err := validateJSONObject(discoveryResult, "discovery_result", c.limits.MaxJSONContainer); err != nil {
		return nil, err
	}
	return &AgentSearchResponse{Header: header, DiscoveryResult: cloneBytes(discoveryResult)}, nil
}

func encodeAgentSearchReject(output *encoder, message *AgentSearchReject) error {
	if message.Cause < SearchRejectAgentInvalid || message.Cause > SearchRejectInternalError {
		return newProtocolError(ErrorCodeInvalidValue, "reject_cause", -1, ErrInvalidValue)
	}
	if message.FailedField < SearchFieldUnspecified || message.FailedField > SearchFieldProof {
		return newProtocolError(ErrorCodeInvalidValue, "failed_field", -1, ErrInvalidValue)
	}
	output.uint8(uint8(message.Cause))
	output.uint8(uint8(message.FailedField))
	return nil
}

func decodeAgentSearchReject(input *decoder, header Header) (Message, error) {
	cause, err := input.uint8("reject_cause")
	if err != nil {
		return nil, err
	}
	field, err := input.uint8("failed_field")
	if err != nil {
		return nil, err
	}
	message := &AgentSearchReject{Header: header, Cause: SearchRejectCause(cause), FailedField: SearchFailedField(field)}
	if err := encodeAgentSearchReject(&encoder{}, message); err != nil {
		return nil, err
	}
	return message, nil
}
