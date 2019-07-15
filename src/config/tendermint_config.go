package config

import (
	"fmt"

	tmConfig "github.com/tendermint/tendermint/config"
)

var (
	defaultTmDir = fmt.Sprintf("%s/tendermint", DefaultDataDir)
)

type TmConfig struct {
	DataDir    string `mapstructure:"datadir"`
	RealConfig *tmConfig.Config
}

// DefaultTmConfig returns the default configuration for a Babble node
func DefaultTmConfig() *TmConfig {
	var conf = &TmConfig{
		DataDir:    defaultTmDir,
		RealConfig: tmConfig.DefaultConfig().SetRoot(defaultTmDir),
	}
	return conf
}

// SetDataDir updates the tendermint configuration directories if they were set to
// to default values.
func (c *TmConfig) SetDataDir(datadir string) {
	if c.DataDir == defaultTmDir {
		c.DataDir = datadir
		c.RealConfig.SetRoot(datadir)
	}
}

// ToRealTmConfig converts the config to real Tendermint config
func (c *TmConfig) ToRealTmConfig() *tmConfig.Config {
	return c.RealConfig
}
