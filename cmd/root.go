package cmd

import (
	"github.com/spf13/cobra"

	"github.com/certikfoundation/oracle-toolset/oracle"
)

var rootCmd = &cobra.Command{
	Use:   "oracle-toolset",
	Short: "A generator for Cobra based Applications",
	Long: `Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd = &cobra.Command{
		Use:   "certik-oracle",
		Short: "CertiK Chain Oracle Operator",
		Long: `CertiK Oracle Operator listens to the create_task event from CertiK Chain,
		queries the primitives, and pushes the result back the chain.`,
	}

	//initRootCmd(rootCmd)
	rootCmd.AddCommand(oracle.ServeCommand())
	//rootCmd.AddCommand(versionCmd)
}
