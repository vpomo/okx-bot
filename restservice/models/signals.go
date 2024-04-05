package models

import (
	"github.com/aidarkhanov/nanoid"
	"github.com/jinzhu/gorm"
	u "okx-bot/restservice/utils"
)

type TradingViewSignal struct {
	gorm.Model
	IdOrder                string `json:"idOrder"`
	Action                 string `json:"action"`
	MarketPosition         string `json:"marketPosition" sql:"-"`
	PrevMarketPosition     string `json:"prevMarketPosition"`
	MarketPositionSize     string `json:"marketPositionSize"`
	PrevMarketPositionSize string `json:"prevMarketPositionSize"`
	Instrument             string `json:"instrument"`
	SignalToken            string `json:"signalToken"`
	Timestamp              string `json:"timestamp"`
	Amount                 string `json:"amount"`
}

type SignalObject struct {
	gorm.Model
	IdSignal     string `json:"idSignal"`
	NameToken    string `json:"nameToken"`
	TimeInterval string `json:"timeInterval"`
}

func (tradingViewSignal *TradingViewSignal) Save() map[string]interface{} {

	GetDB().Create(tradingViewSignal)

	response := u.Message(true, "TradingViewSignal has been saved")
	response["tradingViewSignal"] = tradingViewSignal
	return response
}

func (signalObject *SignalObject) Create(nameToken string, interval string) map[string]interface{} {
	signalObject.IdSignal = nanoid.New()
	signalObject.NameToken = nameToken
	signalObject.TimeInterval = interval

	GetDB().Create(signalObject)

	response := u.Message(true, "SignalObject has been created")
	response["signalObject"] = signalObject
	return response
}
