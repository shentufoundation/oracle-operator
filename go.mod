module github.com/certikfoundation/oracle-toolset

go 1.16

require (
	github.com/certikfoundation/shentu v1.1.1-0.20210401164204-2763c53c046f
	github.com/cosmos/cosmos-sdk v0.41.0
	github.com/hyperledger/burrow v0.31.0
	github.com/rs/zerolog v1.20.0
	github.com/smartystreets/goconvey v1.6.4
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
	github.com/tendermint/go-amino v0.16.0
	github.com/tendermint/tendermint v0.34.3
	github.com/tendermint/tmlibs v0.9.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
