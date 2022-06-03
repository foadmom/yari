package message

import (
	"time"

	c "github.com/foadmom/yari/common"
	"github.com/foadmom/yari/node"
	"github.com/foadmom/yari/types"
)

// ============================================================================
//
// ============================================================================
func CreateMessageHeader(seqNo uint64,
	actionClass types.ActionClassType,
	action types.ActionType) types.MessageHeader {
	var _uuid string = c.GenerateUuid_V4()

	var _header types.MessageHeader
	_header.Uuid = _uuid
	_header.Timestamp = time.Now()
	_header.SequenceId = seqNo
	_header.SourceNodeId = node.ThisNode().Id
	_header.SourceNodeIp = node.ThisNode().IpAddress
	_header.ActionClass = actionClass
	_header.Action = action
	return _header
}

// ============================================================================
//
// ============================================================================
func CreateHeartBeatPing(responseQueue string) *types.RaftMessage {

	var _message types.RaftMessage
	_message.MessageType = string(types.AT_HeartBeatResponse)
	_message.Header = CreateMessageHeader(1, types.AC_HeartBeat, types.AT_HeartBeatRequest)
	var _payload types.HeartBeatPingPayload = types.HeartBeatPingPayload{}
	_payload.ResponseQueue = responseQueue
	_message.Payload = _payload
	return &_message
}

// ============================================================================
//
// ============================================================================
func CreateHeartBeatResponse(status types.NodeStatus) *types.RaftMessage {

	var _message types.RaftMessage
	_message.Header = CreateMessageHeader(1, types.AC_HeartBeat, types.AT_HeartBeatRequest)
	_message.MessageType = string(types.AT_HeartBeatResponse)
	_message.Header.SourceNodeIp = node.ThisNode().IpAddress
	var _payload types.HeartBeatPongPayload = types.HeartBeatPongPayload{}
	_payload.Status = status
	_message.Payload = _payload
	return &_message
}

// ============================================================================
// ============================================================================
// ===============               Campaign                ======================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// CreateCampaignRequest header
// ============================================================================
func CreateCampaignRequest(actionClass types.ActionClassType, action types.ActionType) *types.RaftMessage {

	var _message types.RaftMessage
	_message.MessageType = string(action)
	_message.Header = CreateMessageHeader(1, actionClass, action)
	return &_message
}

func CreateCampaignMessage(action types.ActionType, candidateNodeId string, voteCount int) *types.RaftMessage {
	var _message *types.RaftMessage

	_message = CreateCampaignRequest(types.AC_Campaign, action)

	var _payload types.CampaignPayload = types.CampaignPayload{}
	_payload.Candidate = candidateNodeId
	_payload.VoteCount = voteCount
	_message.Payload = _payload
	return _message
}

func CreateVotingBallot(candidateNodeId string) *types.RaftMessage {
	var _message *types.RaftMessage

	_message = CreateCampaignMessage(types.AT_VoteLodged, candidateNodeId, 0)
	return _message
}

func CreateNewLeaderMessage(leaderNodeId string, voteCount int) *types.RaftMessage {
	var _message *types.RaftMessage

	_message = CreateCampaignMessage(types.AT_NewLeaderElected, leaderNodeId, voteCount)

	var _payload types.CampaignPayload = types.CampaignPayload{}
	_payload.Candidate = leaderNodeId
	_message.Payload = _payload
	return _message
}

// ========================================================
// create a request to nominate me as the leader
// ========================================================
func CreateALaunchMessage() *types.RaftMessage {
	var _message *types.RaftMessage

	_message = CreateCampaignMessage(types.AT_CampaignLaunch, node.ThisNode().Id, 0)
	return _message

}

// ========================================================
// create a request to reset/restart campaign
// ========================================================
func CreateACampaignResetMessage() *types.RaftMessage {
	var _message *types.RaftMessage

	_message = CreateCampaignMessage(types.AT_ResetCampaign, node.ThisNode().Id, 0)
	return _message

}

// ============================================================================
//
// ============================================================================

// ============================================================================
//
// ============================================================================
// func (message *types.RaftMessage) LogMessage () error {
// 	var 	_json	[]byte;
// 	var  	_err	error;

// 	_json, _err = json.Marshal (message);
// 	if (_err == nil) {
// //		storage.LogMessage (message.Header.Uuid, _json);
// 		log.Println (_json);
// 	}

// 	return _err;
// }
