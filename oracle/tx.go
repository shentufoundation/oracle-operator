package oracle

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/hyperledger/burrow/crypto"
	"github.com/hyperledger/burrow/execution/evm/abi"
	"github.com/hyperledger/burrow/logging"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/certikfoundation/shentu/common"
	"github.com/certikfoundation/shentu/x/cvm/compile"

	"github.com/certikfoundation/oracle-toolset/types"
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

// callContract calls contract on certik-chain.
func callContract(ctx types.Context, callee string, function string, args []string) (bool, string, error) {
	cliCtx := ctx.ClientContext()

	caller := cliCtx.GetFromAddress().String()

	calleeAddr, err := sdk.AccAddressFromBech32(callee)
	if err != nil {
		return false, "", err
	}
	accGetter := authtypes.AccountRetriever{}
	if err := accGetter.EnsureExists(cliCtx, calleeAddr); err != nil {
		return false, "", err
	}

	abiSpec, err := queryAbi(cliCtx, callee)
	if err != nil {
		return false, "", err
	}
	abi := []byte(abiSpec) // TODO
	data, err := parseData(function, abi, args, logging.NewNoopLogger())
	if err != nil {
		return false, "", err
	}

	// Decode abiSpec to check if the called function's type is view or pure.
	// If it is, reroute to query.
	var abiEntries []types.ABIEntry
	err = json.Unmarshal(abi, &abiEntries)
	if err != nil {
		return false, "", err
	}
	for _, entry := range abiEntries {
		if entry.Name != function {
			continue
		}
		if entry.Type != "view" && entry.Type != "pure" {
			return false, "", fmt.Errorf("getInsight function should be view or pure function")
		}
		return queryContract(cliCtx, caller, callee, function, abi, data)
	}
	return false, "", fmt.Errorf("function %s was not found in abi", function)
}

// parseData parses Data for contract on certik chain
func parseData(function string, abiSpec []byte, args []string, logger *logging.Logger) ([]byte, error) {
	params := make([]interface{}, 0)

	if string(abiSpec) == compile.NoABI {
		panic("No ABI registered for this contract. Use --raw flag to submit raw bytecode.")
	}

	for _, arg := range args {
		var argi interface{}
		argi = arg
		for _, prefix := range []string{common.Bech32MainPrefix, common.Bech32PrefixConsAddr, common.Bech32PrefixAccAddr} {
			if strings.HasPrefix(arg, prefix) && ((len(arg) - len(prefix)) == 39) {
				data, err := sdk.GetFromBech32(arg, prefix)
				if err != nil {
					return nil, err
				}
				addr, err := crypto.AddressFromBytes(data)
				if err != nil {
					return nil, err
				}
				argi = addr[:]
				break
			}
		}
		params = append(params, argi)
	}

	data, _, err := abi.EncodeFunctionCall(string(abiSpec), function, logger, params...)
	return data, err
}
