package comms


import (
	"yari/message"
	"yari/translation"
	"yari/types"
	t "yari/types"
	"log"
)


func InitComms (comms *t.CommsType) error {
	var _err 			error;

	comms.Transport.Config = comms.Config.Transport;

	_err = InitTransport (&comms.Transport);
	return _err;
}

// ============================================================================
// ============================================================================
// Campaign
// ============================================================================

func InitCampaign (comms *t.CommsType) error {
	return InitCampaignTransport (&comms.Transport);
}


// ============================================================================
// ============================================================================
// Follower heartbeat
// ============================================================================
// ============================================================================
func InitFollowerHBComms (comms *t.CommsType) error {
	return initFollowerHearbeatTransport (&comms.Transport)
}

func DeInitFollowerHBComms (comms *t.CommsType) error {
	FollowerDeInitTransport (&comms.Transport);
	return nil;
}

func WaitForHeartBeatPing (comms *t.CommsType, channel types.RaftMessageChannel) {
	var  	_message  	*types.RaftMessage;
	var  	_err		error;

	for {
		_message, _err = waitForHeartBeatPing (comms);
		if (_err == nil) {
			channel <- _message;
		} else {
			return;
		}

	}
}



// ============================================================================
// 
// ============================================================================
func waitForHeartBeatPing (comms *t.CommsType) (*types.RaftMessage, error) {
	var _message   	*types.RaftMessage;
	var _buffer   	[]byte;
	
	_buffer, _err := WaitForHeartbeatPingRequest (comms);
	if (_err == nil) && (len (_buffer)>0) {
		_message, _err = translation.TranslateToMessage (_buffer);
	}


	return _message, _err;
}

// ============================================================================
// 
// ============================================================================
func SendHeartbeatPong (comms *t.CommsType, _message	*types.RaftMessage) error {
	var _err		error;
	var _buffer 	[]byte;
	
	_buffer, _err = translation.TranslateFromMessage (_message);
	if (_err == nil) {
		_err = SendHeartbeatPongResponse (&comms.Transport, _buffer);
	}
	return _err;
}






// ============================================================================
// ============================================================================
// ============================================================================
// ========================== leader heart beat ===============================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================

func InitLeaderHBComms (comms *t.CommsType) error {
	var   _err 	error;
	_err = leaderInitTransport (&comms.Transport);

	return _err;
}
func DeInitLeaderHBComms (comms *t.CommsType) {
	leaderDeInitTransport (&comms.Transport);
}

func SendHeartbeatPing (comms *t.CommsType) error {
	var  _err   error;

	_request := message.CreateHeartBeatPing(comms.Transport.Config.HeartbeatPongQueue);
	var _translatedPing []byte;
	_translatedPing, _err = translation.TranslateFromMessage(_request);

	if (_err == nil) {
		_err = sendHeartbeatPing (&comms.Transport, _translatedPing);
		if (_err == nil) {
			log.Print ("send ping. ");
		}
	}
	return _err;
}

// ============================================================================
// 
// ============================================================================
func WaitForHeartBeatPong (comms *t.CommsType) (*types.RaftMessage, error) {
	var _message   	*types.RaftMessage;
	var _buffer   	[]byte;
	
	_buffer, _err := waitForHeartbeatPong (&comms.Transport);
	if (_err == nil) && (len (_buffer)>0) {
		_message, _err = translation.TranslateToMessage (_buffer);
	}

	return _message, _err;
}


// ============================================================================
// 
// ============================================================================
func WaitForHeartbeatPongWithChannel (comms *types.CommsType, _chan types.RaftMessageChannel) {

	for {
		_message, _err := WaitForHeartBeatPong (comms);
		if (_err == nil) {
			_chan <- _message;
		} else {
			log.Fatalf ("error reading from the ping listener. error=\n%v\n", _err);
			return;
		}
	}
}

// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// 
// ============================================================================
func ListenToCampaignBroadcast (comms *t.CommsType, outcome types.RaftMessageChannel) {
	var _buffer  	[]byte;
	var _err		 error;
	var _message     *types.RaftMessage;

	for {
		_buffer, _err = receiveFromCampaignQueue (&comms.Transport);

		if (_err == nil) {
			_message, _err = translation.TranslateToMessage (_buffer);
			_payload := _message.Payload.(types.CampaignPayload);
			log.Printf ("node %s received campaign broadcast: %s", _payload.Candidate, _message.MessageType);
			outcome <- _message;
		} else {
			return;
		}
	}
}

func SendToCampaignBroadcast (comms *t.CommsType, message *types.RaftMessage) error {
	var _buffer  	[]byte;
	var _err		error;

	_buffer, _err = translation.TranslateFromMessage (message);
	if (_err == nil) {
		_err = SendToCampaignQueue (&comms.Transport, _buffer);
	}

	return _err;
}

// ============================================================================
// 
// ============================================================================
func OnClose (commsStruct *t.CommsType) {
	OnCloseTransport (&commsStruct.Transport);
}

