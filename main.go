package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"okx-bot/signalview"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	log.Info("started ...")
	log.Info("======================================")
	log.Info("\n")

	//envParams := make(map[string]string)
	//envParams, err := godotenv.Read()

	//if err != nil {
	//	log.Fatal("Error loading .env file")
	//}

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

	//OKx := okx.New()
	//
	//okx.DefaultHttpCli.SetTimeout(5)
	//
	//_, _, err = OKx.Spot.GetExchangeInfo()
	//if err != nil {
	//	panic(err)
	//}
	//
	//btcUSDTCurrencyPair, err := OKx.Spot.NewCurrencyPair(model.BTC, model.USDT)
	//if err != nil {
	//	panic(err)
	//}
	//
	//okxPrvApi := OKx.Spot.NewPrvApi(
	//	options.WithApiKey(envParams["okx_api_key"]),
	//	options.WithApiSecretKey(envParams["okx_api_secret_key"]),
	//	options.WithPassphrase(envParams["okx_api_passphrase"]))
	//
	//order, _, err := okxPrvApi.GetTicker(btcUSDTCurrencyPair)
	//log.Println(err)
	//log.Println(order)
	//
	//order2, _, err := okxPrvApi.GetSpotHistoryOrders()
	//log.Println(err)
	//log.Println(order2)

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
