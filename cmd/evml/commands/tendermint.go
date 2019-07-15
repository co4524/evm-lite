package commands

import (
	"fmt"
	"path/filepath"

	"github.com/bear987978897/evm-lite/src/consensus/tendermint"
	"github.com/bear987978897/evm-lite/src/engine"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tmConfig "github.com/tendermint/tendermint/config"
)

// AddTendermintFlags adds flags to the Tendermint command
func AddTendermintFlags(cmd *cobra.Command) {
	// cmd.Flags().String("tendermint.home", config.Tendermint.DataDir, "Tendermint home directory")
}

// Viber load config
func tmBindFlagsLoadViper(cmd *cobra.Command) error {
	// cmd.Flags() includes flags from this command and all persistent flags from the parent
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	homeDir := viper.GetString("datadir") + "/tendermint"
	// viper.Set("tendermint.home", homeDir)
	viper.SetConfigName("config")                         // name of config file (without extension)
	viper.AddConfigPath(homeDir)                          // search root directory
	viper.AddConfigPath(filepath.Join(homeDir, "config")) // search root directory /config

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// stderr, so if we redirect output to json file, this doesn't appear
		// fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
		// ignore not found error, return other errors
		return err
	}
	return nil
}

func ParseConfig() (*tmConfig.Config, error) {
	conf := tmConfig.DefaultConfig()
	err := viper.Unmarshal(conf)
	if err != nil {
		return nil, err
	}

	conf.SetRoot(viper.GetString("datadir") + "/tendermint")
	tmConfig.EnsureRoot(conf.RootDir)
	if err = conf.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("Error in config file: %v", err)
	}
	return conf, err
}

// NewTendermintCmd returns the command that starts EVM-Lite with Tendermint consensus
func NewTendermintCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tendermint",
		Short: "Run the evm-lite node with Tendermint consensus",
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {
			config.SetDataDir(config.BaseConfig.DataDir)
			logger.WithFields(logrus.Fields{
				"Tendermint": config.Tendermint,
			}).Debug("Config")

			tmBindFlagsLoadViper(cmd)
			config.Tendermint.RealConfig, _ = ParseConfig()
			return nil
		},
		RunE: runTendermint,
	}

	AddTendermintFlags(cmd)

	return cmd
}

func runTendermint(cmd *cobra.Command, args []string) error {
	tendermint := tendermint.NewTendermint(config.Tendermint, logger)
	engine, err := engine.NewEngine(*config, tendermint, logger)
	if err != nil {
		return fmt.Errorf("Error building Engine: %s", err)
	}

	engine.Run()

	return nil
}
