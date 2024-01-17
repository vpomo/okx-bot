package futures

import (
	"fmt"
	"net/http"
	"net/url"
	"okx-bot/exchange/model"
	"okx-bot/exchange/okx/common"
	"okx-bot/exchange/options"
	"okx-bot/exchange/util"
)

type PrvApi struct {
	*common.Prv
	Isolated *IsolatedPrvApi
	Cross    *CrossPrvApi
}

func NewPrvApi(v5 *common.OKxV5, apiOpts ...options.ApiOption) *PrvApi {
	prvApi := new(PrvApi)
	prvApi.Prv = v5.NewPrvApi(apiOpts...)

	prvApi.Isolated = new(IsolatedPrvApi)
	prvApi.Isolated.PrvApi = prvApi

	prvApi.Cross = new(CrossPrvApi)
	prvApi.Cross.PrvApi = prvApi

	return prvApi
}

func (prv *PrvApi) GetFuturesAccount(coin string) (map[string]model.FuturesAccount, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", prv.OKxV5.UriOpts.Endpoint, prv.OKxV5.UriOpts.GetAccountUri)
	params := url.Values{}
	params.Set("ccy", coin)
	data, responseBody, err := prv.DoAuthRequest(http.MethodGet, reqUrl, &params, nil)
	if err != nil {
		return nil, responseBody, err
	}
	acc, err := prv.OKxV5.UnmarshalOpts.GetFuturesAccountResponseUnmarshaler(data)
	return acc, responseBody, err
}

//func (prv *PrvApi) PlaceOrder(pair model.CurrencyPair, qty, price float64, side model.OrderSide, orderTy model.OrderType, opts ...model.OptionParameter) (*model.Order, []byte, error) {
//	reqUrl := fmt.Sprintf("%s%s", prv.UriOpts.Endpoint, prv.UriOpts.NewOrderUri)
//	params := url.Values{}
//
//	params.Set("instId", pair.Symbol)
//	params.Set("tdMode", "isolated")
//	//params.Set("posSide", "")
//	params.Set("ordType", adaptOrderTypeToSym(orderTy))
//	params.Set("px", util.FloatToString(price, pair.PricePrecision))
//	params.Set("sz", util.FloatToString(qty, pair.QtyPrecision))
//
//	side2, posSide := adaptOrderSideToSym(side)
//	params.Set("side", side2)
//	if posSide != "" {
//		params.Set("posSide", posSide)
//	}
//
//	util.MergeOptionParams(&params, opts...)
//
//	data, responseBody, err := prv.DoAuthRequest(http.MethodPost, reqUrl, &params, nil)
//	if err != nil {
//		logger.Errorf("[CreateOrder] err=%s, response=%s", err.Error(), string(data))
//		return nil, responseBody, err
//	}
//
//	ord, err := prv.UnmarshalOpts.CreateOrderResponseUnmarshaler(data)
//	if err != nil {
//		return nil, responseBody, err
//	}
//
//	ord.Pair = pair
//	ord.Price = price
//	ord.Qty = qty
//	ord.Side = side
//	ord.OrderTy = orderTy
//	ord.Status = model.OrderStatus_Pending
//
//	return ord, responseBody, err
//}

func (prv *PrvApi) GetPositions(pair model.CurrencyPair, opts ...model.OptionParameter) ([]model.FuturesPosition, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", prv.OKxV5.UriOpts.Endpoint, prv.OKxV5.UriOpts.GetPositionsUri)
	params := url.Values{}
	params.Set("instId", pair.Symbol)
	util.MergeOptionParams(&params, opts...)
	data, responseBody, err := prv.DoAuthRequest(http.MethodGet, reqUrl, &params, nil)
	if err != nil {
		return nil, responseBody, err
	}
	positions, err := prv.OKxV5.UnmarshalOpts.GetPositionsResponseUnmarshaler(data)
	return positions, responseBody, err
}

func (prv *PrvApi) GetHistoryOrders(pair model.CurrencyPair, opt ...model.OptionParameter) ([]model.Order, []byte, error) {
	opt = append(opt, model.OptionParameter{
		Key:   "instType",
		Value: "SWAP",
	})
	return prv.Prv.GetHistoryOrders(pair, opt...)
}
