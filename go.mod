module github.com/shentufoundation/oracle-operator

go 1.16

require (
	github.com/cosmos/cosmos-sdk v0.45.9
	github.com/rs/zerolog v1.28.0
	github.com/shentufoundation/shentu/v2 v2.6.0
	github.com/smartystreets/goconvey v1.6.4
	github.com/spf13/cobra v1.5.0
	github.com/spf13/viper v1.13.0
	github.com/tendermint/go-amino v0.16.0
	github.com/tendermint/tendermint v0.34.21
	github.com/tendermint/tmlibs v0.9.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
