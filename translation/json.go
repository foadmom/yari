package translation

import (
	"encoding/json"
	"fmt"

	"github.com/foadmom/yari/types"
)

// ============================================================================
// here we unmarshal the header and let specific message handlers
// 		unmarshal the payload
// ============================================================================
func TranslateToMessage(buffer []byte) (*types.RaftMessage, error) {
	var _err error
	var _message *types.RaftMessage
	var _structure jsonIntermediateMessage
	_err = json.Unmarshal(buffer, &_structure)
	if _err == nil {
		_handler := HandlerList[types.ActionType(_structure.Header.Action)]
		if _handler != nil {
			_message, _err = _handler(&_structure)
		} else {
			_err = fmt.Errorf("unable to find a handler for Action %s", _structure.Header.Action)
		}
	}

	return _message, _err
}

// ============================================================================
// marshal a RaftMessage into json
// ============================================================================
func TranslateFromMessage(message *types.RaftMessage) ([]byte, error) {
	_buffer, _err := json.Marshal(message)

	return _buffer, _err
}

type jsonIntermediateMessage struct {
	MessageType string
	Header      types.MessageHeader
	Payload     json.RawMessage `json:",omitempty"`
}

// ============================================================================
// how a json message coming through network and converted to structure:
// 		CreateMessageFromJson will unmarshall the message header from
//		json to RaftMessage.Header.
//		the action in the header is then used as key in HandlerList to
//		call the function for unmarshalling the payload into proper structure
//		for the request type.
//
// ============================================================================

// ============================================================================
// payload for every message can be different. here we create a
//		map with Action in the header as key and
//		the function to unmarshal the specific payload as the value
//		in the map.
// eg. when the header is unmarshalled and the Action is 'HeartBeatResponse'
//		then the map is searched for the correct function to unmarshal the
//		payload which is 'CreateHeartBeatResponse'
// ============================================================================
type payloadHandlerFunction func(*jsonIntermediateMessage) (*types.RaftMessage, error)

var HandlerList map[types.ActionType]payloadHandlerFunction = map[types.ActionType]payloadHandlerFunction{
	types.AT_HeartBeatRequest:  translateJsonToNoPayload,
	types.AT_HeartBeatResponse: translateJsonToHBResponsePayload,
	types.AT_VoteLodged:        translateJsonToCampaignPayload,
	types.AT_CampaignLaunch:    translateJsonToCampaignPayload,
	types.AT_NewLeaderElected:  translateJsonToCampaignPayload,
	types.AT_ResetCampaign:     translateJsonToCampaignPayload,
}

// ============================================================================
// unmarshal the payload for HeartBeat response and return the RaftMessage struct
// ============================================================================
func translateJsonToHBRequestPayload(structure *jsonIntermediateMessage) (*types.RaftMessage, error) {
	var _err error
	var _message types.RaftMessage
	var _payload types.HeartBeatPingPayload

	if len(structure.Payload) > 0 {
		_err = json.Unmarshal([]byte(structure.Payload), &_payload)
		if _err == nil {
			_message.MessageType = structure.MessageType
			_message.Header = structure.Header
			_message.Payload = _payload
		}
	}

	return &_message, _err
}

// ============================================================================
// unmarshal the payload for HeartBeat response and return the RaftMessage struct
// ============================================================================
func translateJsonToHBResponsePayload(structure *jsonIntermediateMessage) (*types.RaftMessage, error) {
	var _err error
	var _message types.RaftMessage
	var _payload types.HeartBeatPongPayload

	if len(structure.Payload) > 0 {
		_err = json.Unmarshal([]byte(structure.Payload), &_payload)
		if _err == nil {
			_message.MessageType = structure.MessageType
			_message.Header = structure.Header
			_message.Payload = _payload
		}
	}

	return &_message, _err
}

// ============================================================================
// unmarshal the payload for Campaign message
// ============================================================================
func translateJsonToCampaignPayload(structure *jsonIntermediateMessage) (*types.RaftMessage, error) {
	var _err error
	var _message types.RaftMessage
	var _payload types.CampaignPayload

	if len(structure.Payload) > 0 {
		_err = json.Unmarshal([]byte(structure.Payload), &_payload)
		if _err == nil {
			_message.MessageType = structure.MessageType
			_message.Header = structure.Header
			_message.Payload = _payload
		}
	}

	return &_message, _err
}

// ============================================================================
//
// ============================================================================
func translateJsonToNoPayload(structure *jsonIntermediateMessage) (*types.RaftMessage, error) {
	var _err error
	var _message types.RaftMessage

	_message.MessageType = structure.MessageType
	_message.Header = structure.Header
	_message.Payload = nil

	return &_message, _err
}
