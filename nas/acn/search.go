package acn

import (
	"fmt"
	"math"
)

func (c *Codec) encodeAgentSearchRequest(output *encoder, message *AgentSearchRequest) error {
	if err := validateString(message.RequesterAgentID, "requester_agent_id", c.limits.MaxAgentID); err != nil {
		return err
	}
	if len(message.RequiredSkills) < 1 || len(message.RequiredSkills) > math.MaxUint8 {
		return newProtocolError(ErrorCodeInvalidLength, "required_skills", -1, ErrInvalidLength)
	}
	if message.DiscoveryScope != DiscoveryScopeIntraPLMN && message.DiscoveryScope != DiscoveryScopeInterPLMN {
		return newProtocolError(ErrorCodeInvalidValue, "discovery_scope", -1, ErrInvalidValue)
	}
	if message.MaxResults == 0 || message.MaxResults > 100 {
		return newProtocolError(ErrorCodeInvalidValue, "max_results", -1, ErrInvalidValue)
	}
	output.lve([]byte(message.RequesterAgentID))
	output.uint8(uint8(len(message.RequiredSkills)))
	for _, skill := range message.RequiredSkills {
		if err := validateString(skill, "required_skill", c.limits.MaxCapability); err != nil {
			return err
		}
		output.lv([]byte(skill))
	}
	output.uint8(uint8(message.DiscoveryScope))
	output.uint8(message.MaxResults)
	output.uint64(message.Timestamp)
	return c.encodeProof(output, message.Proof, "proof")
}

func (c *Codec) decodeAgentSearchRequest(input *decoder, header Header) (Message, error) {
	message := &AgentSearchRequest{Header: header}
	var err error
	if message.RequesterAgentID, err = input.lveString("requester_agent_id", c.limits.MaxAgentID); err != nil {
		return nil, err
	}
	count, err := input.uint8("required_skill_count")
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, newProtocolError(ErrorCodeInvalidLength, "required_skill_count", input.offset-1, ErrInvalidLength)
	}
	message.RequiredSkills = make([]string, 0, int(count))
	for range int(count) {
		skill, err := input.lvString("required_skill", c.limits.MaxCapability)
		if err != nil {
			return nil, err
		}
		message.RequiredSkills = append(message.RequiredSkills, skill)
	}
	scope, err := input.uint8("discovery_scope")
	if err != nil {
		return nil, err
	}
	message.DiscoveryScope = DiscoveryScope(scope)
	if message.DiscoveryScope != DiscoveryScopeIntraPLMN && message.DiscoveryScope != DiscoveryScopeInterPLMN {
		return nil, newProtocolError(ErrorCodeInvalidValue, "discovery_scope", input.offset-1, ErrInvalidValue)
	}
	if message.MaxResults, err = input.uint8("max_results"); err != nil {
		return nil, err
	}
	if message.MaxResults == 0 || message.MaxResults > 100 {
		return nil, newProtocolError(ErrorCodeInvalidValue, "max_results", input.offset-1, ErrInvalidValue)
	}
	if message.Timestamp, err = input.uint64("timestamp"); err != nil {
		return nil, err
	}
	if message.Proof, err = c.decodeProof(input, "proof"); err != nil {
		return nil, err
	}
	return message, nil
}

func (c *Codec) encodeAgentSearchResponse(output *encoder, message *AgentSearchResponse) error {
	if len(message.Results) > math.MaxUint8 {
		return newProtocolError(ErrorCodeInvalidLength, "results", -1, ErrInvalidLength)
	}
	output.uint8(uint8(len(message.Results)))
	for _, result := range message.Results {
		record, err := c.encodeSearchResult(result)
		if err != nil {
			return err
		}
		if len(record) > math.MaxUint16 {
			return newProtocolError(ErrorCodeInvalidLength, "search_result", -1, ErrInvalidLength)
		}
		output.uint16(uint16(len(record)))
		output.bytes(record)
	}
	output.uint64(message.Timestamp)
	return nil
}

func (c *Codec) encodeSearchResult(result SearchResult) ([]byte, error) {
	if err := validateJSONObject(result.AgentCard, "agent_card", c.limits.MaxJSONContainer); err != nil {
		return nil, err
	}
	if !validPriority(result.Priority) {
		return nil, newProtocolError(ErrorCodeInvalidValue, "agent_card.priority", -1, ErrInvalidValue)
	}
	record := &encoder{}
	record.lve(result.AgentCard)
	record.uint8(uint8(result.Priority))
	return record.data, nil
}

func (c *Codec) decodeAgentSearchResponse(input *decoder, header Header) (Message, error) {
	message := &AgentSearchResponse{Header: header}
	var err error
	count, err := input.uint8("result_count")
	if err != nil {
		return nil, err
	}
	message.Results = make([]SearchResult, 0, int(count))
	for range int(count) {
		length, err := input.uint16("record_length")
		if err != nil {
			return nil, err
		}
		record, err := input.bytes("search_result", int(length))
		if err != nil {
			return nil, err
		}
		result, err := c.decodeSearchResult(record)
		if err != nil {
			return nil, err
		}
		message.Results = append(message.Results, result)
	}
	if message.Timestamp, err = input.uint64("timestamp"); err != nil {
		return nil, err
	}
	return message, nil
}

func (c *Codec) decodeSearchResult(data []byte) (SearchResult, error) {
	input := &decoder{data: data}
	agentCard, err := input.lveBytes("agent_card", c.limits.MaxJSONContainer)
	if err != nil {
		return SearchResult{}, err
	}
	if err := validateJSONObject(agentCard, "agent_card", c.limits.MaxJSONContainer); err != nil {
		return SearchResult{}, err
	}
	priority, err := input.uint8("agent_card.priority")
	if err != nil {
		return SearchResult{}, err
	}
	result := SearchResult{AgentCard: cloneBytes(agentCard), Priority: Priority(priority)}
	if !validPriority(result.Priority) {
		return SearchResult{}, newProtocolError(ErrorCodeInvalidValue, "agent_card.priority", input.offset-1, ErrInvalidValue)
	}
	if err := input.finish("search_result"); err != nil {
		return SearchResult{}, fmt.Errorf("decode search result: %w", err)
	}
	return result, nil
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
