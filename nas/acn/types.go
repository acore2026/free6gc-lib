package acn

import "encoding/json"

const (
	// Version1 is the only ACN protocol version supported by this package.
	Version1 uint8 = 0x01

	// PayloadContainerTypeACN is a private prototype allocation, not a 3GPP
	// payload-container type.
	PayloadContainerTypeACN uint8 = 0x0e
)

// Direction identifies the NAS transport direction allowed for an ACN message.
type Direction uint8

const (
	Uplink Direction = iota + 1
	Downlink
)

func (d Direction) String() string {
	switch d {
	case Uplink:
		return "uplink"
	case Downlink:
		return "downlink"
	default:
		return "unknown"
	}
}

// MessageType identifies an ACN control PDU.
type MessageType uint8

const (
	MessageTypeAgentRegisterRequest               MessageType = 0x01
	MessageTypeAgentRegisterAccept                MessageType = 0x02
	MessageTypeAgentRegisterReject                MessageType = 0x03
	MessageTypeAgentDeregisterRequest             MessageType = 0x04
	MessageTypeAgentDeregisterAccept              MessageType = 0x05
	MessageTypeAgentProfileUpdateRequest          MessageType = 0x06
	MessageTypeAgentProfileUpdateAccept           MessageType = 0x07
	MessageTypeAgentProfileUpdateReject           MessageType = 0x08
	MessageTypeAgentSearchRequest                 MessageType = 0x09
	MessageTypeAgentSearchResponse                MessageType = 0x0a
	MessageTypeAgentSearchReject                  MessageType = 0x0b
	MessageTypeAgentDeregisterReject              MessageType = 0x0c
	MessageTypeAgentGroupingInvitation            MessageType = 0x10
	MessageTypeAgentGroupingInvitationResponse    MessageType = 0x11
	MessageTypeAgentGroupInfoNotification         MessageType = 0x12
	MessageTypeAgentGroupInfoNotificationResponse MessageType = 0x13
)

func (t MessageType) String() string {
	switch t {
	case MessageTypeAgentRegisterRequest:
		return "ACN_AGENT_REGISTER_REQUEST"
	case MessageTypeAgentRegisterAccept:
		return "ACN_AGENT_REGISTER_ACCEPT"
	case MessageTypeAgentRegisterReject:
		return "ACN_AGENT_REGISTER_REJECT"
	case MessageTypeAgentDeregisterRequest:
		return "ACN_AGENT_DEREGISTER_REQUEST"
	case MessageTypeAgentDeregisterAccept:
		return "ACN_AGENT_DEREGISTER_ACCEPT"
	case MessageTypeAgentProfileUpdateRequest:
		return "ACN_AGENT_PROFILE_UPDATE_REQUEST"
	case MessageTypeAgentProfileUpdateAccept:
		return "ACN_AGENT_PROFILE_UPDATE_ACCEPT"
	case MessageTypeAgentProfileUpdateReject:
		return "ACN_AGENT_PROFILE_UPDATE_REJECT"
	case MessageTypeAgentSearchRequest:
		return "ACN_AGENT_SEARCH_REQUEST"
	case MessageTypeAgentSearchResponse:
		return "ACN_AGENT_SEARCH_RESPONSE"
	case MessageTypeAgentSearchReject:
		return "ACN_AGENT_SEARCH_REJECT"
	case MessageTypeAgentDeregisterReject:
		return "ACN_AGENT_DEREGISTER_REJECT"
	case MessageTypeAgentGroupingInvitation:
		return "ACN_AGENT_GROUPING_INVITATION"
	case MessageTypeAgentGroupingInvitationResponse:
		return "ACN_AGENT_GROUPING_INVITATION_RESPONSE"
	case MessageTypeAgentGroupInfoNotification:
		return "ACN_AGENT_GROUPINFO_NOTIFICATION"
	case MessageTypeAgentGroupInfoNotificationResponse:
		return "ACN_AGENT_GROUPINFO_NOTIFICATION_RESPONSE"
	default:
		return "ACN_MESSAGE_UNKNOWN"
	}
}

func (t MessageType) direction() (Direction, bool) {
	switch t {
	case MessageTypeAgentRegisterRequest,
		MessageTypeAgentDeregisterRequest,
		MessageTypeAgentProfileUpdateRequest,
		MessageTypeAgentSearchRequest,
		MessageTypeAgentGroupingInvitationResponse,
		MessageTypeAgentGroupInfoNotificationResponse:
		return Uplink, true
	case MessageTypeAgentRegisterAccept,
		MessageTypeAgentRegisterReject,
		MessageTypeAgentDeregisterAccept,
		MessageTypeAgentDeregisterReject,
		MessageTypeAgentProfileUpdateAccept,
		MessageTypeAgentProfileUpdateReject,
		MessageTypeAgentSearchResponse,
		MessageTypeAgentSearchReject,
		MessageTypeAgentGroupingInvitation,
		MessageTypeAgentGroupInfoNotification:
		return Downlink, true
	default:
		return 0, false
	}
}

// Header is common to every decoded and encoded ACN message. Version and
// message type are derived from the concrete Go type.
type Header struct {
	TransactionID uint8
}

func (h Header) GetHeader() Header {
	return h
}

type messageMarker struct{}

func (messageMarker) isACNMessage() {}

// Message is the closed set of ACN version 1 messages.
type Message interface {
	GetHeader() Header
	MessageType() MessageType
	Direction() Direction
	isACNMessage()
}

type AgentRegisterRequest struct {
	messageMarker
	Header
	Owner           string
	AgentName       string
	PublicKey       []byte
	Description     string
	Timestamp       uint64
	Signature       []byte
	Region          string
	OS              string
	SoftwareVersion string
}

func (*AgentRegisterRequest) MessageType() MessageType { return MessageTypeAgentRegisterRequest }
func (*AgentRegisterRequest) Direction() Direction     { return Uplink }

type AgentRegisterAccept struct {
	messageMarker
	Header
	AgentID string
	VC0     json.RawMessage
}

func (*AgentRegisterAccept) MessageType() MessageType { return MessageTypeAgentRegisterAccept }
func (*AgentRegisterAccept) Direction() Direction     { return Downlink }

type RegisterRejectCause uint8

const (
	RegisterRejectInvalidMandatoryField RegisterRejectCause = 0x01
	RegisterRejectOwnerNotAllowed       RegisterRejectCause = 0x02
	RegisterRejectInvalidPublicKey      RegisterRejectCause = 0x03
	RegisterRejectInvalidSignature      RegisterRejectCause = 0x04
	RegisterRejectInvalidTimestamp      RegisterRejectCause = 0x05
	RegisterRejectAgentIDAllocation     RegisterRejectCause = 0x06
	RegisterRejectInternalError         RegisterRejectCause = 0x07
)

type RegisterFailedField uint8

const (
	RegisterFieldUnspecified     RegisterFailedField = 0x00
	RegisterFieldOwner           RegisterFailedField = 0x01
	RegisterFieldAgentName       RegisterFailedField = 0x02
	RegisterFieldPublicKey       RegisterFailedField = 0x03
	RegisterFieldDescription     RegisterFailedField = 0x04
	RegisterFieldTimestamp       RegisterFailedField = 0x05
	RegisterFieldSignature       RegisterFailedField = 0x06
	RegisterFieldRegion          RegisterFailedField = 0x07
	RegisterFieldOS              RegisterFailedField = 0x08
	RegisterFieldSoftwareVersion RegisterFailedField = 0x09
)

type AgentRegisterReject struct {
	messageMarker
	Header
	Cause       RegisterRejectCause
	FailedField RegisterFailedField
}

func (*AgentRegisterReject) MessageType() MessageType { return MessageTypeAgentRegisterReject }
func (*AgentRegisterReject) Direction() Direction     { return Downlink }

type DeregistrationReason uint8

const (
	DeregistrationReasonNormal        DeregistrationReason = 0x00
	DeregistrationReasonUninstalled   DeregistrationReason = 0x01
	DeregistrationReasonReplaced      DeregistrationReason = 0x02
	DeregistrationReasonUserRequest   DeregistrationReason = 0x03
	DeregistrationReasonSecurityEvent DeregistrationReason = 0x04
	DeregistrationReasonRetired       DeregistrationReason = 0x05
	DeregistrationReasonOther         DeregistrationReason = 0xff
)

type AgentDeregisterRequest struct {
	messageMarker
	Header
	AgentID   string
	Reason    DeregistrationReason
	Timestamp uint64
	Signature []byte
}

func (*AgentDeregisterRequest) MessageType() MessageType {
	return MessageTypeAgentDeregisterRequest
}
func (*AgentDeregisterRequest) Direction() Direction { return Uplink }

type AgentDeregisterAccept struct {
	messageMarker
	Header
}

func (*AgentDeregisterAccept) MessageType() MessageType {
	return MessageTypeAgentDeregisterAccept
}
func (*AgentDeregisterAccept) Direction() Direction { return Downlink }

type DeregisterRejectCause uint8

const (
	DeregisterRejectUnknownAgentID   DeregisterRejectCause = 0x01
	DeregisterRejectAgentNotBound    DeregisterRejectCause = 0x02
	DeregisterRejectNotAuthorized    DeregisterRejectCause = 0x03
	DeregisterRejectInvalidSignature DeregisterRejectCause = 0x04
	DeregisterRejectInvalidTimestamp DeregisterRejectCause = 0x05
	DeregisterRejectBackendFailure   DeregisterRejectCause = 0x06
)

type DeregisterFailedField uint8

const (
	DeregisterFieldUnspecified    DeregisterFailedField = 0x00
	DeregisterFieldAgentID        DeregisterFailedField = 0x01
	DeregisterFieldReason         DeregisterFailedField = 0x02
	DeregisterFieldTimestamp      DeregisterFailedField = 0x03
	DeregisterFieldSignature      DeregisterFailedField = 0x04
	DeregisterFieldUEAgentBinding DeregisterFailedField = 0x05
)

type AgentDeregisterReject struct {
	messageMarker
	Header
	Cause       DeregisterRejectCause
	FailedField DeregisterFailedField
}

func (*AgentDeregisterReject) MessageType() MessageType {
	return MessageTypeAgentDeregisterReject
}
func (*AgentDeregisterReject) Direction() Direction { return Downlink }

type Priority uint8

const (
	PriorityUnspecified Priority = 0x00
	PriorityHigh        Priority = 0x01
	PriorityNormal      Priority = 0x02
	PriorityLow         Priority = 0x03
)

type AgentProfileUpdateRequest struct {
	messageMarker
	Header
	AgentID   string
	Priority  Priority
	Timestamp uint64
	Signature []byte
	VCList    json.RawMessage
}

func (*AgentProfileUpdateRequest) MessageType() MessageType {
	return MessageTypeAgentProfileUpdateRequest
}
func (*AgentProfileUpdateRequest) Direction() Direction { return Uplink }

type AgentProfileUpdateAccept struct {
	messageMarker
	Header
}

func (*AgentProfileUpdateAccept) MessageType() MessageType {
	return MessageTypeAgentProfileUpdateAccept
}
func (*AgentProfileUpdateAccept) Direction() Direction { return Downlink }

type ProfileUpdateRejectCause uint8

const (
	ProfileUpdateRejectInvalidAgentID   ProfileUpdateRejectCause = 0x01
	ProfileUpdateRejectNotAuthorized    ProfileUpdateRejectCause = 0x02
	ProfileUpdateRejectInvalidPriority  ProfileUpdateRejectCause = 0x03
	ProfileUpdateRejectInvalidVCList    ProfileUpdateRejectCause = 0x04
	ProfileUpdateRejectInvalidSignature ProfileUpdateRejectCause = 0x05
	ProfileUpdateRejectInvalidTimestamp ProfileUpdateRejectCause = 0x06
	ProfileUpdateRejectInternalError    ProfileUpdateRejectCause = 0x07
)

type ProfileUpdateFailedField uint8

const (
	ProfileUpdateFieldUnspecified ProfileUpdateFailedField = 0x00
	ProfileUpdateFieldAgentID     ProfileUpdateFailedField = 0x01
	ProfileUpdateFieldPriority    ProfileUpdateFailedField = 0x02
	ProfileUpdateFieldTimestamp   ProfileUpdateFailedField = 0x03
	ProfileUpdateFieldSignature   ProfileUpdateFailedField = 0x04
	ProfileUpdateFieldVCList      ProfileUpdateFailedField = 0x05
)

type AgentProfileUpdateReject struct {
	messageMarker
	Header
	Cause       ProfileUpdateRejectCause
	FailedField ProfileUpdateFailedField
}

func (*AgentProfileUpdateReject) MessageType() MessageType {
	return MessageTypeAgentProfileUpdateReject
}
func (*AgentProfileUpdateReject) Direction() Direction { return Downlink }

type SearchType uint8

const (
	SearchTypeCapabilityDiscovery SearchType = 0x01
	SearchTypeAgentInfo           SearchType = 0x02
	SearchTypeOwnerAgents         SearchType = 0x03
)

type searchQueryMarker struct{}

func (searchQueryMarker) isSearchQuery() {}

type SearchQuery interface {
	SearchType() SearchType
	isSearchQuery()
}

type CapabilityDiscoveryQuery struct {
	searchQueryMarker
	SourceAgentID        string
	TaskID               string
	Timestamp            uint64
	RequiredCapabilities []string
}

func (*CapabilityDiscoveryQuery) SearchType() SearchType {
	return SearchTypeCapabilityDiscovery
}

type AgentInfoQuery struct {
	searchQueryMarker
	TargetAgentID string
}

func (*AgentInfoQuery) SearchType() SearchType { return SearchTypeAgentInfo }

type OwnerAgentsQuery struct {
	searchQueryMarker
	OwnerID string
}

func (*OwnerAgentsQuery) SearchType() SearchType { return SearchTypeOwnerAgents }

type AgentSearchRequest struct {
	messageMarker
	Header
	Query SearchQuery
}

func (*AgentSearchRequest) MessageType() MessageType { return MessageTypeAgentSearchRequest }
func (*AgentSearchRequest) Direction() Direction     { return Uplink }

type AgentStatus uint8

const (
	AgentStatusUnknown AgentStatus = 0x00
	AgentStatusOffline AgentStatus = 0x01
	AgentStatusOnline  AgentStatus = 0x02
	AgentStatusBusy    AgentStatus = 0x03
)

type AgentRecord struct {
	AgentID      string
	AgentName    *string
	Description  *string
	Status       *AgentStatus
	Priority     *Priority
	Capabilities []string
}

type AgentSearchResponse struct {
	messageMarker
	Header
	SearchType SearchType
	Records    []AgentRecord
}

func (*AgentSearchResponse) MessageType() MessageType { return MessageTypeAgentSearchResponse }
func (*AgentSearchResponse) Direction() Direction     { return Downlink }

type SearchRejectCause uint8

const (
	SearchRejectInvalidSearchType SearchRejectCause = 0x01
	SearchRejectInvalidMandatory  SearchRejectCause = 0x02
	SearchRejectNotAuthorized     SearchRejectCause = 0x03
	SearchRejectInvalidSource     SearchRejectCause = 0x04
	SearchRejectInvalidTarget     SearchRejectCause = 0x05
	SearchRejectInternalError     SearchRejectCause = 0x06
)

type SearchFailedField uint8

const (
	SearchFieldUnspecified        SearchFailedField = 0x00
	SearchFieldSourceAgentID      SearchFailedField = 0x01
	SearchFieldTaskID             SearchFailedField = 0x02
	SearchFieldTargetAgentID      SearchFailedField = 0x03
	SearchFieldOwnerID            SearchFailedField = 0x04
	SearchFieldRequiredCapability SearchFailedField = 0x05
	SearchFieldTimestamp          SearchFailedField = 0x06
)

type AgentSearchReject struct {
	messageMarker
	Header
	Cause       SearchRejectCause
	FailedField SearchFailedField
}

func (*AgentSearchReject) MessageType() MessageType { return MessageTypeAgentSearchReject }
func (*AgentSearchReject) Direction() Direction     { return Downlink }

type AgentGroupingInvitation struct {
	messageMarker
	Header
	GroupID       string
	SourceAgentID string
	TargetAgentID string
	TaskID        string
	ExpiresAt     uint64
	Proof         []byte
}

func (*AgentGroupingInvitation) MessageType() MessageType {
	return MessageTypeAgentGroupingInvitation
}
func (*AgentGroupingInvitation) Direction() Direction { return Downlink }

type GroupingDecision uint8

const (
	GroupingDecisionAccept GroupingDecision = iota
	GroupingDecisionReject
)

type GroupingRejectReason uint8

const (
	GroupingRejectReasonNone GroupingRejectReason = iota
	GroupingRejectReasonUserReject
	GroupingRejectReasonBusy
	GroupingRejectReasonCapabilityUnavailable
	GroupingRejectReasonSecurityFailure
	GroupingRejectReasonLocalError
)

type AgentGroupingInvitationResponse struct {
	messageMarker
	Header
	GroupID      string
	AgentID      string
	Decision     GroupingDecision
	RejectReason GroupingRejectReason
	Timestamp    uint64
	Proof        []byte
}

func (*AgentGroupingInvitationResponse) MessageType() MessageType {
	return MessageTypeAgentGroupingInvitationResponse
}
func (*AgentGroupingInvitationResponse) Direction() Direction { return Uplink }

type AgentGroupInfoNotification struct {
	messageMarker
	Header
	GroupID       string
	TargetAgentID string
	GroupConfig   json.RawMessage
}

func (*AgentGroupInfoNotification) MessageType() MessageType {
	return MessageTypeAgentGroupInfoNotification
}
func (*AgentGroupInfoNotification) Direction() Direction { return Downlink }

type GroupInfoApplyResult uint8

const (
	GroupInfoApplyResultSuccess GroupInfoApplyResult = iota
	GroupInfoApplyResultFailure
)

type GroupInfoFailureCause uint8

const (
	GroupInfoFailureCauseNone GroupInfoFailureCause = iota
	GroupInfoFailureCauseConfigInvalid
	GroupInfoFailureCauseProofInvalid
	GroupInfoFailureCauseRelayInvalid
	GroupInfoFailureCauseLocalApplyFailed
)

type AgentGroupInfoNotificationResponse struct {
	messageMarker
	Header
	GroupID      string
	AgentID      string
	Result       GroupInfoApplyResult
	FailureCause GroupInfoFailureCause
}

func (*AgentGroupInfoNotificationResponse) MessageType() MessageType {
	return MessageTypeAgentGroupInfoNotificationResponse
}
func (*AgentGroupInfoNotificationResponse) Direction() Direction { return Uplink }
