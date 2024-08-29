package models

import (
	"time"
)

type ModelChosenPairsTrade struct {
	Ticker1        string    `gorm:"column:ticker_1;type:varchar(10);not null" json:"ticker_1"`
	Ticker2        string    `gorm:"column:ticker_2;type:varchar(10);not null" json:"ticker_2"`
	Spread         float64   `json:"spread"`
	PredictedTrend float64   `json:"predicted_trend"`
	Signal         string    `json:"signal"`
	Position       float64   `json:"position"`
	Budget         float64   `json:"budget"`
	Date           time.Time `json:"date"`
}

func (ModelChosenPairsTrade) TableName() string {
	return "gold.model_trades"
}
