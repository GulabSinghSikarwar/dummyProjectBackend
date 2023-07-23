package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Stock struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty"`
	Ticker               string             `bson:"ticker"`
	Name                 string             `bson:"name"`
	ExchangeShort        interface{}        `bson:"exchange_short"`
	ExchangeLong         interface{}        `bson:"exchange_long"`
	MICCode              string             `bson:"mic_code"`
	Currency             string             `bson:"currency"`
	Price                float32            `bson:"price"`
	DayHigh              float32            `bson:"day_high"`
	DayLow               float32            `bson:"day_low"`
	DayOpen              float32            `bson:"day_open"`
	Week52High           float32            `bson:"52_week_high"`
	Week52Low            float32            `bson:"52_week_low"`
	MarketCap            interface{}        `bson:"market_cap"`
	PreviousClosePrice   float32            `bson:"previous_close_price"`
	PreviousCloseTime    time.Time          `bson:"previous_close_price_time"`
	DayChange            float32            `bson:"day_change"`
	Volume               int32              `bson:"volume"`
	IsExtendedHoursPrice bool               `bson:"is_extended_hours_price"`
	LastTradeTime        time.Time          `bson:"last_trade_time"`
}
