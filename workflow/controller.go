package workflow


import (
	//	"fmt"

	"yari/common"
	"yari/comms"
	"yari/message"
	"yari/node"
	"yari/types"
	"log"
	"time"
)


// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// ============================================================================
// this is the top level flow control
// ============================================================================
func MainFlowManager (components *types.RootStructure) {

	log.Println ("workflow.controller.MainFlowManager: started");
	// ==================================================
	// create the internal channels
	components.Internal.MainToCampaignChannel       = make (types.SignalChannel);
	components.Internal.MainToFLeadHBChannel        = make (types.SignalChannel);
	components.Internal.MainToFolHBChannel          = make (types.SignalChannel);

	log.Printf ("MainLoop: starting campaign goroutine as %s\n", string(components.Internal.Role));
	if (components.Internal.Role == types.ROLE_Follower) {go followerHeartbeatManager (components);}
	if (components.Internal.Role == types.ROLE_Leader)   {go leaderHeartbeatManager   (components);}
	go CampaignManager (components);

	for {
		select {
			// from campaign manager
			case _signal := <- components.Internal.MainToCampaignChannel: {
				switch _signal {
					case types.Ch_NewLeader : {
						adjustToNewLeadership (components);
					}
				}
			}
			// from follower heartbeat manager
			case _signal := <- components.Internal.MainToFolHBChannel : {
				switch _signal {
					case types.Ch_HeartbeatTimeout : {
						log.Printf ("received heartbeat timeout. sending StartCampaign.");
						components.Internal.Role = types.ROLE_Campaign;
						components.Internal.MainToCampaignChannel <- types.Ch_StartCampaign;
					}
				}
			}
			// from leader heartbeat manager
			case _signal := <- components.Internal.MainToFLeadHBChannel : {
				switch _signal {
					case types.Ch_HeartbeatTimeout : {
						log.Printf ("received heartbeat timeout. sending StartCampaign.");
					}
				}
			}
		}
		}

	}

	func adjustToNewLeadership (components *types.RootStructure) {
		if (components.Internal.Role == types.ROLE_Follower) {go followerHeartbeatManager (components);}
		if (components.Internal.Role == types.ROLE_Leader)   {go leaderHeartbeatManager   (components);}
	
	}

// ============================================================================
// ============================================================================
// ============================================================================
// =============================== Campaign ===================================
// ============================================================================
// ============================================================================

type campaignManagerResources struct {
    campaignCommsChannel        types.RaftMessageChannel;
//    commsChannel                types.RaftMessageChannel;
    
    votingTimeoutTimer          *time.Timer;
    timeout                     time.Duration;
    randomWaitTimer             *time.Timer;
    randomWait                  time.Duration;
    
    candidate                   string;
    myNodeId                    string;
    // role                        types.RoleType
    
    msgRecveived                *types.RaftMessage
    ballotBox                   ballotBox;
    consensusReached            bool;
	campaignRetries				int; 	// retry again after unsuccessful campaing period
};

func campaignRandomWait (components *types.RootStructure) time.Duration {
	_randomWaitDuration := time.Duration (common.GenerateRandomInt (components.ConfigRoot.Config.Modules.Campaign.MinWait, 
		components.ConfigRoot.Config.Modules.Campaign.MaxWait))*time.Millisecond;

	return _randomWaitDuration;
}

func resetResources (components *types.RootStructure, resources *campaignManagerResources) {
	resources.consensusReached 	= false;
	resources.randomWait = campaignRandomWait (components);
	resources.randomWaitTimer = time.NewTimer (resources.randomWait);
	resources.votingTimeoutTimer	= time.NewTimer (resources.timeout);
	resources.createANewBallotBox ();
}

func initResources (components *types.RootStructure, resources *campaignManagerResources) {
	resources.campaignRetries = 1;
	resources.myNodeId 			= node.ThisNode().Id;
	resources.consensusReached 	= false;

	resources.campaignCommsChannel = make (types.RaftMessageChannel);
	// initialise the timer values and start-stop timers so they don't
	// cause panic in the select
	resources.randomWait = time.Duration (common.GenerateRandomInt (components.ConfigRoot.Config.Modules.Campaign.MinWait, 
									  components.ConfigRoot.Config.Modules.Campaign.MaxWait))*time.Millisecond;
	resources.randomWaitTimer = time.NewTimer (resources.randomWait);
	resources.randomWaitTimer.Stop ();
	resources.timeout = time.Duration(components.ConfigRoot.Config.Modules.Campaign.MaxWait*2) * time.Millisecond;
	resources.votingTimeoutTimer	= time.NewTimer (resources.timeout);
	resources.votingTimeoutTimer.Stop ();

	resources.createANewBallotBox ();
}

func startANewCampaign (components *types.RootStructure, resources *campaignManagerResources) {
	components.Internal.Role = types.ROLE_Campaign;
	log.Printf ("workflow.controller.startANewCampaign: attempt no %d\n", resources.campaignRetries);
	resetResources (components, resources);
	resources.campaignRetries++;

}

func nominationReceived (components *types.RootStructure, resources *campaignManagerResources) {
	log.Printf ("workflow.controller.nominationReceived: role %v: Candidate %s has nominated itself for leader. randomWait %dms\n", 
	                   components.Internal.Role, resources.candidate, resources.randomWait/1000000);
	components.Internal.Role = types.ROLE_Voting;
	// in case it is a message from another node and not me, stop the randomWait timer
	// and start the overall campaign timeout so we don't get stock waiting for consensus forever
	resources.randomWaitTimer.Stop ();
//	resources.votingTimeoutTimer	= time.NewTimer (resources.timeout);
	voteForThisCandidate (components, resources.candidate);
}


func voteRecieved (components *types.RootStructure, resources *campaignManagerResources) (int, bool) {
	// have received a vote card for a candidate
	var  _candidateHasReachedMajority	bool;
	var  _noOfVotes int;

	_noOfVotes, _candidateHasReachedMajority = resources.ballotBox.addVote (resources.candidate, components.Internal.Consensus);
	
	return _noOfVotes, _candidateHasReachedMajority;
}

func stopCampaign (components *types.RootStructure, resources *campaignManagerResources) {
	resources.randomWaitTimer.Stop ();
	resources.votingTimeoutTimer.Stop ();
}

func CampaignManager (components *types.RootStructure) error {
	log.Println ("workflow.controller.CampaignManager: started");
	var _resources	campaignManagerResources;

	comms.InitCampaign (&components.Comms);
	initResources (components, &_resources);

	go comms.ListenToCampaignBroadcast (&components.Comms, _resources.campaignCommsChannel);
	for {
		select {
			case _signal := <- components.Internal.MainToCampaignChannel :
				log.Printf ("CampaignManager: received signal from rootChannel %v\n", _signal);
				switch _signal {
				case types.Ch_StartCampaign :
					if (components.Internal.Role == types.ROLE_Campaign) {
						startANewCampaign (components, &_resources);
					}
				}
			case _message := <- _resources.campaignCommsChannel :
				log.Printf ("CampaignManager: received message from _campaignChannel %v candidate=%s role=%d\n", 
							_message, _resources.candidate, int (components.Internal.Role));
				switch _message.Header.Action {
					case types.AT_CampaignLaunch : {
						receivedACampaignLaunch (components, &_resources, _message)
					}
					case types.AT_VoteLodged : {
						processANewVote (components, &_resources, _message);
					}
					case types.AT_NewLeaderElected : {
						newLeaderElected (components, &_resources, _message);
					}
					case types.AT_ResetCampaign : {
						startANewCampaign (components, &_resources);
					}
				}
			case <- _resources.randomWaitTimer.C:
				log.Printf ("CampaignManager: random wait timer has triggered. node %s\n", _resources.myNodeId);
				launchMyCampaign (components, &_resources);

			case <- _resources.votingTimeoutTimer.C :
				// send a message to reset the campaign timing
				requestCampaignReset (components, &_resources);
		}
	}
}

func receivedACampaignLaunch (components *types.RootStructure, resources *campaignManagerResources, message *types.RaftMessage) {
	log.Printf ("CampaignManager: received message AT_CampaignLaunch received. role=%d", components.Internal.Role);
	if (components.Internal.Role == types.ROLE_Campaign) {
		log.Printf ("CampaignManager: nominationReceived\n");
		var _msgPayload  types.CampaignPayload = message.Payload.(types.CampaignPayload);
		resources.candidate = _msgPayload.Candidate;
		nominationReceived (components, resources);
	}

}

func requestCampaignReset (components *types.RootStructure, resources *campaignManagerResources) error {
	var _err error;

	log.Printf ("CampaignManager: _votingTimeoutTimer has triggered with no consensus. requesting a campaign reset node %s\n", resources.myNodeId);
	_message := message.CreateACampaignResetMessage ();
	_err = comms.SendToCampaignBroadcast (&components.Comms, _message);
	if (_err != nil) {startANewCampaign (components, resources);}

	return _err;
}


func processANewVote (components *types.RootStructure, resources *campaignManagerResources, message *types.RaftMessage) {
	log.Printf ("CampaignManager:VoteLodged received: role=%d candidate=%s\n", 
	components.Internal.Role, resources.candidate);
	var _msgPayload  types.CampaignPayload = message.Payload.(types.CampaignPayload);

	if (components.Internal.Role == types.ROLE_Voting) && (_msgPayload.Candidate == resources.candidate) {
		// have received a vote card for a candidate
		_noOfVotes, _candidateHasReachedMajority := voteRecieved (components, resources);
		log.Printf ("CampaignManager: vote for %s has reached %d and has majority=%v\n", 
			  resources.candidate, _noOfVotes,_candidateHasReachedMajority);
		if (_candidateHasReachedMajority == true) {
		// _message := message.CreateNewLeaderMessage (_resources.candidate, _noOfVotes);
		// _resources.campaignCommsChannel <- _message;
		resources.votingTimeoutTimer.Stop ();
		// log.Printf ("workflow.controller.CampaignManager: %s has a majority votes %d\n", _resources.candidate, _noOfVotes);
		broadcastANewLeader (components, resources.candidate, _noOfVotes);
		}
	}

}


func newLeaderElected (components *types.RootStructure, resources *campaignManagerResources, message *types.RaftMessage) {
	log.Printf ("CampaignManager: new leader elected %s\n", resources.candidate);
	stopCampaign (components, resources);
	var _msgPayload  types.CampaignPayload = message.Payload.(types.CampaignPayload);
	if (_msgPayload.Candidate == node.ThisNode().Id) {
		components.Internal.Role = types.ROLE_Leader;
	} else {
		components.Internal.Role = types.ROLE_Follower;
	}
	components.Internal.MainToCampaignChannel <- types.Ch_NewLeader;
}

// ============================================================================
// Launch my Campaign
// ============================================================================
func launchMyCampaign (components *types.RootStructure, resources *campaignManagerResources) error {
	// create a Campaign notice. tell everyone you want to be the leader
	var _err 			error;		
	var _launchRequest 	*types.RaftMessage;

	_launchRequest = message.CreateALaunchMessage ();
	_err = comms.SendToCampaignBroadcast (&components.Comms, _launchRequest);

	return _err;
}

// ============================================================================
// Vote
// ============================================================================
func voteForThisCandidate (components *types.RootStructure, candidate string) error {
	var _err  		error;
	var _ballot		*types.RaftMessage;

	_ballot = message.CreateVotingBallot (candidate);
	_err = comms.SendToCampaignBroadcast (&components.Comms, _ballot);

	return _err;
}

// ============================================================================
type candidateCard struct {
	Candidate 	string;
	Votes 		int;
}
type ballotBox map [string] candidateCard;
// ============================================================================
// BallotBox. this keep track of votes for each candidate. in a perfect world
// there should only be one candidate but in rare ocassions 2 candidates can 
// nominate themselves at the same time
// ============================================================================
func (campaignStaff *campaignManagerResources) createANewBallotBox  () {
	campaignStaff.ballotBox = nil;
	campaignStaff.ballotBox = make (ballotBox);
}


// ============================================================================
// add the vote for this candidate to the ballot box and return true if it has 
// reached majority (concensus)
// ============================================================================
func (ballot ballotBox) addVote (candidate string, majorityNeeded int) (int,bool) {
	var _hasTheNumbers 	bool = false;

	_candidateCard, _ok := ballot [candidate];
	if (_ok == false) {
		_candidateCard = candidateCard {candidate, 0};
	}
	_candidateCard.Votes++;
	ballot [candidate] = _candidateCard;
	if (_candidateCard.Votes >= majorityNeeded) {_hasTheNumbers = true;}

	return _candidateCard.Votes,_hasTheNumbers;
}


// ============================================================================
// add the vote for this candidate to the ballot box and return true if it has 
// reached majority (concensus)
// ============================================================================
func broadcastANewLeader (components *types.RootStructure, candidate string, votes int) error {
	var _err   		error;
	var _message  	*types.RaftMessage;

	_message = message.CreateNewLeaderMessage (candidate, votes);
	_err = comms.SendToCampaignBroadcast (&components.Comms, _message);

	return _err;
}



// ============================================================================
// ============================================================================
// ============================================================================
// ========================== Follower Heartbeat ==============================
// ============================================================================
// ============================================================================


type followerHBResources struct {
	hearbeatChannel 	types.RaftMessageChannel;
	timeoutTimer 		*time.Timer;
	timeout				time.Duration;

}

func initFollowerHBManager (components *types.RootStructure, resources *followerHBResources) error {
	var _err  error;

	node.ResetNodeList ();
	_err = comms.InitFollowerHBComms (&components.Comms);
	if (_err == nil) {
		var _timeoutValue int = components.ConfigRoot.Config.Modules.Heartbeat.PingRcvTimeout;
		resources.timeout = time.Duration (_timeoutValue) * time.Millisecond;
		resources.hearbeatChannel = make (types.RaftMessageChannel);
		resources.timeoutTimer = time.NewTimer (resources.timeout);
	} else {
		components.Internal.MainToCampaignChannel <- types.Ch_InitialisationRrror;
	}

	return _err;
}

func deInitFollowerHBManager (components *types.RootStructure, resources *followerHBResources) {
	comms.DeInitFollowerHBComms (&components.Comms);
	resources.timeoutTimer.Stop ();
	resources.timeoutTimer = nil;
}

// ============================================================================
// Follower hearbeat manager
// ============================================================================
func followerHeartbeatManager (components *types.RootStructure) {
	var _err  error;
	var _resources  followerHBResources;
	defer deInitFollowerHBManager (components, &_resources);

	_err = initFollowerHBManager (components, &_resources);
	if (_err == nil) {
		go comms.WaitForHeartBeatPing (&components.Comms, _resources.hearbeatChannel);

		for {
			select {
				case _signal := <- components.Internal.MainToFolHBChannel : {
					switch _signal {
						case types.Ch_StopFollowerHB: {
							return;
						}
					}
				}
				case _message := <- _resources.hearbeatChannel: {
					log.Printf ("FollowerHeartbeatManager: received ping %v\n", _message);
					_response := message.CreateHeartBeatResponse(types.NS_Active);
					_err = comms.SendHeartbeatPong (&components.Comms, _response);
					_resources.timeoutTimer = time.NewTimer (_resources.timeout);
				}
				case <- _resources.timeoutTimer.C : {
					components.Internal.MainToFolHBChannel <- types.Ch_HeartbeatTimeout;
					return;
				}
			}
		}
	}
}


// ============================================================================
// ============================================================================
// ============================================================================
// ============================ Leader Heartbeat ==============================
// ============================================================================
// ============================================================================

type leaderHBResources struct {
	hearbeatChannel 	types.RaftMessageChannel;
	hbInterval			time.Duration;
	heartbeatTimer 		*time.Ticker;
	pongTimeout			time.Duration;
	pongTimeoutTimer	*time.Timer;
}

func initLeaderHBManager (components *types.RootStructure, resources *leaderHBResources) error {
	var   _err   	error;

	node.ResetNodeList ();
	_err = comms.InitLeaderHBComms(&components.Comms);
	if (_err == nil) {
		resources.hearbeatChannel = make (types.RaftMessageChannel);
		_tickerValue := components.ConfigRoot.Config.Modules.Heartbeat.PingInterval;
		resources.hbInterval = time.Duration(_tickerValue) * time.Millisecond;
		resources.heartbeatTimer = time.NewTicker (resources.hbInterval);
		_timerValue := components.ConfigRoot.Config.Modules.Heartbeat.PongRcvTimeout;
		resources.pongTimeout = time.Duration(_timerValue) * time.Millisecond;
		resources.pongTimeoutTimer = time.NewTimer (resources.pongTimeout);
		resources.pongTimeoutTimer.Stop ();
	}
	return _err;
}

func deInitLeaderHBManager (components *types.RootStructure, resources *leaderHBResources) {
	node.ResetNodeList ();
	comms.DeInitLeaderHBComms(&components.Comms);
	resources.heartbeatTimer.Stop ();
}

// ============================================================================
// Leader hearbeat manager
// ============================================================================
func leaderHeartbeatManager (components *types.RootStructure) {
	var _err  error;
	var _resources  leaderHBResources;

	_err = initLeaderHBManager (components, &_resources);
	defer deInitLeaderHBManager (components, &_resources);

	if (_err == nil) {
		go comms.WaitForHeartbeatPongWithChannel (&components.Comms, _resources.hearbeatChannel);
		var _pongCount int;
		for {
			select {
				case <- _resources.heartbeatTimer.C : {
					comms.SendHeartbeatPing (&components.Comms);
					_resources.pongTimeoutTimer = time.NewTimer (_resources.pongTimeout);
				}
				case _response := <- _resources.hearbeatChannel : {
					_pongCount++;
					_responder := _response.Header;
					updateNodesStatus (components, _responder);
					// log.Printf ("leaderHeartbeatManager: received pong from node %s\n", _responder);
				}
				case <- _resources.pongTimeoutTimer.C : {
					// check the list and mark them as inactive if they have not responded
					node.ReportLatestStati ();
					node.ResetAllNodesStati ();
					// log.Printf ("-> received %d pongs\n", _pongCount);
					_pongCount = 0;
				}
			}
		}
	}

}

func updateNodesStatus (components *types.RootStructure, header types.MessageHeader) {
	var _node 	*types.NodeMeta = node.GetNodeData (header.SourceNodeId);
	if (_node == nil) {
		_node = node.CreateNode (header.SourceNodeId, header.SourceNodeIp);
		node.AddNode (_node);
	}
	node.SetNodeStatus (_node, types.NS_Active);
}

