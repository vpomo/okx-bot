package main

import (
	"fmt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"okx-bot/exchange/model"
	"okx-bot/exchange/okx"
	"okx-bot/exchange/options"
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

	envParams := make(map[string]string)
	envParams, err := godotenv.Read()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var ta signalview.TradingView

	//err := ta.Get("BINANCE:BTCUSDT", signalview.Interval15min)
	//err = ta.Get("OKX:XRPUSDT", signalview.Interval15min)
	//err = ta.Get("OKX:TRBUSDT", signalview.Interval15min)
	err = ta.Get("OKX:STARLUSDT", signalview.Interval15min)
	if err != nil {
		log.Error(err)
	}
	log.Info(ta) // Full Data

	// Get the value by key
	recSummary := ta.Recommend.Summary // Summary
	log.Info("recSummary: ", recSummary)

	recOsc := ta.Recommend.Oscillators // Oscillators
	log.Info("recOsc: ", recOsc)

	recMA := ta.Recommend.MA // Moving Averages
	log.Info("recMA: ", recMA)

	switch recSummary {
	case signalview.SignalStrongSell:
		fmt.Println("STRONG_SELL")
	case signalview.SignalSell:
		fmt.Println("SELL")
	case signalview.SignalNeutral:
		fmt.Println("NEUTRAL")
	case signalview.SignalBuy:
		fmt.Println("BUY")
	case signalview.SignalStrongBuy:
		fmt.Println("STRONG_BUY")
	default:
		fmt.Println("An error has occurred")
	}

	OKx := okx.New()

	okx.DefaultHttpCli.SetTimeout(5)

	_, _, err = OKx.Spot.GetExchangeInfo()
	if err != nil {
		panic(err)
	}

	btcUSDTCurrencyPair, err := OKx.Spot.NewCurrencyPair(model.BTC, model.USDT)
	if err != nil {
		panic(err)
	}

	okxPrvApi := OKx.Spot.NewPrvApi(
		options.WithApiKey(envParams["okx_api_key"]),
		options.WithApiSecretKey(envParams["okx_api_secret_key"]),
		options.WithPassphrase(envParams["okx_api_passphrase"]))

	order, _, err := okxPrvApi.GetTicker(btcUSDTCurrencyPair)
	log.Println(err)
	log.Println(order)

	order2, _, err := okxPrvApi.GetSpotHistoryOrders()
	log.Println(err)
	log.Println(order2)

	log.Info("\n")
	log.Info("======================================")
	log.Info("finished")
}
