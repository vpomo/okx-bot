package spot

import (
	"errors"
	"okx-bot/exchange/model"
	"okx-bot/exchange/okx/common"
)

type PrvApi struct {
	*common.Prv
}

func (api *PrvApi) CreateOrder(pair model.CurrencyPair, qty, price float64, side model.OrderSide, orderTy model.OrderType, opts ...model.OptionParameter) (*model.Order, []byte, error) {
	//check params
	if model.Spot_Buy != side && side != model.Spot_Sell {
		return nil, nil, errors.New("spot order side is error")
	}

	opts = append(opts,
		model.OptionParameter{
			Key:   "tdMode",
			Value: "cash",
		})

	return api.Prv.CreateOrder(pair, qty, price, side, orderTy, opts...)
}

func (api *PrvApi) GetHistoryOrders(pair model.CurrencyPair, opt ...model.OptionParameter) ([]model.Order, []byte, error) {
	opt = append(opt, model.OptionParameter{
		Key:   "instType",
		Value: "SPOT",
	})
	return api.Prv.GetHistoryOrders(pair, opt...)
}
