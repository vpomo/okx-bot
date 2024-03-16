package main

import (
	"fmt"
	"okx-bot/exchange/model"
	"okx-bot/exchange/okx"
	"okx-bot/exchange/options"
	"okx-bot/signalview"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	log.Info("started ...")
	log.Info("======================================")
	log.Info("\n")

	envParams := make(map[string]string)
	envParams, err := godotenv.Read()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var ta signalview.TradingView

	//getBuySellInfo(ta, "OKX:BTCUSDT")
	//getBuySellInfo(ta, "OKX:XRPUSDT")
	//getBuySellInfo(ta, "OKX:OKBUSDT")
	//getBuySellInfo(ta, "OKX:TONUSDT")
	//getBuySellInfo(ta, "OKX:TRBUSDT")
	//getBuySellInfo(ta, "OKX:PEPEUSDT")
	getBuySellInfo(ta, "OKX:STARLUSDT")

	//recOsc := ta.Recommend.Oscillators // Oscillators
	//log.Info("recOsc: ", recOsc)
	//
	//recMA := ta.Recommend.MA // Moving Averages
	//log.Info("recMA: ", recMA)

	OKx := okx.New()

	okx.DefaultHttpCli.SetTimeout(5)

	_, _, err = OKx.Spot.GetExchangeInfo()
	if err != nil {
		panic(err)
	}

	//okxPrvApi := OKx.Spot.NewPrvApi(
	//	options.WithApiKey(envParams["okx_api_key"]),
	//	options.WithApiSecretKey(envParams["okx_api_secret_key"]),
	//	options.WithPassphrase(envParams["okx_api_passphrase"]))

	//okxSwapPrvApi := OKx.Futures.NewPrvApi(
	//	options.WithApiKey(envParams["okx_api_key"]),
	//	options.WithApiSecretKey(envParams["okx_api_secret_key"]),
	//	options.WithPassphrase(envParams["okx_api_passphrase"]))

	//opts := model.OptionParameter{
	//	Key:   "contractAlias",
	//	Value: "SWAP",
	//}

	//btcUSDTCurrencyPair, err := OKx.Grid.NewCurrencyPair(model.BTC, model.USDT, opts)
	//if err != nil {
	//	panic(err)
	//}
	//order, _, err := okxSwapPrvApi.GetTicker(btcUSDTCurrencyPair)
	//log.Println(err)
	//log.Println(order)

	//order2, _, err := okxSwapPrvApi.GetHistoryOrders(btcUSDTCurrencyPair, opts)
	//fmt.Println("error: ", err)
	//log.Info("order2: ", order2)

	orderRequest := new(model.PlaceOrderRequest)
	orderRequest.InstId = "SOL-USDT"
	orderRequest.TdMode = "cash"
	orderRequest.Side = "buy"
	orderRequest.OrdType = "limit"
	orderRequest.Sz = "1"
	orderRequest.Px = "1"

	okxSwapPrvApi := OKx.Futures.NewPrvApi(
		options.WithApiKey(envParams["okx_api_key"]),
		options.WithApiSecretKey(envParams["okx_api_secret_key"]),
		options.WithPassphrase(envParams["okx_api_passphrase"]))

	newOrder, _, err := okxSwapPrvApi.Cross.PlaceOrder(*orderRequest)
	if err != nil {
		panic(err)
	}
	log.Info("ordId = ", newOrder.Id)

	amendOrder := new(model.AmendOrderRequest)

	amendOrder.OrdId = newOrder.Id
	amendOrder.InstId = "SOL-USDT"
	amendOrder.NewSz = "9"

	order, _, err := okxSwapPrvApi.Cross.AmendOrder(*amendOrder)
	if err != nil {
		panic(err)
	}
	log.Info("ordId = ", order.Data)

	// posHistoryRequest := new(model.FuturesPositionHistoryRequest)
	// posHistoryRequest.InstId = "BTC-USDT-SWAP"
	// posHistory, _, err := okxSwapPrvApi.GetPositionsHistory(*posHistoryRequest)
	// if err != nil {
	// 	panic(err)
	// }
	// log.Info("posHistory = ", posHistory)
	// log.Info("posHistory.InstId = ", posHistory[0].InstId)
	// log.Info("posHistory.Direction = ", posHistory[0].Direction)
	// log.Info("posHistory.Lever = ", posHistory[0].Lever)
	// log.Info("posHistory.CTime = ", posHistory[0].CTime)
	// log.Info("posHistory.UTime = ", posHistory[0].UTime)
	// log.Info("posHistory.OpenAvgPx = ", posHistory[0].OpenAvgPx)
	// log.Info("posHistory.CloseAvgPx = ", posHistory[0].CloseAvgPx)
	// log.Info("posHistory.Pnl = ", posHistory[0].Pnl)
	// log.Info("posHistory.RealizedPnl = ", posHistory[0].RealizedPnl)

	//minInvestRequest := new(model.ComputeMinInvestmentRequest)
	//minInvestRequest.InstId = "BTC-USDT-SWAP"
	//minInvestRequest.AlgoOrdType = "contract_grid"
	//minInvestRequest.MaxPx = "37000"
	//minInvestRequest.MinPx = "34000"
	//minInvestRequest.GridNum = "40" //strconv.FormatUint(20, 10)
	//minInvestRequest.RunType = "1"
	//minInvestRequest.Direction = "long"
	//minInvestRequest.Lever = "10"
	//
	//var investData = new(model.InvestmentData)
	//investData.Amt = "100"
	//investData.Ccy = "USDT"
	//minInvestRequest.InvestmentData = append(minInvestRequest.InvestmentData, *investData)
	//
	//calcGridMinInvestment, respBody, err := OKx.Grid.GetCompMinInvest(*minInvestRequest)
	//if err != nil {
	//	log.Error(err)
	//	panic(err)
	//}
	//log.Info(string(respBody))
	//log.Info("calcGridMinInvestment", calcGridMinInvestment.SingleAmt)
	//resp := calcGridMinInvestment.InvestmentData
	//log.Info("calcGridMinInvestment: ", resp[0].Amt)
	//log.Info("calcGridMinInvestment: ", resp[0].Ccy)

	//gridAlgoOrderDetailsRequest := new(model.GridAlgoOrderDetailsRequest)
	//gridAlgoOrderDetailsRequest.AlgoOrdType = "contract_grid"
	//gridAlgoOrderDetailsRequest.AlgoId = "665614704722841600"
	//
	//gridAlgoOrderDetailsResponse, respBody, err := okxPrvApi.GetGridAlgoOrderDetails(*gridAlgoOrderDetailsRequest)
	//if err != nil {
	//	log.Error(err)
	//	panic(err)
	//}
	//log.Info("gridAlgoOrderDetailsResponse ", gridAlgoOrderDetailsResponse)

	//newGridOrder := new(model.PlaceGridAlgoOrderRequest)
	//newGridOrder.GridNum = "20"
	//newGridOrder.MinPx = "10"
	//newGridOrder.MaxPx = "12"
	//newGridOrder.InstId = "ORDI-USDT-SWAP"
	//newGridOrder.AlgoOrdType = "contract_grid"
	//newGridOrder.Lever = "5"
	//newGridOrder.Direction = "long"
	//newGridOrder.Sz = "8"
	//
	//placeGridAlgoOrderResponse, respBody, err := okxPrvApi.PlaceGridAlgoOrder(*newGridOrder)
	//algoId := placeGridAlgoOrderResponse.AlgoId
	//log.Info("respBody ", string(respBody))
	//if err != nil {
	//	log.Error(err)
	//}
	//log.Info("placeGridAlgoOrderResponse ", placeGridAlgoOrderResponse)
	//log.Info("algoId ", algoId)

	//algoId := "665614704722841600"
	//stopGridOrder := new(model.StopGridAlgoOrderRequest)
	//
	//stopGridOrder.AlgoId = algoId
	//stopGridOrder.InstId = "ORDI-USDT-SWAP"
	//stopGridOrder.AlgoOrdType = "contract_grid"
	//stopGridOrder.StopType = "1"
	//
	//stopGridAlgoOrderResponse, respBody, err := okxPrvApi.StopGridAlgoOrder(*stopGridOrder)
	//log.Info("respBody ", string(respBody))
	//if err != nil {
	//	log.Error(err)
	//}
	//log.Info("stopGridAlgoOrderResponse ", stopGridAlgoOrderResponse)

	log.Info("\n")
	log.Info("======================================")
	log.Info("finished")
}

func getBuySellInfo(ta signalview.TradingView, namePair string) {
	err := ta.Get("OKX:XRPUSDT", signalview.Interval15min)
	if err != nil {
		log.Error(err)
	}
	ma := ta.Recommend.MA              // Moving Averages
	mACDOsc := ta.Oscillators.MACD     // Oscillators
	mRSIOsc := ta.Oscillators.StochRSI // Oscillators

	var summaryResult = ""

	recSummary := ta.Recommend.Summary
	switch recSummary {
	case signalview.SignalStrongSell:
		summaryResult = "STRONG_SELL"
	case signalview.SignalSell:
		summaryResult = "SELL"
	case signalview.SignalNeutral:
		summaryResult = "NEUTRAL"
	case signalview.SignalBuy:
		summaryResult = "BUY"
	case signalview.SignalStrongBuy:
		summaryResult = "STRONG_BUY"
	default:
		fmt.Println("An error has occurred")
	}
	log.Info(namePair, " - ", summaryResult, " ma: ", ma, " macdOsc: ", mACDOsc, " mRSIOsc: ", mRSIOsc)
}

func convert(b []byte) string {
	s := make([]string, len(b))
	for i := range b {
		s[i] = strconv.Itoa(int(b[i]))
	}
	return strings.Join(s, ",")
}

func placeOrder() {

}
