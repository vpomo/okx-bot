package okx

import (
	"okx-bot/exchange/httpcli"
	"okx-bot/exchange/logger"
	"okx-bot/exchange/okx/futures"
	"okx-bot/exchange/okx/grid"
	"okx-bot/exchange/okx/spot"
	"reflect"
)

type OKx struct {
	Spot    *spot.Spot
	Futures *futures.Futures
	Swap    *futures.Swap
	Grid    *grid.Grid
}

func New() *OKx {
	return &OKx{
		Spot:    spot.New(),
		Futures: futures.New(),
		Swap:    futures.NewSwap(),
		Grid:    grid.New(),
	}
}

func SetDefaultHttpCli(cli httpcli.IHttpClient) {
	logger.Infof("use new http client implement: %s", reflect.TypeOf(cli).Elem().String())
	httpcli.Cli = cli
}

var (
	DefaultHttpCli = httpcli.Cli
)
