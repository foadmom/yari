package comms


import (
	"log"
	"net/url"
	"time"

	"qpid.apache.org/amqp"
	"qpid.apache.org/electron"

	"yari/types"
)

func InitTransport (artemis *types.TransportType) error {
	var _err	error;

	_err = establishConnection (artemis);
	return _err;
}



// ============================================================================
// 
// ============================================================================
func InitHBTransport (artemis *types.TransportType, role types.RoleType) error {
	var _err	error;

	return _err;
}

// ============================================================================
// initialise Campaign queue
// ============================================================================
func InitCampaignTransport (artemis *types.TransportType) error {
	var _err	error;

	DeInitCampaignTransport (artemis, nil);
	artemis.CampaignListener, _err = createListner (artemis.Connection, artemis.Config.CampaignQueue, artemis.Config.ListenerPrefetch);
	if (_err == nil) {
		artemis.CampaignSender, _err = createSender (artemis.Connection, artemis.Config.CampaignQueue);
		// if (_err == nil) {
		// 	artemis.VotingListener, _err = createListner (artemis.Connection, artemis.Config.VotingQueue, artemis.Config.ListenerPrefetch);
			// if (_err == nil) {
			// 	artemis.VotingSender, _err = createSender (artemis.Connection, artemis.Config.VotingQueue);
			// }
		// }
	}
	if (_err != nil) {
		DeInitCampaignTransport (artemis, _err);
	}
	return _err;
}


// ============================================================================
// initialise Campaign queue
// ============================================================================
func DeInitCampaignTransport (artemis *types.TransportType, err error) {
	closeReceiver (artemis.CampaignListener, err);
	closeSender   (artemis.CampaignSender,   err);
	// closeReceiver (artemis.VotingListener,   err);
	// closeSender   (artemis.VotingSender,     err);
}


// ============================================================================
// 
// ============================================================================
// func DeInitHBTransport (artemis *types.TransportType, role types.RoleType) {
// 	if (role == types.ROLE_Follower) {
// 		FollowerDeInitTransport (artemis)
// 	} else if (role == types.ROLE_Leader) {
// 		leaderDeInitArtemis (artemis);
// 	}
// }

// ============================================================================
// ============================================================================
// ============================================================================
// 
// ============================================================================
func initFollowerHearbeatTransport (artemis *types.TransportType)  error {
	var _err	error;

	artemis.HeartbeatPingListener, _err = createListner (artemis.Connection, artemis.Config.HeartbeatPingQueue, 
		artemis.Config.ListenerPrefetch);
	if (_err == nil) {
		artemis.HeartbeatPongSender, _err = createSender (artemis.Connection, artemis.Config.HeartbeatPongQueue);
	}

	return _err;
}

// ============================================================================
// 
// ============================================================================
func FollowerDeInitTransport (artemis *types.TransportType) {

	if (artemis.HeartbeatPingListener != nil) {
		artemis.HeartbeatPingListener.Close (nil);
		artemis.HeartbeatPingListener = nil;
	}
	if (artemis.HeartbeatPongSender != nil) {
		artemis.HeartbeatPongSender.Close (nil);
		artemis.HeartbeatPongSender = nil;
	}
}


// ============================================================================
// 
// ============================================================================
func WaitForHeartbeatPingRequest (comms *types.CommsType) ([]byte, error) {
	var _artemis  *types.TransportType = &comms.Transport;
	var _err				error;
	var _message			[]byte;

	_message, _err = receiveMessage (_artemis.HeartbeatPingListener);

	return _message, _err;
}

// ============================================================================
// 
// ============================================================================
func SendHeartbeatPongResponse (artemis  *types.TransportType, message []byte) error {
	return sendMessage (artemis.HeartbeatPongSender, message, artemis.Config.MessageExpiry);
}


// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// 
// ============================================================================
func leaderInitTransport (artemis *types.TransportType) error {
	var _err	error;

	artemis.HeartbeatPongListener, _err = createListner (artemis.Connection, artemis.Config.HeartbeatPongQueue, 
		                                      artemis.Config.ListenerPrefetch);
	if (_err == nil) {
		artemis.HeartbeatPingSender, _err = createSender (artemis.Connection, artemis.Config.HeartbeatPingQueue);
	}

	return _err;
}

// ============================================================================
// 
// ============================================================================
func leaderDeInitTransport (artemis *types.TransportType) {

	if (artemis.HeartbeatPingListener != nil) {
		artemis.HeartbeatPingListener.Close (nil);
		artemis.HeartbeatPingListener = nil;
	}
	if (artemis.HeartbeatPongSender != nil) {
		artemis.HeartbeatPongSender.Close (nil);
		artemis.HeartbeatPongSender = nil;
	}
}

// ============================================================================
// 
// ============================================================================
func waitForHeartbeatPong (artemis *types.TransportType) ([]byte, error) {
	var _err				error;
	var _message			[]byte;

	_message, _err = receiveMessage (artemis.HeartbeatPongListener);

	return _message, _err;
}

// ============================================================================
// 
// ============================================================================
func sendHeartbeatPing (artemis *types.TransportType, message []byte) error {
	return sendMessage (artemis.HeartbeatPingSender, message, artemis.Config.MessageExpiry);
}




// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// Independent functions to be used with any sender or receiver
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================


// ============================================================================
// Connection
// ============================================================================
func establishConnection (artemis *types.TransportType) error {
	var	_url	*url.URL;
	var _err  	error;

	if (artemis.Connection == nil) {
		_url, _err = amqp.ParseURL (artemis.Config.Host);
		artemis.Container = electron.NewContainer(artemis.Config.ContainerName);
		_connOptContainerId := electron.ContainerId (artemis.Container.Id());
		_connOptUser := electron.User (artemis.Config.User);
		_connOptPass := electron.Password([]byte(artemis.Config.Password));
		_connOptSASAllowInsecure := electron.SASLAllowInsecure (true);
		_connOptSASAllowedMechsPlain := electron.SASLAllowedMechs("PLAIN");
		_connOptSASAllowedMechsAnon  := electron.SASLAllowedMechs("ANONYMOUS");
		log.Printf ("artemis.go:establishConnection: connecting with _url=%v\n", _url);
		artemis.Connection, _err = artemis.Container.Dial("tcp", _url.Host, _connOptContainerId,
				 _connOptSASAllowedMechsAnon,
				 _connOptUser, _connOptPass,
				 _connOptSASAllowedMechsPlain,
				 _connOptSASAllowInsecure);
		if (_err != nil) {
			log.Println (_err);
		}
	}

	return _err;
}

// ============================================================================
// 
// ============================================================================
func closeSender (sender electron.Sender, err error) {
	if (sender != nil) {
		sender.Close (err);
	}
}

// ============================================================================
// 
// ============================================================================
func closeReceiver (receiver electron.Receiver, err error) {
	if (receiver != nil) {
		receiver.Close (err);
	}
}

// ============================================================================
// 
// ============================================================================
func createListner (connection electron.Connection, queueName string, preFetch int, )(electron.Receiver, error) {
	var _opts  		[]electron.LinkOption;
	var _err		error;
	var _receiver   electron.Receiver;

	_opts = []electron.LinkOption{electron.Source(queueName)};
	if preFetch > 0 { // Use a pre-fetch window
		_opts = append(_opts, electron.Capacity (preFetch), electron.Prefetch(true));
	} else { // Grant credit for all expected messages at once
		_opts = append(_opts, electron.Capacity(preFetch), electron.Prefetch(false));
	}
	_receiver, _err = createReceiver (connection, _opts);

	return _receiver, _err;
}


// ============================================================================
// 
// ============================================================================
func createSender (c electron.Connection, addr string) (electron.Sender, error) {
	_sender, _err := c.Sender(electron.Target(addr));
	return _sender, _err;
}

// ============================================================================
// ============================================================================
// ============================================================================
// 
// ============================================================================
func createReceiver (connection electron.Connection, opts []electron.LinkOption) (electron.Receiver, error) {
	_receiver, _err := connection.Receiver(opts...)
	return _receiver, _err;
}

// ============================================================================
// 
// ============================================================================
func sendMessage (electronSender electron.Sender, message []byte, expiryInterval int64) error {
	m := amqp.NewMessage()
	m.Marshal(message);
	if (expiryInterval != -1) {
		_expiryTime := time.Now ().Add (time.Duration(expiryInterval) * time.Millisecond);
		m.SetExpiryTime (_expiryTime);
	}
	_outcome := electronSender.SendSync(m); 
	return _outcome.Error;
}


// ============================================================================
// 
// ============================================================================
func receiveMessage (receiver electron.Receiver) ([]byte, error) {
	var _message		string;
	var _err 			error;
	var _rcvdMessage	electron.ReceivedMessage;

	if (receiver != nil) {
		_rcvdMessage, _err = receiver.Receive();
		if _err == nil {
			_rcvdMessage.Accept();
			_rcvdMessage.Message.Unmarshal(&_message);
		} else if _err == electron.Closed {
			log.Println ("artemis:receiveMessage: electron closed");
		} else {
			log.Fatalf("artemis:receiveMessage: receive error %v", _err);
		}
	
	} else {
		_err = electron.Closed;
	}

	return []byte(_message), _err;
}

// ============================================================================
// 
// ============================================================================
// func receiveMessageWithChannel (receiver electron.Receiver, _chan chan []byte) {
// 	_rcvdMessage, _err := receiveMessage (receiver);
// 	if (_err == nil) {
// 		_chan <- _rcvdMessage;
// 	}
// }


// ============================================================================
// 
// ============================================================================
// func receiveMessageWithTimeout (receiver electron.Receiver, timeout time.Duration) ([]byte, error) {
// 	var  _rcvdMessage 		[]byte;
// 	var  _chan				chan []byte;
// 	var  _err				error;

// 	_chan = make (chan []byte);
// 	defer close (_chan);
// 	_timeout  := time.NewTimer (timeout);
// 	go receiveMessageWithChannel (receiver, _chan);
// 	select {
// 		case _rcvdMessage = <- _chan :
// 			_timeout.Stop();
// 		case <-_timeout.C :
// 			_rcvdMessage = nil;
// 			_err = errors.New ("timeout");
// 	}
// 	return _rcvdMessage, _err;
// }


// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// Campaign
// ============================================================================

func receiveFromCampaignQueue (artemis *types.TransportType) ([]byte, error) {

	return receiveMessage (artemis.CampaignListener);
}

func SendToCampaignQueue (artemis *types.TransportType, request []byte) error {
	var  _err		error;
	_err = sendMessage (artemis.CampaignSender, request, artemis.Config.MessageExpiry);

	return _err;
}

// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// 
// ============================================================================
func OnCloseTransport (artemis *types.TransportType) {
	DeInitCampaignTransport (artemis, nil);
	
	if (artemis.Connection != nil) {
		artemis.Connection.Close (nil);
		artemis.Connection = nil;
	}

}