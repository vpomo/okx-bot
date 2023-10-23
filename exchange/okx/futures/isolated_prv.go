package futures

import (
	"errors"
	"okx-bot/exchange/model"
)

type IsolatedPrvApi struct {
	*PrvApi
}

func (f *IsolatedPrvApi) CreateOrder(pair model.CurrencyPair, qty, price float64, side model.OrderSide, orderTy model.OrderType, opts ...model.OptionParameter) (*model.Order, []byte, error) {
	if side != model.Futures_OpenBuy &&
		side != model.Futures_OpenSell &&
		side != model.Futures_CloseBuy &&
		side != model.Futures_CloseSell {
		return nil, nil, errors.New("futures side only is Futures_OpenBuy or Futures_OpenSell or Futures_CloseBuy or Futures_CloseSell")
	}

	opts = append(opts,
		model.OptionParameter{
			Key:   "tdMode",
			Value: "isolated",
		})

	return f.Prv.CreateOrder(pair, qty, price, side, orderTy, opts...)
}
