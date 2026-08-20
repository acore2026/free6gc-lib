package acn

import "math"

func (c *Codec) encodeAgentGroupingRequest(output *encoder, message *AgentGroupingRequest) error {
	if err := validateString(message.AgentID, "agent_id", c.limits.MaxAgentID); err != nil {
		return err
	}
	if len(message.TargetAgentIDs) < 1 || len(message.TargetAgentIDs) > math.MaxUint8 {
		return newProtocolError(ErrorCodeInvalidLength, "target_agents", -1, ErrInvalidLength)
	}
	output.lve([]byte(message.AgentID))
	output.uint8(uint8(len(message.TargetAgentIDs)))
	for _, agentID := range message.TargetAgentIDs {
		if err := validateString(agentID, "target_agent_id", c.limits.MaxAgentID); err != nil {
			return err
		}
		output.lve([]byte(agentID))
	}
	if err := validateJSONObject(message.GroupConfig, "group_config", c.limits.MaxJSONContainer); err != nil {
		return err
	}
	output.lve(message.GroupConfig)
	output.uint64(message.Timestamp)
	return c.encodeProof(output, message.Proof, "proof")
}

func (c *Codec) decodeAgentGroupingRequest(input *decoder, header Header) (Message, error) {
	message := &AgentGroupingRequest{Header: header}
	var err error
	if message.AgentID, err = input.lveString("agent_id", c.limits.MaxAgentID); err != nil {
		return nil, err
	}
	count, err := input.uint8("target_agent_count")
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, newProtocolError(ErrorCodeInvalidLength, "target_agent_count", input.offset-1, ErrInvalidLength)
	}
	message.TargetAgentIDs = make([]string, 0, int(count))
	for range int(count) {
		agentID, err := input.lveString("target_agent_id", c.limits.MaxAgentID)
		if err != nil {
			return nil, err
		}
		message.TargetAgentIDs = append(message.TargetAgentIDs, agentID)
	}
	groupConfig, err := input.lveBytes("group_config", c.limits.MaxJSONContainer)
	if err != nil {
		return nil, err
	}
	if err := validateJSONObject(groupConfig, "group_config", c.limits.MaxJSONContainer); err != nil {
		return nil, err
	}
	message.GroupConfig = cloneBytes(groupConfig)
	if message.Timestamp, err = input.uint64("timestamp"); err != nil {
		return nil, err
	}
	if message.Proof, err = c.decodeProof(input, "proof"); err != nil {
		return nil, err
	}
	return message, nil
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
	if err := validateJSONObject(message.GroupConfig, "group_config", c.limits.MaxJSONContainer); err != nil {
		return err
	}
	if err := validateJSONObject(message.GroupAdministrator, "group_administrator", c.limits.MaxJSONContainer); err != nil {
		return err
	}
	output.lve(message.GroupConfig)
	output.lve(message.GroupAdministrator)
	return nil
}

func (c *Codec) decodeAgentGroupingInvitation(input *decoder, header Header) (Message, error) {
	groupConfig, err := input.lveBytes("group_config", c.limits.MaxJSONContainer)
	if err != nil {
		return nil, err
	}
	if err := validateJSONObject(groupConfig, "group_config", c.limits.MaxJSONContainer); err != nil {
		return nil, err
	}
	groupAdministrator, err := input.lveBytes("group_administrator", c.limits.MaxJSONContainer)
	if err != nil {
		return nil, err
	}
	if err := validateJSONObject(groupAdministrator, "group_administrator", c.limits.MaxJSONContainer); err != nil {
		return nil, err
	}
	return &AgentGroupingInvitation{
		Header:             header,
		GroupConfig:        cloneBytes(groupConfig),
		GroupAdministrator: cloneBytes(groupAdministrator),
	}, nil
}

func encodeAgentGroupingInvitationResponse(output *encoder, message *AgentGroupingInvitationResponse) error {
	if message.Decision != GroupingDecisionAccept && message.Decision != GroupingDecisionReject {
		return newProtocolError(ErrorCodeInvalidValue, "decision", -1, ErrInvalidValue)
	}
	output.uint8(uint8(message.Decision))
	return nil
}

func decodeAgentGroupingInvitationResponse(input *decoder, header Header) (Message, error) {
	decision, err := input.uint8("decision")
	if err != nil {
		return nil, err
	}
	message := &AgentGroupingInvitationResponse{Header: header, Decision: GroupingDecision(decision)}
	if err := encodeAgentGroupingInvitationResponse(&encoder{}, message); err != nil {
		return nil, err
	}
	return message, nil
}

func (c *Codec) encodeAgentGroupInfoNotification(output *encoder, message *AgentGroupInfoNotification) error {
	if err := validateString(message.Version, "version", c.limits.MaxSoftwareVersion); err != nil {
		return err
	}
	if err := validateString(message.GroupID, "group_id", c.limits.MaxGroupID); err != nil {
		return err
	}
	if err := validateJSONObject(message.Members, "members", c.limits.MaxJSONContainer); err != nil {
		return err
	}
	output.lve([]byte(message.Version))
	output.uint64(message.Timestamp)
	output.lve([]byte(message.GroupID))
	output.lve(message.Members)
	return c.encodeProof(output, message.Proof, "proof")
}

func (c *Codec) decodeAgentGroupInfoNotification(input *decoder, header Header) (Message, error) {
	message := &AgentGroupInfoNotification{Header: header}
	var err error
	if message.Version, err = input.lveString("version", c.limits.MaxSoftwareVersion); err != nil {
		return nil, err
	}
	if message.Timestamp, err = input.uint64("timestamp"); err != nil {
		return nil, err
	}
	if message.GroupID, err = input.lveString("group_id", c.limits.MaxGroupID); err != nil {
		return nil, err
	}
	members, err := input.lveBytes("members", c.limits.MaxJSONContainer)
	if err != nil {
		return nil, err
	}
	if err := validateJSONObject(members, "members", c.limits.MaxJSONContainer); err != nil {
		return nil, err
	}
	message.Members = cloneBytes(members)
	if message.Proof, err = c.decodeProof(input, "proof"); err != nil {
		return nil, err
	}
	return message, nil
}

func (c *Codec) encodeAgentGroupInfoNotificationResponse(output *encoder, message *AgentGroupInfoNotificationResponse) error {
	if err := validateString(message.GroupID, "group_id", c.limits.MaxGroupID); err != nil {
		return err
	}
	if err := validateString(message.AgentID, "agent_id", c.limits.MaxAgentID); err != nil {
		return err
	}
	if message.Status != GroupInfoStatusAccepted && message.Status != GroupInfoStatusRejected {
		return newProtocolError(ErrorCodeInvalidValue, "status", -1, ErrInvalidValue)
	}
	if message.Detail != "" {
		if err := validateString(message.Detail, "detail", c.limits.MaxDescription); err != nil {
			return err
		}
	}
	output.lve([]byte(message.GroupID))
	output.lve([]byte(message.AgentID))
	output.uint8(uint8(message.Status))
	output.lve([]byte(message.Detail))
	output.uint64(message.Timestamp)
	return c.encodeProof(output, message.Proof, "proof")
}

func (c *Codec) decodeAgentGroupInfoNotificationResponse(input *decoder, header Header) (Message, error) {
	message := &AgentGroupInfoNotificationResponse{Header: header}
	var err error
	if message.GroupID, err = input.lveString("group_id", c.limits.MaxGroupID); err != nil {
		return nil, err
	}
	if message.AgentID, err = input.lveString("agent_id", c.limits.MaxAgentID); err != nil {
		return nil, err
	}
	status, err := input.uint8("status")
	if err != nil {
		return nil, err
	}
	message.Status = GroupInfoStatus(status)
	if message.Status != GroupInfoStatusAccepted && message.Status != GroupInfoStatusRejected {
		return nil, newProtocolError(ErrorCodeInvalidValue, "status", input.offset-1, ErrInvalidValue)
	}
	if message.Detail, err = input.lveOptionalString("detail", c.limits.MaxDescription); err != nil {
		return nil, err
	}
	if message.Timestamp, err = input.uint64("timestamp"); err != nil {
		return nil, err
	}
	if message.Proof, err = c.decodeProof(input, "proof"); err != nil {
		return nil, err
	}
	return message, nil
}
