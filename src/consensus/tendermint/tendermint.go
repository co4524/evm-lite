package tendermint

import (
	"fmt"
	"os"

	"github.com/bear987978897/evm-lite/src/config"
	"github.com/bear987978897/evm-lite/src/service"
	"github.com/bear987978897/evm-lite/src/state"
	"github.com/sirupsen/logrus"
	"github.com/tendermint/tendermint/libs/log"
	node "github.com/tendermint/tendermint/node"
	p2p "github.com/tendermint/tendermint/p2p"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

type Tendermint struct {
	config     *config.TmConfig
	node       *node.Node
	ethService *service.Service
	ethState   *state.State
	logger     *logrus.Logger
	rpcclient  *rpcclient.HTTP
}

func NewTendermint(config *config.TmConfig, logger *logrus.Logger) *Tendermint {
	return &Tendermint{
		config: config,
		logger: logger,
	}
}

/*******************************************************************************
IMPLEMENT CONSENSUS INTERFACE
*******************************************************************************/

func (t *Tendermint) Init(state *state.State, service *service.Service) error {

	t.logger.Debug("INIT")

	realConfig := t.config.ToRealTmConfig()
	abciApp := NewABCIProxy(state, t.logger)
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	node, err := DefaultNewNodeWithApp(realConfig, abciApp, logger)
	if err != nil {
		return err
	}

	t.ethState = state
	t.ethService = service
	t.node = node
	t.rpcclient = rpcclient.NewHTTP(realConfig.RPC.ListenAddress, "/websocket")

	return nil
}

func (t *Tendermint) Run() error {
	submitCh := t.ethService.GetSubmitCh()
	if err := t.node.Start(); err != nil {
		fmt.Println(err)
		return fmt.Errorf("Failed to start node: %v", err)
	}
	t.logger.Info("Started node", "nodeInfo", t.node.Switch().NodeInfo())

	for {
		select {
		case tx := <-submitCh:
			_, err := t.rpcclient.BroadcastTxSync(tx)

			if err != nil {
				fmt.Errorf("Failed to broacast transaction: %v", err)
			}
		}
	}

	return nil
}

func (t *Tendermint) Info() (map[string]string, error) {
	tmInfo := t.node.NodeInfo().(p2p.DefaultNodeInfo)
	info := map[string]string{
		"type":      "tendermint",
		"id":        string(tmInfo.ID()),
		"laddr":     tmInfo.ListenAddr,
		"network":   tmInfo.Network,
		"version":   tmInfo.Version,
		"moniker":   tmInfo.Moniker,
		"tx_index":  tmInfo.Other.TxIndex,
		"rpc_laddr": tmInfo.Other.RPCAddress,
	}
	return info, nil
}
