package oracle

import (
	"github.com/cosmos/cosmos-sdk/client/tx"

	oracletypes "github.com/shentufoundation/shentu/v2/x/oracle/types"

	"github.com/shentufoundation/oracle-operator/types"
)

// Push pushes MsgInquiryEvent to shentu chain.
func Push(ctx types.Context, ctkMsgChan <-chan interface{}, errorChan chan<- error) {
	for {
		select {
		case <-ctx.Context().Done():
			return
		case msg := <-ctkMsgChan:
			switch m := msg.(type) {
			case *oracletypes.MsgTaskResponse:
				go PushMsgTaskResponse(ctx.WithLoggerLabels("type", "MsgTaskResponse"), *m)
			}
		}
	}
}

// PushMsgTaskResponse pushes MsgTaskResponse message to Shentu Chain.
func PushMsgTaskResponse(ctx types.Context, msg oracletypes.MsgTaskResponse) {
	logger := ctx.Logger()
	if err := msg.ValidateBasic(); err != nil {
		ctx.Logger().Error(err.Error())
		return
	}
	err := tx.BroadcastTx(ctx.ClientContext(), ctx.TxBuilder(), &msg)
	if err != nil {
		ctx.Logger().Error(err.Error())
		return
	}

	logger.Debug("Finished pushing task response back to shentu-chain")
}
