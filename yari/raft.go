package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/foadmom/yari/comms"
	n "github.com/foadmom/yari/node"
	s "github.com/foadmom/yari/storage"
	"github.com/foadmom/yari/types"

	w "github.com/foadmom/yari/workflow"

	"gopkg.in/yaml.v2"
)

var roleString []string = []string{"", "Leader", "Follower"}

var configFile = flag.String("configFile", "/data/workspaces/go/src/yari/config.yaml", "usage:  -configFile=<full path and filename>")

// var configFile  = flag.String ("configFile", "./config.yaml", "usage:  -configFile=<full path and filename>");
var initialRole = flag.Int("initialRole", 2, "usage:  -initialRole=1/2")
var nodeId = flag.String("nodeId", "", "usage:  -nodeId=<node id")

func main() {
	var _err error
	var _components types.RootStructure

	defer onClose(&_components)
	flag.Parse()
	fmt.Printf("YARI Starting node %s as a %s \n", *nodeId, roleString[*initialRole])

	if *nodeId == "" {
		log.Println("-nodeId argument is mandatory. ")
		os.Exit(0)
	}
	_err = InitConfigs(&_components, *configFile)

	_components.Internal.Role = types.RoleType(*initialRole)
	_components.Internal.NodeId = *nodeId

	initialise(&_components)

	handleCntlC(&_components)

	if _err == nil {
		_err = s.InitStorage(&_components.Storage, &_components.ConfigRoot.Config.Modules.Storage)
		n.ThisNodeInit(*nodeId)

		if _err == nil {
			initialise(&_components)
			w.MainFlowManager(&_components)
		}
	}

	fmt.Println("Exiting")
}

func InitConfigs(components *types.RootStructure, configFileName string) error {
	var _buffer []byte
	var _err error
	//	var _config	types.ConfigType;

	_buffer, _err = readConfigFile(configFileName)
	if _err == nil {
		//		fmt.Println (string(_buffer));
		//		_err = json.Unmarshal(_buffer, &components.ConfigRoot);
		_err = yaml.Unmarshal(_buffer, &components.ConfigRoot)
	}

	return _err
}

func readConfigFile(configFileName string) ([]byte, error) {
	var _buffer []byte

	_configFile, _err := os.Open(configFileName)
	defer _configFile.Close()
	if _err == nil {
		_buffer, _err = ioutil.ReadAll(_configFile)
	}
	return _buffer, _err
}

// ============================================================================
//
// ============================================================================
func initialise(components *types.RootStructure) error {
	var _err error

	components.Internal.NodeCount = components.ConfigRoot.Config.General.NodeCount
	components.Internal.Consensus = components.ConfigRoot.Config.General.Consensus

	components.Comms.Config = &components.ConfigRoot.Config.Modules.Networks
	_err = comms.InitComms(&components.Comms)

	return _err
}

func handleCntlC(components *types.RootStructure) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		onClose(components)
		os.Exit(1)
	}()
}

func onClose(components *types.RootStructure) {
	fmt.Println("cleaning up before exit")
	comms.OnClose(&components.Comms)
}

func tempScaffold_nodeList(thisNode *types.NodeMeta) {
	n.AddNode(thisNode)
	fmt.Println(thisNode)
	tempScaffold_addNodes()
}

func tempScaffold_addNodes() {
	n.AddNode(tempScaffold_createNode("Node_EFGH"))
	n.AddNode(tempScaffold_createNode("Node_IJKL"))
	fmt.Println(n.GetNodeList())
}

func tempScaffold_createNode(id string) *types.NodeMeta {
	var _newNode types.NodeMeta = types.NodeMeta{Id: id}
	return &_newNode
}
