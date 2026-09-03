package acn

func (c *Codec) encodeAgentGroupingRequest(output *encoder, message *AgentGroupingRequest) error {
	if err := validateJSONObject(message.GroupCreation, "group_creation", c.limits.MaxJSONContainer); err != nil {
		return err
	}
	output.lve(message.GroupCreation)
	return nil
}

func (c *Codec) decodeAgentGroupingRequest(input *decoder, header Header) (Message, error) {
	groupCreation, err := input.lveBytes("group_creation", c.limits.MaxJSONContainer)
	if err != nil {
		return nil, err
	}
	if err := validateJSONObject(groupCreation, "group_creation", c.limits.MaxJSONContainer); err != nil {
		return nil, err
	}
	return &AgentGroupingRequest{Header: header, GroupCreation: cloneBytes(groupCreation)}, nil
}

func (c *Codec) encodeAgentGroupingAccept(output *encoder, message *AgentGroupingAccept) error {
	if err := validateString(message.GroupID, "group_id", c.limits.MaxGroupID); err != nil {
		return err
	}
	output.lve([]byte(message.GroupID))
	return nil
}

func (c *Codec) decodeAgentGroupingAccept(input *decoder, header Header) (Message, error) {
	groupID, err := input.lveString("group_id", c.limits.MaxGroupID)
	if err != nil {
		return nil, err
	}
	return &AgentGroupingAccept{Header: header, GroupID: groupID}, nil
}

func (c *Codec) encodeAgentGroupingReject(output *encoder, message *AgentGroupingReject) error {
	if message.Cause < GroupingRejectSourceAgentInvalid || message.Cause > GroupingRejectInternalError {
		return newProtocolError(ErrorCodeInvalidValue, "reject_cause", -1, ErrInvalidValue)
	}
	if message.FailedField < GroupingFieldUnspecified || message.FailedField > GroupingFieldProof {
		return newProtocolError(ErrorCodeInvalidValue, "failed_field", -1, ErrInvalidValue)
	}
	if message.RelatedAgentID != "" {
		if err := validateString(message.RelatedAgentID, "related_agent_id", c.limits.MaxAgentID); err != nil {
			return err
		}
	}
	output.uint8(uint8(message.Cause))
	output.uint8(uint8(message.FailedField))
	output.lve([]byte(message.RelatedAgentID))
	return nil
}

func (c *Codec) decodeAgentGroupingReject(input *decoder, header Header) (Message, error) {
	cause, err := input.uint8("reject_cause")
	if err != nil {
		return nil, err
	}
	field, err := input.uint8("failed_field")
	if err != nil {
		return nil, err
	}
	relatedAgentID, err := input.lveOptionalString("related_agent_id", c.limits.MaxAgentID)
	if err != nil {
		return nil, err
	}
	message := &AgentGroupingReject{
		Header: header, Cause: GroupingRejectCause(cause), FailedField: GroupingFailedField(field), RelatedAgentID: relatedAgentID,
	}
	if err := c.encodeAgentGroupingReject(&encoder{}, message); err != nil {
		return nil, err
	}
	return message, nil
}

func (c *Codec) encodeAgentGroupingInvitation(output *encoder, message *AgentGroupingInvitation) error {
	if err := validateJSONObject(message.GroupInvitation, "group_invitation", c.limits.MaxJSONContainer); err != nil {
		return err
	}
	output.lve(message.GroupInvitation)
	return nil
}

func (c *Codec) decodeAgentGroupingInvitation(input *decoder, header Header) (Message, error) {
	groupInvitation, err := input.lveBytes("group_invitation", c.limits.MaxJSONContainer)
	if err != nil {
		return nil, err
	}
	if err := validateJSONObject(groupInvitation, "group_invitation", c.limits.MaxJSONContainer); err != nil {
		return nil, err
	}
	return &AgentGroupingInvitation{Header: header, GroupInvitation: cloneBytes(groupInvitation)}, nil
}

func (c *Codec) encodeAgentGroupingInvitationResponse(output *encoder, message *AgentGroupingInvitationResponse) error {
	if message.Decision != GroupingDecisionAccept && message.Decision != GroupingDecisionReject {
		return newProtocolError(ErrorCodeInvalidValue, "decision", -1, ErrInvalidValue)
	}
	if err := validateString(message.GroupID, "group_id", c.limits.MaxGroupID); err != nil {
		return err
	}
	output.uint8(uint8(message.Decision))
	output.lve([]byte(message.GroupID))
	return nil
}

func (c *Codec) decodeAgentGroupingInvitationResponse(input *decoder, header Header) (Message, error) {
	decision, err := input.uint8("decision")
	if err != nil {
		return nil, err
	}
	groupID, err := input.lveString("group_id", c.limits.MaxGroupID)
	if err != nil {
		return nil, err
	}
	message := &AgentGroupingInvitationResponse{Header: header, Decision: GroupingDecision(decision), GroupID: groupID}
	if err := c.encodeAgentGroupingInvitationResponse(&encoder{}, message); err != nil {
		return nil, err
	}
	return message, nil
}

func (c *Codec) encodeAgentGroupInfoNotification(output *encoder, message *AgentGroupInfoNotification) error {
	if err := validateJSONObject(message.GroupConfig, "group_config", c.limits.MaxJSONContainer); err != nil {
		return err
	}
	output.lve(message.GroupConfig)
	return nil
}

func (c *Codec) decodeAgentGroupInfoNotification(input *decoder, header Header) (Message, error) {
	groupConfig, err := input.lveBytes("group_config", c.limits.MaxJSONContainer)
	if err != nil {
		return nil, err
	}
	if err := validateJSONObject(groupConfig, "group_config", c.limits.MaxJSONContainer); err != nil {
		return nil, err
	}
	return &AgentGroupInfoNotification{Header: header, GroupConfig: cloneBytes(groupConfig)}, nil
}

func (c *Codec) encodeAgentGroupInfoNotificationResponse(output *encoder, message *AgentGroupInfoNotificationResponse) error {
	if message.ApplyResult != GroupInfoApplyResultACK && message.ApplyResult != GroupInfoApplyResultReject {
		return newProtocolError(ErrorCodeInvalidValue, "apply_result", -1, ErrInvalidValue)
	}
	if err := validateString(message.GroupID, "group_id", c.limits.MaxGroupID); err != nil {
		return err
	}
	output.uint8(uint8(message.ApplyResult))
	output.lve([]byte(message.GroupID))
	return nil
}

func (c *Codec) decodeAgentGroupInfoNotificationResponse(input *decoder, header Header) (Message, error) {
	applyResult, err := input.uint8("apply_result")
	if err != nil {
		return nil, err
	}
	groupID, err := input.lveString("group_id", c.limits.MaxGroupID)
	if err != nil {
		return nil, err
	}
	message := &AgentGroupInfoNotificationResponse{Header: header, ApplyResult: GroupInfoApplyResult(applyResult), GroupID: groupID}
	if err := c.encodeAgentGroupInfoNotificationResponse(&encoder{}, message); err != nil {
		return nil, err
	}
	return message, nil
}
