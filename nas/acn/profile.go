package acn

func (c *Codec) encodeAgentProfileUpdateRequest(output *encoder, message *AgentProfileUpdateRequest) error {
	if err := validateJSONObject(message.ProfileUpdate, "profile_update", c.limits.MaxJSONContainer); err != nil {
		return err
	}
	output.lve(message.ProfileUpdate)
	return nil
}

func (c *Codec) decodeAgentProfileUpdateRequest(input *decoder, header Header) (Message, error) {
	profileUpdate, err := input.lveBytes("profile_update", c.limits.MaxJSONContainer)
	if err != nil {
		return nil, err
	}
	if err := validateJSONObject(profileUpdate, "profile_update", c.limits.MaxJSONContainer); err != nil {
		return nil, err
	}
	return &AgentProfileUpdateRequest{Header: header, ProfileUpdate: cloneBytes(profileUpdate)}, nil
}

func (c *Codec) encodeAgentProfileUpdateAccept(output *encoder, message *AgentProfileUpdateAccept) error {
	if err := validateJSONObject(message.ProfileUpdateResult, "profile_update_result", c.limits.MaxJSONContainer); err != nil {
		return err
	}
	output.lve(message.ProfileUpdateResult)
	return nil
}

func (c *Codec) decodeAgentProfileUpdateAccept(input *decoder, header Header) (Message, error) {
	profileUpdateResult, err := input.lveBytes("profile_update_result", c.limits.MaxJSONContainer)
	if err != nil {
		return nil, err
	}
	if err := validateJSONObject(profileUpdateResult, "profile_update_result", c.limits.MaxJSONContainer); err != nil {
		return nil, err
	}
	return &AgentProfileUpdateAccept{Header: header, ProfileUpdateResult: cloneBytes(profileUpdateResult)}, nil
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
	output.uint16(message.RelatedItemIndex)
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
	relatedItemIndex, err := input.uint16("related_item_index")
	if err != nil {
		return nil, err
	}
	message := &AgentProfileUpdateReject{
		Header: header, Cause: ProfileUpdateRejectCause(cause), FailedField: ProfileUpdateFailedField(field),
		RelatedItemIndex: relatedItemIndex,
	}
	if err := encodeAgentProfileUpdateReject(&encoder{}, message); err != nil {
		return nil, err
	}
	return message, nil
}
