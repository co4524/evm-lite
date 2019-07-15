package main

import (
	cmd "github.com/bear987978897/evm-lite/cmd/evml/commands"
)

func main() {

	rootCmd := cmd.RootCmd

	rootCmd.AddCommand(
		cmd.NewSoloCmd(),
		cmd.NewBabbleCmd(),
		cmd.NewRaftCmd(),
		cmd.NewTendermintCmd(),
		cmd.VersionCmd)

	//Do not print usage when error occurs
	rootCmd.SilenceUsage = true

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
