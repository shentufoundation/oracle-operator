package oracle

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"

	cvmtypes "github.com/certikfoundation/shentu/x/cvm/types"
)

// queryAbi queries ABI from certik chain
func queryAbi(cliCtx client.Context, addr string) (string, error) {
	queryClient := cvmtypes.NewQueryClient(cliCtx)

	res, err := queryClient.Abi(context.Background(), &cvmtypes.QueryAbiRequest{Address: addr})
	if err != nil {
		return "", err
	}

	return res.Abi, nil
}

// queryContract queries contract on certik-chain.
func queryContract(cliCtx client.Context, caller, callee, fname string, abiSpec, data []byte) (bool, string, error) {
	queryClient := cvmtypes.NewQueryClient(cliCtx)

	res, err := queryClient.View(
		context.Background(),
		&cvmtypes.QueryViewRequest{
			Caller: caller,
			Callee: callee,
			AbiSpec: abiSpec,
			FunctionName: fname,
			Data: data,
		})
	if err != nil {
		return false, "", fmt.Errorf("querying security primitive contract: %v", err)
	}

	ret := res.ReturnVars
	retBool, err := strconv.ParseBool(ret[0].Value)
	if err != nil {
		return false, "", fmt.Errorf("decoding function return: %v", err)
	}

	return retBool, ret[1].Value, nil
}
