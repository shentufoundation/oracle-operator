package oracle

import (
	"bufio"
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CompleteAndBroadcastTx is adopted from auth.CompleteAndBroadcastTxCLI. The original function prints out response.
func CompleteAndBroadcastTx(cliCtx client.Context, txf tx.Factory, msgs []sdk.Msg) (sdk.TxResponse, error) {
	txf, err := tx.PrepareFactory(cliCtx, txf)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	if txf.SimulateAndExecute() || cliCtx.Simulate {
		_, adjusted, err := tx.CalculateGas(cliCtx.QueryWithData, txf, msgs...)
		if err != nil {
			return sdk.TxResponse{}, err
		}

		txf = txf.WithGas(adjusted)
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", tx.GasEstimateResponse{GasEstimate: txf.Gas()})
	}

	if cliCtx.Simulate {
		return sdk.TxResponse{}, nil
	}

	unsignedTx, err := tx.BuildUnsignedTx(txf, msgs...)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	if !cliCtx.SkipConfirm {
		out, err := cliCtx.TxConfig.TxJSONEncoder()(unsignedTx.GetTx())
		if err != nil {
			return sdk.TxResponse{}, err
		}

		_, _ = fmt.Fprintf(os.Stderr, "%s\n\n", out)

		buf := bufio.NewReader(os.Stdin)
		ok, err := input.GetConfirmation("confirm transaction before signing and broadcasting", buf, os.Stderr)

		if err != nil || !ok {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", "canceled transaction")
			return sdk.TxResponse{}, err
		}
	}

	err = tx.Sign(txf, cliCtx.GetFromName(), unsignedTx, true)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	txBytes, err := cliCtx.TxConfig.TxEncoder()(unsignedTx.GetTx())
	if err != nil {
		return sdk.TxResponse{}, err
	}

	// broadcast to a Tendermint node
	res, err := cliCtx.BroadcastTx(txBytes)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	return *res, nil
}
