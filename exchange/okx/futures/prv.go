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

func (prv *PrvApi) GetPositionsHistory(request model.FuturesPositionHistoryRequest, opts ...model.OptionParameter) ([]model.FuturesPositionHistory, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", prv.OKxV5.UriOpts.Endpoint, prv.OKxV5.UriOpts.GetPositionsHistoryUri)
	params := url.Values{}
	params.Set("instType", request.InstType)
	params.Set("instId", request.InstId)
	params.Set("mgnMode", request.MgnMode)
	params.Set("type", request.Type)
	params.Set("posId", request.PosId)
	params.Set("after", request.After)
	params.Set("before", request.Before)
	params.Set("limit", request.Limit)
	util.MergeOptionParams(&params, opts...)
	data, responseBody, err := prv.DoAuthRequest(http.MethodGet, reqUrl, &params, nil)
	if err != nil {
		return nil, responseBody, err
	}
	positionsHistory, err := prv.OKxV5.UnmarshalOpts.GetPositionsHistoryResponseUnmarshaler(data)
	return positionsHistory, responseBody, err
}

func (prv *PrvApi) GetHistoryOrders(pair model.CurrencyPair, opt ...model.OptionParameter) ([]model.Order, []byte, error) {
	opt = append(opt, model.OptionParameter{
		Key:   "instType",
		Value: "SWAP",
	})
	return prv.Prv.GetHistoryOrders(pair, opt...)
}
