package main

import (
	"fmt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"okx-bot/exchange/model"
	"okx-bot/exchange/okx"
	"okx-bot/exchange/options"
	"okx-bot/signalview"
	"strconv"
	"strings"
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

	//btcUSDTCurrencyPair, err := OKx.Spot.NewCurrencyPair(model.BTC, model.USDT)
	//if err != nil {
	//	panic(err)
	//}

	okxPrvApi := OKx.Spot.NewPrvApi(
		options.WithApiKey(envParams["okx_api_key"]),
		options.WithApiSecretKey(envParams["okx_api_secret_key"]),
		options.WithPassphrase(envParams["okx_api_passphrase"]))

	//order, _, err := okxPrvApi.GetTicker(btcUSDTCurrencyPair)
	//log.Println(err)
	//log.Println(order)

	order2, _, err := okxPrvApi.GetSpotHistoryOrders()
	log.Info(err)
	log.Info(order2)

	minInvestRequest := new(model.ComputeMinInvestmentRequest)
	minInvestRequest.InstId = "BTC-USDT-SWAP"
	minInvestRequest.AlgoOrdType = "contract_grid"
	minInvestRequest.MaxPx = "37000"
	minInvestRequest.MinPx = "34000"
	minInvestRequest.GridNum = "40" //strconv.FormatUint(20, 10)
	minInvestRequest.RunType = "1"
	minInvestRequest.Direction = "long"
	minInvestRequest.Lever = "10"

	var investData = new(model.InvestmentData)
	investData.Amt = "100"
	investData.Ccy = "USDT"
	minInvestRequest.InvestmentData = append(minInvestRequest.InvestmentData, *investData)

	calcGridMinInvestment, respBody, err := OKx.Grid.GetCompMinInvest(*minInvestRequest)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	log.Info(string(respBody))
	log.Info("calcGridMinInvestment", calcGridMinInvestment.SingleAmt)
	resp := calcGridMinInvestment.InvestmentData
	log.Info("calcGridMinInvestment: ", resp[0].Amt)
	log.Info("calcGridMinInvestment: ", resp[0].Ccy)

	gridAlgoOrderDetailsRequest := new(model.GridAlgoOrderDetailsRequest)
	gridAlgoOrderDetailsRequest.AlgoOrdType = "contract_grid"
	gridAlgoOrderDetailsRequest.AlgoId = "642028702938959872"

	gridAlgoOrderDetailsResponse, respBody, err := okxPrvApi.GetGridAlgoOrderDetails(*gridAlgoOrderDetailsRequest)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	log.Info(string(respBody))
	log.Info("gridAlgoOrderDetailsResponse: ", gridAlgoOrderDetailsResponse)

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
