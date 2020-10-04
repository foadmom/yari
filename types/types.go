package types

import (
	"encoding/json"
	"net"

	"log"
	"time"

	"qpid.apache.org/electron"
)

type Buffer 	[]byte;

type RoleType	int;
const (
	ROLE_Leader				RoleType = 1;
	ROLE_Follower			RoleType = 2;
	ROLE_Campaign			RoleType = 3;
	ROLE_Voting 			RoleType = 4;
	ROLE_WaitingForLeader	RoleType = 5;
)

type Ch_SignalType int;
const (
	Ch_Cancel 				Ch_SignalType = -1;
	Ch_InitialisationRrror  Ch_SignalType = -2;

	Ch_StartFollowerHB   	Ch_SignalType = 100;
	Ch_StopFollowerHB 		Ch_SignalType = 101;
	Ch_StartLeaderHB   		Ch_SignalType = 102;
	Ch_StopLeaderHB 		Ch_SignalType = 103;
	Ch_ReceivedPing      	Ch_SignalType = 109;
	Ch_HeartbeatTimeout		  Ch_SignalType = 110;
	Ch_ErrorReceivingReq 	  Ch_SignalType = 112;
	Ch_ErrorReceivingResp 	  Ch_SignalType = 113;
	Ch_ErrorSendingResp 	  Ch_SignalType = 114;
	Ch_ErrorSendingReq  	  Ch_SignalType = 115;

	Ch_StartCampaign		Ch_SignalType = 200;
	Ch_CampaignMode			Ch_SignalType = 201;
	Ch_NewLeader			Ch_SignalType = 202;
	Ch_ConsensusNotReached  Ch_SignalType = 203;
);
type SignalChannel		chan Ch_SignalType;
type RaftMessageChannel	chan *RaftMessage;




type ConfigType struct {
	Config ConfigTreeType     `yaml:"Config"`;
}

type ConfigTreeType struct {
	General ConfigGeneralType		 `yaml:"General"` ;
	Modules ConfigModuleType         `yaml:"Modules"` ;
}

type ConfigGeneralType struct {
	Debug     bool    `yaml:"Debug"`
	Version   string  `yaml:"Version"`
	Build     string  `yaml:"Build"`
	Consensus int     `yaml:"Consensus"`

	Role			RoleType;  // Follower or Leader
	NodeId			string;
	NodeCount		int;	// total number of nodes in the system
}

// ==============
type ConfigModuleType	struct {
	Networks ConfigNetworksType  	`yaml:"Networks"`;
	Heartbeat ConfigHeartbeatType	`yaml:"Heartbeat"`;
	Campaign ConfigCampaignType  	`yaml:"Campaign"`;
	Security ConfigSecurityType  	`yaml:"Security"`;
	Storage  ConfigStorageType    	`yaml:"Storage"`;
}

type ConfigNetworksType struct {
	Transport ConfigTransportType		`yaml:"Transport"`;
}

		type ConfigTransportType struct {
			ContainerName      string `yaml:"ContainerName"`;
			Host               string `yaml:"Host"`;
			Port               string `yaml:"Port"`;
			User               string `yaml:"User"`;
			Password           string `yaml:"Password"`;
			HeartbeatPingQueue string `yaml:"HeartbeatPingQueue"`;
			HeartbeatPongQueue string `yaml:"HeartbeatPongQueue"`;
			CampaignQueue      string `yaml:"CampaignQueue"`;
			VotingQueue        string `yaml:"VotingQueue"`;
			ListenerPrefetch   int    `yaml:"ListenerPrefetch"`;
			MessageExpiry      int64  `yaml:"MessageExpiry"`;
		}

		type ConfigHeartbeatType struct {
			PingInterval   int `yaml:"PingInterval"`;
			PongRcvTimeout int `yaml:"PongRcvTimeout"`;
			PingRcvTimeout int `yaml:"PingRcvTimeout"`;
		}
		
		

type ConfigCampaignType struct {
	MinWait    int `yaml:"MinWait"`;
	MaxWait    int `yaml:"MaxWait"`;
    RetryWait  int `yaml:"RetryWait"`;
}

type ConfigSecurityType struct {
	SSL bool `yaml:"SSL"` ;
}

type ConfigStorageType struct {
	LogDir string 	`yaml:"LogDir"`;
}

// ==================================================================
// ==================================================================
// ==================================================================
// internal structures used
// ==================================================================
// ==================================================================
// ==================================================================


// ==================================================================
// this is the top level structure that should include all other 
// configs and structures used
type RootStructure struct {
	ConfigRoot		ConfigType
	Internal		InternalType;
	Comms 			CommsType;
	Storage			StorageType;
}

type InternalType struct {
	Role			    RoleType;
	NodeId				string;
	NodeCount			int;
	Consensus			int;

	// internal channels
	MainToCampaignChannel			SignalChannel;
	MainToFolHBChannel				SignalChannel;
	MainToFLeadHBChannel			SignalChannel;
	HBFollowerChannel	RaftMessageChannel;
	HBLeaderChannel		RaftMessageChannel;
	CampaignChannel		RaftMessageChannel;
}

type StorageType struct {
	Config			*ConfigStorageType;
	SystemLog		*log.Logger;
}

type CommsType	struct {
	Config			*ConfigNetworksType;
	Transport		TransportType;
}


type TransportType struct {
	Config  		    	ConfigTransportType;
	Container				electron.Container;
	Connection			    electron.Connection
	// for follower
	HeartbeatPingListener	electron.Receiver;
	HeartbeatPongSender  	electron.Sender;
	// and for leader
	HeartbeatPingSender 	electron.Sender;
	HeartbeatPongListener	electron.Receiver;

	CampaignSender			electron.Sender;
	CampaignListener		electron.Receiver;

	VotingSender			electron.Sender;
	VotingListener			electron.Receiver;
}



type ActionClassType 	string;
const (
	AC_HeartBeat			ActionClassType = "HeartBeat";
	AC_Campaign				ActionClassType = "Campaign";
);

type ActionType		string;
const (
	// hearbeat
	AT_HeartBeatRequest		ActionType = "HeartBeatPing";
	AT_HeartBeatResponse	ActionType = "HeartBeatPong";
	// maintenance of node list. add, remove, change status of nodes
	AT_NodeAdded			ActionType = "NodeAdded";
	AT_NodeRemoved			ActionType = "NodeRemoved";
	AT_NodeStatus			ActionType = "NodeStatus";
	// campaign
	AT_CampaignLaunch		ActionType = "CampaignLaunch";
	AT_LeadershipVote		ActionType = "LeadershipVote";
	AT_VoteLodged     		ActionType = "VoteLodged";
	AT_NewLeaderElected		ActionType = "NewLeaderElected";
	AT_NoElectionResult		ActionType = "NoElectionResult";
	AT_ResetCampaign		ActionType = "ResetCamaign";
);


type MessageHeader struct {
	Uuid			string;		// generated for every message
	RequestId		string      // this requestId is generated by the requester and should be 
								// returned in the same field of the response message.
								// think of it as a txn id which can span over 
								// multiple messages.
	Timestamp		time.Time;	// in the form of nanoseconds (time.Now().UnixNano())
	SequenceId		uint64;		// generated by the leader. not sure how it will be used yet.
	SourceNodeIp	string;		// for future use. not sure how to use it.
								// can a node change it's hardware and run on a new machine while
								// maintaining the same nodeId ?
	SourceNodeId	string; 	// default node id is the mac address to make it unique.
								// multiple nodes running on docker on the same machine will
								// produce different unique mac addresses as default.
								// how nodeId is generated can be over-ridden.
	ActionClass		ActionClassType;
	Action			ActionType;

};

type MessagePayload interface{};
type RaftMessage struct {
	MessageType	string;
	Header		MessageHeader;
	Payload		MessagePayload			`json:",omitempty"`;
};

type JsonIntermediateMessage struct {
	MessageType	string;
	Header		MessageHeader;
	Payload		json.RawMessage 		`json:",omitempty"`;
};

type NullPayload struct {

};


// ==============================================
// heartbeat
// ==============================================
type HeartBeatPingPayload struct {
	ResponseQueue 		string;	
};
type HeartBeatPongPayload struct {
	NodeId 		string;
	Status 		NodeStatus;	
};

type HeartbeatStatus struct {
	NodeId 				string;
	Status 				string;
	LastPongTimestamp	time.Time;
}


// ==============================================
// Campaign
// ==============================================
// type VoteForACandidatePayload struct {
// 	Candidate 		string;
// }

type CampaignPayload struct {
	Action 			ActionType;
	Candidate		string; // nodeId of the candidate. it can also be my own nodeId
	VoteCount		int;
}
type CampaignChannel     chan *RaftMessage;


// ==============================================
// Node
// ==============================================
type NodeIdType		string;
type NodeGroupType  string;

type NodeStatus 	string ;
const (
	NS_Unknown	NodeStatus = "Unknown";
	NS_Paused   NodeStatus = "Paused";
	NS_Active   NodeStatus = "Active";
);

const (
	NT_MinTimout	int = 50;
	NT_MaxTimout	int = 500;
)

// this is the definition of node in this raft implementation.
// structures, and interfaces to be defined here

type NodeMeta struct {
	Id					string;
	Group				NodeGroupType;		// for grouping of certain nodes
	Status				NodeStatus;
	IpAddress			string;
	MacAddress			net.HardwareAddr;
	NetInterfaceName	string;
}

 type NodeInt interface {
 }
