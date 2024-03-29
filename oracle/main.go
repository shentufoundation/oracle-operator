// Package oracle defines oracle-operator
package oracle

import (
	"os"

	"github.com/spf13/cobra"

	tmconfig "github.com/tendermint/tendermint/config"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/server"

	"github.com/shentufoundation/oracle-operator/types"
)

// starts the oracle operator.
func start(ctx types.Context) {
	ctx = ctx.WithLoggerLabels("module", "oracle-operator")
	fatalError := make(chan error)
	ctkMsgChan := make(chan interface{}, 1000)
	go Listen(ctx.WithLoggerLabels("protocol", "shentu", "submodule", "listener"), ctkMsgChan, fatalError)
	go Push(ctx.WithLoggerLabels("protocol", "shentu", "submodule", "pusher"), ctkMsgChan, fatalError)
	// exit on fatal error
	err := <-fatalError
	ctx.Logger().Error(err.Error())
	os.Exit(1)
}

// ServeCommand will start the oracle operator as a blocking process.
func ServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start oracle operator",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf := tx.NewFactoryCLI(cliCtx, cmd.Flags()).WithTxConfig(cliCtx.TxConfig)

			cliCtx.SkipConfirm = true // TODO: new cosmos version

			ctx, err := types.NewContextWithDefaultConfigAndLogger()
			if err != nil {
				return err
			}
			ctx = ctx.WithClientContext(&cliCtx).WithTxFactory(&txf)

			if err := serve(ctx); err != nil {
				return err
			}
			return nil
		},
	}

	return registerFlags(cmd)
}

// registerFlags registers additional flags to the command.
func registerFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().String(flags.FlagChainID, "", "The network chain ID")
	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().Uint(types.FlagRPCReadTimeout, 10, "RPC read timeout (in seconds)")
	cmd.Flags().Uint(types.FlagRPCWriteTimeout, 10, "RPC write timeout (in seconds)")
	cmd.Flags().String(types.FlagLogLevel, tmconfig.DefaultLogLevel, "Log level")
	cmd.Flags().String(types.FlagConfigFile, types.DefaultConfigFileName, "Name of the config file")
	return cmd
}

// serve sets up oracle operator running environment.
func serve(ctx types.Context) error {
	done := make(chan struct{})
	panicChan := make(chan interface{}, 1)

	server.TrapSignal(func() {
		done <- struct{}{}
	})

	ctx.Logger().Info("Starting Oracle Operator...")
	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()
		start(ctx)
	}()

	defer close(done)
	select {
	case p := <-panicChan:
		panic(p)
	case <-done:
		ctx.Logger().Info("Shutting Down Oracle Operator...")
	case <-ctx.Context().Done():
		ctx.Logger().Info("Ending Oracle Operator...")
	}
	return nil
}
