package acn

import "encoding/json"

const (
	Version1                uint8 = 0x01
	PayloadContainerTypeACN uint8 = 0x0e
)

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
	MessageTypeAgentGroupingRequest               MessageType = 0x0d
	MessageTypeAgentGroupingAccept                MessageType = 0x0e
	MessageTypeAgentGroupingReject                MessageType = 0x0f
	MessageTypeAgentGroupingInvitation            MessageType = 0x10
	MessageTypeAgentGroupingInvitationResponse    MessageType = 0x11
	MessageTypeAgentGroupInfoNotification         MessageType = 0x12
	MessageTypeAgentGroupInfoNotificationResponse MessageType = 0x13
	MessageTypeAgentNetworkAbilityRequest         MessageType = 0x14
	MessageTypeAgentNetworkAbilityResponse        MessageType = 0x15
	MessageTypeAgentNetworkAbilityReject          MessageType = 0x16
	MessageTypeAgentPublishRequest                MessageType = 0x17
	MessageTypeAgentPublishAccept                 MessageType = 0x18
	MessageTypeAgentPublishReject                 MessageType = 0x19
)

func (t MessageType) String() string {
	names := map[MessageType]string{
		MessageTypeAgentRegisterRequest:               "ACN_AGENT_REGISTER_REQUEST",
		MessageTypeAgentRegisterAccept:                "ACN_AGENT_REGISTER_ACCEPT",
		MessageTypeAgentRegisterReject:                "ACN_AGENT_REGISTER_REJECT",
		MessageTypeAgentDeregisterRequest:             "ACN_AGENT_DEREGISTER_REQUEST",
		MessageTypeAgentDeregisterAccept:              "ACN_AGENT_DEREGISTER_ACCEPT",
		MessageTypeAgentDeregisterReject:              "ACN_AGENT_DEREGISTER_REJECT",
		MessageTypeAgentProfileUpdateRequest:          "ACN_AGENT_PROFILE_UPDATE_REQUEST",
		MessageTypeAgentProfileUpdateAccept:           "ACN_AGENT_PROFILE_UPDATE_ACCEPT",
		MessageTypeAgentProfileUpdateReject:           "ACN_AGENT_PROFILE_UPDATE_REJECT",
		MessageTypeAgentSearchRequest:                 "ACN_AGENT_SEARCH_REQUEST",
		MessageTypeAgentSearchResponse:                "ACN_AGENT_SEARCH_RESPONSE",
		MessageTypeAgentSearchReject:                  "ACN_AGENT_SEARCH_REJECT",
		MessageTypeAgentGroupingRequest:               "ACN_AGENT_GROUPING_REQUEST",
		MessageTypeAgentGroupingAccept:                "ACN_AGENT_GROUPING_ACCEPT",
		MessageTypeAgentGroupingReject:                "ACN_AGENT_GROUPING_REJECT",
		MessageTypeAgentGroupingInvitation:            "ACN_AGENT_GROUPING_INVITATION",
		MessageTypeAgentGroupingInvitationResponse:    "ACN_AGENT_GROUPING_INVITATION_RESPONSE",
		MessageTypeAgentGroupInfoNotification:         "ACN_AGENT_GROUPINFO_NOTIFICATION",
		MessageTypeAgentGroupInfoNotificationResponse: "ACN_AGENT_GROUPINFO_NOTIFICATION_RESPONSE",
		MessageTypeAgentNetworkAbilityRequest:         "ACN_AGENT_NETWORK_ABILITY_REQUEST",
		MessageTypeAgentNetworkAbilityResponse:        "ACN_AGENT_NETWORK_ABILITY_RESPONSE",
		MessageTypeAgentNetworkAbilityReject:          "ACN_AGENT_NETWORK_ABILITY_REJECT",
		MessageTypeAgentPublishRequest:                "ACN_AGENT_PUBLISH_REQUEST",
		MessageTypeAgentPublishAccept:                 "ACN_AGENT_PUBLISH_ACCEPT",
		MessageTypeAgentPublishReject:                 "ACN_AGENT_PUBLISH_REJECT",
	}
	if name, ok := names[t]; ok {
		return name
	}
	return "ACN_MESSAGE_UNKNOWN"
}

func (t MessageType) direction() (Direction, bool) {
	switch t {
	case MessageTypeAgentRegisterRequest,
		MessageTypeAgentDeregisterRequest,
		MessageTypeAgentProfileUpdateRequest,
		MessageTypeAgentSearchRequest,
		MessageTypeAgentGroupingRequest,
		MessageTypeAgentGroupingInvitationResponse,
		MessageTypeAgentGroupInfoNotificationResponse,
		MessageTypeAgentNetworkAbilityRequest,
		MessageTypeAgentPublishRequest:
		return Uplink, true
	case MessageTypeAgentRegisterAccept,
		MessageTypeAgentRegisterReject,
		MessageTypeAgentDeregisterAccept,
		MessageTypeAgentDeregisterReject,
		MessageTypeAgentProfileUpdateAccept,
		MessageTypeAgentProfileUpdateReject,
		MessageTypeAgentSearchResponse,
		MessageTypeAgentSearchReject,
		MessageTypeAgentGroupingAccept,
		MessageTypeAgentGroupingReject,
		MessageTypeAgentGroupingInvitation,
		MessageTypeAgentGroupInfoNotification,
		MessageTypeAgentNetworkAbilityResponse,
		MessageTypeAgentNetworkAbilityReject,
		MessageTypeAgentPublishAccept,
		MessageTypeAgentPublishReject:
		return Downlink, true
	default:
		return 0, false
	}
}

type Header struct {
	TransactionID uint8
}

func (h Header) GetHeader() Header { return h }

type messageMarker struct{}

func (messageMarker) isACNMessage() {}

type Message interface {
	GetHeader() Header
	MessageType() MessageType
	Direction() Direction
	isACNMessage()
}

type AgentRegisterRequest struct {
	messageMarker
	Header
	Owner       string
	AgentName   string
	PublicKey   []byte
	Description string
	Timestamp   uint64
	Signature   []byte
	Metadata    json.RawMessage
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
	RegisterRejectRequestValidationFailed RegisterRejectCause = iota + 1
	RegisterRejectOwnerNotAllowed
	RegisterRejectPublicKeyAlgorithmUnsupported
	RegisterRejectPublicKeyInvalid
	RegisterRejectSignatureInvalid
	RegisterRejectTimestampInvalid
	RegisterRejectIDAllocationFailed
	RegisterRejectInternalError
)

type RegisterFailedField uint8

const (
	RegisterFieldUnspecified RegisterFailedField = iota
	RegisterFieldOwner
	RegisterFieldAgentName
	RegisterFieldPublicKey
	RegisterFieldDescription
	RegisterFieldTimestamp
	RegisterFieldSignature
	RegisterFieldMetadata
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
	DeregistrationReasonNormal DeregistrationReason = iota
	DeregistrationReasonUninstalled
	DeregistrationReasonReplaced
	DeregistrationReasonUserRequest
	DeregistrationReasonSecurityEvent
	DeregistrationReasonRetired
	DeregistrationReasonOther DeregistrationReason = 0xff
)

type AgentDeregisterRequest struct {
	messageMarker
	Header
	AgentID   string
	Reason    DeregistrationReason
	Timestamp uint64
	Signature []byte
}

func (*AgentDeregisterRequest) MessageType() MessageType { return MessageTypeAgentDeregisterRequest }
func (*AgentDeregisterRequest) Direction() Direction     { return Uplink }

type AgentDeregisterAccept struct {
	messageMarker
	Header
}

func (*AgentDeregisterAccept) MessageType() MessageType { return MessageTypeAgentDeregisterAccept }
func (*AgentDeregisterAccept) Direction() Direction     { return Downlink }

type DeregisterRejectCause uint8

const (
	DeregisterRejectAgentNotFound DeregisterRejectCause = iota + 1
	DeregisterRejectSignatureInvalid
	DeregisterRejectRequestValidationFailed
	DeregisterRejectNotAuthorized
	DeregisterRejectInternalError
)

type DeregisterFailedField uint8

const (
	DeregisterFieldUnspecified DeregisterFailedField = iota
	DeregisterFieldAgentID
	DeregisterFieldReason
	DeregisterFieldTimestamp
	DeregisterFieldSignature
)

type AgentDeregisterReject struct {
	messageMarker
	Header
	Cause       DeregisterRejectCause
	FailedField DeregisterFailedField
}

func (*AgentDeregisterReject) MessageType() MessageType { return MessageTypeAgentDeregisterReject }
func (*AgentDeregisterReject) Direction() Direction     { return Downlink }

type NetworkAbilityRejectCause uint8

const (
	NetworkAbilityRejectAgentInvalid NetworkAbilityRejectCause = iota + 1
	NetworkAbilityRejectProofInvalid
	NetworkAbilityRejectRequestValidationFailed
	NetworkAbilityRejectNotAuthorized
	NetworkAbilityRejectInternalError
)

type NetworkAbilityFailedField uint8

const (
	NetworkAbilityFieldUnspecified NetworkAbilityFailedField = iota
	NetworkAbilityFieldAgentID
	NetworkAbilityFieldTimestamp
	NetworkAbilityFieldProof
)

type AgentNetworkAbilityRequest struct {
	messageMarker
	Header
	AgentID   string
	Timestamp uint64
	Proof     json.RawMessage
}

func (*AgentNetworkAbilityRequest) MessageType() MessageType {
	return MessageTypeAgentNetworkAbilityRequest
}
func (*AgentNetworkAbilityRequest) Direction() Direction { return Uplink }

type AgentNetworkAbilityResponse struct {
	messageMarker
	Header
	Timestamp uint64
	VC1       json.RawMessage
}

func (*AgentNetworkAbilityResponse) MessageType() MessageType {
	return MessageTypeAgentNetworkAbilityResponse
}
func (*AgentNetworkAbilityResponse) Direction() Direction { return Downlink }

type AgentNetworkAbilityReject struct {
	messageMarker
	Header
	Cause       NetworkAbilityRejectCause
	FailedField NetworkAbilityFailedField
}

func (*AgentNetworkAbilityReject) MessageType() MessageType {
	return MessageTypeAgentNetworkAbilityReject
}
func (*AgentNetworkAbilityReject) Direction() Direction { return Downlink }

type Priority uint8

const (
	PriorityUnspecified Priority = iota
	PriorityHigh
	PriorityNormal
	PriorityLow
)

type AgentPublishRequest struct {
	messageMarker
	Header
	AgentID   string
	Priority  Priority
	Timestamp uint64
	Signature []byte
	VCList    json.RawMessage
}

func (*AgentPublishRequest) MessageType() MessageType { return MessageTypeAgentPublishRequest }
func (*AgentPublishRequest) Direction() Direction     { return Uplink }

type AgentPublishAccept struct {
	messageMarker
	Header
}

func (*AgentPublishAccept) MessageType() MessageType { return MessageTypeAgentPublishAccept }
func (*AgentPublishAccept) Direction() Direction     { return Downlink }

type PublishRejectCause uint8

const (
	PublishRejectAgentInvalid PublishRejectCause = iota + 1
	PublishRejectVCInvalid
	PublishRejectSignatureInvalid
	PublishRejectRequestValidationFailed
	PublishRejectInternalError
)

type PublishFailedField uint8

const (
	PublishFieldUnspecified PublishFailedField = iota
	PublishFieldAgentID
	PublishFieldPriority
	PublishFieldTimestamp
	PublishFieldSignature
	PublishFieldVCList
)

type AgentPublishReject struct {
	messageMarker
	Header
	Cause       PublishRejectCause
	FailedField PublishFailedField
}

func (*AgentPublishReject) MessageType() MessageType { return MessageTypeAgentPublishReject }
func (*AgentPublishReject) Direction() Direction     { return Downlink }

type AgentProfileUpdateRequest struct {
	messageMarker
	Header
	RequestID   string
	AgentID     string
	UpdateItems json.RawMessage
	Credentials json.RawMessage
	Timestamp   uint64
	Proof       json.RawMessage
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
	ProfileUpdateRejectAgentInvalid ProfileUpdateRejectCause = iota + 1
	ProfileUpdateRejectUpdateItemInvalid
	ProfileUpdateRejectCredentialInvalid
	ProfileUpdateRejectProofInvalid
	ProfileUpdateRejectRequestValidationFailed
	ProfileUpdateRejectInternalError
)

type ProfileUpdateFailedField uint8

const (
	ProfileUpdateFieldUnspecified ProfileUpdateFailedField = iota
	ProfileUpdateFieldRequestID
	ProfileUpdateFieldAgentID
	ProfileUpdateFieldUpdateItems
	ProfileUpdateFieldCredentials
	ProfileUpdateFieldTimestamp
	ProfileUpdateFieldProof
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

type DiscoveryScope uint8

const (
	DiscoveryScopeIntraPLMN DiscoveryScope = iota + 1
	DiscoveryScopeInterPLMN
)

type AgentSearchRequest struct {
	messageMarker
	Header
	TaskID          string
	AgentID         string
	TaskDescription string
	RequiredSkills  []string
	DiscoveryScope  DiscoveryScope
	MaxResults      uint8
	Timestamp       uint64
	Proof           json.RawMessage
}

func (*AgentSearchRequest) MessageType() MessageType { return MessageTypeAgentSearchRequest }
func (*AgentSearchRequest) Direction() Direction     { return Uplink }

type SearchResult struct {
	AgentCard json.RawMessage
	Priority  Priority
}

type AgentSearchResponse struct {
	messageMarker
	Header
	TaskID          string
	TaskDescription string
	Results         []SearchResult
	Timestamp       uint64
}

func (*AgentSearchResponse) MessageType() MessageType { return MessageTypeAgentSearchResponse }
func (*AgentSearchResponse) Direction() Direction     { return Downlink }

type SearchRejectCause uint8

const (
	SearchRejectAgentInvalid SearchRejectCause = iota + 1
	SearchRejectSkillInvalid
	SearchRejectScopeInvalid
	SearchRejectProofInvalid
	SearchRejectRequestValidationFailed
	SearchRejectInternalError
)

type SearchFailedField uint8

const (
	SearchFieldUnspecified SearchFailedField = iota
	SearchFieldTaskID
	SearchFieldAgentID
	SearchFieldTaskDescription
	SearchFieldRequiredSkills
	SearchFieldDiscoveryScope
	SearchFieldMaxResults
	SearchFieldTimestamp
	SearchFieldProof
)

type AgentSearchReject struct {
	messageMarker
	Header
	Cause       SearchRejectCause
	FailedField SearchFailedField
}

func (*AgentSearchReject) MessageType() MessageType { return MessageTypeAgentSearchReject }
func (*AgentSearchReject) Direction() Direction     { return Downlink }

type AgentGroupingRequest struct {
	messageMarker
	Header
	AgentID        string
	TargetAgentIDs []string
	GroupConfig    json.RawMessage
	Timestamp      uint64
	Proof          json.RawMessage
}

func (*AgentGroupingRequest) MessageType() MessageType { return MessageTypeAgentGroupingRequest }
func (*AgentGroupingRequest) Direction() Direction     { return Uplink }

type AgentGroupingAccept struct {
	messageMarker
	Header
	GroupID string
}

func (*AgentGroupingAccept) MessageType() MessageType { return MessageTypeAgentGroupingAccept }
func (*AgentGroupingAccept) Direction() Direction     { return Downlink }

type GroupingRejectCause uint8

const (
	GroupingRejectSourceAgentInvalid GroupingRejectCause = iota + 1
	GroupingRejectTargetAgentInvalid
	GroupingRejectTargetAgentRejected
	GroupingRejectTargetAgentTimeout
	GroupingRejectRequestValidationFailed
	GroupingRejectProofInvalid
	GroupingRejectInternalError
)

type GroupingFailedField uint8

const (
	GroupingFieldUnspecified GroupingFailedField = iota
	GroupingFieldSourceAgentID
	GroupingFieldTargetAgents
	GroupingFieldGroupConfig
	GroupingFieldTimestamp
	GroupingFieldProof
)

type AgentGroupingReject struct {
	messageMarker
	Header
	Cause          GroupingRejectCause
	FailedField    GroupingFailedField
	RelatedAgentID string
}

func (*AgentGroupingReject) MessageType() MessageType { return MessageTypeAgentGroupingReject }
func (*AgentGroupingReject) Direction() Direction     { return Downlink }

type AgentGroupingInvitation struct {
	messageMarker
	Header
	GroupConfig        json.RawMessage
	GroupAdministrator json.RawMessage
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

type AgentGroupingInvitationResponse struct {
	messageMarker
	Header
	Decision GroupingDecision
}

func (*AgentGroupingInvitationResponse) MessageType() MessageType {
	return MessageTypeAgentGroupingInvitationResponse
}
func (*AgentGroupingInvitationResponse) Direction() Direction { return Uplink }

type AgentGroupInfoNotification struct {
	messageMarker
	Header
	Version   string
	Timestamp uint64
	GroupID   string
	Members   json.RawMessage
	Proof     json.RawMessage
}

func (*AgentGroupInfoNotification) MessageType() MessageType {
	return MessageTypeAgentGroupInfoNotification
}
func (*AgentGroupInfoNotification) Direction() Direction { return Downlink }

type GroupInfoStatus uint8

const (
	GroupInfoStatusAccepted GroupInfoStatus = iota
	GroupInfoStatusRejected
)

type AgentGroupInfoNotificationResponse struct {
	messageMarker
	Header
	GroupID   string
	AgentID   string
	Status    GroupInfoStatus
	Detail    string
	Timestamp uint64
	Proof     json.RawMessage
}

func (*AgentGroupInfoNotificationResponse) MessageType() MessageType {
	return MessageTypeAgentGroupInfoNotificationResponse
}
func (*AgentGroupInfoNotificationResponse) Direction() Direction { return Uplink }
