package common

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"okx-bot/exchange/httpcli"
	"okx-bot/exchange/logger"
	"okx-bot/exchange/model"
	"okx-bot/exchange/options"
	"okx-bot/exchange/util"
	"strings"
	"time"
)

type Prv struct {
	*OKxV5
	apiOpts options.ApiOptions
}

func (prv *Prv) GetAccount(coin string) (map[string]model.Account, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", prv.UriOpts.Endpoint, prv.UriOpts.GetAccountUri)
	params := url.Values{}
	params.Set("ccy", coin)
	data, responseBody, err := prv.DoAuthRequest(http.MethodGet, reqUrl, &params, nil)
	if err != nil {
		return nil, responseBody, err
	}
	acc, err := prv.UnmarshalOpts.GetAccountResponseUnmarshaler(data)
	return acc, responseBody, err
}

func (prv *Prv) CreateOrder(pair model.CurrencyPair, qty, price float64, side model.OrderSide, orderTy model.OrderType, opts ...model.OptionParameter) (*model.Order, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", prv.UriOpts.Endpoint, prv.UriOpts.NewOrderUri)
	params := url.Values{}

	params.Set("instId", pair.Symbol)
	//params.Set("tdMode", "cash")
	//params.Set("posSide", "")
	params.Set("ordType", adaptOrderTypeToSym(orderTy))
	params.Set("px", util.FloatToString(price, pair.PricePrecision))
	params.Set("sz", util.FloatToString(qty, pair.QtyPrecision))

	side2, posSide := adaptOrderSideToSym(side)
	params.Set("side", side2)
	if posSide != "" {
		params.Set("posSide", posSide)
	}

	util.MergeOptionParams(&params, opts...)

	data, responseBody, err := prv.DoAuthRequest(http.MethodPost, reqUrl, &params, nil)
	if err != nil {
		logger.Errorf("[CreateOrder] err=%s, response=%s", err.Error(), string(data))
		return nil, responseBody, err
	}

	ord, err := prv.UnmarshalOpts.CreateOrderResponseUnmarshaler(data)
	if err != nil {
		return nil, responseBody, err
	}

	ord.Pair = pair
	ord.Price = price
	ord.Qty = qty
	ord.Side = side
	ord.OrderTy = orderTy
	ord.Status = model.OrderStatus_Pending

	return ord, responseBody, err
}

func (prv *Prv) GetOrderInfo(pair model.CurrencyPair, id string, opt ...model.OptionParameter) (*model.Order, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", prv.UriOpts.Endpoint, prv.UriOpts.GetOrderUri)
	params := url.Values{}
	params.Set("instId", pair.Symbol)
	params.Set("ordId", id)

	util.MergeOptionParams(&params, opt...)

	data, responseBody, err := prv.DoAuthRequest(http.MethodGet, reqUrl, &params, nil)
	if err != nil {
		return nil, responseBody, err
	}

	ord, err := prv.UnmarshalOpts.GetOrderInfoResponseUnmarshaler(data[1 : len(data)-1])
	if err != nil {
		return nil, responseBody, err
	}

	ord.Pair = pair

	return ord, responseBody, nil
}

func (prv *Prv) GetPendingOrders(pair model.CurrencyPair, opt ...model.OptionParameter) ([]model.Order, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", prv.UriOpts.Endpoint, prv.UriOpts.GetPendingOrdersUri)
	params := url.Values{}
	params.Set("instId", pair.Symbol)

	util.MergeOptionParams(&params, opt...)

	data, responseBody, err := prv.DoAuthRequest(http.MethodGet, reqUrl, &params, nil)
	if err != nil {
		return nil, responseBody, err
	}

	orders, err := prv.UnmarshalOpts.GetPendingOrdersResponseUnmarshaler(data)
	return orders, responseBody, err
}

func (prv *Prv) GetHistoryOrders(pair model.CurrencyPair, opt ...model.OptionParameter) ([]model.Order, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", prv.UriOpts.Endpoint, prv.UriOpts.GetHistoryOrdersUri)
	params := url.Values{}
	params.Set("instId", pair.Symbol)
	params.Set("limit", "50")

	util.MergeOptionParams(&params, opt...)

	data, responseBody, err := prv.DoAuthRequest(http.MethodGet, reqUrl, &params, nil)
	if err != nil {
		return nil, responseBody, err
	}

	orders, err := prv.UnmarshalOpts.GetHistoryOrdersResponseUnmarshaler(data)
	return orders, responseBody, err
}

func (prv *Prv) GetSpotHistoryOrders(opt ...model.OptionParameter) ([]model.Order, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", prv.UriOpts.Endpoint, prv.UriOpts.GetHistoryOrdersUri)
	params := url.Values{}
	params.Set("instType", "SPOT")

	util.MergeOptionParams(&params, opt...)

	data, responseBody, err := prv.DoAuthRequest(http.MethodGet, reqUrl, &params, nil)
	if err != nil {
		return nil, responseBody, err
	}

	orders, err := prv.UnmarshalOpts.GetHistoryOrdersResponseUnmarshaler(data)
	return orders, responseBody, err
}

func (prv *Prv) CancelOrder(pair model.CurrencyPair, id string, opt ...model.OptionParameter) ([]byte, error) {
	reqUrl := fmt.Sprintf("%s%s", prv.UriOpts.Endpoint, prv.UriOpts.CancelOrderUri)
	params := url.Values{}
	params.Set("instId", pair.Symbol)
	params.Set("ordId", id)
	util.MergeOptionParams(&params, opt...)

	data, responseBody, err := prv.DoAuthRequest(http.MethodPost, reqUrl, &params, nil)
	if data != nil && len(data) > 0 {
		return responseBody, prv.UnmarshalOpts.CancelOrderResponseUnmarshaler(data)
	}

	return responseBody, err
}

func (prv *Prv) GetGridAlgoOrderDetails(req model.GridAlgoOrderDetailsRequest, opt ...model.OptionParameter) (model.GridAlgoOrderDetailsResponse, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", prv.UriOpts.Endpoint, prv.UriOpts.GetAlgoOrderDetails)

	params := url.Values{}
	params.Set("algoOrdType", req.AlgoOrdType)
	params.Set("algoId", req.AlgoId)

	util.MergeOptionParams(&params, opt...)

	data, responseBody, err := prv.DoAuthRequest(http.MethodGet, reqUrl, &params, nil)
	if err != nil {
		return model.GridAlgoOrderDetailsResponse{}, responseBody, err
	}

	logger.Info("responseBody", string(responseBody))
	logger.Info("data", string(data))

	details, err := prv.UnmarshalOpts.GetAlgoOrderDetailsResponseUnmarshaler(data)

	return details, responseBody, err
}

func (prv *Prv) PlaceGridAlgoOrder(req model.PlaceGridAlgoOrderRequest, opt ...model.OptionParameter) (model.PlaceGridAlgoOrderResponse, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", prv.UriOpts.Endpoint, prv.UriOpts.PostPlaceGridAlgoOrderUri)

	params := url.Values{}
	params.Set("instId", req.InstId)
	params.Set("algoOrdType", req.AlgoOrdType)
	params.Set("maxPx", req.MaxPx)
	params.Set("minPx", req.MinPx)
	params.Set("gridNum", req.GridNum)
	params.Set("runType", req.RunType)
	params.Set("sz", req.Sz)
	params.Set("direction", req.Direction)
	params.Set("lever", req.Lever)
	params.Set("tag", req.Tag)
	params.Set("basePos", req.BasePos)

	util.MergeOptionParams(&params, opt...)

	data, responseBody, err := prv.DoAuthRequest(http.MethodPost, reqUrl, &params, nil)
	if err != nil {
		logger.Errorf("[PlaceOrder] err=%s, response=%s", err.Error(), string(data))
		return model.PlaceGridAlgoOrderResponse{}, responseBody, err
	}

	logger.Info("responseBody", string(responseBody))
	logger.Info("data", string(data))

	details, err := prv.UnmarshalOpts.PlaceGridAlgoOrderResponseUnmarshaler(data)

	return details, responseBody, err
}

func (prv *Prv) StopGridAlgoOrder(req model.StopGridAlgoOrderRequest, opt ...model.OptionParameter) (model.StopGridAlgoOrderResponse, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", prv.UriOpts.Endpoint, prv.UriOpts.PostStopGridAlgoOrderUri)

	params := url.Values{}
	params.Set("algoId", req.AlgoId)
	params.Set("instId", req.InstId)
	params.Set("algoOrdType", req.AlgoOrdType)
	params.Set("stopType", req.StopType)

	util.MergeOptionParams(&params, opt...)

	data, responseBody, err := prv.DoAuthRequest(http.MethodPost, reqUrl, &params, nil)
	if err != nil {
		return model.StopGridAlgoOrderResponse{}, responseBody, err
	}

	logger.Info("responseBody", string(responseBody))
	logger.Info("data", string(data))

	details, err := prv.UnmarshalOpts.StopGridAlgoOrderResponseUnmarshaler(data)

	return details, responseBody, err
}

func (prv *Prv) PlaceOrder(req model.PlaceOrderRequest, opt ...model.OptionParameter) (*model.Order, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", prv.UriOpts.Endpoint, prv.UriOpts.NewOrderUri)

	params := url.Values{}
	params.Set("instId", req.InstId)
	params.Set("tdMode", req.TdMode)
	params.Set("side", req.Side)
	params.Set("posSide", req.PosSide)
	params.Set("ordType", req.OrdType)
	params.Set("sz", req.Sz)
	params.Set("px", req.Px)

	util.MergeOptionParams(&params, opt...)

	data, responseBody, err := prv.DoAuthRequest(http.MethodPost, reqUrl, &params, nil)
	if err != nil {
		return nil, responseBody, err
	}

	logger.Info("responseBody", string(responseBody))
	logger.Info("data", string(data))

	details, err := prv.UnmarshalOpts.CreateOrderResponseUnmarshaler(data)
	logger.Info("details", details)

	return details, responseBody, err
}

func (prv *Prv) AmendOrder(req model.AmendOrderRequest, opt ...model.OptionParameter) (model.AmendOrderResponse, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", prv.UriOpts.Endpoint, prv.UriOpts.AmendOrderUri)

	params := url.Values{}
	params.Set("instId", req.InstId)
	params.Set("ordId", req.OrdId)
	params.Set("clOrdId", req.ClOrdId)
	params.Set("reqId", req.ReqId)
	params.Set("newSz", req.NewSz)
	params.Set("newPx", req.NewPx)
	params.Set("newPxUsd", req.NewPxUsd)
	params.Set("NewPxVol", req.NewPxVol)

	util.MergeOptionParams(&params, opt...)

	data, responseBody, err := prv.DoAuthRequest(http.MethodPost, reqUrl, &params, nil)

	if err != nil {
		fmt.Println(err)
		return model.AmendOrderResponse{}, responseBody, err
	}

	logger.Info("responseBody", string(responseBody))
	logger.Info("data", string(data))

	details, err := prv.UnmarshalOpts.AmendOrderResponseUmarshaler(data)
	logger.Info("details", details)

	return details, responseBody, err
}

func (prv *Prv) DoSignParam(httpMethod, apiUri, apiSecret, reqBody string) (signStr, timestamp string) {
	timestamp = time.Now().UTC().Format("2006-01-02T15:04:05.000Z") //iso time style
	payload := fmt.Sprintf("%s%s%s%s", timestamp, strings.ToUpper(httpMethod), apiUri, reqBody)
	signStr, _ = util.HmacSHA256Base64Sign(apiSecret, payload)
	return
}

func (prv *Prv) DoAuthRequest(httpMethod, reqUrl string, params *url.Values, headers map[string]string) ([]byte, []byte, error) {
	var (
		reqBodyStr string
		reqUri     string
	)

	if http.MethodGet == httpMethod {
		reqUrl += "?" + params.Encode()
	}

	if http.MethodPost == httpMethod {
		params.Set("tag", "86d4a3bf87bcBCDE")
		reqBody, _ := util.ValuesToJson(*params)
		reqBodyStr = string(reqBody)
	}

	_url, _ := url.Parse(reqUrl)
	reqUri = _url.RequestURI()
	signStr, timestamp := prv.DoSignParam(httpMethod, reqUri, prv.apiOpts.Secret, reqBodyStr)
	logger.Debugf("[DoAuthRequest] sign base64: %s, timestamp: %s", signStr, timestamp)

	headers = map[string]string{
		"Content-Type": "application/json; charset=UTF-8",
		//"Accept":               "application/json",
		"x-simulated-trading":  "1",
		"OK-ACCESS-KEY":        prv.apiOpts.Key,
		"OK-ACCESS-PASSPHRASE": prv.apiOpts.Passphrase,
		"OK-ACCESS-SIGN":       signStr,
		"OK-ACCESS-TIMESTAMP":  timestamp}

	respBody, err := httpcli.Cli.DoRequest(httpMethod, reqUrl, reqBodyStr, headers)
	if err != nil {
		return nil, respBody, err
	}
	logger.Debugf("[DoAuthRequest] response body: %s", string(respBody))

	var baseResp BaseResp
	err = prv.OKxV5.UnmarshalOpts.ResponseUnmarshaler(respBody, &baseResp)
	if err != nil {
		return nil, respBody, err
	}

	if baseResp.Code != 0 {
		return nil, respBody, errors.New(baseResp.Msg)
	}

	return baseResp.Data, respBody, nil
}

func NewPrvApi(opts ...options.ApiOption) *Prv {
	var api = new(Prv)
	for _, opt := range opts {
		opt(&api.apiOpts)
	}
	return api
}
