package oracle

import (
	"fmt"
	"strings"
	"sync"

	abciTypes "github.com/tendermint/tendermint/abci/types"
	ctkClient "github.com/tendermint/tendermint/rpc/client/http"
	tendermintTypes "github.com/tendermint/tendermint/types"

	oracletypes "github.com/certikfoundation/shentu/v2/x/oracle/types"

	"github.com/certikfoundation/oracle-operator/types"
)

// Listen listens for events from CertiK chain.
func Listen(ctx types.Context, ctkMsgChan chan<- interface{}, fatalError chan<- error) {
	// load configuration and logger
	logger := ctx.Logger()
	node := ctx.ClientContext().NodeURI
	logger.Info("start to listen to certik-chain", "node", node)

	// initialize client
	client, err := ctkClient.New(ctx.ClientContext().NodeURI, "/websocket")
	if err != nil {
		logger.Error("ctkClient dialing", "error", err.Error())
		fatalError <- err
		return
	}

	// start the listener
	err = client.Start()
	if err != nil {
		logger.Error("ctkClient subscribing", "error", err.Error())
		fatalError <- err
		return
	}
	defer client.Stop()

	// subscribe the TXs according to the query
	txChan, err := client.Subscribe(ctx.Context(), "", "tm.event='Tx'", 1000) // TODO
	if err != nil {
		logger.Error("ctkClient subscribing", "error", err.Error())
		fatalError <- err
		return
	}

	for {
		select {
		case <-ctx.Context().Done():
			logger.Info("stop listening...")
			return
		case tx := <-txChan:
			// get tendermint transaction data in struct of ResponseDeliverTx
			txData, ok := tx.Data.(tendermintTypes.EventDataTx)
			if !ok {
				logger.Error("received non-event tx", "tx", tx.Data)
			}
			for _, event := range txData.Result.Events {
				switch event.Type {
				case "create_task":
					logger.Info("Received event", "type", "create_task")
					go handleMsgCreateTask(ctx.WithLoggerLabels("type", "create_task"), event, ctkMsgChan)
				case "compute_score":
					logger.Info("Received event", "type", "compute_score")
					go handleMsgComputeScoreTask(ctx.WithLoggerLabels("type", "compute_score"), event, ctkMsgChan)
				}
			}
		}
	}
}

// handleMsgCreateTask parses MsgCreateTask TX data and passes organized message to endpoint querier.
func handleMsgCreateTask(ctx types.Context, event abciTypes.Event, ctkMsgChan chan<- interface{}) {
	logger := ctx.Logger()
	// parse event
	msgCreateTask, err := parseMsgCreateTask(event)
	if err != nil {
		logger.Error("parsing event", "type", "create_task", "error", err.Error(), event)
		return
	}
	logger.Debug("parsed event", "type", "create_task", "msg", msgCreateTask)
	// get payload
	payload, err := getPrimitivePayload(msgCreateTask)
	if err != nil {
		logger.Error("getting task payload", "type", "create_task", "error", err.Error(), "msg", msgCreateTask)
		return
	}
	logger.Debug("task payload", "type", "create_task", "payload", payload)

	// get aggregation strategy
	var strategy types.Strategy
	var ok bool
	strategy, ok = ctx.Config().Strategy[payload.Client]
	if !ok {
		// use default strategy if specific client is not found in config
		strategy, ok = ctx.Config().Strategy["default"]
		logger.Info("Chain configuration has not been set. Use default chain configuration instead.")
		if !ok {
			logger.Error("target client chain strategy not specified", "client", payload.Client, "payload", payload)
			return
		}
	}
	aggregator, err := NewAggregation(strategy)
	if err != nil {
		logger.Error("aggregation strategy not defined", "type", strategy.Type)
		return
	}
	// get primitive socres
	var wg sync.WaitGroup
	primitiveScores := make(chan types.PrimitiveScore, len(strategy.Primitives))
	wg.Add(len(strategy.Primitives))
	for _, primitive := range strategy.Primitives {
		go queryPrimitive(
			ctx.WithLoggerLabels("primitive", primitive),
			primitive,
			payload,
			primitiveScores,
			&wg,
		)
	}
	wg.Wait()
	close(primitiveScores)
	// aggregate primitive scores
	score, err := aggregator.Aggregate(primitiveScores)
	if err != nil {
		logger.Error("aggregation failed", "type", strategy.Type, "error", err.Error(), "payload", payload)
		return
	}
	logger.Info(
		"aggregation result",
		"type", strategy.Type,
		"score", score,
		"payload", payload,
	)
	// push back
	ctkMsgChan <- oracletypes.NewMsgTaskResponse(
		msgCreateTask.Contract,
		msgCreateTask.Function,
		int64(score),
		ctx.ClientContext().GetFromAddress(),
	)
}

// parseMsgCreateTask parses TX data of creating tasks.
func parseMsgCreateTask(event abciTypes.Event) (oracletypes.MsgCreateTask, error) {
	var contract, function string
	for _, v := range event.GetAttributes() {
		switch string(v.GetKey()) {
		case "contract":
			contract = string(v.GetValue())
			if contract == "" {
				return oracletypes.MsgCreateTask{}, fmt.Errorf("missing contract in event content")
			}
		case "function":
			function = string(v.GetValue())
			if function == "" {
				return oracletypes.MsgCreateTask{}, fmt.Errorf("missing function in event content")
			}
		}
	}
	msgCreateTask := oracletypes.MsgCreateTask{
		Contract: contract,
		Function: function,
	}
	return msgCreateTask, nil
}

// parseMsgCreateTaskContract parses MsgCreateTaskContract contract field.
func parseMsgCreateTaskContract(contract string) (types.Client, string, error) {
	if !strings.Contains(contract, ":") {
		return types.DefaultClient, contract, nil
	}
	seg := strings.Split(contract, ":")
	if len(seg) <= 1 {
		return "", "", fmt.Errorf(contract)
	}
	return types.Client(strings.Join(seg[:len(seg)-1], ":")), seg[len(seg)-1], nil
}

// handleMsgEndpointQuery parses MsgEndpointQuery TX data and passes organized message to endpoint querier.
func handleMsgComputeScoreTask(ctx types.Context, event abciTypes.Event, ctkMsgChan chan<- interface{}) {

}
