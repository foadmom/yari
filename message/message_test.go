package message


import (
	n "yari/node"
	"yari/types"
	"testing"

	"github.com/stretchr/testify/assert"
)


// ============================================================================
// 
// ============================================================================
func testStandardRaftMessage (t *testing.T, _message *types.RaftMessage, messageTitle string, 
	                          nodeId string, actionClass types.ActionClassType, action types.ActionType) {
	assert := assert.New (t);
	assert.Equal (_message.MessageType, messageTitle, "message type must be "+messageTitle);
	assert.NotEmpty (_message.Header.Uuid, "the created message must have Uuid");
	assert.NotEqual (_message.Header.Timestamp, nil, "Timestamp must not be nil");
	assert.Equal (_message.Header.SourceNodeId, nodeId, "the request must have the nodeId " + nodeId);
	assert.Equal (_message.Header.ActionClass, actionClass, "the ActionClass must be "+actionClass);
	assert.Equal (_message.Header.Action, action, "the Action must be "+action);
}

// ============================================================================
// 
// ============================================================================
func testHeartBeatResponse (t *testing.T, _message *types.RaftMessage) {
	assert := assert.New (t);

	testStandardRaftMessage (t, _message, string(types.AT_HeartBeatResponse), n.ThisNode().Id, types.AC_HeartBeat, types.AT_HeartBeatResponse);
	assert.Equal (_message.Header.SequenceId, uint64(2), "SequenceId should always be 2 for the heartbeat request");
}


// ============================================================================
// 
// ============================================================================
func testHeartBeatRequest (t *testing.T, _message *types.RaftMessage) {
	assert := assert.New (t);

	testStandardRaftMessage (t, _message, string(types.AT_HeartBeatRequest), n.ThisNode().Id, types.AC_HeartBeat, types.AT_HeartBeatRequest);

	assert.Equal (_message.Header.SequenceId, uint64(2), "SequenceId should always be 1 for the heartbeat request");
}




func TestCreateHeartBeatResponse (t *testing.T) {
	_message := CreateHeartBeatResponse (types.NS_Active);
	testHeartBeatResponse (t, _message);
}

// func TestCreateMessageFromJson (t *testing.T) {
// 	var _buffer string = `{"MessageType":"HeartBeatRequest","Header":{"Uuid":"9876df7gsd78fgy","RequestId":"","Timestamp":"2020-05-08T14:00:27.115342853+02:00","SequenceId":666,"SourceNodeIp":"192.168.1.189","SourceNodeId":"192.168.1.189","ActionClass":"HeartBeat","Action":"HeartBeatRequest"},"Payload":{}}`;
// 	var _bytes []byte = []byte(_buffer);
// 	var _config c.Config = c.Config{};
// 	_message, _ := CreateMessageFromJson (&_config, &_bytes);
// 	testStandardRaftMessage (t, _message, "HeartBeatRequest", "192.168.1.189", AC_HeartBeat, AT_HeartBeatRequest);
// }