package models

import "github.com/jinzhu/gorm"

type Signal struct {
	gorm.Model
	Id                     string `json:"id"`
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
