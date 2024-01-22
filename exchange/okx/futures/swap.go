package futures

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"okx-bot/exchange/model"
	"okx-bot/exchange/okx/common"
	"okx-bot/exchange/options"
)

type Swap struct {
	*common.OKxV5
	currencyPairM map[string]model.CurrencyPair
}

func NewSwap() *Swap {
	var currencyPairM = make(map[string]model.CurrencyPair, 64)
	return &Swap{
		OKxV5:         common.New(),
		currencyPairM: currencyPairM}
}

func (f *Swap) GetExchangeInfo() (map[string]model.CurrencyPair, []byte, error) {
	m, b, er := f.OKxV5.GetExchangeInfo("SWAP")
	f.currencyPairM = m
	return m, b, er
}

func (f *Swap) NewCurrencyPair(baseSym, quoteSym string, opts ...model.OptionParameter) (model.CurrencyPair, error) {
	currencyPair := f.currencyPairM[baseSym+quoteSym]
	log.Info("currencyPair", currencyPair)
	if currencyPair.Symbol == "" {
		return currencyPair, errors.New("not found currency pair")
	}
	return currencyPair, nil
}

func (f *Swap) NewPrvApi(apiOpts ...options.ApiOption) *PrvApi {
	return NewPrvApi(f.OKxV5, apiOpts...)
}
