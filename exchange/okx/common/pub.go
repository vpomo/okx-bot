package common

import (
	"fmt"
	"net/http"
	"net/url"
	"okx-bot/exchange/httpcli"
	"okx-bot/exchange/logger"
	"okx-bot/exchange/model"
	"okx-bot/exchange/util"
	"strconv"
)

func (okx *OKxV5) GetName() string {
	return "okx.com"
}

func (okx *OKxV5) GetDepth(pair model.CurrencyPair, size int, opt ...model.OptionParameter) (*model.Depth, []byte, error) {
	params := url.Values{}
	params.Set("instId", pair.Symbol)
	params.Set("sz", fmt.Sprint(size))
	util.MergeOptionParams(&params, opt...)

	data, responseBody, err := okx.DoNoAuthRequest("GET", okx.UriOpts.Endpoint+okx.UriOpts.DepthUri, &params)
	if err != nil {
		return nil, responseBody, err
	}

	dep, err := okx.UnmarshalOpts.DepthUnmarshaler(data)
	if err != nil {
		return nil, data, err
	}

	dep.Pair = pair

	return dep, responseBody, err
}

func (okx *OKxV5) GetTicker(pair model.CurrencyPair, opt ...model.OptionParameter) (*model.Ticker, []byte, error) {
	params := url.Values{}
	params.Set("instId", pair.Symbol)

	data, responseBody, err := okx.DoNoAuthRequest("GET", okx.UriOpts.Endpoint+okx.UriOpts.TickerUri, &params)
	if err != nil {
		return nil, data, err
	}

	tk, err := okx.UnmarshalOpts.TickerUnmarshaler(data)
	if err != nil {
		return nil, nil, err
	}

	tk.Pair = pair

	return tk, responseBody, err
}

func (okx *OKxV5) GetKline(pair model.CurrencyPair, period model.KlinePeriod, opt ...model.OptionParameter) ([]model.Kline, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", okx.UriOpts.Endpoint, okx.UriOpts.KlineUri)
	param := url.Values{}
	param.Set("instId", pair.Symbol)
	param.Set("bar", AdaptKlinePeriodToSymbol(period))
	param.Set("limit", "100")
	util.MergeOptionParams(&param, opt...)

	data, responseBody, err := okx.DoNoAuthRequest(http.MethodGet, reqUrl, &param)
	if err != nil {
		return nil, nil, err
	}
	klines, err := okx.UnmarshalOpts.KlineUnmarshaler(data)
	return klines, responseBody, err
}

func (okx *OKxV5) GetExchangeInfo(instType string, opt ...model.OptionParameter) (map[string]model.CurrencyPair, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", okx.UriOpts.Endpoint, okx.UriOpts.GetExchangeInfoUri)
	param := url.Values{}
	param.Set("instType", instType)
	util.MergeOptionParams(&param, opt...)

	data, responseBody, err := okx.DoNoAuthRequest(http.MethodGet, reqUrl, &param)
	if err != nil {
		return nil, responseBody, err
	}

	currencyPairMap, err := okx.UnmarshalOpts.GetExchangeInfoResponseUnmarshaler(data)

	return currencyPairMap, responseBody, err
}

func (okx *OKxV5) GetCompMinInvest(minInvestReq model.ComputeMinInvestmentRequest, opt ...model.OptionParameter) (model.ComputeMinInvestmentResponse, []byte, error) {
	reqUrl := fmt.Sprintf("%s%s", okx.UriOpts.Endpoint, okx.UriOpts.PostComputeMinInvestment)
	params := url.Values{}
	params.Set("instId", minInvestReq.InstId)
	params.Set("algoOrdType", minInvestReq.AlgoOrdType)
	params.Set("maxPx", minInvestReq.MaxPx)
	params.Set("minPx", minInvestReq.MinPx)
	params.Set("gridNum", minInvestReq.GridNum)
	params.Set("runType", minInvestReq.RunType)
	params.Set("direction", minInvestReq.Direction)
	params.Set("lever", minInvestReq.Lever)
	params.Set("basePos", strconv.FormatBool(minInvestReq.BasePos))
	//params.Set("investmentData", minInvestReq.InvestmentData)

	util.MergeOptionParams(&params, opt...)

	data, responseBody, err := okx.DoNoAuthRequest(http.MethodPost, reqUrl, &params)
	if err != nil {
		return model.ComputeMinInvestmentResponse{}, responseBody, err
	}

	logger.Info("responseBody", string(responseBody))
	logger.Info("data", string(data))

	minInvestment, err := okx.UnmarshalOpts.GetCompMinInvestResponseUnmarshaler(data)

	return minInvestment, responseBody, err
}

func (okx *OKxV5) DoNoAuthRequest(httpMethod, reqUrl string, params *url.Values) ([]byte, []byte, error) {
	reqBody := ""

	if http.MethodGet == httpMethod {
		reqUrl += "?" + params.Encode()
	}

	if http.MethodPost == httpMethod {
		params.Set("tag", "86d4a3bf87bcBCDE")
		reqBodyByte, _ := util.ValuesToJson(*params)
		reqBody = string(reqBodyByte)
	}

	responseBody, err := httpcli.Cli.DoRequest(httpMethod, reqUrl, reqBody, nil)
	if err != nil {
		return nil, responseBody, err
	}

	var baseResp BaseResp
	err = okx.UnmarshalOpts.ResponseUnmarshaler(responseBody, &baseResp)
	if err != nil {
		return responseBody, responseBody, err
	}

	if baseResp.Code == 0 {
		logger.Debugf("[DoNoAuthRequest] response=%s", string(responseBody))
		return baseResp.Data, responseBody, nil
	}

	logger.Debugf("[DoNoAuthRequest] error=%s", baseResp.Msg)
	return nil, responseBody, err
}
