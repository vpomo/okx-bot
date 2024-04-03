package models

import (
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

func (tradingViewSignal *TradingViewSignal) Save() map[string]interface{} {

	GetDB().Create(tradingViewSignal)

	response := u.Message(true, "TradingViewSignal has been saved")
	response["tradingViewSignal"] = tradingViewSignal
	return response
}
