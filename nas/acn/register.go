package acn

func (c *Codec) encodeAgentRegisterRequest(output *encoder, message *AgentRegisterRequest) error {
	for _, field := range []struct {
		value   string
		name    string
		maximum int
	}{
		{message.Owner, "owner", c.limits.MaxOwnerID},
		{message.AgentName, "agent_name", c.limits.MaxAgentName},
		{message.Description, "description", c.limits.MaxDescription},
	} {
		if err := validateString(field.value, field.name, field.maximum); err != nil {
			return err
		}
	}
	if err := validateBytes(message.PublicKey, "public_key", c.limits.MaxPublicKey); err != nil {
		return err
	}
	if err := validateBytes(message.Signature, "signature", c.limits.MaxSignature); err != nil {
		return err
	}
	if err := validateJSONObject(message.Metadata, "metadata", c.limits.MaxJSONContainer); err != nil {
		return err
	}
	output.lve([]byte(message.Owner))
	output.lve([]byte(message.AgentName))
	output.lve(message.PublicKey)
	output.lve([]byte(message.Description))
	output.uint64(message.Timestamp)
	output.lve(message.Signature)
	output.lve(message.Metadata)
	return nil
}

func (c *Codec) decodeAgentRegisterRequest(
	input *decoder,
	header Header,
) (Message, error) {
	message := &AgentRegisterRequest{Header: header}
	var err error
	if message.Owner, err = input.lveString("owner", c.limits.MaxOwnerID); err != nil {
		return nil, err
	}
	if message.AgentName, err = input.lveString("agent_name", c.limits.MaxAgentName); err != nil {
		return nil, err
	}
	publicKey, err := input.lveBytes("public_key", c.limits.MaxPublicKey)
	if err != nil {
		return nil, err
	}
	message.PublicKey = cloneBytes(publicKey)
	if message.Description, err = input.lveString(
		"description",
		c.limits.MaxDescription,
	); err != nil {
		return nil, err
	}
	if message.Timestamp, err = input.uint64("timestamp"); err != nil {
		return nil, err
	}
	signature, err := input.lveBytes("signature", c.limits.MaxSignature)
	if err != nil {
		return nil, err
	}
	message.Signature = cloneBytes(signature)
	metadata, err := input.lveBytes("metadata", c.limits.MaxJSONContainer)
	if err != nil {
		return nil, err
	}
	if err := validateJSONObject(metadata, "metadata", c.limits.MaxJSONContainer); err != nil {
		return nil, err
	}
	message.Metadata = cloneBytes(metadata)
	return message, nil
}

func (c *Codec) encodeAgentRegisterAccept(output *encoder, message *AgentRegisterAccept) error {
	if err := validateString(message.AgentID, "agent_id", c.limits.MaxAgentID); err != nil {
		return err
	}
	if err := validateJSONObject(message.VC0, "vc0", c.limits.MaxJSONContainer); err != nil {
		return err
	}
	output.lve([]byte(message.AgentID))
	output.lve(message.VC0)
	return nil
}

func (c *Codec) decodeAgentRegisterAccept(
	input *decoder,
	header Header,
) (Message, error) {
	message := &AgentRegisterAccept{Header: header}
	var err error
	if message.AgentID, err = input.lveString("agent_id", c.limits.MaxAgentID); err != nil {
		return nil, err
	}
	vc0, err := input.lveBytes("vc0", c.limits.MaxJSONContainer)
	if err != nil {
		return nil, err
	}
	if err := validateJSONObject(vc0, "vc0", c.limits.MaxJSONContainer); err != nil {
		return nil, err
	}
	message.VC0 = cloneBytes(vc0)
	return message, nil
}

func encodeAgentRegisterReject(output *encoder, message *AgentRegisterReject) error {
	if !validRegisterRejectCause(message.Cause) {
		return newProtocolError(ErrorCodeInvalidValue, "reject_cause", -1, ErrInvalidValue)
	}
	if !validRegisterFailedField(message.FailedField) {
		return newProtocolError(ErrorCodeInvalidValue, "failed_field", -1, ErrInvalidValue)
	}
	output.uint8(uint8(message.Cause))
	output.uint8(uint8(message.FailedField))
	return nil
}

func decodeAgentRegisterReject(input *decoder, header Header) (Message, error) {
	cause, err := input.uint8("reject_cause")
	if err != nil {
		return nil, err
	}
	if !validRegisterRejectCause(RegisterRejectCause(cause)) {
		return nil, newProtocolError(ErrorCodeInvalidValue, "reject_cause", input.offset-1, ErrInvalidValue)
	}
	failedField, err := input.uint8("failed_field")
	if err != nil {
		return nil, err
	}
	if !validRegisterFailedField(RegisterFailedField(failedField)) {
		return nil, newProtocolError(ErrorCodeInvalidValue, "failed_field", input.offset-1, ErrInvalidValue)
	}
	return &AgentRegisterReject{
		Header:      header,
		Cause:       RegisterRejectCause(cause),
		FailedField: RegisterFailedField(failedField),
	}, nil
}

func validRegisterRejectCause(value RegisterRejectCause) bool {
	return value >= RegisterRejectRequestValidationFailed && value <= RegisterRejectInternalError
}

func validRegisterFailedField(value RegisterFailedField) bool {
	return value >= RegisterFieldUnspecified && value <= RegisterFieldMetadata
}
