package pelldvs

import (
	"fmt"

	"github.com/0xPellNetwork/pelldvs-libs/log"
	aggRPC "github.com/0xPellNetwork/pelldvs/aggregator/rpc"
	avsitypes "github.com/0xPellNetwork/pelldvs/avsi/types"
	"github.com/0xPellNetwork/pelldvs/config"
	"github.com/0xPellNetwork/pelldvs/node"
	"github.com/0xPellNetwork/pelldvs/p2p"
	"github.com/0xPellNetwork/pelldvs/privval"
	"github.com/0xPellNetwork/pelldvs/proxy"
	rpclocal "github.com/0xPellNetwork/pelldvs/rpc/client/local"
)

type Node struct {
	agg     *aggRPC.RPCClientAggregator
	nodeCfg *config.Config
	logger  log.Logger
	node    *node.Node
}

func (n *Node) Start() error {
	if n.node == nil {
		return fmt.Errorf("node is nil")
	}

	return n.node.Start()
}

func NewNode(
	logger log.Logger,
	app avsitypes.Application,
	nodeCfg *config.Config,
) (*Node, error) {
	var n = &Node{
		logger:  logger,
		nodeCfg: nodeCfg,
	}

	var err error
	n.agg, err = aggRPC.NewRPCClientAggregator(nodeCfg.Pell.AggregatorRPCURL)
	if err != nil {
		logger.Error("Failed to create aggregator client", "error", err)
		return nil, fmt.Errorf("failed to create aggregator client: %v", err)
	}

	// Load or generate node key
	nodeKey, err := p2p.LoadOrGenNodeKey(n.nodeCfg.NodeKeyFile())
	if err != nil {
		logger.Error("Failed to load node key", "error", err)
		return nil, fmt.Errorf("failed to load node key: %v", err)
	}

	// Load or generate private validator
	pv, err := privval.LoadOrGenFilePV(n.nodeCfg.PrivValidatorKeyFile())
	if err != nil {
		return nil, fmt.Errorf("failed to load or gen file PV: %v", err)
	}

	n.node, err = node.NewNode(n.nodeCfg,
		pv,
		nodeKey,
		proxy.NewLocalClientCreator(app),
		config.DefaultDBProvider,
		n.agg,
		node.DefaultMetricsProvider(n.nodeCfg.Instrumentation),
		logger,
	)
	if err != nil {
		logger.Error("NewClient Failed to create node", "error", err)
		return nil, fmt.Errorf("failed to create node: %v", err)
	}
	return n, nil
}

func (n *Node) GetLocalClient() *rpclocal.Local {

	return rpclocal.New(n.node)
}
