package types

import (
	"context"

	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
)

// Context define the type of data for the operator node.
type Context struct {
	CLIContext *client.Context
	ctx        context.Context
	txFactory  *tx.Factory
	config     *Config
	logger     log.Logger
}

// Context returns the internal context object.
func (c Context) Context() context.Context { return c.ctx }

// ClientContext returns a copy of the certik chain client context object value.
func (c Context) ClientContext() client.Context { return *c.CLIContext }

// Codec returns a reference to the client codec.
func (c Context) Codec() *amino.Codec { return c.CLIContext.LegacyAmino.Amino }

// TxBuilder returns a copy of the certik chain transaction builder object.
func (c Context) TxBuilder() tx.Factory { return *c.txFactory }

// Config returns a copy of the oracle operator node global configuration.
func (c Context) Config() Config { return *c.config }

// Logger returns the logger for oracle node internal use.
func (c Context) Logger() log.Logger { return c.logger }

// NewContextWithDefaultConfigAndLogger returns a new context with global configuration set from a config file.
func NewContextWithDefaultConfigAndLogger() (Context, error) {
	if err := initConfig(); err != nil {
		return Context{}, err
	}
	if err := initLogger(); err != nil {
		return Context{}, err
	}
	return NewContext(&config, logger), nil
}

// NewContext creates a new context.
func NewContext(config *Config, logger log.Logger) Context {
	return Context{
		CLIContext: &client.Context{},
		ctx:        context.Background(),
		txFactory:  &tx.Factory{},
		config:     config,
		logger:     logger,
	}
}

// WithContext returns a copy of the context with an updated internal context.
func (c Context) WithContext(ctx context.Context) Context {
	c.ctx = ctx
	return c
}

// WithClientContext returns a copy of the context with an updated CosmoSDK client context.
func (c Context) WithClientContext(ctx *client.Context) Context {
	c.CLIContext = ctx
	return c
}

// WithTxFactory returns a copy of the context with an updated tx builder.
func (c Context) WithTxFactory(txFactory *tx.Factory) Context {
	c.txFactory = txFactory
	return c
}

// WithConfig returns a copy of the context with an updated configuration setting.
func (c Context) WithConfig(config *Config) Context {
	c.config = config
	return c
}

// WithLogger returns a copy of the context with an updated logger.
func (c Context) WithLogger(logger log.Logger) Context {
	c.logger = logger
	return c
}

// WithLoggerLabels returns a copy of the context with updated logger labels.
func (c Context) WithLoggerLabels(keyvals ...interface{}) Context {
	c.logger = c.logger.With(keyvals...)
	return c
}

// WithValue returns a copy of the context with an extra key-value information.
func (c Context) WithValue(key, value interface{}) Context {
	c.ctx = context.WithValue(c.ctx, key, value)
	return c
}
