package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Stock struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty"`
	Ticker             string             `bson:"ticker"`
	Name               string             `bson:"name"`
	ExchangeShort      interface{}        `bson:"exchange_short"`
	ExchangeLong       interface{}        `bson:"exchange_long"`
	MICCode            string             `bson:"mic_code"`
	Currency           string             `bson:"currency"`
	Price              float64            `bson:"price"`
	DayHigh            float64            `bson:"day_high"`
	DayLow             float64            `bson:"day_low"`
	DayOpen            float64            `bson:"day_open"`
	Week52High         float64            `bson:"52_week_high"`
	Week52Low          float64            `bson:"52_week_low"`
	MarketCap          interface{}        `bson:"market_cap"`
	PreviousClosePrice float64            `bson:"previous_close_price"`
	// PreviousClosePriceTime time.Time         `bson:"previous_close_price_time" time_format:"2006-01-02T15:04:05.000000"`
	DayChange            float64 `bson:"day_change"`
	Volume               int64   `bson:"volume"`
	IsExtendedHoursPrice bool    `bson:"is_extended_hours_price"`
	// LastTradeTime        time.Time `bson:"last_trade_time"`
}

// 253
