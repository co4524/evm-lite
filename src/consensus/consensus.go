package consensus

import (
	"github.com/bear987978897/evm-lite/src/service"
	"github.com/bear987978897/evm-lite/src/state"
)

// Consensus is the interface that abstracts the consensus system
type Consensus interface {
	Init(*state.State, *service.Service) error
	Run() error
	Info() (map[string]string, error)
}
