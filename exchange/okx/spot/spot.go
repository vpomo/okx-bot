package spot

import (
	"okx-bot/exchange/model"
	"okx-bot/exchange/okx/common"
	"okx-bot/exchange/options"
)

type Spot struct {
	*common.OKxV5
	currencyPairM map[string]model.CurrencyPair
}

func New() *Spot {
	v5 := common.New()
	currencyPairCacheMap := make(map[string]model.CurrencyPair, 64)
	return &Spot{v5, currencyPairCacheMap}
}

func (s *Spot) NewPrvApi(apiOps ...options.ApiOption) *PrvApi {
	prv := new(PrvApi)
	prv.Prv = s.OKxV5.NewPrvApi(apiOps...)
	prv.Prv.OKxV5 = s.OKxV5
	return prv
}
