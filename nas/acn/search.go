package acn

import (
	"fmt"
	"math"
)

const (
	agentRecordNamePresent         uint8 = 1 << 0
	agentRecordDescriptionPresent  uint8 = 1 << 1
	agentRecordStatusPresent       uint8 = 1 << 2
	agentRecordPriorityPresent     uint8 = 1 << 3
	agentRecordCapabilitiesPresent uint8 = 1 << 4
	agentRecordReservedMask        uint8 = 0xe0

	minEncodedAgentRecordLength = 4
	minFramedAgentRecordLength  = 2 + minEncodedAgentRecordLength
)

func (c *Codec) encodeAgentSearchRequest(
	output *encoder,
	message *AgentSearchRequest,
) error {
	if typedNil(message.Query) {
		return newProtocolError(ErrorCodeInvalidValue, "query", -1, ErrInvalidValue)
	}

	switch query := message.Query.(type) {
	case *CapabilityDiscoveryQuery:
		if query == nil {
			return newProtocolError(ErrorCodeInvalidValue, "query", -1, ErrInvalidValue)
		}
		if err := validateString(
			query.SourceAgentID,
			"source_agent_id",
			c.limits.MaxAgentID,
		); err != nil {
			return err
		}
		if err := validateString(query.TaskID, "task_id", c.limits.MaxTaskID); err != nil {
			return err
		}
		if len(query.RequiredCapabilities) < 1 || len(query.RequiredCapabilities) > math.MaxUint8 {
			return newProtocolError(
				ErrorCodeInvalidLength,
				"required_capabilities",
				-1,
				ErrInvalidLength,
			)
		}
		for _, capability := range query.RequiredCapabilities {
			if err := validateString(
				capability,
				"required_capability",
				c.limits.MaxCapability,
			); err != nil {
				return err
			}
		}
		output.uint8(uint8(SearchTypeCapabilityDiscovery))
		output.lve([]byte(query.SourceAgentID))
		output.lve([]byte(query.TaskID))
		output.uint64(query.Timestamp)
		output.uint8(uint8(len(query.RequiredCapabilities)))
		for _, capability := range query.RequiredCapabilities {
			output.lv([]byte(capability))
		}
	case *AgentInfoQuery:
		if query == nil {
			return newProtocolError(ErrorCodeInvalidValue, "query", -1, ErrInvalidValue)
		}
		if err := validateString(
			query.TargetAgentID,
			"target_agent_id",
			c.limits.MaxAgentID,
		); err != nil {
			return err
		}
		output.uint8(uint8(SearchTypeAgentInfo))
		output.lve([]byte(query.TargetAgentID))
	case *OwnerAgentsQuery:
		if query == nil {
			return newProtocolError(ErrorCodeInvalidValue, "query", -1, ErrInvalidValue)
		}
		if err := validateString(query.OwnerID, "owner_id", c.limits.MaxOwnerID); err != nil {
			return err
		}
		output.uint8(uint8(SearchTypeOwnerAgents))
		output.lve([]byte(query.OwnerID))
	default:
		return newProtocolError(ErrorCodeInvalidValue, "query", -1, ErrInvalidValue)
	}
	return nil
}

func (c *Codec) decodeAgentSearchRequest(
	input *decoder,
	header Header,
) (Message, error) {
	rawSearchType, err := input.uint8("search_type")
	if err != nil {
		return nil, err
	}
	message := &AgentSearchRequest{Header: header}

	switch SearchType(rawSearchType) {
	case SearchTypeCapabilityDiscovery:
		query := &CapabilityDiscoveryQuery{}
		if query.SourceAgentID, err = input.lveString(
			"source_agent_id",
			c.limits.MaxAgentID,
		); err != nil {
			return nil, err
		}
		if query.TaskID, err = input.lveString("task_id", c.limits.MaxTaskID); err != nil {
			return nil, err
		}
		if query.Timestamp, err = input.uint64("timestamp"); err != nil {
			return nil, err
		}
		count, err := input.uint8("capability_count")
		if err != nil {
			return nil, err
		}
		if count == 0 {
			return nil, newProtocolError(
				ErrorCodeInvalidLength,
				"capability_count",
				input.offset-1,
				ErrInvalidLength,
			)
		}
		query.RequiredCapabilities = make([]string, 0, int(count))
		for range int(count) {
			capability, err := input.lvString("required_capability", c.limits.MaxCapability)
			if err != nil {
				return nil, err
			}
			query.RequiredCapabilities = append(query.RequiredCapabilities, capability)
		}
		message.Query = query
	case SearchTypeAgentInfo:
		targetAgentID, err := input.lveString("target_agent_id", c.limits.MaxAgentID)
		if err != nil {
			return nil, err
		}
		message.Query = &AgentInfoQuery{TargetAgentID: targetAgentID}
	case SearchTypeOwnerAgents:
		ownerID, err := input.lveString("owner_id", c.limits.MaxOwnerID)
		if err != nil {
			return nil, err
		}
		message.Query = &OwnerAgentsQuery{OwnerID: ownerID}
	default:
		return nil, newProtocolError(
			ErrorCodeInvalidValue,
			"search_type",
			input.offset-1,
			ErrInvalidValue,
		)
	}
	return message, nil
}

func (c *Codec) encodeAgentSearchResponse(
	output *encoder,
	message *AgentSearchResponse,
) error {
	if !validSearchType(message.SearchType) {
		return newProtocolError(ErrorCodeInvalidValue, "search_type", -1, ErrInvalidValue)
	}
	if message.SearchType == SearchTypeAgentInfo && len(message.Records) > 1 {
		return newProtocolError(ErrorCodeCountMismatch, "records", -1, ErrCountMismatch)
	}
	if len(message.Records) > math.MaxUint16 {
		return newProtocolError(ErrorCodeInvalidLength, "records", -1, ErrInvalidLength)
	}

	output.uint8(uint8(message.SearchType))
	output.uint16(uint16(len(message.Records)))
	for index := range message.Records {
		if len(output.data) > c.limits.MaxACNPayload-minFramedAgentRecordLength {
			return newProtocolError(
				ErrorCodePayloadTooLarge,
				"records",
				-1,
				ErrPayloadTooLarge,
			)
		}
		recordBytes, err := c.encodeAgentRecord(&message.Records[index])
		if err != nil {
			return err
		}
		if len(recordBytes) > c.limits.MaxACNPayload-len(output.data)-2 {
			return newProtocolError(
				ErrorCodePayloadTooLarge,
				"agent_record",
				-1,
				ErrPayloadTooLarge,
			)
		}
		output.uint16(uint16(len(recordBytes)))
		output.bytes(recordBytes)
	}
	return nil
}

func (c *Codec) encodeAgentRecord(record *AgentRecord) ([]byte, error) {
	if err := validateString(record.AgentID, "agent_record.agent_id", c.limits.MaxAgentID); err != nil {
		return nil, err
	}
	presence := uint8(0)
	if record.AgentName != nil {
		if err := validateString(
			*record.AgentName,
			"agent_record.agent_name",
			c.limits.MaxAgentName,
		); err != nil {
			return nil, err
		}
		presence |= agentRecordNamePresent
	}
	if record.Description != nil {
		if err := validateString(
			*record.Description,
			"agent_record.description",
			c.limits.MaxDescription,
		); err != nil {
			return nil, err
		}
		presence |= agentRecordDescriptionPresent
	}
	if record.Status != nil {
		if !validAgentStatus(*record.Status) {
			return nil, newProtocolError(
				ErrorCodeInvalidValue,
				"agent_record.status",
				-1,
				ErrInvalidValue,
			)
		}
		presence |= agentRecordStatusPresent
	}
	if record.Priority != nil {
		if !validPriority(*record.Priority) {
			return nil, newProtocolError(
				ErrorCodeInvalidValue,
				"agent_record.priority",
				-1,
				ErrInvalidValue,
			)
		}
		presence |= agentRecordPriorityPresent
	}
	if record.Capabilities != nil {
		if len(record.Capabilities) < 1 || len(record.Capabilities) > math.MaxUint8 {
			return nil, newProtocolError(
				ErrorCodeInvalidLength,
				"agent_record.capabilities",
				-1,
				ErrInvalidLength,
			)
		}
		for _, capability := range record.Capabilities {
			if err := validateString(
				capability,
				"agent_record.capability",
				c.limits.MaxCapability,
			); err != nil {
				return nil, err
			}
		}
		presence |= agentRecordCapabilitiesPresent
	}

	output := &encoder{data: make([]byte, 0, 64)}
	output.lve([]byte(record.AgentID))
	output.uint8(presence)
	if record.AgentName != nil {
		output.lve([]byte(*record.AgentName))
	}
	if record.Description != nil {
		output.lve([]byte(*record.Description))
	}
	if record.Status != nil {
		output.uint8(uint8(*record.Status))
	}
	if record.Priority != nil {
		output.uint8(uint8(*record.Priority))
	}
	if record.Capabilities != nil {
		output.uint8(uint8(len(record.Capabilities)))
		for _, capability := range record.Capabilities {
			output.lv([]byte(capability))
		}
	}
	if len(output.data) > math.MaxUint16 {
		return nil, newProtocolError(
			ErrorCodeInvalidLength,
			"agent_record",
			-1,
			fmt.Errorf("%w: record length %d", ErrInvalidLength, len(output.data)),
		)
	}
	return output.data, nil
}

func (c *Codec) decodeAgentSearchResponse(
	input *decoder,
	header Header,
) (Message, error) {
	rawSearchType, err := input.uint8("search_type")
	if err != nil {
		return nil, err
	}
	searchType := SearchType(rawSearchType)
	if !validSearchType(searchType) {
		return nil, newProtocolError(ErrorCodeInvalidValue, "search_type", input.offset-1, ErrInvalidValue)
	}
	count, err := input.uint16("result_count")
	if err != nil {
		return nil, err
	}
	if searchType == SearchTypeAgentInfo && count > 1 {
		return nil, newProtocolError(ErrorCodeCountMismatch, "result_count", input.offset-2, ErrCountMismatch)
	}
	if int(count) > input.remaining()/minFramedAgentRecordLength {
		return nil, newProtocolError(
			ErrorCodeCountMismatch,
			"result_count",
			input.offset-2,
			ErrCountMismatch,
		)
	}

	message := &AgentSearchResponse{
		Header:     header,
		SearchType: searchType,
		Records:    make([]AgentRecord, 0, int(count)),
	}
	for range int(count) {
		recordLength, err := input.uint16("record_length")
		if err != nil {
			return nil, err
		}
		if recordLength < minEncodedAgentRecordLength {
			return nil, newProtocolError(
				ErrorCodeInvalidLength,
				"record_length",
				input.offset-2,
				ErrInvalidLength,
			)
		}
		recordBytes, err := input.bytes("agent_record", int(recordLength))
		if err != nil {
			return nil, err
		}
		record, err := c.decodeAgentRecord(recordBytes)
		if err != nil {
			return nil, err
		}
		message.Records = append(message.Records, record)
	}
	return message, nil
}

func (c *Codec) decodeAgentRecord(data []byte) (AgentRecord, error) {
	input := &decoder{data: data}
	record := AgentRecord{}
	var err error
	if record.AgentID, err = input.lveString(
		"agent_record.agent_id",
		c.limits.MaxAgentID,
	); err != nil {
		return AgentRecord{}, err
	}
	presence, err := input.uint8("agent_record.presence")
	if err != nil {
		return AgentRecord{}, err
	}
	if presence&agentRecordReservedMask != 0 {
		return AgentRecord{}, newProtocolError(
			ErrorCodeInvalidPresenceBitmap,
			"agent_record.presence",
			input.offset-1,
			ErrInvalidPresenceBitmap,
		)
	}
	if presence&agentRecordNamePresent != 0 {
		value, err := input.lveString("agent_record.agent_name", c.limits.MaxAgentName)
		if err != nil {
			return AgentRecord{}, err
		}
		record.AgentName = &value
	}
	if presence&agentRecordDescriptionPresent != 0 {
		value, err := input.lveString("agent_record.description", c.limits.MaxDescription)
		if err != nil {
			return AgentRecord{}, err
		}
		record.Description = &value
	}
	if presence&agentRecordStatusPresent != 0 {
		value, err := input.uint8("agent_record.status")
		if err != nil {
			return AgentRecord{}, err
		}
		status := AgentStatus(value)
		if !validAgentStatus(status) {
			return AgentRecord{}, newProtocolError(
				ErrorCodeInvalidValue,
				"agent_record.status",
				input.offset-1,
				ErrInvalidValue,
			)
		}
		record.Status = &status
	}
	if presence&agentRecordPriorityPresent != 0 {
		value, err := input.uint8("agent_record.priority")
		if err != nil {
			return AgentRecord{}, err
		}
		priority := Priority(value)
		if !validPriority(priority) {
			return AgentRecord{}, newProtocolError(
				ErrorCodeInvalidValue,
				"agent_record.priority",
				input.offset-1,
				ErrInvalidValue,
			)
		}
		record.Priority = &priority
	}
	if presence&agentRecordCapabilitiesPresent != 0 {
		count, err := input.uint8("agent_record.capability_count")
		if err != nil {
			return AgentRecord{}, err
		}
		if count == 0 {
			return AgentRecord{}, newProtocolError(
				ErrorCodeInvalidLength,
				"agent_record.capability_count",
				input.offset-1,
				ErrInvalidLength,
			)
		}
		record.Capabilities = make([]string, 0, int(count))
		for range int(count) {
			value, err := input.lvString(
				"agent_record.capability",
				c.limits.MaxCapability,
			)
			if err != nil {
				return AgentRecord{}, err
			}
			record.Capabilities = append(record.Capabilities, value)
		}
	}
	if err := input.finish("agent_record"); err != nil {
		return AgentRecord{}, err
	}
	return record, nil
}

func validSearchType(value SearchType) bool {
	return value >= SearchTypeCapabilityDiscovery && value <= SearchTypeOwnerAgents
}

func validAgentStatus(value AgentStatus) bool {
	return value <= AgentStatusBusy
}

func encodeAgentSearchReject(output *encoder, message *AgentSearchReject) error {
	if !validSearchRejectCause(message.Cause) {
		return newProtocolError(ErrorCodeInvalidValue, "reject_cause", -1, ErrInvalidValue)
	}
	if !validSearchFailedField(message.FailedField) {
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
	if !validSearchRejectCause(SearchRejectCause(cause)) {
		return nil, newProtocolError(ErrorCodeInvalidValue, "reject_cause", input.offset-1, ErrInvalidValue)
	}
	failedField, err := input.uint8("failed_field")
	if err != nil {
		return nil, err
	}
	if !validSearchFailedField(SearchFailedField(failedField)) {
		return nil, newProtocolError(ErrorCodeInvalidValue, "failed_field", input.offset-1, ErrInvalidValue)
	}
	return &AgentSearchReject{
		Header:      header,
		Cause:       SearchRejectCause(cause),
		FailedField: SearchFailedField(failedField),
	}, nil
}

func validSearchRejectCause(value SearchRejectCause) bool {
	return value >= SearchRejectInvalidSearchType && value <= SearchRejectInternalError
}

func validSearchFailedField(value SearchFailedField) bool {
	return value >= SearchFieldUnspecified && value <= SearchFieldTimestamp
}
