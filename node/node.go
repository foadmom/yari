package node

import (
	"fmt"
	"log"
	"net"

	u "github.com/foadmom/yari/common"
	"github.com/foadmom/yari/types"
)

type NodeListType map[string]*types.NodeMeta

var NodeList NodeListType = make(map[string]*types.NodeMeta)

var thisNode types.NodeMeta = types.NodeMeta{}

func ThisNode() *types.NodeMeta {
	return &thisNode
}

func ThisNodeInit(nodeId string) {
	var _node *types.NodeMeta = ThisNode()

	SetNodeNetIntDetails(_node)
	//	CreateId(_node);
	_node.Id = nodeId
	// _node.CampaignDelay = u.GenerateRandomInt(types.NT_MinTimout, types.NT_MaxTimout);
	_node.Status = types.NS_Active
}

func CreateId(node *types.NodeMeta) {
	node.Id = node.IpAddress
}

func SetNodeNetIntDetails(node *types.NodeMeta) {
	var _netIP net.IP
	_netIP, node.NetInterfaceName, node.MacAddress = u.GetMyNetInterfaceDetails()
	node.IpAddress = _netIP.String()
}

func GetMyNodeStatus() types.NodeStatus {
	var _node *types.NodeMeta = ThisNode()
	return _node.Status
}

//=============================================================================
// create a NodeMeta object
//=============================================================================
func CreateNode(nodeId string, nodeIp string) *types.NodeMeta {
	var _node types.NodeMeta
	_node.Id = nodeId
	_node.IpAddress = nodeIp

	return &_node
}

//=============================================================================
// set NodeStatus
//=============================================================================
func SetNodeStatus(node *types.NodeMeta, status types.NodeStatus) {
	node.Status = status
}

//=============================================================================
// add/update a node. if the node exists already, it will be overwritten
//=============================================================================
func AddNode(node *types.NodeMeta) error {
	var _err error = nil
	if len(node.Id) > 0 {
		NodeList[node.Id] = node
	} else {
		_err = fmt.Errorf("can NOT AddNode with node.Id=nil")
	}
	return _err
}

func RemoveNode(id string) {
	delete(NodeList, id)
}

func GetNodeList() NodeListType {
	return NodeList
}

func GetNodeData(nodeId string) *types.NodeMeta {
	var _node *types.NodeMeta = NodeList[nodeId]

	return _node
}

func ResetNodeList() {
	NodeList = make(map[string]*types.NodeMeta)
}

//=============================================================================
// this is used at the start of heartbeat so we can check the status of each
// node at the end of the heartbeat period to see which one has not responded.
//=============================================================================
func ResetAllNodesStati() {
	for _, _node := range NodeList {
		if _node.Status == types.NS_Active {
			_node.Status = types.NS_Unknown
		} else {
			// log.Printf ("node %s status is %s\n", _nodeId, string(_node.Status));
		}
	}
}

func ReportLatestStati() {
	for _, _node := range NodeList {
		log.Printf("node %s status is %s\n", _node.Id, string(_node.Status))
	}
}
