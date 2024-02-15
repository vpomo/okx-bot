package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"okx-bot/exchange/logger"
	"okx-bot/exchange/model"
	"time"

	"github.com/buger/jsonparser"
	"github.com/spf13/cast"
)

type RespUnmarshaler struct {
}

func (un *RespUnmarshaler) UnmarshalDepth(data []byte) (*model.Depth, error) {
	var (
		dep model.Depth
		err error
	)

	err = jsonparser.ObjectEach(data[1:len(data)-1],
		func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			switch string(key) {
			case "ts":
				dep.UTime = time.UnixMilli(cast.ToInt64(string(value)))
			case "asks":
				items, _ := un.unmarshalDepthItem(value)
				dep.Asks = items
			case "bids":
				items, _ := un.unmarshalDepthItem(value)
				dep.Bids = items
			}
			return nil
		})

	return &dep, err
}

func (un *RespUnmarshaler) unmarshalDepthItem(data []byte) (model.DepthItems, error) {
	var items model.DepthItems
	_, err := jsonparser.ArrayEach(data, func(asksItemData []byte, dataType jsonparser.ValueType, offset int, err error) {
		item := model.DepthItem{}
		i := 0
		_, err = jsonparser.ArrayEach(asksItemData, func(itemVal []byte, dataType jsonparser.ValueType, offset int, err error) {
			valStr := string(itemVal)
			switch i {
			case 0:
				item.Price = cast.ToFloat64(valStr)
			case 1:
				item.Amount = cast.ToFloat64(valStr)
			}
			i += 1
		})
		items = append(items, item)
	})
	return items, err
}

func (un *RespUnmarshaler) UnmarshalTicker(data []byte) (*model.Ticker, error) {
	var tk = &model.Ticker{}

	var open float64
	_, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		err = jsonparser.ObjectEach(value, func(key []byte, val []byte, dataType jsonparser.ValueType, offset int) error {
			valStr := string(val)
			switch string(key) {
			case "last":
				tk.Last = cast.ToFloat64(valStr)
			case "askPx":
				tk.Sell = cast.ToFloat64(valStr)
			case "bidPx":
				tk.Buy = cast.ToFloat64(valStr)
			case "vol24h":
				tk.Vol = cast.ToFloat64(valStr)
			case "high24h":
				tk.High = cast.ToFloat64(valStr)
			case "low24h":
				tk.Low = cast.ToFloat64(valStr)
			case "ts":
				tk.Timestamp = cast.ToInt64(valStr)
			case "open24h":
				open = cast.ToFloat64(valStr)
			}
			return nil
		})
	})

	if err != nil {
		logger.Errorf("[UnmarshalTicker] %s", err.Error())
		return nil, err
	}

	tk.Percent = (tk.Last - open) / open * 100

	return tk, nil
}

func (un *RespUnmarshaler) UnmarshalGetKlineResponse(data []byte) ([]model.Kline, error) {
	var klines []model.Kline
	_, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var (
			k model.Kline
			i int
		)
		_, err = jsonparser.ArrayEach(value, func(val []byte, dataType jsonparser.ValueType, offset int, err error) {
			valStr := string(val)
			switch i {
			case 0:
				k.Timestamp = cast.ToInt64(valStr)
			case 1:
				k.Open = cast.ToFloat64(valStr)
			case 2:
				k.High = cast.ToFloat64(valStr)
			case 3:
				k.Low = cast.ToFloat64(valStr)
			case 4:
				k.Close = cast.ToFloat64(valStr)
			case 5:
				k.Vol = cast.ToFloat64(valStr)
			}
			i += 1
		})
		klines = append(klines, k)
	})

	return klines, err
}

func (un *RespUnmarshaler) UnmarshalCreateOrderResponse(data []byte) (*model.Order, error) {
	var ord = new(model.Order)
	err := jsonparser.ObjectEach(data[1:len(data)-1], func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		valStr := string(value)
		switch string(key) {
		case "ordId":
			ord.Id = valStr
		case "clOrdId":
			ord.CId = valStr
		}
		return nil
	})
	return ord, err
}

func (un *RespUnmarshaler) UnmarshalGetPendingOrdersResponse(data []byte) ([]model.Order, error) {
	var (
		orders []model.Order
		err    error
	)

	_, err = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		ord, err := un.UnmarshalGetOrderInfoResponse(value)
		if err != nil {
			return
		}
		orders = append(orders, *ord)
	})

	return orders, err
}

func (un *RespUnmarshaler) UnmarshalGetHistoryOrdersResponse(data []byte) ([]model.Order, error) {
	return un.UnmarshalGetPendingOrdersResponse(data)
}

func (un *RespUnmarshaler) UnmarshalGetOrderInfoResponse(data []byte) (ord *model.Order, err error) {
	var side, posSide string
	var utime int64
	ord = new(model.Order)

	err = jsonparser.ObjectEach(data, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		valStr := string(value)
		switch string(key) {
		case "ordId":
			ord.Id = valStr
		case "px":
			ord.Price = cast.ToFloat64(valStr)
		case "sz":
			ord.Qty = cast.ToFloat64(valStr)
		case "cTime":
			ord.CreatedAt = cast.ToInt64(valStr)
		case "avgPx":
			ord.PriceAvg = cast.ToFloat64(valStr)
		case "accFillSz":
			ord.ExecutedQty = cast.ToFloat64(valStr)
		case "fee":
			ord.Fee = cast.ToFloat64(valStr)
		case "feeCcy":
			ord.FeeCcy = valStr
		case "clOrdId":
			ord.CId = valStr
		case "side":
			side = valStr
		case "posSide":
			posSide = valStr
		case "ordType":
			ord.OrderTy = adaptSymToOrderTy(valStr)
		case "state":
			ord.Status = adaptSymToOrderStatus(valStr)
		case "uTime":
			utime = cast.ToInt64(valStr)
		}
		return nil
	})

	ord.Side = adaptSymToOrderSide(side, posSide)
	if ord.Status == model.OrderStatus_Canceled {
		ord.CanceledAt = utime
		if ord.ExecutedQty > 0 {
			ord.FinishedAt = utime
		}
	}

	if ord.Status == model.OrderStatus_Finished {
		ord.FinishedAt = utime
	}

	return
}

func (un *RespUnmarshaler) UnmarshalGetAccountResponse(data []byte) (map[string]model.Account, error) {
	var accMap = make(map[string]model.Account, 2)

	_, err := jsonparser.ArrayEach(data[1:len(data)-1], func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var acc model.Account
		err = jsonparser.ObjectEach(value, func(key []byte, accData []byte, dataType jsonparser.ValueType, offset int) error {
			valStr := string(accData)
			switch string(key) {
			case "ccy":
				acc.Coin = valStr
			case "availEq":
				acc.AvailableBalance = cast.ToFloat64(valStr)
			case "eq":
				acc.Balance = cast.ToFloat64(valStr)
			case "frozenBal":
				acc.FrozenBalance = cast.ToFloat64(valStr)
			}
			return err
		})

		if err != nil {
			return
		}

		accMap[acc.Coin] = acc
	}, "details")

	return accMap, err
}

func (un *RespUnmarshaler) UnmarshalGetFuturesAccountResponse(data []byte) (map[string]model.FuturesAccount, error) {
	var accMap = make(map[string]model.FuturesAccount, 2)

	_, err := jsonparser.ArrayEach(data[1:len(data)-1], func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var acc model.FuturesAccount
		err = jsonparser.ObjectEach(value, func(key []byte, accData []byte, dataType jsonparser.ValueType, offset int) error {
			valStr := string(accData)
			switch string(key) {
			case "ccy":
				acc.Coin = valStr
			case "availEq":
				acc.AvailEq = cast.ToFloat64(valStr)
			case "eq":
				acc.Eq = cast.ToFloat64(valStr)
			case "frozenBal":
				acc.FrozenBal = cast.ToFloat64(valStr)
			case "upl":
				acc.Upl = cast.ToFloat64(valStr)
			case "mgnRatio":
				acc.MgnRatio = cast.ToFloat64(valStr)
			}
			return err
		})

		if err != nil {
			return
		}

		accMap[acc.Coin] = acc
	}, "details")

	return accMap, err
}

func (un *RespUnmarshaler) UnmarshalCancelOrderResponse(data []byte) error {
	sCodeData, _, _, err := jsonparser.Get(data[1:len(data)-1], "sCode")
	if err != nil {
		return err
	}

	if cast.ToInt64(string(sCodeData)) == 0 {
		return nil
	}

	return errors.New(string(data))
}

func (un *RespUnmarshaler) UnmarshalGetPositionsResponse(data []byte) ([]model.FuturesPosition, error) {
	var (
		positions []model.FuturesPosition
		err       error
	)

	_, err = jsonparser.ArrayEach(data, func(posData []byte, dataType jsonparser.ValueType, offset int, err error) {
		var pos model.FuturesPosition
		err = jsonparser.ObjectEach(posData, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			valStr := string(value)
			switch string(key) {
			case "availPos":
				pos.AvailQty = cast.ToFloat64(valStr)
			case "avgPx":
				pos.AvgPx = cast.ToFloat64(valStr)
			case "pos":
				pos.Qty = cast.ToFloat64(valStr)
			case "posSide":
				if valStr == "long" {
					pos.PosSide = model.Futures_OpenBuy
				}
				if valStr == "short" {
					pos.PosSide = model.Futures_OpenSell
				}
			case "upl":
				pos.Upl = cast.ToFloat64(valStr)
			case "uplRatio":
				pos.UplRatio = cast.ToFloat64(valStr)
			case "lever":
				pos.Lever = cast.ToFloat64(valStr)
			}
			return nil
		})
		positions = append(positions, pos)
	})

	return positions, err
}

func (un *RespUnmarshaler) UnmarshalGetPositionsHisotoryResponse(data []byte) ([]model.FuturesPositionHistory, error) {
	var (
		positionsHistory []model.FuturesPositionHistory
		err              error
	)

	_, err = jsonparser.ArrayEach(data, func(posData []byte, dataType jsonparser.ValueType, offset int, err error) {
		var posHistory model.FuturesPositionHistory
		_ = jsonparser.ObjectEach(posData, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			valStr := string(value)
			switch string(key) {
			case "direction":
				posHistory.Direction = cast.ToString(valStr)
			case "type":
				posHistory.Type = cast.ToInt8(valStr)
			case "cTime":
				posHistory.CTime = cast.ToInt32(valStr)
			case "uTime":
				posHistory.UTime = cast.ToInt32(valStr)
			case "realizedPnl":
				posHistory.RealizedPnl = cast.ToFloat64(valStr)
			}
			return nil
		})
		positionsHistory = append(positionsHistory, posHistory)
	})
	return positionsHistory, err
}

func (un *RespUnmarshaler) UnmarshalGetExchangeInfoResponse(data []byte) (map[string]model.CurrencyPair, error) {
	var (
		err             error
		currencyPairMap = make(map[string]model.CurrencyPair, 20)
	)

	_, err = jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		var (
			currencyPair model.CurrencyPair
			instTy       string
			ctValCcy     string
			settleCcy    string
		)

		err = jsonparser.ObjectEach(value, func(key []byte, val []byte, dataType jsonparser.ValueType, offset int) error {
			valStr := string(val)
			switch string(key) {
			case "instType":
				instTy = valStr
			case "instId":
				currencyPair.Symbol = valStr
			case "minSz":
				currencyPair.MinQty = cast.ToFloat64(valStr)
			case "tickSz":
				currencyPair.PricePrecision = AdaptQtyOrPricePrecision(valStr)
			case "lotSz":
				currencyPair.QtyPrecision = AdaptQtyOrPricePrecision(valStr)
			case "baseCcy":
				currencyPair.BaseSymbol = valStr
			case "quoteCcy":
				currencyPair.QuoteSymbol = valStr
			case "ctValCcy":
				ctValCcy = valStr
				currencyPair.ContractValCurrency = valStr
			case "ctVal":
				currencyPair.ContractVal = cast.ToFloat64(valStr)
			case "settleCcy":
				settleCcy = valStr
				currencyPair.SettlementCurrency = valStr
			case "alias":
				currencyPair.ContractAlias = valStr
			case "expTime":
				currencyPair.ContractDeliveryDate = cast.ToInt64(valStr)
			}
			return nil
		})

		if instTy == "SWAP" {
			currencyPair.BaseSymbol = ctValCcy
			currencyPair.QuoteSymbol = settleCcy
		}

		//adapt
		if instTy == "FUTURES" {
			currencyPair.BaseSymbol = settleCcy
			currencyPair.QuoteSymbol = ctValCcy
		}

		k := fmt.Sprintf("%s%s%s", currencyPair.BaseSymbol, currencyPair.QuoteSymbol, currencyPair.ContractAlias)
		currencyPairMap[k] = currencyPair
	})

	return currencyPairMap, err
}

func (un *RespUnmarshaler) UnmarshalResponse(data []byte, res interface{}) error {
	return json.Unmarshal(data, res)
}

func (un *RespUnmarshaler) UnmarshalGetComputeMinInvestmentResponse(data []byte) (model.ComputeMinInvestmentResponse, error) {
	var minInvestment = new(model.ComputeMinInvestmentResponse)
	var investData = new(model.InvestmentData)

	_, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		err = jsonparser.ObjectEach(value, func(key []byte, respData []byte, dataType jsonparser.ValueType, offset int) error {
			investmentDataStr := string(respData)
			switch string(key) {
			case "minInvestmentData":
				_, _ = jsonparser.ArrayEach(respData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					err = jsonparser.ObjectEach(value, func(key []byte, respData []byte, dataType jsonparser.ValueType, offset int) error {
						newInvestmentDataStr := string(respData)
						switch string(key) {
						case "amt":
							investData.Amt = newInvestmentDataStr
						case "ccy":
							investData.Ccy = newInvestmentDataStr
						}
						return err
					})
					if err != nil {
						return
					}
				})
			case "singleAmt":
				minInvestment.SingleAmt = investmentDataStr
			}
			return err
		})

		if err != nil {
			return
		}
	})
	minInvestment.InvestmentData = append(minInvestment.InvestmentData, *investData)

	return *minInvestment, err
}

func (un *RespUnmarshaler) UnmarshalGetAlgoOrderDetailsResponse(data []byte) (model.GridAlgoOrderDetailsResponse, error) {
	var details = new(model.GridAlgoOrderDetailsResponse)
	var rebateTrans = new(model.RebateTrans)
	var triggerParams = new(model.TriggerParams)

	_, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		err = jsonparser.ObjectEach(value, func(key []byte, respData []byte, dataType jsonparser.ValueType, offset int) error {
			detailsStr := string(respData)
			switch string(key) {
			case "rebateTrans":
				_, _ = jsonparser.ArrayEach(respData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					err = jsonparser.ObjectEach(value, func(key []byte, respData []byte, dataType jsonparser.ValueType, offset int) error {
						detailsStr := string(respData)
						switch string(key) {
						case "rebate":
							rebateTrans.Rebate = detailsStr
						case "ccy":
							rebateTrans.RebateCcy = detailsStr
						}
						return err
					})
					if err != nil {
						return
					}
					details.RebateTrans = append(details.RebateTrans, *rebateTrans)
				})
			case "triggerParams":
				_, _ = jsonparser.ArrayEach(respData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					err = jsonparser.ObjectEach(value, func(key []byte, respData []byte, dataType jsonparser.ValueType, offset int) error {
						detailsStr := string(respData)
						switch string(key) {
						case "triggerAction":
							triggerParams.TriggerAction = detailsStr
						case "triggerStrategy":
							triggerParams.TriggerStrategy = detailsStr
						case "delaySeconds":
							triggerParams.DelaySeconds = detailsStr
						case "triggerTime":
							triggerParams.TriggerTime = detailsStr
						case "triggerType":
							triggerParams.TriggerType = detailsStr
						case "timeframe":
							triggerParams.Timeframe = detailsStr
						case "thold":
							triggerParams.Thold = detailsStr
						case "triggerCond":
							triggerParams.TriggerCond = detailsStr
						case "timePeriod":
							triggerParams.TimePeriod = detailsStr
						case "triggerPx":
							triggerParams.TriggerPx = detailsStr
						case "stopType":
							triggerParams.StopType = detailsStr
						}
						return err
					})
					if err != nil {
						return
					}
					details.TriggerParams = append(details.TriggerParams, *triggerParams)
				})
			case "algoId":
				details.AlgoId = detailsStr
			case "algoClOrdId":
				details.AlgoClOrdId = detailsStr
			case "instType":
				details.InstType = detailsStr
			case "instId":
				details.InstId = detailsStr
			case "cTime":
				details.CTime = detailsStr
			case "uTime":
				details.UTime = detailsStr
			case "algoOrdType":
				details.AlgoOrdType = detailsStr
			case "state":
				details.State = detailsStr
			case "maxPx":
				details.MaxPx = detailsStr
			case "minPx":
				details.MinPx = detailsStr
			case "gridNum":
				details.GridNum = detailsStr
			case "runType":
				details.RunType = detailsStr
			case "tpTriggerPx":
				details.TpTriggerPx = detailsStr
			case "slTriggerPx":
				details.SlTriggerPx = detailsStr
			case "tradeNum":
				details.TradeNum = detailsStr
			case "arbitrageNum":
				details.ArbitrageNum = detailsStr
			case "singleAmt":
				details.SingleAmt = detailsStr
			case "perMinProfitRate":
				details.PerMinProfitRate = detailsStr
			case "perMaxProfitRate":
				details.PerMaxProfitRate = detailsStr
			case "runPx":
				details.RunPx = detailsStr
			case "totalPnl":
				details.TotalPnl = detailsStr
			case "pnlRatio":
				details.PnlRatio = detailsStr
			case "investment":
				details.Investment = detailsStr
			case "gridProfit":
				details.GridProfit = detailsStr
			case "floatProfit":
				details.FloatProfit = detailsStr
			case "totalAnnualizedRate":
				details.TotalAnnualizedRate = detailsStr
			case "annualizedRate":
				details.AnnualizedRate = detailsStr
			case "cancelType":
				details.CancelType = detailsStr
			case "stopType":
				details.StopType = detailsStr
			case "activeOrdNum":
				details.ActiveOrdNum = detailsStr
			case "quoteSz":
				details.QuoteSz = detailsStr
			case "baseSz":
				details.BaseSz = detailsStr
			case "curQuoteSz":
				details.CurQuoteSz = detailsStr
			case "curBaseSz":
				details.CurBaseSz = detailsStr
			case "profit":
				details.Profit = detailsStr
			case "stopResult":
				details.StopResult = detailsStr
			case "direction":
				details.Direction = detailsStr
			case "basePos":
				details.BasePos = detailsStr
			case "sz":
				details.Sz = detailsStr
			case "lever":
				details.Lever = detailsStr
			case "actualLever":
				details.ActualLever = detailsStr
			case "liqPx":
				details.LiqPx = detailsStr
			case "uly":
				details.Uly = detailsStr
			case "instFamily":
				details.InstFamily = detailsStr
			case "ordFrozen":
				details.OrdFrozen = detailsStr
			case "availEq":
				details.AvailEq = detailsStr
			case "eq":
				details.Eq = detailsStr
			case "tag":
				details.Tag = detailsStr
			case "profitSharingRatio":
				details.ProfitSharingRatio = detailsStr
			case "copyType":
				details.CopyType = detailsStr
			}
			return err
		})

		if err != nil {
			return
		}
	})

	return *details, err
}

func (un *RespUnmarshaler) UnmarshalPostPlaceGridAlgoOrder(data []byte) (model.PlaceGridAlgoOrderResponse, error) {
	var details = new(model.PlaceGridAlgoOrderResponse)

	_, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		err = jsonparser.ObjectEach(value, func(key []byte, respData []byte, dataType jsonparser.ValueType, offset int) error {
			detailsStr := string(respData)
			switch string(key) {
			case "algoId":
				details.AlgoId = detailsStr
			case "algoClOrdId":
				details.AlgoClOrdId = detailsStr
			case "sCode":
				details.SCode = detailsStr
			case "sMsg":
				details.SMsg = detailsStr
			case "tag":
				details.Tag = detailsStr
			}
			return err
		})

		if err != nil {
			return
		}
	})

	return *details, err
}

func (un *RespUnmarshaler) UnmarshalPlaceOrder(respPlaceOrderData []byte) (model.PlaceOrderResponse, error) {
	var placeOrderResponse = new(model.PlaceOrderResponse)
	var placeOrderResponseData = new(model.PlaceOrderResponseData)
	logger.Info("********* respPlaceOrderData ", string(respPlaceOrderData))

	_, err := jsonparser.ArrayEach(respPlaceOrderData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		err = jsonparser.ObjectEach(value, func(key []byte, respData []byte, dataType jsonparser.ValueType, offset int) error {
			detailsStr := string(respData)
			logger.Info("!!!!!    key ", string(key), " = ", detailsStr)
			switch string(key) {
			case "code":
				placeOrderResponse.Code = detailsStr
			case "msg":
				placeOrderResponse.Msg = detailsStr
			case "inTime":
				placeOrderResponse.InTime = detailsStr
			case "outTime":
				placeOrderResponse.OutTime = detailsStr
			case "data":
				_, _ = jsonparser.ArrayEach(respData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					err = jsonparser.ObjectEach(value, func(key []byte, respData []byte, dataType jsonparser.ValueType, offset int) error {
						detailsStr := string(respData)
						switch string(key) {
						case "ordId":
							placeOrderResponseData.OrdId = detailsStr
						case "clOrdId":
							placeOrderResponseData.ClOrdId = detailsStr
						case "tag":
							placeOrderResponseData.Tag = detailsStr
						case "sCode":
							placeOrderResponseData.SCode = detailsStr
						case "sMsg":
							logger.Info("++++++++++ sMsg - ", detailsStr)
							placeOrderResponseData.SMsg = detailsStr
						}
						return err
					})
					if err != nil {
						return
					}
					placeOrderResponse.Data = append(placeOrderResponse.Data, *placeOrderResponseData)
				})
			}
			return err
		})

		if err != nil {
			return
		}
	})

	logger.Info("@@@@@@@@ placeOrderResponse ", *placeOrderResponse)
	return *placeOrderResponse, err
}

func (un *RespUnmarshaler) UnmarshalPostStopGridAlgoOrder(data []byte) (model.StopGridAlgoOrderResponse, error) {
	var details = new(model.StopGridAlgoOrderResponse)

	_, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		err = jsonparser.ObjectEach(value, func(key []byte, respData []byte, dataType jsonparser.ValueType, offset int) error {
			detailsStr := string(respData)
			switch string(key) {
			case "algoId":
				details.AlgoId = detailsStr
			case "algoClOrdId":
				details.AlgoClOrdId = detailsStr
			case "sCode":
				details.SCode = detailsStr
			case "sMsg":
				details.SMsg = detailsStr
			case "tag":
				details.Tag = detailsStr
			}
			return err
		})

		if err != nil {
			return
		}
	})

	return *details, err
}
