package cmd

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	tmcfg "github.com/tendermint/tendermint/config"
	tmcli "github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/server"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/shentufoundation/shentu/v2/app"
	"github.com/shentufoundation/shentu/v2/app/params"
	"github.com/shentufoundation/shentu/v2/common"

	"github.com/shentufoundation/oracle-operator/oracle"
)

// NewRootCmd creates a new root command. It is called once in the
// main function.
func NewRootCmd() (*cobra.Command, params.EncodingConfig) {
	encodingConfig := app.MakeEncodingConfig()
	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithHomeDir(app.DefaultNodeHome)

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(common.Bech32PrefixAccAddr, common.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(common.Bech32PrefixValAddr, common.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(common.Bech32PrefixConsAddr, common.Bech32PrefixConsPub)
	config.Seal()

	rootCmd := &cobra.Command{
		Use:   "oracle-operator",
		Short: "CertiK Chain Oracle Operator",
		Long:  `CertiK Oracle Operator listens to create_task events from CertiK Chain, queries primitives, and pushes result back to the chain.`,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}
			template, customAppConfig := initAppConfig()
			return server.InterceptConfigsPreRunHandler(cmd, template, customAppConfig)
		},
	}

	initRootCmd(rootCmd)

	return rootCmd, encodingConfig
}

func initRootCmd(rootCmd *cobra.Command) {
	rootCmd.AddCommand(
		oracle.ServeCommand(),
		keys.Commands(app.DefaultNodeHome),
		version.NewVersionCommand(),
	)
}

func Execute(rootCmd *cobra.Command, defaultHome string) error {
	// Create and set a client.Context on the command's Context. During the pre-run
	// of the root command, a default initialized client.Context is provided to
	// seed child command execution with values such as AccountRetriver, Keyring,
	// and a Tendermint RPC. This requires the use of a pointer reference when
	// getting and setting the client.Context. Ideally, we utilize
	// https://github.com/spf13/cobra/pull/1118.
	srvCtx := server.NewDefaultContext()
	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &client.Context{})
	ctx = context.WithValue(ctx, server.ServerContextKey, srvCtx)

	rootCmd.PersistentFlags().String(flags.FlagLogLevel, zerolog.InfoLevel.String(), "The logging level (trace|debug|info|warn|error|fatal|panic)")
	rootCmd.PersistentFlags().String(flags.FlagLogFormat, tmcfg.LogFormatPlain, "The logging format (json|plain)")

	executor := tmcli.PrepareBaseCmd(rootCmd, "", defaultHome)
	return executor.ExecuteContext(ctx)
}

// initAppConfig helps to override default appConfig template and configs.
// return "", nil if no custom configuration is required for the application.
func initAppConfig() (string, interface{}) {
	type CustomAppConfig struct {
		serverconfig.Config
	}

	// Optionally allow the chain developer to overwrite the SDK's default
	// server config.
	srvCfg := serverconfig.DefaultConfig()
	srvCfg.API.Enable = true
	srvCfg.StateSync.SnapshotInterval = 1500
	srvCfg.StateSync.SnapshotKeepRecent = 2
	srvCfg.MinGasPrices = "0.025" + common.MicroCTKDenom
	srvCfg.IAVLDisableFastNode = false

	AppTemplate := serverconfig.DefaultConfigTemplate

	return AppTemplate, CustomAppConfig{*srvCfg}
}
