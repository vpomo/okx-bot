package grid

import (
	"errors"
	"okx-bot/exchange/model"
	"okx-bot/exchange/okx/common"
	"okx-bot/exchange/options"
)

type Grid struct {
	*common.OKxV5
	currencyPairM map[string]model.CurrencyPair
}

func New() *Grid {
	currencyPairM := make(map[string]model.CurrencyPair, 64)
	return &Grid{OKxV5: common.New(), currencyPairM: currencyPairM}
}

func (g *Grid) NewPrvApi(apiOps ...options.ApiOption) *PrvApi {
	prv := new(PrvApi)
	prv.Prv = g.OKxV5.NewPrvApi(apiOps...)
	prv.Prv.OKxV5 = g.OKxV5
	return prv
}

func (g *Grid) NewCurrencyPair(baseSym, quoteSym string, opts ...model.OptionParameter) (model.CurrencyPair, error) {
	if len(opts) >= 1 && opts[0].Key == "contractAlias" {
		contractAlias := opts[0].Value
		currencyPair := g.currencyPairM[baseSym+quoteSym+contractAlias]
		if currencyPair.Symbol != "" {
			return currencyPair, nil
		}
		return currencyPair, errors.New("not found currency pair")
	}
	return model.CurrencyPair{}, errors.New("please input contract alias option parameter")
}
