package acn

func (c *Codec) encodeAgentGroupingInvitation(
	output *encoder,
	message *AgentGroupingInvitation,
) error {
	if err := validateString(message.GroupID, "group_id", c.limits.MaxGroupID); err != nil {
		return err
	}
	if err := validateString(message.SourceAgentID, "source_agent_id", c.limits.MaxAgentID); err != nil {
		return err
	}
	if err := validateString(message.TargetAgentID, "target_agent_id", c.limits.MaxAgentID); err != nil {
		return err
	}
	if err := validateString(message.TaskID, "task_id", c.limits.MaxTaskID); err != nil {
		return err
	}
	if err := validateBytes(message.Proof, "proof", c.limits.MaxSignature); err != nil {
		return err
	}
	output.lve([]byte(message.GroupID))
	output.lve([]byte(message.SourceAgentID))
	output.lve([]byte(message.TargetAgentID))
	output.lve([]byte(message.TaskID))
	output.uint64(message.ExpiresAt)
	output.lve(message.Proof)
	return nil
}

func (c *Codec) decodeAgentGroupingInvitation(
	input *decoder,
	header Header,
) (Message, error) {
	message := &AgentGroupingInvitation{Header: header}
	var err error
	if message.GroupID, err = input.lveString("group_id", c.limits.MaxGroupID); err != nil {
		return nil, err
	}
	if message.SourceAgentID, err = input.lveString("source_agent_id", c.limits.MaxAgentID); err != nil {
		return nil, err
	}
	if message.TargetAgentID, err = input.lveString("target_agent_id", c.limits.MaxAgentID); err != nil {
		return nil, err
	}
	if message.TaskID, err = input.lveString("task_id", c.limits.MaxTaskID); err != nil {
		return nil, err
	}
	if message.ExpiresAt, err = input.uint64("expires_at"); err != nil {
		return nil, err
	}
	proof, err := input.lveBytes("proof", c.limits.MaxSignature)
	if err != nil {
		return nil, err
	}
	message.Proof = cloneBytes(proof)
	return message, nil
}

func (c *Codec) encodeAgentGroupingInvitationResponse(
	output *encoder,
	message *AgentGroupingInvitationResponse,
) error {
	if err := validateString(message.GroupID, "group_id", c.limits.MaxGroupID); err != nil {
		return err
	}
	if err := validateString(message.AgentID, "agent_id", c.limits.MaxAgentID); err != nil {
		return err
	}
	if !validGroupingDecision(message.Decision, message.RejectReason) {
		return newProtocolError(ErrorCodeInvalidValue, "decision", -1, ErrInvalidValue)
	}
	if err := validateBytes(message.Proof, "proof", c.limits.MaxSignature); err != nil {
		return err
	}
	output.lve([]byte(message.GroupID))
	output.lve([]byte(message.AgentID))
	output.uint8(uint8(message.Decision))
	output.uint8(uint8(message.RejectReason))
	output.uint64(message.Timestamp)
	output.lve(message.Proof)
	return nil
}

func (c *Codec) decodeAgentGroupingInvitationResponse(
	input *decoder,
	header Header,
) (Message, error) {
	message := &AgentGroupingInvitationResponse{Header: header}
	var err error
	if message.GroupID, err = input.lveString("group_id", c.limits.MaxGroupID); err != nil {
		return nil, err
	}
	if message.AgentID, err = input.lveString("agent_id", c.limits.MaxAgentID); err != nil {
		return nil, err
	}
	decision, err := input.uint8("decision")
	if err != nil {
		return nil, err
	}
	reason, err := input.uint8("reject_reason")
	if err != nil {
		return nil, err
	}
	message.Decision = GroupingDecision(decision)
	message.RejectReason = GroupingRejectReason(reason)
	if !validGroupingDecision(message.Decision, message.RejectReason) {
		return nil, newProtocolError(ErrorCodeInvalidValue, "decision", input.offset-2, ErrInvalidValue)
	}
	if message.Timestamp, err = input.uint64("timestamp"); err != nil {
		return nil, err
	}
	proof, err := input.lveBytes("proof", c.limits.MaxSignature)
	if err != nil {
		return nil, err
	}
	message.Proof = cloneBytes(proof)
	return message, nil
}

func validGroupingDecision(decision GroupingDecision, reason GroupingRejectReason) bool {
	if decision == GroupingDecisionAccept {
		return reason == GroupingRejectReasonNone
	}
	return decision == GroupingDecisionReject &&
		reason >= GroupingRejectReasonUserReject && reason <= GroupingRejectReasonLocalError
}

func (c *Codec) encodeAgentGroupInfoNotification(
	output *encoder,
	message *AgentGroupInfoNotification,
) error {
	if err := validateString(message.GroupID, "group_id", c.limits.MaxGroupID); err != nil {
		return err
	}
	if err := validateString(message.TargetAgentID, "target_agent_id", c.limits.MaxAgentID); err != nil {
		return err
	}
	if err := validateJSONObject(message.GroupConfig, "group_config", c.limits.MaxJSONContainer); err != nil {
		return err
	}
	output.lve([]byte(message.GroupID))
	output.lve([]byte(message.TargetAgentID))
	output.lve(message.GroupConfig)
	return nil
}

func (c *Codec) decodeAgentGroupInfoNotification(
	input *decoder,
	header Header,
) (Message, error) {
	message := &AgentGroupInfoNotification{Header: header}
	var err error
	if message.GroupID, err = input.lveString("group_id", c.limits.MaxGroupID); err != nil {
		return nil, err
	}
	if message.TargetAgentID, err = input.lveString("target_agent_id", c.limits.MaxAgentID); err != nil {
		return nil, err
	}
	groupConfig, err := input.lveBytes("group_config", c.limits.MaxJSONContainer)
	if err != nil {
		return nil, err
	}
	if err := validateJSONObject(groupConfig, "group_config", c.limits.MaxJSONContainer); err != nil {
		return nil, err
	}
	message.GroupConfig = cloneBytes(groupConfig)
	return message, nil
}

func (c *Codec) encodeAgentGroupInfoNotificationResponse(
	output *encoder,
	message *AgentGroupInfoNotificationResponse,
) error {
	if err := validateString(message.GroupID, "group_id", c.limits.MaxGroupID); err != nil {
		return err
	}
	if err := validateString(message.AgentID, "agent_id", c.limits.MaxAgentID); err != nil {
		return err
	}
	if !validGroupInfoResult(message.Result, message.FailureCause) {
		return newProtocolError(ErrorCodeInvalidValue, "result", -1, ErrInvalidValue)
	}
	output.lve([]byte(message.GroupID))
	output.lve([]byte(message.AgentID))
	output.uint8(uint8(message.Result))
	output.uint8(uint8(message.FailureCause))
	return nil
}

func (c *Codec) decodeAgentGroupInfoNotificationResponse(
	input *decoder,
	header Header,
) (Message, error) {
	message := &AgentGroupInfoNotificationResponse{Header: header}
	var err error
	if message.GroupID, err = input.lveString("group_id", c.limits.MaxGroupID); err != nil {
		return nil, err
	}
	if message.AgentID, err = input.lveString("agent_id", c.limits.MaxAgentID); err != nil {
		return nil, err
	}
	result, err := input.uint8("result")
	if err != nil {
		return nil, err
	}
	cause, err := input.uint8("failure_cause")
	if err != nil {
		return nil, err
	}
	message.Result = GroupInfoApplyResult(result)
	message.FailureCause = GroupInfoFailureCause(cause)
	if !validGroupInfoResult(message.Result, message.FailureCause) {
		return nil, newProtocolError(ErrorCodeInvalidValue, "result", input.offset-2, ErrInvalidValue)
	}
	return message, nil
}

func validGroupInfoResult(result GroupInfoApplyResult, cause GroupInfoFailureCause) bool {
	if result == GroupInfoApplyResultSuccess {
		return cause == GroupInfoFailureCauseNone
	}
	return result == GroupInfoApplyResultFailure &&
		cause >= GroupInfoFailureCauseConfigInvalid && cause <= GroupInfoFailureCauseLocalApplyFailed
}
